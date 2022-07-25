package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID          string `json:"id" validate:"len:36|regexp:\\d+"`
		Name        string
		Age         int      `validate:"min:18|max:50"`
		Email       string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role        UserRole `validate:"in:admin,stuff|regexp:\\W+"`
		Phones      []string `validate:"len:11"`
		meta        json.RawMessage
		Application App `validate:"nested"`
	}

	App struct {
		Version   string `validate:"len:5"`
		UserToken Token  `validate:"nested"`
	}

	Token struct {
		Header      []byte
		Payload     []byte
		Signature   []byte
		GetResponse Response `validate:"nested"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

// TestValidate проверка работоспособности.
func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "rwerwer",
				Name:   "sdfsdf",
				Age:    12,
				Email:  "324234",
				Role:   "324234",
				Phones: []string{"qewqeqweqweqwe", "qwertyuiop["},
				meta:   []byte{},
				Application: App{
					Version: "123123123",
					UserToken: Token{
						Header:    []byte{},
						Payload:   []byte{},
						Signature: []byte{},
						GetResponse: Response{
							Code: 300,
						},
					},
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrStringLength,
				},
				ValidationError{
					Field: "ID",
					Err:   ErrStringRegexp,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrNumberMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrStringRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrStringIn,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrStringRegexp,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrStringLength,
				},
				ValidationError{
					Field: "Version",
					Err:   ErrStringLength,
				},
				ValidationError{
					Field: "Code",
					Err:   ErrNumberIn,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			var valErr ValidationErrors
			require.ErrorAs(t, err, &valErr)
			for _, error := range strings.Split(tt.expectedErr.Error(), "\n") {
				require.ErrorContains(t, err, error)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	t.Run("error ErrInterfaceType", func(t *testing.T) {
		err := Validate("")
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrInterfaceType.Error())
	})

	t.Run("error ErrInvalidEmptyTag", func(t *testing.T) {
		testStruct := struct {
			Field1 string `validate:"len:"`
		}{
			"qwerty",
		}

		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrInvalidEmptyTag.Error())
	})

	t.Run("error ErrStringLength", func(t *testing.T) {
		testStruct := struct {
			Field1 string `validate:"len:1"`
		}{
			"qwerty",
		}
		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrStringLength.Error())
	})

	t.Run("error ErrStringRegexp", func(t *testing.T) {
		testStruct := struct {
			Field1 string `validate:"regexp:\\d+"`
		}{
			"qwerty",
		}
		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrStringRegexp.Error())
	})

	t.Run("error ErrStringIn", func(t *testing.T) {
		testStruct := struct {
			Field1 string `validate:"regexp:\\d+"`
		}{
			"qwerty",
		}
		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrStringRegexp.Error())
	})

	t.Run("error ErrNumberMax", func(t *testing.T) {
		testStruct := struct {
			Field1 int64 `validate:"max:0"`
		}{
			1,
		}
		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrNumberMax.Error())
	})

	t.Run("error ErrNumberMin", func(t *testing.T) {
		testStruct := struct {
			Field1 int64 `validate:"min:0"`
		}{
			-1,
		}
		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrNumberMin.Error())
	})

	t.Run("error ErrNumberIn", func(t *testing.T) {
		testStruct := struct {
			Field1 int64 `validate:"in:0"`
		}{
			-1,
		}
		err := Validate(testStruct)
		var valErr ValidationErrors
		require.ErrorAs(t, err, &valErr)
		require.ErrorContains(t, valErr, ErrNumberIn.Error())
	})
}
