[![Build Status](https://travis-ci.org/lcaballero/hitman.svg?branch=master)](https://travis-ci.org/lcaballero/hitman)

# Introduction

`hitman` is a small library useful for cleaning up `go routines`
refered to as `Targets`.

A `Target` is required to implement a minimal interface for use with
`Targets.Add(NamedTarget)` or `Targets.AddTarget(Target)`, namely a
`Target` must have a `Start() hitman.KillChannel` and a `NamedTarget`
must have both `Name() string` and `Start() hitman.KillChannel`.  (The
difference is provided for convenience where a developer doesn't truly
care about naming the `Target` for logging purposes.)

See the example below which first creates a `hitman.Targets` instance
and then adds `Targets` using `Add`.  When the application needs to
stop the collection of routines and allow those go routines to
shutdown gracefully it calls `Targets.Close()` which will send a
`hitman.KillSignal` to each routine allowing them to clean up and
shutdown gracefully.

## Example

```Go

// Components implement the Target interface and Start() KillChannel

type Component struct {}

func (w *Component) Start() KillChannel {
	kill := NewKillChannel()
	go func(done KillChannel) {
		select {
		case cleaner := <-done:
			cleaner.WaitGroup.Done()
			return
		}
	}(kill)
	return kill
}

// Component lifecycles are then managed with a Targets collection:
func main() {
	targets := NewTargets()
	targets.Add(&Component{})
	targets.Add(&Component{})
	targets.Close()
}

```

## Using Death

If you use the [Death][Death] library then the example above can be
written like so:


```Go

import (
    DEATH "github.com/vrecan/death"
    SYS "syscall"
)

type Service struct {}
func NewService() (*Server, error) {
  s := &Service{}
  // Do something that might produce error
  _, err := ioutil.ReadFile("doh.txt")
  return s, err
}

func main() {
	targets := NewTargets()
	targets.AddOrPanic(NewComponent())
	targets.AddOrPanic(NewComponent())

	//pass the signals you want to end your application
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	// when you want to block for shutdown signals
	// this will finish when a signal of your type is sent to your application
	death.WaitForDeath(targets) 
}
```

## Changes

### 2016-04-06
1. Revised the readme.

### 2016-03-01

1. Added `AddOrPanic` for services that return errors on construction
   (convenience).
1. Partitioned the `Target` interface into `Named` and `Target` so
   that `Name()` was optional.
1. Added Put that can take a given name instead of requiring the
   `Named` interface.
1. Added tests for the above changes.

## License

See license file.

The use and distribution terms for this software are covered by the
[Eclipse Public License 1.0][EPL-1], which can be found in the file
'license' at the root of this distribution. By using this software in
any fashion, you are agreeing to be bound by the terms of this
license. You must not remove this notice, or any other, from this
software.


[EPL-1]: http://opensource.org/licenses/eclipse-1.0.txt
[DEATH]: https://github.com/vrecan/death
