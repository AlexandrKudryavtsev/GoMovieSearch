package elastic

import (
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

const (
	_defaultAddress      = "http://localhost:9200"
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Elastic struct {
	addresses    []string
	username     string
	password     string
	connAttempts int
	connTimeout  time.Duration

	Client *elasticsearch.Client
}

func New(opts ...Option) (*Elastic, error) {
	e := &Elastic{
		addresses:    []string{_defaultAddress},
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(e)
	}

	var err error
	cfg := elasticsearch.Config{
		Addresses: e.addresses,
		Username:  e.username,
		Password:  e.password,
	}

	for e.connAttempts > 0 {
		e.Client, err = elasticsearch.NewClient(cfg)
		if err == nil {
			_, err = e.Client.Ping()
			if err == nil {
				break
			}
		}

		log.Printf("Elasticsearch is trying to connect, attempts left: %d", e.connAttempts)
		time.Sleep(e.connTimeout)
		e.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("elastic - NewElastic - connAttempts == 0: %w", err)
	}

	return e, nil
}
