package main

import (
	"math/rand"
	"time"
)

const (
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Numerals     = "1234567890"
	Alphanumeric = Alphabet + Numerals
)

// RandomInt64Range returns a random big integer in the range from min to max.
func RandomInt64Range(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func RandomFloat64Range(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomDate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	sec := RandomInt64Range(min, max)
	return time.Unix(sec, 0)
}

// RandomString returns a random string n characters long, composed of entities
// from charset.
func RandomString(n int, charset string) string {
	randstr := make([]byte, n) // Random string to return
	charlen := len(charset)
	for i := 0; i < n; i++ {
		r := rand.Intn(charlen)
		randstr[i] = charset[r]
	}
	return string(randstr)
}

// RandomStringRange returns a random string at least min and no more than max
// characters long, composed of entitites from charset.
func RandomStringRange(min, max int, charset string) string {
	strlen := int(RandomInt64Range(int64(min), int64(max)))
	return RandomString(strlen, charset)
}

// RandomAlphaStringRange returns a random alphanumeric string at least min and no more
// than max characters long.
func RandomAlphaStringRange(min, max int) string {
	return RandomStringRange(min, max, Alphanumeric)
}

// RandomAlphaString returns a random alphanumeric string n characters long.
func RandomAlphaString(n int) string {
	return RandomStringRange(n, n+1, Alphanumeric)
}
