package server

import (
	"bytes"
	"encoding/json"
	"os"
)

type PrettyJSONWriter struct {
	w *os.File
}

func (pjw *PrettyJSONWriter) Write(data []byte) (n int, err error) {
	var buf bytes.Buffer
	err = json.Indent(&buf, data, "", "  ") // Indent with 2 spaces
	if err != nil {
		return 0, err
	}
	buf.WriteByte('\n') // Add a newline after each JSON object
	return pjw.w.Write(buf.Bytes())
}
