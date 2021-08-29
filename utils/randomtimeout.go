package utils

import (
	"log"
	"math/rand"
	"time"
)

func RandomTimeOut() {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 5000
	randomNum := rand.Intn(max-min+1) + min
	log.Printf("Force random wait for: %d millisecond", randomNum)
	time.Sleep(time.Duration(randomNum) * time.Millisecond)
}
