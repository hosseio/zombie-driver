// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"sync"
)

var (
	lockDriverRepositoryMockSave sync.RWMutex
)

// Ensure, that DriverRepositoryMock does implement DriverRepository.
// If this is not the case, regenerate this file with moq.
var _ DriverRepository = &DriverRepositoryMock{}

// DriverRepositoryMock is a mock implementation of DriverRepository.
//
//     func TestSomethingThatUsesDriverRepository(t *testing.T) {
//
//         // make and configure a mocked DriverRepository
//         mockedDriverRepository := &DriverRepositoryMock{
//             SaveFunc: func(in1 Driver) error {
// 	               panic("mock out the Save method")
//             },
//         }
//
//         // use mockedDriverRepository in code that requires DriverRepository
//         // and then make assertions.
//
//     }
type DriverRepositoryMock struct {
	// SaveFunc mocks the Save method.
	SaveFunc func(in1 Driver) error

	// calls tracks calls to the methods.
	calls struct {
		// Save holds details about calls to the Save method.
		Save []struct {
			// In1 is the in1 argument value.
			In1 Driver
		}
	}
}

// Save calls SaveFunc.
func (mock *DriverRepositoryMock) Save(in1 Driver) error {
	if mock.SaveFunc == nil {
		panic("DriverRepositoryMock.SaveFunc: method is nil but DriverRepository.Save was just called")
	}
	callInfo := struct {
		In1 Driver
	}{
		In1: in1,
	}
	lockDriverRepositoryMockSave.Lock()
	mock.calls.Save = append(mock.calls.Save, callInfo)
	lockDriverRepositoryMockSave.Unlock()
	return mock.SaveFunc(in1)
}

// SaveCalls gets all the calls that were made to Save.
// Check the length with:
//     len(mockedDriverRepository.SaveCalls())
func (mock *DriverRepositoryMock) SaveCalls() []struct {
	In1 Driver
} {
	var calls []struct {
		In1 Driver
	}
	lockDriverRepositoryMockSave.RLock()
	calls = mock.calls.Save
	lockDriverRepositoryMockSave.RUnlock()
	return calls
}
