package main

import "testing"

func TestDollarMulti(t *testing.T) {
	five := Dollar{5}
	ten := five.Times(2)
	if ten.Amount != 10 {
		t.Errorf("failed, expected 10, got: %d\n", ten.Amount)
	}
}

func TestMoneyMulti(t *testing.T) {
	five := Money{Amount: 5, Currency: "eup"}
	ten := five.Times(2)
	if ten.Amount != 10 || ten.Currency !=  "eup" {
		t.Errorf("failed, expected 10eup, got: %d%s\n", ten.Amount, ten.Currency)
	}
}

type Dollar struct {
	Amount int
}

func (d Dollar) Times(x int) Dollar {
	return Dollar{d.Amount * x}
}

type Money struct {
	Amount int
	Currency string
}

func (m Money) Times(x int) Money {
	return Money{
		Amount: m.Amount * x,
		Currency: m.Currency,
	}
}