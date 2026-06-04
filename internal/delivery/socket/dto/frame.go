package dto

import (
	"errors"
	"fmt"
	"strings"
)

type Command string

const (
	CmdSend       Command = "SEND"
	CmdSubscribe  Command = "SUBSCRIBE"
	CmdNext       Command = "NEXT"
	CmdDisconnect Command = "DISCONNECT"
	CmdDashboard  Command = "DASHBOARD"

	CmdMessage       Command = "MESSAGE"
	CmdError         Command = "ERROR"
	CmdConnected     Command = "CONNECTED"
	CmdDashboardData Command = "DASHBOARD_DATA"
)

type Frame struct {
	Command Command
	Headers map[string]string
	Body    string
}

func NewFrame(cmd Command, body string, headers map[string]string) *Frame {
	if headers == nil {
		headers = make(map[string]string)
	}
	return &Frame{Command: cmd, Headers: headers, Body: body}
}

func (f *Frame) Header(key string) (string, bool) {
	v, ok := f.Headers[key]
	if !ok || v == "" {
		return "", ok
	}
	return v, ok
}

func (f *Frame) MustHeader(key string) (string, error) {
	v, ok := f.Headers[key]
	if !ok || v == "" {
		return "", fmt.Errorf("missing required header %q", key)
	}
	return v, nil
}

// Encode serialises to wire format:
//
//	COMMAND\n
//	key:value\n
//	\n
//	body\x00
func (f *Frame) Encode() string {
	var sb strings.Builder
	sb.WriteString(string(f.Command))
	sb.WriteByte('\n')
	for k, v := range f.Headers {
		sb.WriteString(k)
		sb.WriteByte(':')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	sb.WriteByte('\n')
	sb.WriteString(f.Body)
	sb.WriteByte('\x00')
	return sb.String()
}

func ParseFrame(raw string) (*Frame, error) {
	raw = strings.TrimRight(raw, "\x00")
	if raw == "" {
		return nil, errors.New("empty frame")
	}

	parts := strings.SplitN(raw, "\n\n", 2)
	headerSection := parts[0]
	body := ""
	if len(parts) == 2 {
		body = parts[1]
	}

	lines := strings.Split(headerSection, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return nil, errors.New("frame missing command")
	}

	cmd := Command(strings.TrimSpace(lines[0]))
	if cmd == "" {
		return nil, errors.New("frame command is empty")
	}

	headers := make(map[string]string)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("malformed header line: %q", line)
		}
		headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	return &Frame{Command: cmd, Headers: headers, Body: body}, nil
}
