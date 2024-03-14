package v1

import (
	"encoding/json"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Service interface {
	Transactions(apiKey string, options ...TxOption) (*Response[Transaction], error)
}

type service struct {
	apiKey string
	logger *slog.Logger
	cache  *ttlcache.Cache[string, Response[Transaction]]
}

func NewService(logger *slog.Logger, cacheDuration time.Duration) Service {
	cache := ttlcache.New[string, Response[Transaction]](
		ttlcache.WithTTL[string, Response[Transaction]](cacheDuration),
		//ttlcache.WithLoader[string, Response[Transaction]](loadTX), //TODO think about it
	)
	return &service{logger: logger, cache: cache}
}

func (s *service) Transactions(apiKey string, options ...TxOption) (*Response[Transaction], error) {
	s.apiKey = apiKey
	filter := TxFilter{}

	var data *Response[Transaction]

	cacheItem := s.cache.Get("transactions")

	if cacheItem != nil && !cacheItem.IsExpired() {
		s.logger.Debug("hit the cache", "expiration", cacheItem.ExpiresAt())
		data = ptr[Response[Transaction]](cacheItem.Value())
	} else {
		d, err := s.get(1, filter)
		if err != nil {
			return nil, err
		}
		if len(d.Errors) > 0 {
			return nil, d.Errors[0]
		}
		d, err = s.get(d.Meta.TotalCount, filter)
		if err != nil {
			return nil, err
		}
		s.cache.Set("transactions", *d, ttlcache.DefaultTTL)
		data = d
	}

	for _, option := range options {
		option(&filter)
	}

	if data.Data != nil {
		filteredData := filter.Apply(data.Data)
		data.Meta.Count = len(filteredData)
		data.Data = filteredData
	}

	return data, nil
}

func (s *service) get(pageSize int, filter TxFilter) (*Response[Transaction], error) {
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

func ptr[T any](t T) *T {
	return &t
}
