package main

import (
	"log"
	"plumber/pkg/mux"
	"plumber/pkg/pipe"
	"time"
)

// TODO: add dataplane
// TODO: add rest component

// NOTE: test run
func main() {
	log.Println("start plumber")
	cfg1 := pipe.Config{
		Name:                    "first pipe",
		InputChannelBufferSize:  10,
		OutputChannelBufferSize: 10,
		NumWorkers:              5,
		Task: func(input interface{}) interface{} {
			in := input.(int)
			return in * in
		},
	}
	p1 := pipe.New(cfg1)
	go p1.Flow()

	cfg2 := pipe.Config{
		Name:                    "second pipe",
		InputChannelBufferSize:  10,
		OutputChannelBufferSize: 10,
		NumWorkers:              2,
		Task: func(input interface{}) interface{} {
			in := input.(int)
			return in + in
		},
	}

	p2 := pipe.New(cfg2)
	go p2.Flow()

	inputChann := make(chan interface{}, 10)
	ways := []*mux.Way{
		{
			State: func(input interface{}) bool {
				in := input.(int)
				return in%2 == 0
			},
			To: p1.InputChannel(),
		},
		{
			State: func(input interface{}) bool {
				in := input.(int)
				return in%2 == 1
			},
			To: p2.InputChannel(),
		},
	}
	mux := mux.New(inputChann, ways, 4)
	go mux.Flow()

	for num := range 10 {
		inputChann <- num
		time.Sleep(1 * time.Second)
	}

}
