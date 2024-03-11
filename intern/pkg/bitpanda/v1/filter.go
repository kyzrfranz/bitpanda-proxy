package v1

import (
	"fmt"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"time"
)

//TODO make this generic

type TxOption func(*TxFilter)

type MetaFilter struct {
	PageSize int
}

type TxFilter struct {
	Status     string
	Type       string
	Cursor     string
	CoinSymbol string
	From       int
	To         int
	IsSavings  bool
}

func (f TxFilter) AppendToReq(r *http.Request) {
	q := r.URL.Query()
	if f.Status != "" {
		q.Add("status", f.Status)
	}

	if f.Type != "" {
		q.Add("type", f.Type)
	}

	if f.Cursor != "" {
		q.Add("cursor", f.Cursor)
	}

	if f.CoinSymbol != "" {
		q.Add("coin_symbol", f.CoinSymbol)
	}

	if f.IsSavings != false {
		q.Add("is_savings", strconv.FormatBool(f.IsSavings))
	}

	if f.From != 0 && f.To != 0 {
		q.Add("from", fmt.Sprintf("%d", f.From))
		q.Add("to", fmt.Sprintf("%d", f.To))
	}
	r.URL.RawQuery = q.Encode()
}

func (f TxFilter) Apply(data []Transaction) []Transaction {
	return lo.Filter(data, func(item Transaction, index int) bool {
		pass := true
		if f.CoinSymbol != "" {
			pass = pass && item.Attributes.CryptocoinSymbol == f.CoinSymbol
		}
		if pass && f.IsSavings {
			pass = pass && item.Attributes.Trade.Attributes.IsSavings
		}
		if pass && (f.From != 0 || f.To != 0) {
			if f.To == 0 {
				f.To = int(time.Now().Unix())
			}
			dt, err := strconv.Atoi(item.Attributes.Time.Unix)
			pass = err == nil && pass && dt >= f.From && dt <= f.To
		}
		return pass
	})
}

func FilterFromReq(r *http.Request) []TxOption {
	var options []TxOption

	query := r.URL.Query()

	// For each filter option, check if it's present and set it in the filter
	if status := query.Get("status"); status != "" {
		options = append(options, WithStatus(status))
	}
	if txType := query.Get("type"); txType != "" {
		options = append(options, WithType(txType))
	}
	if cursor := query.Get("cursor"); cursor != "" {
		options = append(options, WithCursor(cursor))
	}
	if coinSymbol := query.Get("coin_symbol"); coinSymbol != "" {
		options = append(options, WithCoinSymbol(coinSymbol))
	}
	if isSavings := query.Get("is_savings"); isSavings != "" {
		if bVal, err := strconv.ParseBool(isSavings); err == nil {
			options = append(options, WithIsSavings(bVal))
		}
	}

	var msFrom, msTo int
	if from := query.Get("from"); from != "" {
		if fromInt, err := strconv.Atoi(from); err == nil {
			msFrom = fromInt
		}
	}
	if to := query.Get("to"); to != "" {
		if toInt, err := strconv.Atoi(to); err == nil {
			msTo = toInt
		}
	}

	options = append(options, WithDateRange(msFrom, msTo))

	return options
}

func WithStatus(status string) TxOption {
	return func(filter *TxFilter) {
		filter.Status = status
	}
}

func WithIsSavings(isSavings bool) TxOption {
	return func(filter *TxFilter) {
		filter.IsSavings = isSavings
	}
}

func WithType(txType string) TxOption {
	return func(filter *TxFilter) {
		filter.Type = txType
	}
}

func WithCursor(cursor string) TxOption {
	return func(filter *TxFilter) {
		filter.Cursor = cursor
	}
}

func WithCoinSymbol(coinSymbol string) TxOption {
	return func(filter *TxFilter) {
		filter.CoinSymbol = coinSymbol
	}
}

func WithDateRange(from, to int) TxOption {
	return func(filter *TxFilter) {
		filter.From = from
		filter.To = to
	}
}
