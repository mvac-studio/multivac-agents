package processors

type ProcessorContext[T any, C any] struct {
	Value   T
	Context C
}
