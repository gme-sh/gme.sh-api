package short

import (
	"math/rand"
	"strings"
	"time"
)

var availableCharacters = []byte("ABCDEFGHKLMNPRSTUVWXYZabcdefghkmnprstuvwxyz0123455689+-!")

func generateShortID(accept func(id string) bool, try uint64) string {
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(10-5) + 5

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

		return generateShortID(accept, try)
	}

	return resultStr
}

func GenerateShortID(accept func(id string) bool) string {
	return generateShortID(accept, 0)
}
