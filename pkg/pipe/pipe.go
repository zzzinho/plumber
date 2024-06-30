package pipe

import (
	"fmt"
	"log"
)

type Interface interface {
	Flow()
	InputChannel() chan<- interface{}
	OutputChannel() <-chan interface{}
}

type Pipe struct {
	name       string
	input      chan interface{}
	output     chan interface{}
	numWorkers int // TODO: need suitable name
	task       func(interface{}) interface{}
}

type Config struct {
	Name                    string
	InputChannelBufferSize  int
	OutputChannelBufferSize int
	NumWorkers              int
	Task                    func(interface{}) interface{}
}

func New(cfg Config) Interface {
	var inputChan chan interface{}
	if cfg.InputChannelBufferSize < 0 {
		log.Fatal("input channel buffer size should be 0 or more")
	} else if cfg.InputChannelBufferSize == 0 {
		inputChan = make(chan interface{})
	} else {
		inputChan = make(chan interface{}, cfg.InputChannelBufferSize)
	}

	var outputChan chan interface{}
	if cfg.OutputChannelBufferSize < 0 {
		outputChan = nil
	} else if cfg.OutputChannelBufferSize == 0 {
		outputChan = make(chan interface{})
	} else {
		outputChan = make(chan interface{}, cfg.OutputChannelBufferSize)
	}
	return &Pipe{
		name:       cfg.Name,
		numWorkers: cfg.NumWorkers,
		input:      inputChan,
		output:     outputChan,
		task:       cfg.Task,
	}
}

func (p *Pipe) InputChannel() chan<- interface{} {
	return p.input
}

func (p *Pipe) OutputChannel() <-chan interface{} {
	return p.output
}

func (p *Pipe) Flow() {
	for n := range p.numWorkers {
		go func() {
			fmt.Printf("Pipe %s_%d is running\n", p.name, n)
			for i := range p.input {
				fmt.Printf("[Pipe_%d] input: %d\n", n, i.(int))
				o := p.task(i)
				fmt.Printf("[Pipe_%d] output: %d\n", n, o.(int))
				if p.output != nil {
					p.output <- o
				}
			}
			fmt.Printf("Pipe %s_%d is closed\n", p.name, n)
		}()
	}
}
