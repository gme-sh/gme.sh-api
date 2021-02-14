package short

import (
	"math/rand"
	"strings"
	"time"
)

var availableCharacters = []byte("ABCDEFGHKLMNPRSTUVWXYZabcdefghkmnprstuvwxyz0123455689+-!")

func GenerateID(length int, accept func(id string) bool, try uint64) string {
	rand.Seed(time.Now().UnixNano())

	var res strings.Builder
	for i := 0; i < length; i++ {
		randByte := availableCharacters[rand.Intn(len(availableCharacters))]
		res.WriteByte(randByte)
	}

	resultStr := res.String()

	if !accept(resultStr) {
		try++
		if try > 5 {
			return ""
		}

		return GenerateID(length, accept, try)
	}

	return resultStr
}

func GenerateShortID(accept func(id string) bool) string {
	length := rand.Intn(10-5) + 5
	return GenerateID(length, accept, 0)
}

func AlwaysTrue(s string) bool {
	return true
}
