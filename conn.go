package p4

import (
	"time"
)

const defaultTimeout = 15 * time.Second

// ConnOptions Conn is an interface to the Conn command line client.
type ConnOptions struct {
	address  string
	binary   string
	username string
	password string
	client   string
}

type ConnOptionFunc func(*Conn)

type Conn struct {
	ConnOptions
	env     []string
	timeout time.Duration
}

func NewConn(address, username, password string, options ...ConnOptionFunc) (conn *Conn, err error) {
	conn = &Conn{
		ConnOptions: ConnOptions{
			binary:   "p4",
			address:  address,
			username: username,
			password: password,
		},
		timeout: defaultTimeout,
	}
	for _, opt := range options {
		opt(conn)
	}
	if err = conn.Login(); err != nil {
		return
	}
	return
}

func NewClientConn(address, username, password, client string) (conn *Conn, err error) {
	if conn, err = NewConn(address, username, password, WithClient(client)); err != nil {
		return
	}
	return
}

func WithClient(client string) ConnOptionFunc {
	return func(conn *Conn) {
		conn.SetClient(client)
	}
}

func WithTimeout(timeout time.Duration) ConnOptionFunc {
	return func(conn *Conn) {
		if timeout > 0 {
			conn.timeout = timeout
		} else {
			conn.timeout = defaultTimeout
		}
	}
}
