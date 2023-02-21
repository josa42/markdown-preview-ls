package control

type Channels struct {
	Open   chan bool
	Close  chan bool
	Update chan string
}

type PreviewChannels struct {
	Close  chan bool
	Update chan string
}
