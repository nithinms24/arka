package util

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/adwitiyaio/arka/dependency"
	"github.com/gofrs/uuid"
)

const DependencyAppUtil = "app_util"

// AppUtil ... Service providing utility functionality throughout the application
type AppUtil interface {
	// GetCurrentTime ... Get the current time from the system
	GetCurrentTime() time.Time

	//GenerateOTP ... Generate a one time password
	GenerateOTP(length int) string
	// GenerateUniqueToken ... Generate a unique token
	GenerateUniqueToken() string
	// GetExpiryTimeForDuration ... Get an expiry time based on the duration (in hours) passed
	GetExpiryTimeForDuration(duration int) time.Time

	// CompareSlices ... Find the elements in one array of string but not in the other
	CompareSlices(a, b []string) []string

	// ParseStringForTime ... Parse string into time.RFC3339 format
	ParseStringForTime(date string) (time.Time, error)
	//FormatDate ... Format date to get day of month with suffix
	FormatDate(t time.Time) string
	// IsTimeExpired ... Validate if the specified time has expired based on the current time
	IsTimeExpired(t time.Time) bool
}

// Bootstrap ... Bootstraps the application utility functionality
func Bootstrap() {
	a := &simpleAppUtil{}
	dependency.GetManager().Register(DependencyAppUtil, a)
}

type simpleAppUtil struct{}

func (as *simpleAppUtil) GetCurrentTime() time.Time {
	return time.Now()
}

func (as *simpleAppUtil) GenerateOTP(length int) string {
	const digits = "1234567890"
	if length <= 0 {
		return ""
	}
	buffer := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, buffer, length)
	if n != length || err != nil {
		return ""
	}
	for i := 0; i < len(buffer); i++ {
		buffer[i] = digits[int(buffer[i])%len(digits)]
	}
	return string(buffer)
}

func (as *simpleAppUtil) GenerateUniqueToken() string {
	code, _ := uuid.NewV4()
	return code.String()
}

func (as *simpleAppUtil) GetExpiryTimeForDuration(duration int) time.Time {
	t := as.GetCurrentTime().Add(time.Hour*time.Duration(duration) + time.Minute*0 + time.Second*0)
	return t
}

func (as *simpleAppUtil) CompareSlices(a, b []string) (diff []string) {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func (as *simpleAppUtil) ParseStringForTime(date string) (time.Time, error) {
	tm, err := time.Parse(time.RFC3339, date)
	return tm, err
}

func (as *simpleAppUtil) IsTimeExpired(t time.Time) bool {
	if time.Now().Sub(t).Seconds() <= 0 {
		return false
	}
	return true
}

func (as *simpleAppUtil) FormatDate(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.Format("2" + suffix + " Jan")
}
