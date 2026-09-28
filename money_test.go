package main

import (
	"testing"
	"tdd/stocks"
)

func TestMoneyMulti(t *testing.T) {
	currency := "EUR"
	five := stocks.NewMoney(10, currency)
	actual := five.Times(2)
	expected := stocks.NewMoney(20, currency)
	assetEqual(t, expected, actual)
}

func TestMoney_Divide(t *testing.T) {
	currency := "KRW"
	m := stocks.NewMoney(4002, currency)
	actual := m.Divide(4)
	expected := stocks.NewMoney(1000.5, currency)
	assetEqual(t, expected, actual)
}

func TestAddtionOfUSDAndEUR(t *testing.T) {
	fiveUSD := stocks.NewMoney(5, "USD")
	tenEUR := stocks.NewMoney(10, "EUR")
	var portfolio stocks.Portfolio
	portfolio = portfolio.Add(fiveUSD)
	portfolio = portfolio.Add(tenEUR)
	rst := portfolio.Evaluate("USD")

	expected := stocks.NewMoney(17, "USD")

	assetEqual(t, expected, rst)
}

func TestAddtionOfDollarsAndWons(t *testing.T) {
	oneDollar := stocks.NewMoney(1, "USD")
	elvenHundredWon := stocks.NewMoney(1100, "KRW")
	var portfolio stocks.Portfolio
	portfolio = portfolio.Add(oneDollar)
	portfolio = portfolio.Add(elvenHundredWon)
	rst := portfolio.Evaluate("KRW")

	expected := stocks.NewMoney(2200, "KRW")

	assetEqual(t, expected, rst)
}

func assetEqual(t *testing.T, expected, actual stocks.Money) {
	if actual != expected {
		t.Errorf("failed, expected %+v, got: %+v\n", expected, actual)
	}
}
