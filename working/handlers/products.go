package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/lptnkv/microservice_example/data"
)

type ProductsHandler struct {
	l *log.Logger
}

func NewProductsHandler(l *log.Logger) *ProductsHandler {
	return &ProductsHandler{l}
}

func (p *ProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	} else if r.Method == http.MethodPost {
		return
	} else if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProduct(id, w, r)
	}
}

func (p *ProductsHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marhsal json", http.StatusInternalServerError)
	}
}

func (p *ProductsHandler) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling POST Product")
	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarhsal json", http.StatusInternalServerError)
	}

	p.l.Printf("Prod: %#v\n", prod)
	data.AddProduct(prod)
}

func (p *ProductsHandler) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling PUT Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusInternalServerError)
	}
	data.UpdateProduct(id, prod)
}
