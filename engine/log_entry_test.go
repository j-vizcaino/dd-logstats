package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLogEntry(t *testing.T) {
	line := "127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 9996662326"

	entry, err := NewLogEntry(line)
	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", entry.ClientIP)
	assert.Equal(t, "user-identifier", entry.UserIdentifier)
	assert.Equal(t, "frank", entry.UserID)
	assert.Equal(t, "Tue, 10 Oct 2000 13:55:36 -0700", entry.Timestamp.Format(time.RFC1123Z))
	assert.Equal(t, HTTPInfo{
		Method:     "GET",
		Version:    "HTTP/1.0",
		ReturnCode: 200,
	}, entry.HTTP)
	assert.Equal(t, "/apache_pb.gif", entry.URL)
	assert.EqualValues(t, 9996662326, entry.SizeBytes)

}

func TestNewLogEntryErrors(t *testing.T) {
	_, err := NewLogEntry("127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0")
	assert.EqualError(t, err, "log line has invalid format (line='127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0')")

	_, err = NewLogEntry("127.0.0.1 user-identifier frank [40/Xxx/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326")
	assert.EqualError(t, err, "log line has invalid timestamp format (parsing time \"40/Xxx/2000:13:55:36 -0700\" as \"02/Jan/2006:15:04:05 -0700\": cannot parse \"Xxx/2000:13:55:36 -0700\" as \"Jan\")")
}

func TestLogEntryURLSection(t *testing.T) {
	l := LogEntry{URL: "/foo/bar/file.png"}
	assert.Equal(t, "foo", l.URLSection())

	l = LogEntry{URL: "/file.png"}
	assert.Equal(t, "file.png", l.URLSection())

	l = LogEntry{URL: "/"}
	assert.Empty(t, l.URLSection())
}
