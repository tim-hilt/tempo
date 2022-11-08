package set

import "testing"

func TestAdd(t *testing.T) {
	s := New[string]()
	s.Add("item")
	s.Add("item2")

	if l := len(s.items); l != 2 {
		t.Errorf("len(s.items) should be 2, was %d", l)
	}

	if ci1 := s.Contains("item"); ci1 == false {
		t.Errorf("s.Contains(\"item\") should be true, was %t", ci1)
	}

	if ci2 := s.Contains("item2"); ci2 == false {
		t.Errorf("s.Contains(\"item2\") should be true, was %t", ci2)
	}
}
