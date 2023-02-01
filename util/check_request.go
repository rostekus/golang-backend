package util

import (
	"fmt"
	"regexp"
)

func CheckEmail(email string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if re.MatchString(email) {
		return nil
	}
	return fmt.Errorf("%q is not a valid email address", email)
}
