package ilua

func init() {
	// set std logger
	SetLogger(&stdLog{})
}
