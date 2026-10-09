package stocks

import "errors"

type Bank struct {
	exchangeRates map[string]float64
}

func NewBank() Bank {
	return Bank{ exchangeRates: make(map[string]float64) }
}

func (b Bank) AddExchangeRates(currencyFrom, currencyTo string, rate float64)  {
	key := currencyFrom + "->" + currencyTo
	b.exchangeRates[key] = rate
}

func (b Bank) Convert(m Money, currency string) (*Money, error) {
	if m.currency == currency {
		return &m, nil
	}

	key := m.currency + "->" + currency
	
	if rate, ok := b.exchangeRates[key]; ok {
		money := NewMoney(m.amount * rate, currency)
		return &money, nil
	} else {
		return nil, errors.New(key)
	}
}