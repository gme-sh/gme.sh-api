package short

import (
	"math/rand"
	"strings"
	"time"
)

var availableCharacters = []byte("ABCDEFGHKLMNPRSTUVWXYZabcdefghkmnprstuvwxyz0123455689")

// GenerateID generates a string of length {length} using the characters from availableCharacters.
// This generated string is checked against the {accept} function and,
// if this function returns true, it will return the generated string.
// Otherwise, the process is tried a total of 5 more times until the {accept} function returns true,
// or the attempts are exhausted. After that an empty string is returned.
func GenerateID(length int, accept func(id *ShortID) bool, try uint64) ShortID {
	rand.Seed(time.Now().UnixNano())

	var res strings.Builder
	for i := 0; i < length; i++ {
		randByte := availableCharacters[rand.Intn(len(availableCharacters))]
		res.WriteByte(randByte)
	}

	result := ShortID(res.String())

	if !accept(&result) {
		try++
		if try > 5 {
			return ""
		}

		return GenerateID(length, accept, try)
	}

	return result
}

// GenerateShortID generates a 5 - 10 long string and checks it against the {accept} function.
// More information: GenerateID()
func GenerateShortID(accept func(id *ShortID) bool) ShortID {
	length := rand.Intn(10-5) + 5
	return GenerateID(length, accept, 0)
}

// AlwaysTrue always returns true. Used for GenerateID and GenerateShortID.
func AlwaysTrue(_ *ShortID) bool {
	return true
}
