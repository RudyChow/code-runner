package channels

var RemoveContainerChan chan string = make(chan string, 1000)
var RemoveFileChan chan string = make(chan string, 1000)
