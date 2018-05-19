package random

import (
	"crypto/rand"
	"math/big"
)

// Set so we can stub out for unittests.
var seeder = rand.Reader

// Alphanumeric generator. Ensures all values are A-Z, a-z, or 0-9.
func Alphanumeric(length int) string {
	if length <= 0 {
		return ""
	}
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return genString(chars, length)
}

// AlphanumericLowercase generator. Ensures all values are a-z or 0-9.
func AlphanumericLowercase(length int) string {
	if length <= 0 {
		return ""
	}
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	return genString(chars, length)
}

// SafeAlphanumeric generator. Ensures all values are A-Z (except for I, L, O) and 2-9.
// Does not return lowercase.
func SafeAlphanumeric(length int) string {
	if length <= 0 {
		return ""
	}
	const chars = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"
	return genString(chars, length)
}

func genString(chars string, length int) string {
	str := make([]byte, length)
	for i := range str {
		n, err := rand.Int(seeder, big.NewInt(int64(len(chars))-1))
		if err != nil {
			panic("random: genstring: " + err.Error())
		}
		str[i] = (chars[n.Int64()])
	}
	return string(str)
}
