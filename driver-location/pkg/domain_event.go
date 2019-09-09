package pkg

type DomainEvent interface {
	AggregateID() string
}

//go:generate moq -out event_dispatcher_mock.go . EventDispatcher
type EventDispatcher interface {
	Dispatch(domainEvent []DomainEvent)
}
