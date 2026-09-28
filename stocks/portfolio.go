package stocks

type Portfolio []Money

func (p Portfolio) Add(m Money) Portfolio {
	return append(p, m)
}

func (p Portfolio) Evaluate(currency string) Money {
	rate := 1.2
	total := 0.0
	for _, m := range p {
		if m.currency != currency {
			total += m.amount * rate
		} else {
			total += m.amount
		}
	}
	return Money{total, currency}
}
