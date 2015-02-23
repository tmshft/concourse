// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc/scheduler"
	"github.com/pivotal-golang/lager"
)

type FakeScanner struct {
	ScanStub        func(lager.Logger, string) error
	scanMutex       sync.RWMutex
	scanArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
	}
	scanReturns struct {
		result1 error
	}
}

func (fake *FakeScanner) Scan(arg1 lager.Logger, arg2 string) error {
	fake.scanMutex.Lock()
	fake.scanArgsForCall = append(fake.scanArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
	}{arg1, arg2})
	fake.scanMutex.Unlock()
	if fake.ScanStub != nil {
		return fake.ScanStub(arg1, arg2)
	} else {
		return fake.scanReturns.result1
	}
}

func (fake *FakeScanner) ScanCallCount() int {
	fake.scanMutex.RLock()
	defer fake.scanMutex.RUnlock()
	return len(fake.scanArgsForCall)
}

func (fake *FakeScanner) ScanArgsForCall(i int) (lager.Logger, string) {
	fake.scanMutex.RLock()
	defer fake.scanMutex.RUnlock()
	return fake.scanArgsForCall[i].arg1, fake.scanArgsForCall[i].arg2
}

func (fake *FakeScanner) ScanReturns(result1 error) {
	fake.ScanStub = nil
	fake.scanReturns = struct {
		result1 error
	}{result1}
}

var _ scheduler.Scanner = new(FakeScanner)
