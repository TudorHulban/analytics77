package appanalytics

import (
	"math/rand"
	"strings"
)

func randomNumberString(minDigits, maxDigits int) string {
	if minDigits <= 0 {
		minDigits = 1
	}

	if maxDigits < minDigits {
		maxDigits = minDigits
	}

	// Determine a random length between minDigits and maxDigits (inclusive)
	length := minDigits + rand.Intn(maxDigits-minDigits+1) //nolint:gosec

	var sb strings.Builder
	sb.Grow(length)

	// First digit must be from '1' to '9' to avoid leading zeros
	firstDigit := byte('1' + rand.Intn(9)) //nolint:gosec
	sb.WriteByte(firstDigit)

	// Remaining digits can be from '0' to '9'
	for ix := 1; ix < length; ix++ {
		digit := byte('0' + rand.Intn(10)) //nolint:gosec

		sb.WriteByte(digit)
	}

	return sb.String()
}
