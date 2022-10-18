package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shijiayun/atomicshm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("./atomicshmdemo <shmName>")
		os.Exit(1)
	}
	u, err := atomicshm.OpenOrCreateUint64(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open shm uint64, err: %v", err)
	}
	delta := uint64(1000)
	totalDelta := uint64(0)
	for {
		totalDelta += delta
		newVal := u.Add(delta)
		log.Printf("totalDelta: %d, curVal: %d", totalDelta, newVal)
		if newVal > 10000000 {
			u.Store(0)
		}
		time.Sleep(time.Millisecond)
	}
}
