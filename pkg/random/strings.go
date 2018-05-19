package random

import (
	"crypto/rand"
	"math/big"
)

// Set so we can stub out for unittests.
var seeder = rand.Reader

// AlphanumericString generator. Ensures all values are A-Z, a-z, or 0-9.
func AlphanumericString(length int) string {
	if length <= 0 {
		return ""
	}
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return genString(chars, length)
}

// AlphanumericLowercaseString generator. Ensures all values are a-z or 0-9.
func AlphanumericLowercaseString(length int) string {
	if length <= 0 {
		return ""
	}
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	return genString(chars, length)
}

// SafeAlphanumericString generator. Ensures all values are A-Z (except for I, L, O) and 2-9.
// Does not return lowercase.
func SafeAlphanumericString(length int) string {
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
			panic("SafeAlphanumericString: " + err.Error())
		}
		str[i] = (chars[n.Int64()])
	}
	return string(str)
}
