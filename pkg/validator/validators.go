package validator

import (
	"errors"
	"forum/internal/schemas"
	"regexp"
	"strings"
	"unicode"

	"github.com/gofrs/uuid"
)

var (
	ErrorInvalidPassword      = errors.New("Invalid password")
	ErrorInvalidEmail         = errors.New("Invalid email adress")
	ErrorPasswordMatch        = errors.New("Password and confirmed password do not match")
	ErrorEmptyField           = errors.New("Empty field")
	ErrorInvalidPostTtitleLen = errors.New("the title should be 5 - 20 characters in length")
	ErrorInvalidPostBodyLen   = errors.New("the body  should be 20 - 250 characters in length")
	ErrorInvalidUUID          = errors.New("invalid uuid")
)

// следует ли эту функцию сделать методом для user?
func ValidateRegisterInput(user schemas.CreateUser) error {
	if user.Username == "" || user.Email == "" || user.Password == "" || user.PasswordConfirm == "" {
		return ErrorEmptyField
	}
	if !validateEmail(user.Email) {
		return ErrorInvalidEmail
	}

	// if err := validatePassword(user.Password); err != nil {
	// 	return ErrorInvalidPassword
	// }

	if !validatePasswordConfirmed(user.Password, user.PasswordConfirm) {
		return ErrorPasswordMatch
	}

	return nil
}

func ValidateSignInInput(user schemas.AuthUser) error {
	if user.Email == "" || user.Password == "" {
		return ErrorEmptyField
	}
	if !validateEmail(user.Email) {
		return ErrorInvalidEmail
	}

	return nil
}

func ValidateCreatePostInput(post schemas.CreatePost) error {
	title := strings.TrimSpace(post.Title)
	body := strings.TrimSpace(post.Body)

	if len(title) > 20 && len(title) < 5 {
		return ErrorInvalidPostTtitleLen
	}
	if len(body) < 20 && len(body) > 250 {
		return ErrorInvalidPostBodyLen
	}
	return nil
}

func ValidateUpdatePostInput(post schemas.UpdatePost) error {
	title := strings.TrimSpace(post.CreatePost.Title)
	body := strings.TrimSpace(post.CreatePost.Body)

	if len(title) > 20 && len(title) < 5 {
		return ErrorInvalidPostTtitleLen
	}
	if len(body) < 20 && len(body) > 250 {
		return ErrorInvalidPostBodyLen
	}
	return nil
}

func validatePasswordConfirmed(password, PasswordConfirm string) bool {
	if password != PasswordConfirm {
		return false
	}
	return true
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrorInvalidPassword
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !validateSpecialChar(password) {
		return ErrorInvalidPassword
	}

	return nil
}

func validateSpecialChar(password string) bool {
	specialCharsSet := "!@$%^&*()_-+"

	for _, char := range password {
		if strings.ContainsRune(specialCharsSet, char) {
			return true
		}
	}

	return false
}

func validateEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func ValidateUUID(input string) bool {
	_, err := uuid.FromString(input)
	return err == nil
}
