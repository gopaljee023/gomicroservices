package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	//delete method $curl localhost:9090 -XDELETE $curl localhost:9090 -XDELETE -v
	if req.Method == http.MethodDelete {
		p.deleteProducts(rw, req)
		return
	}

	//put method
	if req.Method == http.MethodPut {
		//extrac the id from the URI..  using mux .. it is easy
		reg := regexp.MustCompile(`/([0-9]+)`)
		groups := reg.FindAllStringSubmatch(req.URL.Path, -1)

		if len(groups) != 1 {
			p.l.Println("Novalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		fmt.Printf("capture groups %q\n", groups)
		//Q. why there will be two capture group
		if len(groups[0]) != 2 {
			p.l.Println("Novalid URI more than capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := groups[0][1]

		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Not able to ")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		fmt.Println("received id is ", id)
		p.updateProducts(id, rw, req)
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
	data.AddProduct(prod)

}
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle a update/put request")
	prod := &data.Product{} // correct
	bytes, _ := ioutil.ReadAll(r.Body)

	p.l.Println("resceivied body", string(bytes))

	//err := prod.FromJSON(r.Body) //not working ..don't know why

	err := json.Unmarshal(bytes, prod)
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
