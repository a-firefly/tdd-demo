package main

import "testing"

func TestMultiInDollar(t *testing.T) {
	usd := "USD"
	five := Money{5, usd}
	ten := five.Times(2)
	if ten.Amount != 10 || ten.Currency != usd {
		t.Errorf("failed, expected 10, got: %f\n", ten.Amount)
	}
}

func TestMoneyMulti(t *testing.T) {
	eur := "EUR"
	five := Money{Amount: 5, Currency: eur}
	ten := five.Times(2)
	if ten.Amount != 10 || ten.Currency !=  eur {
		t.Errorf("failed, expected 10eup, got: %f%s\n", ten.Amount, ten.Currency)
	}
}

func TestMoney_Divide(t *testing.T) {
	currency := "KRW"
	m := Money{Amount: 4002, Currency: currency}
	rst := m.Divide(4)
	if rst.Amount != 1000.5 || rst.Currency !=  currency {
		t.Errorf("failed, expected 1000.5 %s, got: %f%s\n", currency, rst.Amount, rst.Currency)
	}
}
