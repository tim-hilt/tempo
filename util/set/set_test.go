package set

import "testing"

func TestAdd(t *testing.T) {
	// Given
	s := New[string]()

	// When
	s.Add("item")
	s.Add("item2")

	// Then
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

func TestDelete(t *testing.T) {
	// Given
	s := New[string]()
	s.Add("item")
	s.Add("item2")

	// When
	s.Delete("item")

	// Then
	if l := len(s.items); l != 1 {
		t.Errorf("len(s.items) should be 1, was %d", l)
	}
}

func TestContains(t *testing.T) {
	// Given
	s := New[string]()

	// When
	s.Add("item")
	s.Add("item2")

	// Then
	if ci1 := s.Contains("item"); ci1 == false {
		t.Errorf("s.Contains(\"item\") should be true, was %t", ci1)
	}

	if ci2 := s.Contains("item2"); ci2 == false {
		t.Errorf("s.Contains(\"item2\") should be true, was %t", ci2)
	}
}

func TestReset(t *testing.T) {
	// Given
	s := New[string]()
	s.Add("item")
	s.Add("item2")

	// When
	s.Reset()

	// Then
	if l := len(s.items); l != 0 {
		t.Errorf("len(s.items) should be 0, was %d", l)
	}
}

func TestItems(t *testing.T) {
	// Given
	s := New[string]()
	s.Add("item")
	s.Add("item2")

	// When/Then
	if l := len(s.Items()); l != 2 {
		t.Errorf("len(s.items) should be 2, was %d", l)
	}
}
