package validator

import (
	"encoding/json"
	"net/http"
)

var cachedBreeds []Breed

type Breed struct {
	Id   string `json:id`
	Name string `json:name`
}

func FetchBreeds() ([]Breed, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}

	cachedBreeds = breeds

	return breeds, nil
}

func ValidateBreed(breed string) (bool, error) {
	if cachedBreeds != nil {
		for _, b := range cachedBreeds {
			if b.Name == breed {
				return true, nil
			}
		}
	}
	breeds, err := FetchBreeds()
	if err != nil {
		return false, err
	}

	for _, b := range breeds {
		if b.Name == breed {
			return true, nil
		}
	}

	return false, err
}
