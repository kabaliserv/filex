package auth

import "regexp"

var (
	regexEmail    = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	regexUserName = regexp.MustCompile(`^[-a-zA-Z0-9]{3,}$`)
)

func ValidEmail(v string) bool {

	if res := regexEmail.FindString(v); res != "" {
		return true
	}

	return false
}

func ValidUserName(v string) bool {

	if result := regexUserName.FindString(v); result != "" {
		return true
	}

	return false
}
