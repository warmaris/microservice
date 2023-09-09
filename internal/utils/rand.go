package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(length int) string {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length+2)
	src.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
