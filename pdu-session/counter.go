package main

import (
	// "log"
	"sync/atomic"
)

func IncrementActiveRequest(){
	atomic.AddInt64(&activeRequests, 1)
	// log.Printf("active=%d", newVal)
}

func DecrementActiveRequest(){
	atomic.AddInt64(&activeRequests, -1)
	// log.Printf("active=%d", newVal)
}

func GetActiveRequests() int64 {
	return atomic.LoadInt64(&activeRequests)
}