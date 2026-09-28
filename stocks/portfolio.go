package stocks

import (
	"errors"
	"strings"
)

type Portfolio []Money

func (p Portfolio) Add(m Money) Portfolio {
	return append(p, m)
}

func (p Portfolio) Evaluate(currency string) (Money, error) {
	// 	expectedErreMsg := "Missing exchange rate(s): [USD->Kalganid,EUR->Kalganid,KRW->Kalganid,]"

	failures := []string{}
	total := 0.0
	for _, m := range p {
		amount, ok := convert(m, currency)
		if ok {
			total += amount
		} else {
			failures = append(failures, m.currency+"->"+currency)
		}
	}

	if len(failures) == 0 {
		return Money{total, currency}, nil
	}

	missing := strings.Join(failures, ",")
	msg := "Missing exchange rate(s): " +"[" + missing + "]"

	return Money{0, ""}, errors.New(msg)
}

// convert 把 Money 转成指定币种对应的金额
func convert(m Money, currency string) (float64, bool) {
	if m.currency == currency {
		return m.amount, true
	}

	exchangeRates := map[string]float64{
		"EUR->USD": 1.2,
		"USD->KRW": 1100,
	}

	key := m.currency + "->" + currency
	rate, ok := exchangeRates[key]

	return rate * m.amount, ok
}
