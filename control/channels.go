package control

type Channels struct {
	Open    chan bool
	Update  chan string
	Started chan bool
}
