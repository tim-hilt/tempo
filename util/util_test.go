package util

import "testing"

func TestDivmod(t *testing.T) {
	if q, r := Divmod(3600, 60); q != 360 && r != 0 {
		t.Errorf("q should be 360, was %d; r should be 0, was %d", q, r)
	}
}

func TestMax(t *testing.T) {
	if max := Max(1, 4, 1, 5, 6, 42); max != 42 {
		t.Errorf("max should be 42, was %d", max)
	}
}

func TestContains(t *testing.T) {
	if i := Contains([]string{"foo", "bar", "baz"}, "baz"); i != true {
		t.Errorf("i should be true, was %t", i)
	}

	if i := Contains([]string{"foo", "bar", "baz"}, "fizz"); i != false {
		t.Errorf("i should be false, was %t", i)
	}
}

func TestSlicesEqual(t *testing.T) {
	if i := SlicesEqual([]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}); i != true {
		t.Errorf("i should be true, was %t", i)
	}

	if i := SlicesEqual([]string{"foo", "bar", "fizz"}, []string{"foo", "bar", "baz"}); i != false {
		t.Errorf("i should be false, was %t", i)
	}
}
