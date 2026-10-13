package main

import (
	"reflect"
	"tdd/stocks"
	"testing"
)

var bank stocks.Bank

func init() {
	initBankWithExchangeRates()
}

func initBankWithExchangeRates() {
	bank = stocks.NewBank()
	bank.AddExchangeRates("EUR", "USD", 1.2)
	bank.AddExchangeRates("USD", "KRW", 1100)
}

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
	initBankWithExchangeRates()

	fiveUSD := stocks.NewMoney(5, "USD")
	tenEUR := stocks.NewMoney(10, "EUR")
	var portfolio stocks.Portfolio
	portfolio = portfolio.Add(fiveUSD)
	portfolio = portfolio.Add(tenEUR)
	rst, err := portfolio.Evaluate(bank, "USD")
	asserNil(t, err)
	expected := stocks.NewMoney(17, "USD")

	assetEqual(t, expected, *rst)
}

func TestAddtionOfDollarsAndWons(t *testing.T) {
	initBankWithExchangeRates()

	oneDollar := stocks.NewMoney(1, "USD")
	elvenHundredWon := stocks.NewMoney(1100, "KRW")
	var portfolio stocks.Portfolio
	portfolio = portfolio.Add(oneDollar)
	portfolio = portfolio.Add(elvenHundredWon)
	rst, err := portfolio.Evaluate(bank, "KRW")
	asserNil(t, err)
	expected := stocks.NewMoney(2200, "KRW")

	assetEqual(t, expected, *rst)
}

func TestAddtionWithMissingExchangeRates(t *testing.T) {
	initBankWithExchangeRates()

	oneDollar := stocks.NewMoney(1, "USD")
	oneEuro := stocks.NewMoney(1, "EUR")
	oneWon := stocks.NewMoney(1, "KRW")
	var portfolio stocks.Portfolio
	portfolio = portfolio.Add(oneDollar)
	portfolio = portfolio.Add(oneEuro)
	portfolio = portfolio.Add(oneWon)

	expectedErreMsg := "Missing exchange rate(s): [USD->Kalganid,EUR->Kalganid,KRW->Kalganid]"
	value, actualErr := portfolio.Evaluate(bank, "Kalganid")
	asserNil(t, value)
	if actualErr.Error() != expectedErreMsg {
		t.Errorf("failed, expected %q, actual %q", expectedErreMsg, actualErr.Error())
	}

	assetEqual(t, expectedErreMsg, actualErr.Error())
}

// TestConversionWithDifferentRates 测试修改汇率前后转换货币是不是正确
func TestConversionWithDifferentRates(t *testing.T) {
	initBankWithExchangeRates()

	tenEuros := stocks.NewMoney(10, "EUR")
	actual, err := bank.Convert(tenEuros, "USD")
	expected := stocks.NewMoney(12, "USD")
	asserNil(t, err)
	assetEqual(t, *actual, expected)

	// 修改汇率后再兑换货币
	bank.AddExchangeRates("EUR", "USD", 1.3)
	actual, err = bank.Convert(tenEuros, "USD")
	expected = stocks.NewMoney(13, "USD")
	asserNil(t, err)
	assetEqual(t, *actual, expected)
}

func TestWhatIsCoversionRateFromEURToUSD(t *testing.T) {
	initBankWithExchangeRates()

	tenEuros := stocks.NewMoney(10, "EUR")
	actual, err := bank.Convert(tenEuros, "USD")
	expected := stocks.NewMoney(12, "USD")
	asserNil(t, err)
	assetEqual(t, *actual, expected)
}

func TestConversionWithMissingExchangeRate(t *testing.T) {
	initBankWithExchangeRates()

	tenEuros := stocks.NewMoney(10, "EUR")
	convertedMoney, err := bank.Convert(tenEuros, "Kalganid")
	asserNil(t, convertedMoney)
	expected := "EUR->Kalganid"
	assetEqual(t, expected, err.Error())
}

// asserNil 断言 actual 为 nil，包括接口中装着的带类型的 nil。
// any 是 interface{} 的别名，可以接收任意类型；reflect 用来检查其动态类型和值。
// 0、空字符串和已初始化的空切片等零值或空值不等于 nil，断言应失败。
func asserNil(t *testing.T, actual any) {
	// 将断言失败的位置定位到调用处。
	t.Helper()
	// 接口可以理解为“动态类型 + 动态值”，两者都没有时接口本身才等于 nil。
	// 先处理这种情况，因为 reflect.ValueOf(nil) 返回无效值，调用 IsNil 会 panic。
	if actual == nil {
		return
	}

	// 例如 var p *int = nil，传入 any 后动态类型为 *int、动态值为 nil，
	// 此时 actual == nil 为 false，需要用 IsNil 检查接口内的实际值。
	value := reflect.ValueOf(actual)
	// 只有允许为 nil 的类型才能调用 IsNil；int、string、struct 等类型会 panic。
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		if value.IsNil() {
			return
		}
	}

	t.Errorf("failed, expected nil, got: %+v", actual)
}

func assetEqual[T comparable](t *testing.T, expected, actual T) {
	if actual != expected {
		t.Errorf("failed, expected %+v, got: %+v\n", expected, actual)
	}
}
