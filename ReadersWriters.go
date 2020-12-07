package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var NR int = 3
var NW int = 3
var N int = NR + NW

var readers int = 0
var writers int = 0
var mutex sync.RWMutex
var maxCountMutex sync.Mutex
var maxReaders int = 2
var locked int = 0

func sleep() {
	time.Sleep(time.Duration(rand.Intn(1500)+100) * time.Millisecond)
}

func reader(i int) {
	for {
		sleep()
		fmt.Printf("WANT_R! (ID =:%d), readers=%d, writers=%d\n", i, readers, writers)
		mutex.RLock()
		if readers >= maxReaders {
			fmt.Println("BLOCKED!!!")
			locked++
			maxCountMutex. Lock()
		}
		readers++
		fmt.Printf("++readers (ID =:%d), readers=%d, writers=%d\n", i, readers, writers)
		sleep()
		readers--
		if locked > 0 {
			locked--
			fmt.Println("UNBLOCKED")
			maxCountMutex.Unlock()
		}
		fmt.Printf("--readers (ID =:%d), readers=%d, writers=%d\n", i, readers, writers)
		mutex.RUnlock()
	}
}

func writer(i int) {
	for {
		fmt.Printf("WANT_W! (ID =:%d), readers=%d, writers=%d\n", i, readers, writers)
		mutex.Lock()
		writers++
		fmt.Printf("++writers (ID =:%d), readers=%d, writers=%d\n", i, readers, writers)
		sleep()
		writers--
		fmt.Printf("--writers (ID =:%d), readers=%d, writers=%d\n", i, readers, writers)
		mutex.Unlock()
	}
}

func main() {
	var wg = sync.WaitGroup{}
	maxCountMutex.Lock()
	wg.Add(1)
	for i := NW; i < N; i++ {
		go reader(i)
	}
	for i := 0; i < NW; i++ {
		go writer(i)
	}
	wg.Wait()
}
