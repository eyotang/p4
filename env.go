package p4

func (conn *Conn) AddEnv(key, value string) {
	conn.env = append(conn.env, key+"="+value)
}
