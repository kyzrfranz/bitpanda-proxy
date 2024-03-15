package v1

import (
	"crypto/md5"
	"encoding/hex"
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

// TxFilter provides various filter options for BP transactions.
// Some of them are handled by the BP API, some will be handled by this proxy.
type TxFilter struct {
	Status     string // handled by BP
	Type       string // handled by BP
	Cursor     string // handled by BP
	CoinSymbol string // custom
	From       int64  // custom
	To         int64  // custom
	IsSavings  bool   // custom
}

// AppendToReq builds the query for the BP Api for the params it will allow.
func (f *TxFilter) AppendToReq(r *http.Request) {
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

	r.URL.RawQuery = q.Encode()
}

func (f *TxFilter) hash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s-%s-%s", f.Cursor, f.Status, f.Type)))
	return hex.EncodeToString(hash[:])
}

// Apply will do a filter on params the BP API does not provide on the results.
func (f *TxFilter) Apply(data []Transaction) []Transaction {
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
				f.To = time.Now().Unix()
			}
			dt := item.Attributes.Time.Unix
			pass = pass && dt >= f.From && dt <= f.To
		}
		return pass
	})
}

// FilterFromReq parses the request to this proxy and extracs all the passed params into the filter.
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

	var msFrom, msTo int64
	if from := query.Get("from"); from != "" {
		if fromInt, err := strconv.ParseInt(from, 10, 64); err == nil {
			msFrom = fromInt
		}
	}
	if to := query.Get("to"); to != "" {
		if toInt, err := strconv.ParseInt(to, 10, 64); err == nil {
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

func WithDateRange(from, to int64) TxOption {
	return func(filter *TxFilter) {
		filter.From = from
		filter.To = to
	}
}
