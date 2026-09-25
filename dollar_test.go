package main

import "testing"

func TestMultiInDollar(t *testing.T) {
	usd := "USD"
	five := Money{5, usd}
	ten := five.Times(2)
	if ten.Amount != 10 || ten.Currency != usd {
		t.Errorf("failed, expected 10, got: %d\n", ten.Amount)
	}
}

func TestMoneyMulti(t *testing.T) {
	eur := "EUR"
	five := Money{Amount: 5, Currency: eur}
	ten := five.Times(2)
	if ten.Amount != 10 || ten.Currency !=  eur {
		t.Errorf("failed, expected 10eup, got: %d%s\n", ten.Amount, ten.Currency)
	}
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