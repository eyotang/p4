package p4

func (conn *Conn) WriteTriggers(content []byte) (out []byte, err error) {
	return conn.Input([]string{"triggers", "-i"}, content)
}

func (conn *Conn) Triggers() (out []byte, err error) {
	return conn.Output([]string{"triggers", "-o"})
}
