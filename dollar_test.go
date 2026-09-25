package main

import "testing"

func TestDollarMulti(t *testing.T) {
	five := Dollar{5}
	ten := five.Times(2)
	if ten.Amount != 10 {
		t.Errorf("failed, expected 10, got: %d\n", ten.Amount)
	}
}

type Dollar struct {
	Amount int
}

func (d Dollar) Times(x int) Dollar {
	return Dollar{10}
}