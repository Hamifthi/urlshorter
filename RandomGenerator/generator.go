package RandomGenerator

import (
	"math/rand"
	"time"
)

type RandomService struct {
	NumberOfDigits int
}

func (r *RandomService) GenerateRandomString() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, r.NumberOfDigits)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
