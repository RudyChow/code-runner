package common

var (
	RemoveContainerChan chan string = make(chan string, 1000)
	RemoveFileChan      chan string = make(chan string, 1000)
)
