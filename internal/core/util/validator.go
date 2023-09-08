package util

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func ValidateString(value string, min, max int) error {
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf("must contain from %d-%d characters", min, max)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits or underscores")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 8, 200)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not an valid email address")
	}
	return nil
}

func ValidateID(value int64) error {
	if value <= 0 {
		return fmt.Errorf("id must be positive integer")
	}
	return nil
}

func ValidateURL(value string) error {
	_, err := url.ParseRequestURI(value)
	return err
}
