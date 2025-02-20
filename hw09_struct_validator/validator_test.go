package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
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

	// Additional test structures.
	Numbers struct {
		Values []int `validate:"min:10|max:20"`
	}

	ComplexValidation struct {
		String string `validate:"len:5|regexp:\\d+"`
		Int    int    `validate:"min:10|max:20|in:15,16,17"`
	}

	UnknownValidation struct {
		Some int `validate:"unknown:10"`
	}

	InvalidMin struct {
		Some int `validate:"min:aa"`
	}
)

func TestError(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr *error
	}{
		{
			name:        "non-struct validation",
			in:          "not a struct",
			expectedErr: &ErrNotStruct,
		},
		{
			name: "unknown validation",
			in: UnknownValidation{
				Some: 1,
			},
			expectedErr: &ErrUnknownValidator,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)

			require.ErrorAs(t, err, tt.expectedErr, "actual error %q", err)
		})
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			name: "valid user",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John",
				Age:    25,
				Email:  "test@test.te",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid user",
			in: User{
				ID:     "123",
				Age:    15,
				Email:  "invalid-email",
				Role:   "unknown",
				Phones: []string{"123", "456"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errors.New("string length does not match required length: length 3 does not match length 36")},
				{Field: "Age", Err: errors.New("value is less than min: 15 is less than min 18")},
				{Field: "Email", Err: errors.New("string does not match pattern: ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: errors.New("value is not in allowed set: unknown is not in set admin,stuff")},
				//nolint:lll
				{Field: "Phones[0]", Err: errors.New("string length does not match required length: length 3 does not match length 11")},
				//nolint:lll
				{Field: "Phones[1]", Err: errors.New("string length does not match required length: length 3 does not match length 11")},
			},
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 418,
				Body: "I'm a teapot",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errors.New("value is not in allowed set: 418 is not in allowed set 200,404,500")},
			},
		},
		{
			name: "valid numbers",
			in: Numbers{
				Values: []int{15, 12, 18},
			},
			expectedErr: nil,
		},
		{
			name: "invalid numbers",
			in: Numbers{
				Values: []int{5, 25, 15},
			},
			expectedErr: ValidationErrors{
				{Field: "Values[0]", Err: errors.New("value is less than min: 5 is less than min 10")},
				{Field: "Values[1]", Err: errors.New("value is greater than max: 25 is greater than max 20")},
			},
		},
		{
			name: "valid complex validation",
			in: ComplexValidation{
				String: "12345",
				Int:    16,
			},
			expectedErr: nil,
		},
		{
			name: "invalid complex validation",
			in: ComplexValidation{
				String: "abcde", // matches len:5 but not regexp:\d+
				Int:    12,      // matches min/max but not in:15,16,17
			},
			expectedErr: ValidationErrors{
				{Field: "String", Err: errors.New("string does not match pattern: \\d+")},
				{Field: "Int", Err: errors.New("value is not in allowed set: 12 is not in allowed set 15,16,17")},
			},
		},
		{
			name: "invalid min value",
			in: InvalidMin{
				Some: 1,
			},
			expectedErr: ValidationErrors{
				{Field: "Some", Err: errors.New("invalid min value: aa")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	errs := ValidationErrors{
		{Field: "field1", Err: errors.New("error1")},
		{Field: "field2", Err: errors.New("error2")},
	}

	expected := "field1: error1; field2: error2"
	require.Equal(t, expected, errs.Error())
}

func TestValidationErrors_ErrorEmpty(t *testing.T) {
	errs := ValidationErrors{}
	require.Equal(t, "", errs.Error())
}
