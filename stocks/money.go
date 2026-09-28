package stocks

type Money struct {
	amount float64
	currency string
}

func NewMoney(amount float64, currency string) Money {
	return Money{ amount: amount, currency: currency }
}

func (m Money) Times(x float64) Money {
	return Money{
		amount: m.amount * x,
		currency: m.currency,
	}
}

func (m Money) Divide(x float64) Money {
	return Money{
		amount: m.amount / x,
		currency: m.currency,
	}
}