package main

type Money struct {
	Amount float32
	Currency string
}

func (m Money) Times(x float32) Money {
	return Money{
		Amount: m.Amount * x,
		Currency: m.Currency,
	}
}

func (m Money) Divide(x float32) Money {
	return Money{
		Amount: m.Amount / x,
		Currency: m.Currency,
	}
}