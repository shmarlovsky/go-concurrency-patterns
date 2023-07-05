package synctwo

import (
	"fmt"
	"math/rand"
	"time"
)

type Msg struct {
	Str  string
	Wait chan bool
}

// returns receive only channel
// has wait channel inside which sync execution between go routines
// each "speaker" is able to continue only when main goroutine is ready to listen
func BoringGenerator(msg string, milt int) <-chan Msg {
	ch := make(chan Msg)
	// shared across all messages
	waitForIt := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			ch <- Msg{fmt.Sprintf("%v %v", msg, i), waitForIt}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)*milt))
			<-waitForIt
		}
	}()
	return ch
}

func UseGen() {
	ch := BoringGenerator("Hey Kiryl", 10)
	ch2 := BoringGenerator("Hey Kate", 5)

	for i := 0; i < 5; i++ {
		fmt.Printf("Got: %v\n", <-ch)
		fmt.Printf("Got: %v\n", <-ch2)
	}

	fmt.Println("Enough")
}

// listens for both input channels and put in resulting channel which is returned
func Multiplex(in1, in2 <-chan Msg) <-chan Msg {
	aggr := make(chan Msg)
	go func() {
		for {
			aggr <- <-in1
		}
	}()
	go func() {
		for {
			aggr <- <-in2
		}
	}()

	return aggr
}

func UseMultiplex() {
	ch := Multiplex(BoringGenerator("Kiryl", 10), BoringGenerator("Kate", 3))

	for i := 0; i < 20; i++ {

		msg1 := <-ch
		fmt.Println(msg1.Str)
		msg2 := <-ch
		fmt.Println(msg2.Str)

		// indicates that listener is ready to receive next messages
		msg1.Wait <- true
		msg2.Wait <- true
	}

	fmt.Println("Enough")
}
