package stocks

type Portfolio []Money

func (p Portfolio) Add(m Money) Portfolio {
	return append(p, m)
}

func (p Portfolio) Evaluate(currency string) Money {
	total := 0.0
	for _, m := range p {
		total += convert(m, currency)
	}
	return Money{total, currency}
}

// convert 把 Money 转成指定币种对应的金额
func convert(m Money, currency string) float64 {
	if m.currency == currency {
		return m.amount
	}

	exchangeRates := map[string]float64 {
		"EUR->USD": 1.2,
		"USD->KRW": 1100,
	}

	key := m.currency + "->" + currency
	rate := exchangeRates[key]

	return rate * m.amount
}
