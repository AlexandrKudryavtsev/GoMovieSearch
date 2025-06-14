package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/entity"
	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/usecase"
	elastic "github.com/AlexandrKudryavtsev/GoMovieSearch/pkg/elasticsearch"
)

type moviesRepo struct {
	es *elastic.Elastic
}

func NewMoviesRepo(es *elastic.Elastic) usecase.MoviesRepo {
	return &moviesRepo{es: es}
}

func (r *moviesRepo) Index(ctx context.Context, data []entity.Movie) error {
	if err := r.CreateIndex(ctx); err != nil {
		return fmt.Errorf("failed to ensure index exists: %w", err)
	}

	var buf bytes.Buffer
	for _, movie := range data {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "movies",
				"_id":    movie.DetailsURL,
			},
		}

		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("failed to encode meta: %w", err)
		}

		if err := json.NewEncoder(&buf).Encode(movie); err != nil {
			return fmt.Errorf("failed to encode movie: %w", err)
		}
	}

	res, err := r.es.Client.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.es.Client.Bulk.WithContext(ctx),
		r.es.Client.Bulk.WithIndex("movies"),
	)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("bulk response error: %s, body: %s", res.Status(), string(body))
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if raw["errors"].(bool) {
		return fmt.Errorf("some documents failed to index: %v", raw["items"])
	}

	return nil
}

func (r *moviesRepo) Autocomplete(ctx context.Context, query string) ([]entity.Movie, error) {
	if len(query) < 2 {
		return nil, nil
	}

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"russianTitle^2", "originalTitle"},
				"fuzziness": "AUTO",
				"operator":  "or",
			},
		},
		"size": 5,
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := r.es.Client.Search(
		r.es.Client.Search.WithContext(ctx),
		r.es.Client.Search.WithIndex("movies"),
		r.es.Client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.Status())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source entity.Movie `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	movies := make([]entity.Movie, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		movies = append(movies, hit.Source)
	}

	return movies, nil
}

func (r *moviesRepo) Search(ctx context.Context, query string) ([]entity.Movie, error) {
	if len(query) < 2 {
		return nil, nil
	}

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"russianTitle": map[string]interface{}{
								"query":     query,
								"boost":     2,
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"match": map[string]interface{}{
							"originalTitle": map[string]interface{}{
								"query":     query,
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"match": map[string]interface{}{
							"year": query,
						},
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := r.es.Client.Search(
		r.es.Client.Search.WithContext(ctx),
		r.es.Client.Search.WithIndex("movies"),
		r.es.Client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.Status())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source entity.Movie `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	movies := make([]entity.Movie, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		movies = append(movies, hit.Source)
	}

	return movies, nil
}

func (r *moviesRepo) CreateIndex(ctx context.Context) error {
	mapping := `{
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 0,
            "analysis": {
                "analyzer": {
                    "ru_en_analyzer": {
                        "type": "custom",
                        "tokenizer": "standard",
                        "filter": ["lowercase", "russian_stop", "english_stop", "russian_stemmer", "english_stemmer"]
                    }
                },
                "filter": {
                    "russian_stop": {
                        "type": "stop",
                        "stopwords": "_russian_"
                    },
                    "english_stop": {
                        "type": "stop",
                        "stopwords": "_english_"
                    },
                    "russian_stemmer": {
                        "type": "stemmer",
                        "language": "russian"
                    },
                    "english_stemmer": {
                        "type": "stemmer",
                        "language": "english"
                    }
                }
            }
        },
        "mappings": {
            "properties": {
                "russianTitle": {
                    "type": "text",
                    "analyzer": "ru_en_analyzer"
                },
                "originalTitle": {
                    "type": "text",
                    "analyzer": "ru_en_analyzer"
                },
                "year": {
                    "type": "keyword"
                },
                "detailsUrl": {
                    "type": "keyword"
                },
                "posterUrl": {
                    "type": "keyword"
                }
            }
        }
    }`

	r.es.Client.Indices.Delete([]string{"movies"}, r.es.Client.Indices.Delete.WithIgnoreUnavailable(true))

	res, err := r.es.Client.Indices.Create("movies",
		r.es.Client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error creating index: %s, response: %s", res.Status(), string(body))
	}

	return nil
}
