// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package driver_zombie

import (
	"sync"
)

var (
	lockDistanceCalculatorMockCalculate sync.RWMutex
)

// Ensure, that DistanceCalculatorMock does implement DistanceCalculator.
// If this is not the case, regenerate this file with moq.
var _ DistanceCalculator = &DistanceCalculatorMock{}

// DistanceCalculatorMock is a mock implementation of DistanceCalculator.
//
//     func TestSomethingThatUsesDistanceCalculator(t *testing.T) {
//
//         // make and configure a mocked DistanceCalculator
//         mockedDistanceCalculator := &DistanceCalculatorMock{
//             CalculateFunc: func(driverID string, lastMinutes int) (int, error) {
// 	               panic("mock out the Calculate method")
//             },
//         }
//
//         // use mockedDistanceCalculator in code that requires DistanceCalculator
//         // and then make assertions.
//
//     }
type DistanceCalculatorMock struct {
	// CalculateFunc mocks the Calculate method.
	CalculateFunc func(driverID string, lastMinutes int) (int, error)

	// calls tracks calls to the methods.
	calls struct {
		// Calculate holds details about calls to the Calculate method.
		Calculate []struct {
			// DriverID is the driverID argument value.
			DriverID string
			// LastMinutes is the lastMinutes argument value.
			LastMinutes int
		}
	}
}

// Calculate calls CalculateFunc.
func (mock *DistanceCalculatorMock) Calculate(driverID string, lastMinutes int) (int, error) {
	if mock.CalculateFunc == nil {
		panic("DistanceCalculatorMock.CalculateFunc: method is nil but DistanceCalculator.Calculate was just called")
	}
	callInfo := struct {
		DriverID    string
		LastMinutes int
	}{
		DriverID:    driverID,
		LastMinutes: lastMinutes,
	}
	lockDistanceCalculatorMockCalculate.Lock()
	mock.calls.Calculate = append(mock.calls.Calculate, callInfo)
	lockDistanceCalculatorMockCalculate.Unlock()
	return mock.CalculateFunc(driverID, lastMinutes)
}

// CalculateCalls gets all the calls that were made to Calculate.
// Check the length with:
//     len(mockedDistanceCalculator.CalculateCalls())
func (mock *DistanceCalculatorMock) CalculateCalls() []struct {
	DriverID    string
	LastMinutes int
} {
	var calls []struct {
		DriverID    string
		LastMinutes int
	}
	lockDistanceCalculatorMockCalculate.RLock()
	calls = mock.calls.Calculate
	lockDistanceCalculatorMockCalculate.RUnlock()
	return calls
}