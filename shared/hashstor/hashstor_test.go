package hashstor

import (
	"testing"
	"time"
)

func TestAddGetHash(t *testing.T) {

	tests := []struct {
		in      []string
		want    []string
		wantErr bool
	}{
		{[]string{"foo", "foo"}, []string{"foo"}, false},
		{[]string{"foo0", "foo1", "foo2"}, []string{"foo2", "foo1", "foo0"}, false},
		{[]string{}, []string{"foo"}, true},
	}

	for i, tt := range tests {

		var hs HashStor

		for _, s := range tt.in {
			hs.Add(s)
		}

		got, err := hs.GetHash()
		if err == nil && tt.wantErr {
			t.Errorf("AddGetHash() test %d: err %v, want error", i, err)
		}

		if err != nil {
			continue
		}

		found := false

		for _, s := range tt.want {

			if s == got {
				found = true
			}
		}

		if found == false {
			t.Errorf("AddGetHash() test %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestRemove(t *testing.T) {

	tests := []struct {
		in     []string
		remove string
	}{
		{[]string{"foo0", "foo1", "foo2"}, "foo0"},
		{[]string{"foo0", "foo1", "foo2"}, "foo1"},
		{[]string{"foo0", "foo1", "foo2"}, "foo2"},
	}

	for i, tt := range tests {

		var hs HashStor

		for _, s := range tt.in {
			hs.Add(s)
		}

		found := false

		for i := range hs.Hashes {

			if tt.remove == hs.Hashes[i].Hash {
				found = true
			}
		}

		if found == false {
			t.Errorf("Remove() test %d: check Add()", i)
			return
		}

		hs.Remove(tt.remove)

		found = false

		for i := range hs.Hashes {

			if tt.remove == hs.Hashes[i].Hash {
				found = true
			}
		}

		if found {
			t.Errorf("Remove() test %d: found %s", i, tt.remove)
		}
	}
}

func TestRemoveAll(t *testing.T) {

	tests := []struct {
		in []string
	}{
		{[]string{"foo0", "foo1", "foo2"}},
	}

	for i, tt := range tests {

		var hs HashStor

		for _, s := range tt.in {
			hs.Add(s)
		}

		if len(tt.in) != len(hs.Hashes) {
			t.Errorf("RemoveAll() test %d: check Add()", i)
			return
		}

		hs.RemoveAll()

		got := len(hs.Hashes)

		if got != 0 {
			t.Errorf("RemoveAll() test %d: got %d, want 0", i, got)
		}
	}
}

func TestSetAsExpired(t *testing.T) {

	tests := []struct {
		in     []string
		expire string
		want   string
	}{
		{[]string{"foo0", "foo1"}, "", "foo0"},
		{[]string{"foo0", "foo1"}, "foo0", "foo1"},
	}

	for i, tt := range tests {

		var hs HashStor

		for _, s := range tt.in {
			hs.Add(s)
		}

		if len(tt.in) != len(hs.Hashes) {
			t.Errorf("SetAsExpired() test %d: check Add()", i)
			return
		}

		hs.SetAsExpired(tt.expire, time.Minute*1)

		got, err := hs.GetHash()
		if err != nil {
			t.Errorf("SetAsExpired() test %d: err %v,", i, err)
			return
		}

		if got != tt.want {
			t.Errorf("SetAsExpired() test %d: got %s, want %s", i, got, tt.want)
		}
	}
}
