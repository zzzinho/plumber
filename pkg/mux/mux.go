package mux

import "fmt"

// TODO: need new name

type Way struct {
	State func(interface{}) bool
	To    chan<- interface{}
}

type Mux struct {
	inputChan  <-chan interface{}
	ways       []*Way
	numWorkers int
}

func New(inputChan <-chan interface{}, ways []*Way, numWorkers int) *Mux {
	return &Mux{
		inputChan:  inputChan,
		ways:       ways,
		numWorkers: numWorkers,
	}
}

func (m *Mux) Flow() {
	for n := range m.numWorkers {
		go func() {
			fmt.Printf("mux_%d is running\n", n)
			for input := range m.inputChan {
				fmt.Printf("[Mux_%d] input: %d\n", n, input.(int))
				for _, way := range m.ways {
					if way.State(input) {
						way.To <- input
					}
				}
			}
		}()
	}
}
