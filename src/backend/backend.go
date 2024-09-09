package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	db     *sql.DB
	addr   string
	router *mux.Router
}

type Product struct {
	id        int
	name      string
	inventory int
	price     int
}

func New(kind, path, addr string) (*Backend, error) {
	db, err := sql.Open(kind, path)
	if err != nil {
		return &Backend{}, err
	}
	router := mux.NewRouter()
	return &Backend{db: db, addr: addr, router: router}, nil
}

func (b *Backend) InitializeRoutes() {
	b.router.HandleFunc("/products", b.allProducts).Methods("GET")
	b.router.HandleFunc("/product/{id}", b.fetchProduct).Methods("GET")
	b.router.HandleFunc("/products", b.newProduct).Methods("POST")
	b.router.HandleFunc("/orders", b.allOrders).Methods("GET")
	b.router.HandleFunc("/order/{id}/products", b.allProductsOfOrder).Methods("GET")
	b.router.HandleFunc("/orders", b.newOrder).Methods("POST")
}

func (b *Backend) allProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(b.db)
	if err != nil {
		fmt.Printf("getProducts error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (b *Backend) fetchProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p product
	p.ID, _ = strconv.Atoi(id)
	err := p.getProduct(b.db)
	if err != nil {
		fmt.Printf("getProduct error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (b *Backend) newProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var p product
	json.Unmarshal(reqBody, &p)

	err := p.createProduct(b.db)
	if err != nil {
		fmt.Printf("createProduct error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (b *Backend) allOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := getOrders(b.db)
	if err != nil {
		fmt.Printf("getOrders: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, orders)
}

func (b *Backend) allProductsOfOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var o order
	o.ID, _ = strconv.Atoi(id)
	products, err := o.getProducts(b.db)
	if err != nil {
		fmt.Printf("allProductsOfOrder error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (b *Backend) newOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var o order
	err := json.Unmarshal(reqBody, &o)
	if err != nil {
		fmt.Printf("newOrder unmarshal error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = o.newOrder(b.db)
	if err != nil {
		fmt.Printf("newOrder error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, o)
}

func (b *Backend) Run() {
	fmt.Println("Server started and listening on port ", b.addr)
	log.Fatal(http.ListenAndServe(b.addr, b.router))
}

func probe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running\n")
}

// Helper functions
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
