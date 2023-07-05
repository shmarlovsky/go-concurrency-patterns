package generator

import (
	"fmt"
	"math/rand"
	"time"
)

// returns receive only channel
func BoringGenerator(msg string, milt int) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%v %v", msg, i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)*milt))
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
func Multiplex(in1, in2 <-chan string) <-chan string {
	aggr := make(chan string)
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
		fmt.Printf("Got: %v\n", <-ch)
	}

	fmt.Println("Enough")
}
