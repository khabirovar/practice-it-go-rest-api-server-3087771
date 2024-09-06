package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	db   *sql.DB
	addr string
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
	return &Backend{db: db, addr: addr}, nil
}

func probe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running\n")
}

func (b *Backend) Run() {
	http.HandleFunc("/probe", probe)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Server started and listening on port ", b.addr)
	log.Fatal(http.ListenAndServe(b.addr, nil))
}
