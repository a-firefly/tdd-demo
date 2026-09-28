package main

type Portfolie []Money

func (p Portfolie) Add(m Money) Portfolie {
	return append(p, m)
}

func (p Portfolie) Evaluate(currency string) Money {
	rate := float32(1.2)
	var total float32
	for _, m := range p {
		if m.Currency != currency {
			total += m.Amount * rate
		} else {
			total += m.Amount
		}
	}
	return Money{total, currency}
}