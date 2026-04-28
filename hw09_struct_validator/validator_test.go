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

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
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

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			
			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}
				
				switch tt.expectedErr.(type) {
				case ValidationErrors:
					var validationErrors ValidationErrors
					if !errors.As(err, &validationErrors) {
						t.Errorf("expected ValidationErrors, got: %T", err)
					}
				default:
					if !errors.Is(err, tt.expectedErr) {
						t.Errorf("expected error %v, got: %v", tt.expectedErr, err)
					}
				}
			}
		})
	}
}
