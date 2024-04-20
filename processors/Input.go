package processors

type Input[IN any] struct {
	input chan IN
}

func NewInputProcessor[IN any]() *Input[IN] {
	return &Input[IN]{
		input: make(chan IN),
	}
}
