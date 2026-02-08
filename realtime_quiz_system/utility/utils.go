package utility

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const sessionCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Exclude similar characters like 0/O, 1/I

func String2Int64(s string) (i int64) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return
}

// GenerateSessionCode generates a unique 6-character alphanumeric session code
func GenerateSessionCode() string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = sessionCodeChars[rand.Intn(len(sessionCodeChars))]
	}
	return string(code)
}

// NormalizeSessionCode converts session code to uppercase for case-insensitive comparison
func NormalizeSessionCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
