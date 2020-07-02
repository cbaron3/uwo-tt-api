package worker

import (
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Trim removes all irrelevant whitespace
func Trim(str string) string {
	// Trim the ends of whitespace
	str = strings.TrimSpace(str)

	// Trim the inner text of whitespace
	space := regexp.MustCompile(`\s+`)
	str = space.ReplaceAllString(str, " ")

	return str
}

// SleepRandom sleeps for a random amount of time in seconds between min and max
func SleepRandom(min int, max int) {
	time.Sleep(time.Duration(rand.Intn(max-min+1)+min) * time.Second)
}

// CreateData is a helper to create string map for form data
func CreateData(subject string) map[string][]string {
	data := url.Values{}
	data.Set("subject", subject)
	data.Set("command", "search")

	return data
}
