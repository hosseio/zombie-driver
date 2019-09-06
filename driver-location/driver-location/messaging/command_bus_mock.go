package messaging

import (
	"sync"

	"github.com/chiguirez/cromberbus"
)

var (
	lockCommandBusMockDispatch sync.RWMutex
)

var _ cromberbus.CommandBus = &CommandBusMock{}

type CommandBusMock struct {
	// DispatchFunc mocks the Dispatch method.
	DispatchFunc func(command cromberbus.Command) error

	// calls tracks calls to the methods.
	calls struct {
		// Dispatch holds details about calls to the Dispatch method.
		Dispatch []struct {
			// Command is the command argument value.
			Command cromberbus.Command
		}
	}
}

// Dispatch calls DispatchFunc.
func (mock *CommandBusMock) Dispatch(command cromberbus.Command) error {
	if mock.DispatchFunc == nil {
		panic("CommandBusMock.DispatchFunc: method is nil but CommandBus.Dispatch was just called")
	}
	callInfo := struct {
		Command cromberbus.Command
	}{
		Command: command,
	}
	lockCommandBusMockDispatch.Lock()
	mock.calls.Dispatch = append(mock.calls.Dispatch, callInfo)
	lockCommandBusMockDispatch.Unlock()
	return mock.DispatchFunc(command)
}

// DispatchCalls gets all the calls that were made to Dispatch.
// Check the length with:
//     len(mockedCommandBus.DispatchCalls())
func (mock *CommandBusMock) DispatchCalls() []struct {
	Command cromberbus.Command
} {
	var calls []struct {
		Command cromberbus.Command
	}
	lockCommandBusMockDispatch.RLock()
	calls = mock.calls.Dispatch
	lockCommandBusMockDispatch.RUnlock()
	return calls
}
