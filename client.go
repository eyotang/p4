package p4

import (
	"bytes"
	"strings"
	"text/template"
)

type Client struct {
	Client        string   `json:"client"`
	Owner         string   `json:"owner"`
	Host          string   `json:"host"`
	Description   string   `json:"description"`
	Root          string   `json:"root"`
	Options       string   `json:"options"`
	SubmitOptions string   `json:"submitOptions"`
	Stream        string   `json:"stream"`
	Type          string   `json:"type"`
	View          []string `json:"view"`
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

var _clientTemplate = template.New("p4 client template")
var _clientTemplateTxt = `Client: {{.Client}}
Owner:  {{.Owner}}
Root:   {{.Root}}
Options:        noallwrite noclobber nocompress unlocked nomodtime normdir
SubmitOptions:  submitunchanged
LineEnd:        local
Stream: {{.Stream}}
Type:   {{.Type}}
View:
{{- range .View }}
        {{.}}
{{- end }}
`

func (conn *Conn) CreateClient(clientInfo Client) (message string, err error) {
	var (
		out        []byte
		contentBuf = bytes.NewBuffer(nil)
	)
	if _, err = _clientTemplate.Parse(_clientTemplateTxt); err != nil {
		return
	}
	if err = _clientTemplate.Execute(contentBuf, clientInfo); err != nil {
		return
	}
	if out, err = conn.Input([]string{"client", "-i"}, contentBuf.Bytes()); err != nil {
		return
	}
	message = strings.TrimSpace(string(out))
	return
}

// https://www.perforce.com/manuals/cmdref/Content/CmdRef/p4_client.html#p4_client
const (
	ClientTypeWriteable   = "writeable" // default
	ClientTypeReadonly    = "readonly"
	ClientTypePartitioned = "partitioned"
)

func (conn *Conn) CreatePartitionClient(clientInfo Client) (message string, err error) {
	clientInfo.Type = ClientTypePartitioned
	return conn.CreateClient(clientInfo)
}

func (conn *Conn) Client(name string) (client *Client, err error) {
	var (
		ok     bool
		result []Result
	)
	if result, err = conn.RunMarshaled("client", []string{"-o", name}); err != nil {
		return
	}
	for idx := range result {
		if client, ok = result[idx].(*Client); ok {
			break
		}
	}
	return
}
