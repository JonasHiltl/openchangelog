package analytics

func NewNoopEmitter() Emitter {
	return noopEmitter{}
}

type noopEmitter struct{}

func (e noopEmitter) Emit(Event) {}
