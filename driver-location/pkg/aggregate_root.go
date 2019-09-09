package pkg

type AggregateRoot interface {
	Record(event DomainEvent)
	Uncommited() []DomainEvent
	ClearEvents()
}

type BaseAggregateRoot struct {
	Events []DomainEvent
}

func (a *BaseAggregateRoot) Record(event DomainEvent) {
	a.Events = append(a.Events, event)
}

func (a *BaseAggregateRoot) Uncommited() []DomainEvent {
	return a.Events
}

func (a *BaseAggregateRoot) ClearEvents() {
	a.Events = []DomainEvent{}
}
