package model

import (
	"fmt"
	"testing"
)

type userData struct {
	Name     string
	Email    string
	Password string
}

func TestUserValid(t *testing.T) {
	tests := []struct {
		desc  string
		in    *userData
		valid bool
	}{
		{
			desc:  "All valid",
			in:    &userData{Name: "Mark", Email: "aaa@test.com", Password: "Mark@123456789"},
			valid: true,
		},
		{
			desc:  "Name is requried",
			in:    &userData{Name: "", Email: "aaa@test.com", Password: "Test@123456789"},
			valid: false,
		},
		{
			desc:  "Email is required",
			in:    &userData{Name: "Mark", Email: "", Password: "Test@123456789"},
			valid: false,
		},
		{
			desc:  "Password is required",
			in:    &userData{Name: "Mark", Email: "aaa@test.com", Password: ""},
			valid: false,
		},
		{
			desc:  "Email format",
			in:    &userData{Name: "Mark", Email: "aaa#test.com", Password: "Test@123456789"},
			valid: false,
		},
		{
			desc:  "Email format valid",
			in:    &userData{Name: "Mark", Email: "abc123@outlook.com", Password: "Test@123456789"},
			valid: true,
		},
		{
			desc:  "Name format",
			in:    &userData{Name: "@", Email: "aaa@test.com", Password: "Test@123456789"},
			valid: false,
		},
		{
			desc:  "Password length",
			in:    &userData{Name: "Mark", Email: "aaa@test.com", Password: "Abc@123"},
			valid: false,
		},
		{
			desc:  "Password format - missing special char",
			in:    &userData{Name: "Mark", Email: "aaa@test.com", Password: "Mark123456789"},
			valid: false,
		},
		{
			desc:  "Password format - missing uppercase",
			in:    &userData{Name: "Mark", Email: "aaa@test.com", Password: "mark@123456789"},
			valid: false,
		},
		{
			desc:  "Password format - missing lowercase",
			in:    &userData{Name: "Mark", Email: "aaa@test.com", Password: "MARK@123456789"},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			u := &User{
				Name:     tt.in.Name,
				Email:    tt.in.Email,
				Password: tt.in.Password,
			}
			err := u.Valid(true)
			// if err != nil {
			// 	fmt.Println("err: ", err)
			// 	fmt.Println("err is ErrValidUserFailed: ", errors.Is(err, ErrValidUserFailed))
			// }
			got := err == nil
			want := tt.valid

			if got != want {
				t.Errorf("user: %+v \nvalidate result should be %t, but got %t, error: %v", tt.in, want, got, err)
			}
		})
	}
}

func TestExtractNameFromEmail(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			"abc@test.com",
			"abc",
		},
		{
			"$abc#@test.com",
			"abc",
		},
		{
			"a@test.com",
			"a",
		},
		{
			"abc.efg@test.com",
			"abc.efg",
		},
		{
			".abc_@test.com",
			"abc",
		},
		{
			"a#b^c@test.com",
			"a.b.c",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("extract name from %s", tt.in), func(t *testing.T) {
			got := ExtractNameFromEmail(tt.in)
			if got != tt.want {
				t.Errorf("want %s, but got %s", tt.want, got)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{
			in:   "abc@test.com",
			want: true,
		},
		{
			in:   "abc.123.def@test.com",
			want: true,
		},
		{
			in:   "123@test.com.cn",
			want: true,
		},
		{
			in:   "abc@test.com'org",
			want: false,
		},
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.in)
		var got = false
		if err == nil {
			got = true
		}

		if got != tt.want {
			t.Errorf("%s validate result should be %t, but got %t", tt.in, tt.want, got)
		}
	}
}

func TestValidPassword(t *testing.T) {
	tests := []struct {
		desc  string
		pwd   string
		valid bool
	}{
		{
			desc:  "Valid password with all requirements",
			pwd:   "Admin@123456789",
			valid: true,
		},
		{
			desc:  "Valid password with different special char",
			pwd:   "User1!234567890",
			valid: true,
		},
		{
			desc:  "Password too short",
			pwd:   "Abc@12345",
			valid: false,
		},
		{
			desc:  "Missing uppercase letter",
			pwd:   "user@123456789",
			valid: false,
		},
		{
			desc:  "Missing lowercase letter",
			pwd:   "USER@123456789",
			valid: false,
		},
		{
			desc:  "Missing number",
			pwd:   "User@Password!",
			valid: false,
		},
		{
			desc:  "Missing special character",
			pwd:   "User123456789",
			valid: false,
		},
		{
			desc:  "Contains invalid special character",
			pwd:   "User-123456789",
			valid: false,
		},
		{
			desc:  "Old weak password should fail",
			pwd:   "admin123456",
			valid: false,
		},
		{
			desc:  "Old weak password should fail",
			pwd:   "moderator123",
			valid: false,
		},
		{
			desc:  "Old weak password should fail",
			pwd:   "user123456",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := ValidPassword(tt.pwd)
			got := err == nil
			if got != tt.valid {
				t.Errorf("password '%s' validation should be %t, but got %t, error: %v", tt.pwd, tt.valid, got, err)
			}
		})
	}
}
