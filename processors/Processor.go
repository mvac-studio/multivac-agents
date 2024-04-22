package processors

import (
	"log"
)

type ProcessorAction[IN any, OUT any | *any] interface {
	Process(input IN) OUT
}

type Processor[IN any, OUT any] struct {
	ProcessorAction[IN, OUT]
	*Input[IN]
	*Output[OUT]
	process func(input IN) (OUT, error)
}

func NewProcessor[IN any, OUT any | *any](processor func(input IN) (OUT, error)) Processor[IN, OUT] {
	p := Processor[IN, OUT]{
		process: processor,
		Input:   NewInputProcessor[IN](),
		Output:  NewOutputProcessor[OUT](),
	}
	p.initialize()
	return p
}

func (o Processor[IN, OUT]) initialize() {
	go func() {
		for {
			select {
			case input := <-o.input:
				result, err := o.process(input)

				if err != nil {
					log.Println(err)
					continue
				}
				o.output <- result

			}
		}
	}()
}
