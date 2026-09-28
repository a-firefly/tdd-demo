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
	rate := 1.2
	if m.currency == currency {
		return m.amount
	}

	return rate * m.amount
}
