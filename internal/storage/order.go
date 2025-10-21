package storage

import (
	"phone_book/internal/models"
)

type orderRecord struct {
	Order_id    int    `json:"order_id"`
	Customer_id int    `json:"customer_id"`
	Shelf_life  string `json:"shelf_life"`
	Issued      bool   `json:"issued"`
	Issued_date string `json:"issued_date"`
	Deleted     bool   `json:"deleted"`
	Returned    bool   `json:"returned"`
}

func (o orderRecord) toDomain() models.Order {
	return models.Order{
		Order_id:    o.Order_id,
		Customer_id: o.Customer_id,
		Shelf_life:  o.Shelf_life,
		Issued:      o.Issued,
		Issued_date: o.Issued_date,
		Deleted:     o.Deleted,
		Returned:    o.Returned,
	}
}

func transform(order models.Order) orderRecord {
	return orderRecord{
		Order_id:    int(order.Order_id),
		Customer_id: int(order.Customer_id),
		Shelf_life:  string(order.Shelf_life),
		Issued:      bool(order.Issued),
		Issued_date: string(order.Issued_date),
		Deleted:     bool(order.Deleted),
		Returned:    bool(order.Returned),
	}
}
