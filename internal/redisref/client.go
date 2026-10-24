package redisref

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	connection net.Conn
	reader     *bufio.Reader
}

func Dial(address string) (*Client, error) {
	connection, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{connection: connection, reader: bufio.NewReader(connection)}, nil
}

func (c *Client) Close() error { return c.connection.Close() }

func (c *Client) Command(args ...string) (any, error) {
	if err := c.connection.SetDeadline(time.Now().Add(time.Minute)); err != nil {
		return nil, err
	}
	var request strings.Builder
	fmt.Fprintf(&request, "*%d\r\n", len(args))
	for _, arg := range args {
		fmt.Fprintf(&request, "$%d\r\n%s\r\n", len(arg), arg)
	}
	if _, err := io.WriteString(c.connection, request.String()); err != nil {
		return nil, err
	}
	return c.read()
}

func (c *Client) read() (any, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 3 || !strings.HasSuffix(line, "\r\n") {
		return nil, fmt.Errorf("invalid RESP line %q", line)
	}
	payload := line[1 : len(line)-2]
	switch line[0] {
	case '+':
		return payload, nil
	case '-':
		return nil, fmt.Errorf("redis: %s", payload)
	case ':':
		return strconv.ParseInt(payload, 10, 64)
	case '$', '*':
		n, err := strconv.Atoi(payload)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid RESP length %q", payload)
		}
		if line[0] == '$' {
			data := make([]byte, n+2)
			if _, err := io.ReadFull(c.reader, data); err != nil {
				return nil, err
			}
			if string(data[n:]) != "\r\n" {
				return nil, fmt.Errorf("invalid bulk terminator")
			}
			return string(data[:n]), nil
		}
		values := make([]any, n)
		for i := range values {
			values[i], err = c.read()
			if err != nil {
				return nil, err
			}
		}
		return values, nil
	default:
		return nil, fmt.Errorf("unsupported RESP marker %q", line[0])
	}
}
