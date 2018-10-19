package util

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func Round64(input float64) int64 {
	if input < 0 {
		return int64(input - 0.5)
	}
	return int64(input + 0.5)
}

func TimeToMs(date time.Time) int64 {
	return date.UnixNano() / 1000000
}

func MsToTime(ms int64) time.Time {
	return time.Unix(0, ms*1000000)
}

func NewGuid() (string, error) {
	byteArray := make([]byte, 16)

	_, err := rand.Read(byteArray)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(byteArray), nil
}

func NewInviteId() (string, error) {
	byteArray := make([]byte, 4)

	_, err := rand.Read(byteArray)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(byteArray), nil
}
