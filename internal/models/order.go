package models

type Order struct {
	Order_id    int
	Customer_id int
	Shelf_life  string
	Issued      bool
	Deleted     bool
	Returned    bool
}
