package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

type testCase struct {
	in          interface{}
	expectedErr error
}

func getTestCases() []testCase {
	return []testCase{
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"79991234567"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "short",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"79991234567"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    15,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"79991234567"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    60,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"79991234567"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    25,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"79991234567"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    25,
				Email:  "test@example.com",
				Role:   "guest",
				Phones: []string{"79991234567"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "Not Found",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 500,
				Body: "Internal Server Error",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 403,
				Body: "Forbidden",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in:          "not a struct",
			expectedErr: ErrNotStruct,
		},
		{
			in:          123,
			expectedErr: ErrNotStruct,
		},
		{
			in:          nil,
			expectedErr: ErrNotStruct,
		},
	}
}

func checkError(t *testing.T, err, expectedErr error) {
	t.Helper()

	if expectedErr == nil {
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	var validationErrors ValidationErrors
	if errors.As(expectedErr, &validationErrors) {
		if !errors.As(err, &validationErrors) {
			t.Errorf("expected ValidationErrors, got: %T", err)
		}
		return
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got: %v", expectedErr, err)
	}
}

func TestValidate(t *testing.T) {
	tests := getTestCases()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			checkError(t, err, tt.expectedErr)
		})
	}
}
