package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HTTPInfo holds information regarding the HTTP request/response
type HTTPInfo struct {
	Method     string
	Version    string
	ReturnCode uint
}

// LogEntry holds all the information contained in a Common Log Format line
type LogEntry struct {
	ClientIP       string
	UserIdentifier string
	UserID         string
	Timestamp      time.Time
	HTTP           HTTPInfo
	URL            string
	SizeBytes      uint64
}

// Log line regexp. NOTE: capture group names are there for readability
var logFormat = regexp.MustCompile("^(?P<ClientIP>[^ ]+) " +
	"(?P<UserIdentifier>[^ ]+) " +
	"(?P<UserID>[^ ]+) " +
	"\\[(?P<Timestamp>[^ ]+ [^ ]+)\\] " +
	"\"(?P<HTTPMethod>[A-Za-z]+) (?P<URL>/[^ ]*) (?P<HTTPVersion>[^\"]+)\" " +
	"(?P<HTTPReturnCode>[1-5]\\d{2}) " +
	"(?P<SizeBytes>\\d+)")

const logEltCount = 9

const timestampFormat = "02/Jan/2006:15:04:05 -0700"

// NewLogEntry creates a LogEntry from the given line.
// Line is expected to be formatted according to the
// Common Log Format [https://en.wikipedia.org/wiki/Common_Log_Format]
func NewLogEntry(line string) (*LogEntry, error) {
	elts := logFormat.FindStringSubmatch(line)
	// NOTE: elts[0] contains the whole match, that's why we need +1
	if (elts == nil) || (len(elts) != logEltCount+1) {
		return nil, fmt.Errorf("log line has invalid format (line='%s')", line)
	}

	ts, err := time.Parse(timestampFormat, elts[4])
	if err != nil {
		return nil, fmt.Errorf("log line has invalid timestamp format (%s)", err)
	}

	httpRet, _ := strconv.ParseUint(elts[8], 10, 32)
	sizeBytes, _ := strconv.ParseUint(elts[9], 10, 64)

	return &LogEntry{
		ClientIP:       elts[1],
		UserIdentifier: elts[2],
		UserID:         elts[3],
		Timestamp:      ts.UTC(),
		HTTP: HTTPInfo{
			Method:     elts[5],
			Version:    elts[7],
			ReturnCode: uint(httpRet),
		},
		URL:       elts[6],
		SizeBytes: sizeBytes,
	}, nil
}

// URLSection returns the section of the URL
// (ie. URL component between the first and the second '/')
func (l *LogEntry) URLSection() string {
	// NOTE: skip leading '/' (ensured by regexp)
	return strings.SplitN(l.URL[1:], "/", 2)[0]
}
