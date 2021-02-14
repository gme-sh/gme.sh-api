package short

import (
	"math/rand"
	"strings"
	"time"
)

var availableCharacters = []byte("ABCDEFGHKLMNPRSTUVWXYZabcdefghkmnprstuvwxyz0123455689+-!")

func GenerateID(length int, accept func(id ShortID) bool, try uint64) ShortID {
	rand.Seed(time.Now().UnixNano())

	var res strings.Builder
	for i := 0; i < length; i++ {
		randByte := availableCharacters[rand.Intn(len(availableCharacters))]
		res.WriteByte(randByte)
	}

	result := ShortID(res.String())

	if !accept(result) {
		try++
		if try > 5 {
			return ""
		}

		return GenerateID(length, accept, try)
	}

	return result
}

func GenerateShortID(accept func(id ShortID) bool) ShortID {
	length := rand.Intn(10-5) + 5
	return GenerateID(length, accept, 0)
}

func AlwaysTrue(_ ShortID) bool {
	return true
}
