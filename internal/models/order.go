package models

type Order struct {
	Order_id    int
	Customer_id int
	Shelf_life  string
	Issued      bool
	Issued_date string
	Deleted     bool
	Returned    bool
}
