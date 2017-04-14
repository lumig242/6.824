package raft

import (
	"log"
	"time"
	"math/rand"
)

// Debugging
const (
	Debug = 0
	TIMEOUT_DURATION = time.Duration(200) * time.Millisecond
)


func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

func GetRandTimeoutDuration() (time.Duration){
	return time.Duration(rand.Int63() % 300 + 300) * time.Millisecond
}


