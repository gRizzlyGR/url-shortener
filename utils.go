package main

import (
	"encoding/base64"
	"log"
	"net/url"
	"strconv"
	"time"
)

// IsURLValid checks if the provided URL is valid or not. The URL must be
// complete with protocol and www. If some error is thrown while parsing, the
// validation fails
func IsURLValid(rawURL string) bool {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		log.Printf("Cannot parse url: %s - Error: %v\n", rawURL, err)
		return false
	}

	log.Printf("Parsed url: %s \n", parsed)
	return true
}

// GenerateKey returns an id based on the current timestamp in milliseconds, then
// encoded in base64. You can choose to pad the encoded string or not, using
// withPadding param
func GenerateKey(withPadding bool) string {
	var padding rune
	if withPadding == true {
		padding = base64.StdPadding
	} else {
		padding = base64.NoPadding
	}

	millis := time.Now().UnixNano() / 1000000
	stringMillis := strconv.FormatInt(millis, 10)
	return base64.StdEncoding.WithPadding(padding).EncodeToString([]byte(stringMillis))
}
