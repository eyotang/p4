package p4

import (
	"strings"
)

func (conn *Conn) ChangeUser(userName, ticket string) {
	foundUser := false
	foundPass := false
	for i, v := range conn.env {
		if strings.HasPrefix(v, "P4USER=") {
			conn.env[i] = "P4USER=" + userName
			foundUser = true
		}
		if strings.HasPrefix(v, "P4PASSWD=") {
			conn.env[i] = "P4PASSWD=" + ticket
			foundPass = true
		}
	}

	if !foundUser {
		conn.env = append(conn.env, "P4USER="+userName)
	}
	if !foundPass {
		conn.env = append(conn.env, "P4PASSWD="+ticket)
	}
}

func (conn *Conn) GetUserTicket() (userName, ticket string) {
	for _, v := range conn.env {
		if strings.HasPrefix(v, "P4USER=") {
			userName = strings.TrimPrefix(v, "P4USER=")
		}
		if strings.HasPrefix(v, "P4PASSWD=") {
			ticket = strings.TrimPrefix(v, "P4PASSWD=")
		}
	}
	return
}
