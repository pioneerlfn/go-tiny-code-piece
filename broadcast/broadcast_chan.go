package main

import (
	"fmt"
	"sync"
)

// Channel in golang does not broadcast, value in channel can only be consumed once.
// in this file a broadcast service is implemented.
// see https://science.mroman.ch/gobroadcastchannels.html

type BroadcastService struct {
	// This is the channel the service will listen on...
	chBroadcast chan int
	// and forward it to these.
	chListeners []chan int
	// Requests for new listeners to be added...
	chNewRequests chan (chan int)
	// Requests for listeners to be removed...
	chRemoveRequests chan (chan int)
}

// Create a new BroadcastService.
func NewBroadcastService() *BroadcastService {
	return &BroadcastService{
		chBroadcast:      make(chan int),
		chListeners:      make([]chan int, 3),
		chNewRequests:    make(chan (chan int)),
		chRemoveRequests: make(chan (chan int)),
	}
}

// This creates a new listener and returns the channel a goroutine
// should listen on.
func (bs *BroadcastService) Listener() chan int {
	ch := make(chan int)
	bs.chNewRequests <- ch
	return ch
}

// This removes a listener.
func (bs *BroadcastService) RemoveListener(ch chan int) {
	bs.chRemoveRequests <- ch
}

func (bs *BroadcastService) addListener(ch chan int) {
	for i, v := range bs.chListeners {
		if v == nil {
			bs.chListeners[i] = ch
			return
		}
	}

	bs.chListeners = append(bs.chListeners, ch)
}

func (bs *BroadcastService) removeListener(ch chan int) {
	for i, v := range bs.chListeners {
		if v == ch {
			bs.chListeners[i] = nil
			// important to close! otherwise the goroutine listening on it
			// might block forever!
			close(ch)
			return
		}
	}
}

func (bs *BroadcastService) Run() chan int {
	go func() {
		for {
			// process requests for new listeners or removal of listeners
			select {
			case newCh := <-bs.chNewRequests:
				bs.addListener(newCh)
			case removeCh := <-bs.chRemoveRequests:
				bs.removeListener(removeCh)
			case v, ok := <-bs.chBroadcast:
				// terminate everything if the input channel is closed
				if !ok {
					goto terminate
				}

				// forward the value to all channels
				for _, dstCh := range bs.chListeners {
					if dstCh == nil {
						continue
					}

					dstCh <- v
				}
			}
		}

	terminate:

		// close all listeners
		for _, dstCh := range bs.chListeners {
			if dstCh == nil {
				continue
			}

			close(dstCh)
		}
	}()

	return bs.chBroadcast
}

func main() {
	bs := NewBroadcastService()
	chBroadcast := bs.Run()
	chA := bs.Listener()
	chB := bs.Listener()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		for v := range chA {
			fmt.Println("A", v)
		}
		wg.Done()
	}()

	go func() {
		for v := range chB {
			fmt.Println("B", v)
		}
		wg.Done()
	}()

	for i := 0; i < 3; i++ {
		chBroadcast <- i
	}

	bs.RemoveListener(chA)

	for i := 3; i < 6; i++ {
		chBroadcast <- i
	}

	close(chBroadcast)
	wg.Wait()
}
