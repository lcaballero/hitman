package hitman

import (
	"sync"
	"log"
)

// A KillContract is sent on the KillChannel when a "go routine" should be
// terminated.  It's provided the name in order to verify the contract.
type KillSignal struct {
	Name string
	WaitGroup *sync.WaitGroup
}

// A creator of a kill-able "go routine" should return a KillChannel that
// can be used to kill the routine.
type KillChannel chan KillSignal

// Make a kill-channel that has 0 queue size since a each channel is
// responsible for terminating 1 go routine.
func NewKillChannel() KillChannel {
	return make(KillChannel, 0)
}

// Simple struct that can be used to track a name with the channel that
// can kill the given target.
type Contract struct {
	Name string
	Done KillChannel
}

// Kill takes down a Target, and the routine that receives the signal
// should call Done() on the WaitGroup when the contract is complete.
func (t Contract) Kill(wg *sync.WaitGroup) {
	t.Done <- KillSignal{
		Name: t.Name,
		WaitGroup: wg,
	}
}

// A Target is an interface used to acquire the KillChannel for a given
// target.  The name should be unique.
type Target interface {
	Name() string
	Start() KillChannel
}

// Implements io.Closer.
type Targets map[string]Contract

func NewTargets() Targets {
	return make(Targets, 0)
}

// Saves the given target which can later be terminated.
func (targets Targets) Add(m Target) {
	targets[m.Name()] = Contract{
		Name: m.Name(),
		Done: m.Start(),
	}
}

// Calls Kill for each target added to the collection.
func (targets Targets) Close() error {
	wg := sync.WaitGroup{}
	wg.Add(len(targets))
	for name, target := range targets {
		log.Println("Killing: ", name)
		target.Kill(&wg)
	}
	wg.Wait()
	return nil
}
