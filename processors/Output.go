package processors

type Output[OUT any] struct {
	output chan OUT
}

func NewOutputProcessor[OUT any]() *Output[OUT] {
	return &Output[OUT]{
		output: make(chan OUT),
	}
}

// To adds a processor to the output channel
func (p *Output[OUT]) To(processors ...*Input[OUT]) {
	for _, processor := range processors {
		go func(source *Output[OUT], target *Input[OUT]) {
			for message := range source.output {
				target.input <- message
			}
		}(p, processor)
	}
}
