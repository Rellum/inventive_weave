package json_test

import (
	"bytes"
	"github.com/Rellum/inventive_weave/pkg/json"
	"github.com/Rellum/inventive_weave/svc/creators/types"
	"os"
	"testing"
	"time"
)

func Test_parse_product(t *testing.T) {
	// Given
	inputJson := `{"Products": [{"id": "prod_2O5Yst3NQSf8b6xBQCVgi4KJw6p", "creatorId": "usr_2O5YssnJNfDynPxoCsT5hsulDwM", "createTime": "2023-04-06T21:01:59.752638+02:00"}]}`
	expected := time.Date(2023, 4, 6, 19, 1, 59, 752638000, time.UTC)

	// When
	data, err := json.Decode[types.Data](bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatalf("failed to parse product: %v", err)
	}

	// Then
	if len(data.Products) != 1 {
		t.Errorf("incorrect number of products: %d", len(data.Products))
	}
	if len(data.Creators) != 0 {
		t.Errorf("incorrect number of creators: %d", len(data.Creators))
	}
	if data.Products[0].Id != "prod_2O5Yst3NQSf8b6xBQCVgi4KJw6p" {
		t.Errorf("incorrect product id: %s", data.Products[0].Id)
	}
	if data.Products[0].CreatorId != "usr_2O5YssnJNfDynPxoCsT5hsulDwM" {
		t.Errorf("incorrect product creator id: %s", data.Products[0].CreatorId)
	}
	if expected.Equal(data.Products[0].CreateTime) == false {
		t.Errorf("incorrect product create time: %v != %v", expected, data.Products[0].CreateTime.UTC())
	}
}

func Test_parse_creator(t *testing.T) {
	// Given
	inputJson := `{"Creators": [{"id": "usr_2GFm8GQt4gcSrHrFjMk97NRe3tB", "email": "4gcSrHr@FjMk97NRe3tB.com"}]}`

	// When
	data, err := json.Decode[types.Data](bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatalf("failed to parse product: %v", err)
	}

	// Then
	if len(data.Products) != 0 {
		t.Errorf("incorrect number of products: %d", len(data.Products))
	}
	if len(data.Creators) != 1 {
		t.Errorf("incorrect number of creators: %d", len(data.Creators))
	}
	if data.Creators[0].Id != "usr_2GFm8GQt4gcSrHrFjMk97NRe3tB" {
		t.Errorf("incorrect creator id: %s", data.Creators[0].Id)
	}
	if data.Creators[0].Email != "4gcSrHr@FjMk97NRe3tB.com" {
		t.Errorf("incorrect product creator id: %s", data.Creators[0].Email)
	}
}

func Test_parse_file(t *testing.T) {
	// Given
	f, err := os.Open("../../data/example1.json")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	defer f.Close()

	// When
	data, err := json.Decode[types.Data](f)
	if err != nil {
		t.Fatalf("failed to parse product: %v", err)
	}

	// Then
	if len(data.Products) != 100 {
		t.Errorf("incorrect number of products: %d", len(data.Products))
	}
	if len(data.Creators) != 10 {
		t.Errorf("incorrect number of creators: %d", len(data.Creators))
	}
}
