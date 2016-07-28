package hitman

import (
	"sync"
	"log"
"github.com/pborman/uuid"
)

// A KillSignal is sent on the KillChannel when a "go routine" should be
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

// Named interface provides a way to name a target optionally.
type Named interface {
	Name() string
}

// A Target is an interface used to acquire the KillChannel for a given
// target.  The name should be unique.
type Target interface {
	Start() KillChannel
}

// NamedTarget combines the Named and Target interfaces for those targets
// that care to implement both.
type NamedTarget interface {
	Named
	Target
}

// Implements io.Closer.
type Targets map[string]Contract

// Creates a new Targets mapping.
func NewTargets() Targets {
	return make(Targets, 0)
}

// Add saves the given target which can later be terminated.
func (targets Targets) Add(m NamedTarget) {
	targets.Put(m.Name(), m)
}

// AddTarget checks to see if the given Target is Named and if so uses
// that name, else it generates a new UUID as the Target's name.
func (targets Targets) AddTarget(m Target) {
	named, ok := m.(Named)
	if ok {
		targets.Put(named.Name(), m)
	} else {
		targets.Put(uuid.New(), m)
	}
}

// AddOrPanic adds the given target to the collection of targets so long as
// the provided error is nil, else it panics.  This is useful when
// constructing a new service that might produce an error on construction.
func (targets Targets) AddOrPanic(t Target, err error) {
	if err != nil {
		panic(err)
	}
	targets.AddTarget(t)
}

// AddPair adds target without Target having to implement Named
func (targets Targets) Put(name string, m Target) {
	_, alreadyHasTarget := targets[name]
	if alreadyHasTarget {
		log.Printf("Given name of target already used: %s\n", name)
	}
	targets[name] = Contract{
		Name: name,
		Done: m.Start(),
	}
}

// Calls Kill for each target added to the collection.
func (targets Targets) Close() error {
	wg := sync.WaitGroup{}
	wg.Add(len(targets) * 2)
	for n, c := range targets {
		go func(name string, contract Contract) {
			defer wg.Done()
			log.Println("Killing: ", name)
			contract.Kill(&wg)
			log.Println("Killed: ", name)
		}(n, c)
	}
	wg.Wait()
	return nil
}
