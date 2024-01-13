package handlers

import (
	"errors"
	"log"
	"net/http"
	"product-api/data"
	"product-api/exception"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
	e *exception.Error
}

func NewProducts(l *log.Logger, e *exception.Error) *Products {
	return &Products{l, e}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		p.l.Println("GET")

		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.getProducts(w)
			return
		}

		if len(g) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, _ := strconv.Atoi(g[0][1])

		err := p.getProduct(id, w)

		if errors.Is(err, data.ErrorProductNotFound) {
			err := p.e.ProductNotFoundError().ToJSON(w)
			if err != nil {
				return
			}
			return
		}

	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT")
		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, _ := strconv.Atoi(g[0][1])

		p.updateProduct(id, w, r)
	}

	if r.Method == http.MethodDelete {
		p.l.Println("DELETE")

		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, _ := strconv.Atoi(g[0][1])

		p.deleteProduct(id, w)

	}

}

func (p *Products) getProducts(w http.ResponseWriter) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) getProduct(id int, w http.ResponseWriter) error {
	p.l.Println("Handle GET Product")
	dp := data.GetProduct(id)
	if dp != nil {
		err := dp.ToJSON(w)
		if err != nil {
			return err
		}
	}

	if dp == nil {
		return data.ErrorProductNotFound
	}

	return nil
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Product: %#v", product)
	data.AddProduct(product)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")
	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Product: %#v", product)
	err = data.UpdateProducts(product, id)

	if errors.Is(err, data.ErrorProductNotFound) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {

		http.Error(w, "Invalid parameter", http.StatusInternalServerError)
		return

	}
}

func (p *Products) deleteProduct(id int, w http.ResponseWriter) {

	p.l.Println("handle DELETE Product")

	err := data.DeleteProduct(id)

	if errors.Is(err, data.ErrorProductNotFound) {
		err := p.e.ProductNotFoundError().ToJSON(w)
		if err != nil {
			http.Error(w, "Invalid parameter", http.StatusInternalServerError)
		}
		return
	}

	if err != nil {

		http.Error(w, "Invalid parameter", http.StatusInternalServerError)
		return

	}

}
