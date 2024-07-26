package activity

import (
	"encoding/json"
	"github.com/Rellum/inventive_weave/svc/creators/types"
	"time"
)

func Parse(b []byte) (*types.Data, error) {
	var d types.Data
	err := json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func MostActive(data types.Data) ([]string, error) {
	agg := make(map[string]struct {
		email    string
		products int
		latest   time.Time
	})
	for i := range data.Products {
		agg[data.Products[i].CreatorId]++
	}
}
