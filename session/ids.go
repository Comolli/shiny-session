package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const referenceDate = 1483228800000

var (
	macAddress  [6]byte    // The unique MAC address of the computer running this program.
	lastMutex   sync.Mutex // The mutex which syncs access to the timestamp and counter.
	lastTime    uint64     // The timestamp of the last CCUID.
	lastCounter uint64     // The counter of the last CCUID.
)

func initCUID() {
	// Get a unique MAC address.
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) >= 6 {
			copy(macAddress[:], iface.HardwareAddr)
			break
		}
	}
}

func CUID() string {
	lastMutex.Lock()
	defer lastMutex.Unlock()

	// Initialize the bits with the timestamp.
	now := time.Now()
	timestamp := uint64(now.Unix())*1000 - referenceDate + uint64(now.Nanosecond())/1000000
	timestamp &= (1 << 40) - 1

	// Counter.
	if timestamp == lastTime {
		lastCounter++
	} else {
		lastCounter = 0
	}
	lastTime = timestamp
	counter := uint64(lastCounter & 0xff)

	// MAC address.
	var macHash uint16
	for _, b := range macAddress {
		macHash = (macHash << 5) - macHash // *= 31 (a prime).
		macHash += uint16(b)
	}
	spill := lastCounter >> 8
	if spill != 0 {
		macHash += uint16(spill & 0xffff)
	}
	mac := uint64(macHash)

	// Assemble.
	bits := (timestamp << 24) | (mac << 8) | counter

	// Transform to Base62.
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base := uint64(len(chars))
	var base64 string
	for len := 0; len < 11; len++ {
		base64 = string(chars[bits%base]) + base64
		bits /= base
	}

	return base64
}

func RandomID(length int) (string, error) {
	id := make([]byte, length)
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var b [1]byte
	for length > 0 {
		n, err := rand.Reader.Read(b[:])
		if err != nil {
			return "", err
		}
		if n < 1 {
			return "", errors.New("Unable to generate random number")
		}
		length--
		id[length] = chars[int(b[0])%len(chars)]
	}
	return string(id), nil
}

func generateSesssionID() (string, error) {
	// For more on collisions:
	// https://en.wikipedia.org/wiki/Birthday_problem

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("Could not generate session ID: %s", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
