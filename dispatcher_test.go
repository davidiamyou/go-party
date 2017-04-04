package party

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	d := NewDispatcher(1, 1, time.Second)
	d.Start()

	result := make(chan int, 1)
	task := func() {
		time.Sleep(time.Millisecond * 200)
		result <- 100
	}

	er := d.Run(task)

	select {
	case <-er.PanicChannel:
		t.Error("should not have failed.")
	case <-er.TimeoutChannel:
		t.Error("should not have timed out.")
	case num := <-result:
		assert.Equal(t, 100, num)
	}
}

func TestDispatcherThatPanics(t *testing.T) {
	d := NewDispatcher(1, 1, time.Second)
	d.Start()

	result := make(chan int, 1)
	task := func() {
		time.Sleep(time.Millisecond * 200)
		panic("foo")
		result <- 100
	}

	er := d.Run(task)

	select {
	case r := <-er.PanicChannel:
		assert.Equal(t, "foo", r.(string))
	case <-er.TimeoutChannel:
		t.Error("should not have timed out.")
	case <-result:
		t.Error("should not have received result")
	}
}

func TestDispatcherThatTimesOut(t *testing.T) {
	d := NewDispatcher(1, 1, time.Millisecond*200)
	d.Start()

	result := make(chan int, 1)
	task := func() {
		time.Sleep(time.Second)
		result <- 100
	}

	er := d.Run(task)

	select {
	case <-er.PanicChannel:
		t.Error("should not have failed.")
	case b := <-er.TimeoutChannel:
		assert.True(t, b)
	case <-result:
		t.Error("should not have received result")
	}
}

func benchmarkDispatcher(partySize int, waitList int, iterations int, b *testing.B) {
	d := NewDispatcher(partySize, waitList, time.Hour)
	d.Start()

	b.N = iterations
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result := make(chan int, 1)
			task := func() {
				time.Sleep(200 * time.Millisecond)
				result <- 100
			}

			er := d.Run(task)
			select {
			case r := <-er.PanicChannel:
				b.Fatal(r)
			case <-er.TimeoutChannel:
				b.Fatal("timed out")
			case <-result:
			}
		}
	})
}

func BenchmarkDispatcher_PartySize2_WaitList10_Iteration100(b *testing.B) {
	benchmarkDispatcher(2, 10, 100, b)
}

func BenchmarkDispatcher_PartySize5_WaitList10_Iteration100(b *testing.B) {
	benchmarkDispatcher(5, 10, 100, b)
}

func BenchmarkDispatcher_PartySize10_WaitList10_Iteration100(b *testing.B) {
	benchmarkDispatcher(10, 10, 100, b)
}

func BenchmarkDispatcher_PartySize25_WaitList10_Iteration100(b *testing.B) {
	benchmarkDispatcher(25, 10, 100, b)
}

func BenchmarkDispatcher_PartySize50_WaitList10_Iteration100(b *testing.B) {
	benchmarkDispatcher(50, 10, 100, b)
}
