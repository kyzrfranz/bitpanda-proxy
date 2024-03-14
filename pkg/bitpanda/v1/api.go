package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const BaseUrl = "https://api.bitpanda.com/v1"

type Response[T any] struct {
	Errors []Error `json:"errors,omitempty"`
	Data   []T     `json:"data,omitempty"`
	Meta   *Meta   `json:"meta,omitempty"`
	Links  *Links  `json:"links,omitempty"`
}

type Error struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
	Title  string `json:"title"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Title)
}

type Links struct {
	Next string `json:"next"`
	Last string `json:"last"`
	Self string `json:"self"`
}

type Meta struct {
	TotalCount int    `json:"total_count"`
	Count      int    `json:"count"`
	PageSize   int    `json:"page_size"`
	Page       int    `json:"page"`
	PageNumber int    `json:"page_number"`
	NextCursor string `json:"next_cursor"`
}

type Transaction struct {
	Type       string       `json:"type"`
	Attributes TxAttributes `json:"attributes"`
	Id         string       `json:"id"`
}

type Trade struct {
	Type       string          `json:"type"`
	Attributes TradeAttributes `json:"attributes"`
	Id         string          `json:"id"`
}

type CommonAttributes struct {
	Status           string        `json:"status"`
	Type             string        `json:"type"`
	CryptocoinId     string        `json:"cryptocoin_id"`
	CryptocoinSymbol string        `json:"cryptocoin_symbol"`
	WalletId         string        `json:"wallet_id"`
	Time             Time          `json:"time"`
	IsSavings        bool          `json:"is_savings"`
	Tags             []interface{} `json:"tags"`
	IsCard           bool          `json:"is_card"`
}

type TxAttributes struct {
	CommonAttributes
	Amount            string `json:"amount"`
	Recipient         string `json:"recipient"`
	Confirmations     int    `json:"confirmations"`
	InOrOut           string `json:"in_or_out"`
	AmountEur         string `json:"amount_eur"`
	AmountEurInclFee  string `json:"amount_eur_incl_fee"`
	ConfirmationBy    string `json:"confirmation_by"`
	Confirmed         bool   `json:"confirmed"`
	Trade             Trade  `json:"trade"`
	LastChanged       Time   `json:"last_changed"`
	Fee               string `json:"fee"`
	CurrentFiatId     string `json:"current_fiat_id"`
	CurrentFiatAmount string `json:"current_fiat_amount"`
	IsMetalStorageFee bool   `json:"is_metal_storage_fee"`
	PublicStatus      string `json:"public_status"`
	IsBfc             bool   `json:"is_bfc"`
}

type TradeAttributes struct {
	CommonAttributes
	FiatId           string `json:"fiat_id"`
	AmountFiat       string `json:"amount_fiat"`
	AmountCryptocoin string `json:"amount_cryptocoin"`
	FiatToEurRate    string `json:"fiat_to_eur_rate"`
	FiatWalletId     string `json:"fiat_wallet_id"`
	Price            string `json:"price"`
	IsSwap           bool   `json:"is_swap"`
	BfcUsed          bool   `json:"bfc_used"`
}

type Time struct {
	DateIso8601 time.Time `json:"date_iso8601"`
	Unix        int64     `json:"unix"`
}

func (t *Time) UnmarshalJSON(data []byte) error {
	tmp := struct {
		DateIso8601 time.Time `json:"date_iso8601"`
		Unix        string    `json:"unix"` //TODO for whatever reasons the peeps at bitpanda decided to do this...
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	t.DateIso8601 = tmp.DateIso8601

	unixInt, err := strconv.ParseInt(tmp.Unix, 10, 64)
	if err != nil {
		return err
	}
	t.Unix = unixInt

	return nil
}
