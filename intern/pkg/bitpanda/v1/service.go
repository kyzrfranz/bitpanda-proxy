package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Service interface {
	Transactions(apiKey string, options ...TxOption) (*Response[Transaction], error)
}

type service struct {
	apiKey string
	logger *slog.Logger
}

func NewService(logger *slog.Logger) Service {
	return &service{logger: logger}
}

func (s service) Transactions(apiKey string, options ...TxOption) (*Response[Transaction], error) {
	s.apiKey = apiKey
	filter := TxFilter{}

	for _, option := range options {
		option(&filter)
	}

	data, err := s.get(1, filter)
	if err != nil {
		return nil, err
	}

	if data.Data != nil {
		data, err = s.get(data.Meta.TotalCount, filter)
		if err != nil {
			return nil, err
		}

		filteredData := filter.Apply(data.Data)
		data.Meta.Count = len(filteredData)
		data.Data = filteredData
	}

	return data, nil
}

func (s service) get(pageSize int, filter TxFilter) (*Response[Transaction], error) {
	client := &http.Client{}

	txUrl := fmt.Sprintf("%s/wallets/transactions?page_size=%d", BaseUrl, pageSize)
	s.logger.Debug("calling bitpanda api", "endpoint", "transactions", "url", txUrl)
	req, err := http.NewRequest("GET", txUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", s.apiKey)
	filter.AppendToReq(req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var txData Response[Transaction]
	err = json.Unmarshal(body, &txData)
	return &txData, err
}
