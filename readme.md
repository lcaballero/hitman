# Introduction

`hitman` is a small library useful for cleaning up `go routines`.  The api consists
mainly of `Targets` and workers.  Workers are required to implement a minimal interface
for use with `Targets.Add()`.  When the code needs to shutdown gracefully it can then
call `Targets.Close()` which will send shutdown signals to the workers.

## Changes

### 2016-03-01

1. Added `AddOrPanic` for services that return errors on construction (convenience).
1. Partitioned the `Target` interface into `Named` and `Target` so that `Name()` was optional.
1. Added Put that can take a given name instead of requiring the `Named` interface.
1. Added tests for the above changes.

## Example

```Go

// Workers implement the Target interface and Start() KillChannel

type Worker struct {}

func (w *Worker) Start() KillChannel {
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

// Worker lifecycles are then managed with a Targets collection:
func main() {
	hits := NewTargets()
	hits.Add(&Worker{})
	hits.Add(&Worker{})
	hits.Close()
}

```

## Using Death

If you use the [Death][Death] library then it might it can be written like
the example below.

__note__: [Death][Death] is waiting to update it's interface based on this [PR][PR]
and for the duration you can use this [FORK][FORK].

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
	hits := NewTargets()
	hits.AddOrPanic(NewWorker())
	hits.AddOrPanic(NewWorker())

	//pass the signals you want to end your application
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	// when you want to block for shutdown signals
	// this will finish when a signal of your type is sent to your application
	death.WaitForDeath(hits) 

}

```

## License

See license file.

The use and distribution terms for this software are covered by the
[Eclipse Public License 1.0][EPL-1], which can be found in the file 'license' at the
root of this distribution. By using this software in any fashion, you are
agreeing to be bound by the terms of this license. You must not remove this
notice, or any other, from this software.


[EPL-1]: http://opensource.org/licenses/eclipse-1.0.txt
[DEATH]: https://github.com/vrecan/death
[PR]: https://github.com/vrecan/death/pull/8
[FORK]: https://github.com/lcaballero/death/tree/use-closer