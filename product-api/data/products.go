package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (ps *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ps)
}

// correct
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          1,
		Name:        "Espressoo",
		Description: "Short and strong coffee without milk",
		Price:       3.45,
		SKU:         "bacc4444",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
	fmt.Printf("after addtion: %#v", p)
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProduct(id int, p *Product) (err error) {

	_, index, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[index] = p
	return nil

}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for index, p := range productList {
		if p.ID == id {
			return p, index, nil
		}
	}
	//errors.New() error
	return nil, -1, ErrProductNotFound
}
