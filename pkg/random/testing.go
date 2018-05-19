package random

import (
	"crypto/rand"
	mathrand "math/rand"
)

// EnableTestSeeder enables a consistent pseudo random seeder.
func EnableTestSeeder() {
	seeder = mathrand.New(mathrand.NewSource(0))
}

// DisableTestSeeder disables the test seeder.
func DisableTestSeeder() {
	seeder = rand.Reader
}
