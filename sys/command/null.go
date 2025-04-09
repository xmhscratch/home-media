package command

func NewNullWriter() *NullWriter {
	return &NullWriter{}
}

func (ctx NullWriter) Read(p []byte) (int, error) {
	return len(p), nil
}

func (ctx NullWriter) Write(p []byte) (int, error) {
	// fmt.Println(string(p))
	return len(p), nil
}

func (ctx NullWriter) Close() error {
	return nil
}
