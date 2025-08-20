package validator

import (
	"fmt"
	"testing"
)

func TestFetchBreeds(t *testing.T) {
	breeds, err := FetchBreeds()
	if err != nil {
		t.Fatal(err)
	}

	for _, b := range breeds {
		t.Log(b.Name)
	}

	fmt.Println("Total: ", len(breeds))
}

func TestValidateBreed(t *testing.T) {
	b := "Siam"
	isValid, err := ValidateBreed(b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(b, "is valid: ", isValid)
}
