package utils

import (
	"math/rand"
	"time"
)

func RandomString(userAgentList []string) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userAgentList))
	return userAgentList[randomIndex]
}
