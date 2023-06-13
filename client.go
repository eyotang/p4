package p4

import "strings"

type Client struct {
	Owner       string `json:"owner"`
	Client      string `json:"client"`
	Root        string `json:"root"`
	Host        string `json:"host"`
	Stream      string `json:"stream"`
	Description string `json:"description"`
}

func (c *Client) String() string {
	return ""
}

// Clients path格式：//Stream_Root (没有后面的/...)
func (conn *Conn) Clients(path string) (clients []*Client, err error) {
	var (
		result []Result
	)
	if err = validateLocation(path); err != nil {
		return
	}
	if result, err = conn.RunMarshaled("clients", []string{"-S", path}); err != nil {
		return
	}
	for idx := range result {
		if client, ok := result[idx].(*Client); !ok {
			continue
		} else {
			clients = append(clients, client)
		}
	}
	return
}

// UnloadedClients path格式：//Stream_Root (没有后面的/...)
func (conn *Conn) UnloadedClients(path string) (clients []*Client, err error) {
	var (
		result []Result
	)
	if err = validateLocation(path); err != nil {
		return
	}
	if result, err = conn.RunMarshaled("clients", []string{"-U", "-S", path}); err != nil {
		return
	}
	for idx := range result {
		if client, ok := result[idx].(*Client); !ok {
			continue
		} else {
			clients = append(clients, client)
		}
	}
	return
}

func (conn *Conn) DeleteClient(name string) (message string, err error) {
	var (
		out []byte
	)
	out, err = conn.Output([]string{"client", "-df", name})
	message = strings.TrimSpace(string(out))
	return
}
