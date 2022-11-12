package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

	//How to get value from context if it is already set in middleware handler func
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	fmt.Printf("CHECK CHECK:%#v", prod)
	//nicer using %#v
	p.l.Printf("Prod:%#v", prod)
	data.AddProduct(&prod)

}
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert Id", http.StatusBadRequest)
	}

	p.l.Println("Handle a update/put request")

	//How to get value from context if it is already set in middleware handler func
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return
	}

}

// [IMP]:type KeyProduct struct : missed the {} and took almost 1 hour to figure out the reason
type KeyProduct struct{}

// Small code that's why we are writting in handler otherwise we need to create this is other file
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	//[IMP]:return http.HandlerFunc(rw http.ResponseWriter, r *http.Request) : showing error..missing func(rw )
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		bytes, _ := ioutil.ReadAll(r.Body)

		p.l.Println("receivied body", string(bytes))

		prod := data.Product{}

		//	err := prod.FromJSON(r.Body)
		err := json.Unmarshal(bytes, &prod) //this is working.

		if err != nil {
			p.l.Println("[ERROR] deserializing Product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}
		//add the product to the context ???????/

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

		//[IMP] : previously I was not assigning r.WithContext(ctx) to r, and thus not getting value set with context.WithValue
		r = r.WithContext(ctx)

		//Call the next Handler, which can be another middleware in the chain, or the final Handler
		next.ServeHTTP(rw, r)

	})

}
