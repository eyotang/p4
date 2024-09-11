package p4

func (conn *Conn) ChangeUser(userName, ticket string) {
	conn.env = append(conn.env, "P4USER="+userName)
	conn.env = append(conn.env, "P4PASSWD="+ticket)
}
