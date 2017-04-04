# go-party

Simple 2-tier go worker pools that support timeout and panic recovery

## Install

```
go get github.com/davidiamyou/go-party
```

## Usage

```go
import (
    "time"
    "github.com/davidiamyou/go-party"
)

func main() {
    d := party.NewDispatcher(1, 1, time.Second)
    d.Start()

    // declare the channel to receive your output
    // you can use channel to provide input arguments in a similar way
    result := make(chan int, 1)

    // task definition, just a plain function
    task := func() {
    	time.Sleep(time.Millisecond * 200)
    	result <- 100
    }

    // run the task through the dispatcher
    // immediately returns a report object containing two channels for panic event and timeout event
    er := d.Run(task)

    // based on the event signal, decide your next step
    select {
    case r := <-er.PanicChannel:
    	// task panicked $r
    case <-er.TimeoutChannel:
    	// task timed out
    case num := <-result:
    	// task returned result $num
    }
}
```

