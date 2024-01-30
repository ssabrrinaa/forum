package exceptions

import "fmt"

type BaseError struct {
	StatusCode int
	Message    string
}

type ClientSideError struct {
	BaseError
}

type ValidationError struct {
	ClientSideError
}

type AuthenticationError struct {
	ClientSideError
}

type AuthorizationError struct {
	ClientSideError
}

type ResourceNotFoundError struct {
	ClientSideError
}

type BadRequestError struct {
	ClientSideError
}

type StatusMethodNotAllowed struct {
	ClientSideError
}

type ConnectionError struct {
	ClientSideError
}

type SessionExpiredError struct {
	ClientSideError
}

type StatusConflictError struct {
	ClientSideError
}

type ForbiddenError struct {
	ClientSideError
}

type ServerSideError struct {
	BaseError
}

type InternalServerError struct {
	ServerSideError
}

func (e BaseError) Error() string {
	return fmt.Sprintf("Status Code: %v, Message: %v", e.StatusCode, e.Message)
}

func NewValidationError(msg string) ValidationError {
	return ValidationError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_422,
				Message:    fmt.Sprintf("Validation Error: %v", msg),
			},
		},
	}
}

func NewAuthenticationError(msg string) AuthenticationError {
	return AuthenticationError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_401,
				Message:    fmt.Sprintf("Authentication Error: %v", msg),
			},
		},
	}
}

func NewResourceNotFoundError(msg string) ResourceNotFoundError {
	return ResourceNotFoundError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_404,
				Message:    fmt.Sprintf("Not Found Error: %v", msg),
			},
		},
	}
}

func NewBadRequestError(msg string) BadRequestError {
	return BadRequestError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_400,
				Message:    fmt.Sprintf("Bad Request Error: %v", msg),
			},
		},
	}
}

func NewStatusMethodNotAllowed() StatusMethodNotAllowed {
	return StatusMethodNotAllowed{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_405,
				Message:    fmt.Sprintf("Method Not Allowed error"),
			},
		},
	}
}

func NewSessionExpiredError(msg string) SessionExpiredError {
	return SessionExpiredError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_401,
				Message:    fmt.Sprintf("Session Expired Error: %v", msg),
			},
		},
	}
}

func NewStatusConflictError(msg string) StatusConflictError {
	return StatusConflictError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_409,
				Message:    fmt.Sprintf("Status Conflict Error: %v", msg),
			},
		},
	}
}

func NewForbiddenError(msg string) ForbiddenError {
	return ForbiddenError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: HTTP_403,
				Message:    fmt.Sprintf("Forbidden Error: %v", msg),
			},
		},
	}
}

func NewInternalServerError() InternalServerError {
	return InternalServerError{
		ServerSideError: ServerSideError{
			BaseError: BaseError{
				StatusCode: HTTP_500,
				Message:    fmt.Sprintf("Internal Server Error"),
			},
		},
	}
}
