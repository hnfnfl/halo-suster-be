package util

import (
	"math/rand"
)

type PrefixID string

const (
	ITPrefix    PrefixID = "IT"
	NursePrefix PrefixID = "NS"
)

func UuidGenerator(prefix PrefixID) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randStr := make([]byte, 30)
	for i := range randStr {
		randStr[i] = chars[rand.Intn(len(chars))]
	}

	return string(prefix) + string(randStr)
}
