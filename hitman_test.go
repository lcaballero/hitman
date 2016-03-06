package hitman
   
import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
)

func TestName(t *testing.T) {

	Convey("Should kill them one.", t, func() {
		hits := NewTargets()
		hits.Add(NewWorker())
		hits.Add(NewWorker())
		hits.Add(NewWorker())
		hits.Add(NewWorker())

		So(len(hits), ShouldEqual, 4)
		err := hits.Close()
		So(err, ShouldBeNil)
	})

	Convey("New kill channels should have a length of 0.", t, func() {
		k := NewKillChannel()
		So(len(k), ShouldEqual, 0)
	})

	Convey("Should kill them one.", t, func() {
		hits := NewTargets()
		t := &Worker{}
		hits.Add(t)

		So(len(hits), ShouldEqual, 1)

		err := hits.Close()
		So(err, ShouldBeNil)
	})

	Convey("New Targets collection should have len of zero", t, func() {
		t := NewTargets()
		w := NewWorker()
		t.Add(w)

		v,ok := t[w.Name()]

		So(len(t), ShouldEqual, 1)
		So(v, ShouldNotBeNil)
		So(ok, ShouldBeTrue)
	})

	Convey("New Targets collection should have len of zero", t, func() {
		t := NewTargets()
		So(len(t), ShouldEqual, 0)
	})

	Convey("New Target via literals", t, func() {
		t := NewTargets()
		t.AddTarget(&Unnamed{})
		So(len(t), ShouldEqual, 1)
		t.Close()
	})

	Convey("New Target or panic should be fine if no error", t, func() {
		t := NewTargets()
		t.AddOrPanic(NewService())
		So(len(t), ShouldEqual, 1)
		t.Close()
	})

	Convey("Should panic when adding a service that failed to instantiate without error", t, func() {
		defer func() {
			panicky := recover()
			err, ok := panicky.(error)
			So(ok, ShouldBeTrue)
			So(err, ShouldNotBeNil)
		}()
		NewTargets().AddOrPanic(NewServiceWithError())
	})
}


type Service struct {}
func NewServiceWithError() (*Service, error) {
	s := &Service{}
	return s, fmt.Errorf("Had error while creating new service")
}
func NewService() (*Service, error) {
	s := &Service{}
	return s, nil
}
func (w *Service) Start() KillChannel {
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


type Unnamed struct {}
func (w *Unnamed) Start() KillChannel {
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
