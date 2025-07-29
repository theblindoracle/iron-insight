-- name: GetMeetDataForLifterName :one
select * from opl
where name like $1 || '%'
order by meet_date desc;

