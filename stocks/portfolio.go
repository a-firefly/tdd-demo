package stocks

import (
	"errors"
	"strings"
)

type Portfolio []Money

func (p Portfolio) Add(m Money) Portfolio {
	return append(p, m)
}

func (p Portfolio) Evaluate(bank Bank, currency string) (*Money, error) {
	// 	expectedErreMsg := "Missing exchange rate(s): [USD->Kalganid,EUR->Kalganid,KRW->Kalganid,]"

	failures := []string{}
	total := 0.0
	for _, m := range p {
		if toMoney, err := bank.Convert(m, currency); err == nil {
			total += toMoney.amount
		} else {
			failures = append(failures, m.currency+"->"+currency)
		}
	}

	if len(failures) == 0 {
		return &Money{total, currency}, nil
	}

	missing := strings.Join(failures, ",")
	msg := "Missing exchange rate(s): " + "[" + missing + "]"

	return nil, errors.New(msg)
}
