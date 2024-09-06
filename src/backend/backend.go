package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

func probe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running\n")
}

func (b *Backend) Run() {
	b.router.HandleFunc("/probe", probe)
	b.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var p Product
		rows, err := b.db.Query("SELECT id, name, inventory, price from products")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&p.id, &p.name, &p.inventory, &p.price)
			fmt.Fprintf(w, "Product #%2d: %20s %5d %5d\n", p.id, p.name, p.inventory, p.price)
		}
	})
	http.Handle("/", b.router)
	fmt.Println("Server started and listening on port ", b.addr)
	log.Fatal(http.ListenAndServe(b.addr, b.router))
}
