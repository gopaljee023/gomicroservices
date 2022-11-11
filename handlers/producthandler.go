package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gopaljee023/gomicroservices/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	//lp listof produt
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle a post request")

	prod := &data.Product{} // correct
	bytes, _ := ioutil.ReadAll(r.Body)

	p.l.Println("resceivied body", string(bytes))

	//err := prod.FromJSON(r.Body) //not working ..don't know why

	err := json.Unmarshal(bytes, prod) //this is working.

	if err != nil {
		p.l.Println("Unable to unmarshal json: will report ui")
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	//nicer using %#v
	p.l.Printf("Prod:%#v", prod)
	data.AddProduct(prod)

}
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert Id", http.StatusBadRequest)
	}

	p.l.Println("Handle a update/put request")
	prod := &data.Product{} // correct
	bytes, _ := ioutil.ReadAll(r.Body)

	p.l.Println("resceivied body", string(bytes))

	//err := prod.FromJSON(r.Body) //not working ..don't know why

	err = json.Unmarshal(bytes, prod)
	if err != nil {
		p.l.Println("Unable to unmarshal json: will report ui")
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return
	}

}

func (p *Products) deleteProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle a delete request")
}
