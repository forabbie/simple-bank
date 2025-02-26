package util

import (
	"math/rand"
	"strings"
	"time"
)
const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min +1) // 0->max-min
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string  {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies) // len = length of the slice
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}

func RandomPassword() string {
	return RandomString(6)
}

type Transfer struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

func RandomTransfer() Transfer {
	return Transfer {
		FromAccountID: RandomInt(1, 100),
		ToAccountID: RandomInt(1, 100),
		Amount: RandomMoney(),
	}
}