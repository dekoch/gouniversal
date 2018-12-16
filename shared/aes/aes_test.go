package aes

import (
	"testing"
)

func TestNewKey(t *testing.T) {

	tests := []struct {
		in        int
		wantError bool
	}{
		{0, true},
		{1, true},
		{16, false},
		{24, false},
		{32, false},
	}

	for _, tt := range tests {

		got, err := NewKey(tt.in)

		l := len(got)

		if l != tt.in {
			t.Errorf("NewKey(): got %d, want %d", l, tt.in)
		}

		if tt.wantError {
			if err == nil {
				t.Errorf("NewKey(): got nil, want error")
			}
		} else {
			if err != nil {
				t.Errorf("NewKey(): got %v, want nil", err)
			}
		}
	}
}

func TestEncrypt(t *testing.T) {

	tests := []struct {
		inKey     []byte
		inText    string
		wantError bool
	}{
		{[]byte{}, "", true},
		{[]byte{}, "foo", true},
		{[]byte{46, 48, 103, 98, 87, 147, 61, 199, 230, 16, 84, 57, 191, 211, 186, 79, 239, 85, 208, 204, 169, 90, 205, 49, 238, 173, 43, 103, 186, 244, 104}, "foo", true},
		{[]byte{46, 48, 103, 98, 87, 147, 61, 199, 230, 16, 84, 57, 191, 211, 186, 79, 239, 85, 208, 204, 169, 90, 205, 49, 238, 173, 43, 103, 186, 244, 104, 190}, "", false},
		{[]byte{46, 48, 103, 98, 87, 147, 61, 199, 230, 16, 84, 57, 191, 211, 186, 79, 239, 85, 208, 204, 169, 90, 205, 49, 238, 173, 43, 103, 186, 244, 104, 190}, "foo", false},
	}

	for _, tt := range tests {

		got, err := Encrypt(tt.inKey, tt.inText)

		if tt.wantError {
			if err == nil {
				t.Errorf("Encrypt(): got nil, want error")
			}
		} else {

			if len(got) == 0 {
				t.Errorf("Encrypt(): no output")
			}

			if err != nil {
				t.Errorf("Encrypt(): got %v, want nil", err)
			}
		}
	}
}

func TestDecrypt(t *testing.T) {

	tests := []struct {
		inKey     []byte
		inText    string
		wantText  string
		wantError bool
	}{
		{[]byte{}, "", "", true},
		{[]byte{}, "foo", "", true},
		{[]byte{46, 48, 103, 98, 87, 147, 61, 199, 230, 16, 84, 57, 191, 211, 186, 79, 239, 85, 208, 204, 169, 90, 205, 49, 238, 173, 43, 103, 186, 244, 104, 190}, "", "", true},
		{[]byte{46, 48, 103, 98, 87, 147, 61, 199, 230, 16, 84, 57, 191, 211, 186, 79, 239, 85, 208, 204, 169, 90, 205, 49, 238, 173, 43, 103, 186, 244, 104, 190}, "MEYcHyFmLJsOWdiP9Oz7QoMVBgDDPEc5drvRSDWsRI4", "foo", false},
	}

	for _, tt := range tests {

		got, err := Decrypt(tt.inKey, tt.inText)

		if got != tt.wantText {
			t.Errorf("Decrypt(): got %s, want %s", got, tt.wantText)
		}

		if tt.wantError {
			if err == nil {
				t.Errorf("Decrypt(): got nil, want error")
			}
		} else {
			if err != nil {
				t.Errorf("Decrypt(): got %v, want nil", err)
			}
		}
	}
}

func TestPackage(t *testing.T) {

	tests := []struct {
		in int
	}{
		{16},
		{24},
		{32},
	}

	for _, tt := range tests {

		key, err := NewKey(tt.in)
		if err != nil {
			t.Errorf("TestPackage: got %v, want nil", err)
		}

		e, err := Encrypt(key, "foo")
		if err != nil {
			t.Errorf("TestPackage: got %v, want nil", err)
		}

		d, err := Decrypt(key, e)
		if err != nil {
			t.Errorf("TestPackage: got %v, want nil", err)
		}

		if d != "foo" {
			t.Errorf("TestPackage: got %s, want foo", d)
		}
	}
}
