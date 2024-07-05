package utils

import (
	"math/rand"
	"time"
)

func RandomString(strings []string) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(strings))
	return strings[randomIndex]
}
