package p4

func (conn *Conn) Print(path string) (content []byte, err error) {
	// -q : suppress header line.
	out, err := conn.Output([]string{"print", "-q", path})
	if err != nil {
		return nil, err
	}
	return out, nil
}
