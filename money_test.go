package main

import "testing"

func TestMultiInDollar(t *testing.T) {
	currency := "USD"
	five := Money{5, currency}
	actual := five.Times(2)
	expected := Money{ Amount: 10, Currency: currency}
	assetEqual(t, expected, actual)
}

func TestMoneyMulti(t *testing.T) {
	currency := "EUR"
	five := Money{Amount: 10, Currency: currency}
	actual := five.Times(2)
	expected := Money{ Amount: 20, Currency: currency}
	assetEqual(t, expected, actual)
}

func TestMoney_Divide(t *testing.T) {
	currency := "KRW"
	m := Money{Amount: 4002, Currency: currency}
	actual := m.Divide(4)
	expected := Money{ Amount: 1000.5, Currency: currency}
	assetEqual(t, expected, actual)
}

func assetEqual(t *testing.T, expected, actual Money) {
	if actual != expected {
		t.Errorf("failed, expected %+v, got: %+v\n", expected, actual)
	}
}
