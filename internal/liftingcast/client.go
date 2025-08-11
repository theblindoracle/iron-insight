package liftingcast

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	heartbeatInterval = 30 * time.Second
)

type Client struct {
	conn     *websocket.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	baseURL  string
	meetID   string
	password string
	apiKey   string

	// Channels
	messages chan []byte
	errors   chan error
}

type Service interface {
	Connect() error
	Start() error
	Stop() error
	Messages() <-chan []byte
	Errors() <-chan error
}

func New(baseURL, meetID, password, apiKey string) Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:      ctx,
		cancel:   cancel,
		baseURL:  baseURL,
		meetID:   meetID,
		password: password,
		apiKey:   apiKey,
		messages: make(chan []byte, 100),
		errors:   make(chan error, 10),
	}
}

// const ws = new WebSocket(
//
//	`${apiBaseUrl}?meetId=${encodeURI(meetId)}&auth=${encodeURI(
//	  btoa(`${meetId}:${password}`)
//	)}&apiKey=${encodeURI(apiKey)}`jjj
//
// );

func (c *Client) Connect() error {
	// Create base64 encoded auth string (meetId:password)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.meetID, c.password)))

	// Parse the base URL
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	// Add query parameters

	query := parsedURL.Query()
	query.Set("apiKey", c.apiKey)
	query.Set("auth", auth)
	query.Set("meetId", c.meetID)
	parsedURL.RawQuery = query.Encode()

	fmt.Println(parsedURL.String())

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	c.conn = conn
	log.Printf("Connected to %s", parsedURL.String())
	return nil
}

func (c *Client) Start() error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	go c.readMessages()
	go c.heartbeat()

	return nil
}

func (c *Client) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}

	return nil
}

func (c *Client) Messages() <-chan []byte {
	return c.messages
}

func (c *Client) Errors() <-chan error {
	return c.errors
}

func (c *Client) readMessages() {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.errors <- err
				return
			}

			select {
			case c.messages <- message:
			case <-c.ctx.Done():
				return
			default:
				log.Printf("Message buffer full, dropping message")
			}
		}
	}
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if c.conn != nil {
				err := c.conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					c.errors <- err
					return
				}
				log.Printf("Sent heartbeat to %s", c.baseURL)
			}
		}
	}
}
