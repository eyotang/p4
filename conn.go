package p4

import (
	"time"
)

const defaultExpTimeout = 15 * time.Second

var defaultTimeout = &Timeout{
	Login: defaultExpTimeout,
	Read:  defaultExpTimeout,
	Write: defaultExpTimeout,
}

type opType int

const (
	opTypeInvalid opType = iota
	opTypeLogin
	opTypeRead
	opTypeWrite
	opTypeMax
)

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
	timeout *Timeout
}

type Timeout struct {
	Login time.Duration `json:"login"`
	Read  time.Duration `json:"read"`
	Write time.Duration `json:"write"`
}

func (t *Timeout) OpTimeout(opt opType) time.Duration {
	if t == nil {
		return defaultExpTimeout
	}
	switch opt {
	case opTypeLogin:
		return t.Login
	case opTypeRead:
		return t.Read
	case opTypeWrite:
		return t.Write
	default:
		return defaultExpTimeout
	}
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

func WithClient(client string) ConnOptionFunc {
	return func(conn *Conn) {
		conn.SetClient(client)
	}
}

func WithTimeout(timeout *Timeout) ConnOptionFunc {
	return func(conn *Conn) {
		if timeout != nil {
			conn.timeout = timeout
		} else {
			conn.timeout = defaultTimeout
		}
	}
}
