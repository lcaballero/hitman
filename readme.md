# Introduction

`hitman` is a small library useful for cleaning up `go routines`.  The api consists
mainly of `Targets` and workers.  Workers are required to implement a minimal interface
for use with `Targets.Add()`.  When the code needs to shutdown gracefully it can then
call `Targets.Close()` which will send shutdown signals to the workers.

## Example

```Go

// Workers implement the Target interface requiring Name() and Start() KillChannel
var id = 0
func newId() int {
	id++
	return id
}

type Worker struct {
	id int
}
func NewWorker() *Worker {
	return &Worker{
		id: newId(),
	}
}
func (w *Worker) Name() string {
	return fmt.Sprintf("Work IT!!, %d", w.id)
}
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
	hits.Add(NewWorker())
	hits.Add(NewWorker())
	hits.Add(NewWorker())
	hits.Add(NewWorker())

	hits.Close()
}

```

## Using Death

If you use the [Death][Death] library then it might it can be written like
the example below.

__note__: [Death][Death] is waiting to update it's interface based on this [PR][PR]
and for the duration you can use this [FORK][FORK]

```Go

import (
    DEATH "github.com/vrecan/death"
    SYS "syscall"
)

func main() {
	hits := NewTargets()
	hits.Add(NewWorker())
	hits.Add(NewWorker())
	hits.Add(NewWorker())
	hits.Add(NewWorker())

	//pass the signals you want to end your application
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	// when you want to block for shutdown signals
	// this will finish when a signal of your type is sent to your application
	death.WaitForDeath(hits) 

}

func main() {
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