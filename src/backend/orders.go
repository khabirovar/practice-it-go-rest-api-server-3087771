package backend

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customer_name"`
	Total        int         `json:"total"`
	Status       string      `json:"status"`
	Products     map[int]int `json:"products"`
}

func getOrders(db *sql.DB) ([]order, error) {
	query := `
		SELECT id, customerName, total, status FROM orders
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []order{}
	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Total, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, err
}

func (o *order) getProducts(db *sql.DB) ([]product, error) {
	query :=
		`
		SELECT products.id,
			   products.productCode,
			   products.name,
			   products.inventory,
			   products.price,
			   products.status
		FROM products
		JOIN order_items
		ON order_items.product_id=products.id
		JOIN orders
		ON order_items.order_id=orders.id
		WHERE order_items.order_id = ?
	`
	rows, err := db.Query(query, o.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, err
}

func (o *order) newOrder(db *sql.DB) error {
	query :=
		`
		INSERT into orders(customerName, total, status)
		VALUES (?,?,?)
	`
	res, err := db.Exec(query, o.CustomerName, o.Total, o.Status)
	if err != nil {
		return nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil
	}
	o.ID = int(id)

	query =
		`
		INSERT INTO order_items(order_id, product_id)
		VALUES(?,?,?)
	`
	for key, val := range o.Products {
		_, err = db.Exec(query, o.ID, key, val)
		if err != nil {
			return err
		}
	}
	return nil
}
