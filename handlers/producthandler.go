package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gopaljee023/gomicroservices/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//Get: $curl localhost:9090/ -v
	if req.Method == http.MethodGet {
		p.getProducts(rw, req)
		return
	}

	//POST:curl -d '{ "id":1,"name":"tea","price": 3.4,"sku":"23211dd"}' localhost:9090
	if req.Method == http.MethodPost {
		p.addProducts(rw, req)

		return
	}

	//$curl localhost:9090 -d "yourdata"
	if req.Method == http.MethodPut {
		p.updateProducts(rw, req)
		return
	}
	//delte: $curl localhost:9090 -XDELETE $curl localhost:9090 -XDELETE -v
	if req.Method == http.MethodDelete {
		p.deleteProducts(rw, req)
		return
	}
	//catch all other
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {

	//lp listof produt
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {

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

}
func (p *Products) updateProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle a update/put request")
}

func (p *Products) deleteProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle a delete request")
}
