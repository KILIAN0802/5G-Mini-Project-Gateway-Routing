package main

import (
	"log"
	"sync"
)

var requestMutex sync.Mutex

func IncrementActiveRequest(){
	requestMutex.Lock()
	activeRequests++
	log.Printf("active=%d", activeRequests)
	requestMutex.Unlock()
}

func DecrementActiveRequest(){
	requestMutex.Lock()
	activeRequests--
	log.Printf("active=%d", activeRequests)
	requestMutex.Unlock()
}

func GetActiveRequests() int {
	requestMutex.Lock()
	defer requestMutex.Unlock()
	return activeRequests
}