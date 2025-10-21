package module

import (
	"fmt"
	"phone_book/internal/models"
	"time"
)

type Storage interface {
	AcceptOrder(order models.Order) error
	ReturnOrder(order models.Order) error
	ListOrders(order models.Order) ([]models.Order, error)
	IssueOrder(ordersIds []int) error
	AcceptReturn(order models.Order) error
	ListReturns(customer_id int, page, pageSize int) ([]models.Order, error)
	// ReWrite(telephones []models.Telephone) error
}

type Deps struct {
	Storage Storage
}

type Module struct {
	Deps
}

func NewModule(d Deps) Module {
	return Module{Deps: d}
}

func (m Module) AcceptOrder(order models.Order) error {
	targetDate, err := time.Parse("2006-01-02", order.Shelf_life)
	if err != nil {
		return fmt.Errorf("некорректный формат даты срока хранения: %v", err)
	}

	currentDate := time.Now().Truncate(24 * time.Hour)
	if targetDate.Before(currentDate) {
		return fmt.Errorf("срок хранения заказа %d в прошлом", order.Order_id)
	}

	existing_orders, err := m.Storage.ListOrders(models.Order{})
	if err != nil {
		return err
	}

	for _, e := range existing_orders {
		if e.Order_id == order.Order_id && e.Customer_id == order.Customer_id && !e.Deleted {
			return fmt.Errorf("заказ %d уже принят этим клиентом", order.Order_id)
		}

		if e.Order_id == order.Order_id && e.Customer_id != order.Customer_id && !e.Deleted {
			return fmt.Errorf("заказ %d уже принят другим клиентом (ID: %d)", order.Order_id, e.Customer_id)
		}
	}
	return m.Storage.AcceptOrder(order)
}

func (m Module) ReturnOrder(order models.Order) error {
	existing_orders, err := m.Storage.ListOrders(models.Order{})
	if err != nil {
		return err
	}

	for _, e := range existing_orders {
		if e.Order_id == order.Order_id {
			if e.Deleted {
				return fmt.Errorf("заказ №%d уже возвращен курьеру", e.Order_id)
			}
			if e.Issued {
				return fmt.Errorf("заказ №%d уже выдан и не может быть возвращен курьеру", order.Order_id)
			}

			targetDate, err := time.Parse("2006-01-02", e.Shelf_life)
			if err != nil {
				return fmt.Errorf("некорректный формат даты срока хранения: %v", err)
			}

			currentDate := time.Now().Truncate(24 * time.Hour)
			if targetDate.After(currentDate) {
				return fmt.Errorf("срок хранения заказа %d ещё не истёк", order.Order_id)
			}

			return m.Storage.ReturnOrder(order)
		}
	}

	return fmt.Errorf("заказ %d не найден", order.Order_id)
}

func (m Module) ListOrders(order models.Order) ([]models.Order, error) {
	return m.Storage.ListOrders(order)
}

func (m Module) IssueOrder(ordersIds []int) error {
	existing_orders, err := m.Storage.ListOrders(models.Order{})
	if err != nil {
		return err
	}

	idToOrder := make(map[int]models.Order, len(existing_orders))
	for _, o := range existing_orders {
		idToOrder[o.Order_id] = o
	}

	var customerId int
	for i, id := range ordersIds {
		order, ok := idToOrder[id]
		if !ok {
			return fmt.Errorf("заказ %d не найден", id)
		}

		if order.Deleted {
			return fmt.Errorf("заказ %d уже возвращен курьеру и не может быть выдан", id)
		}

		if order.Issued {
			return fmt.Errorf("заказ %d уже выдан", id)
		}

		if i == 0 {
			customerId = order.Customer_id
		} else if order.Customer_id != customerId {
			return fmt.Errorf("все заказы должны принадлежать только одному клиенту")
		}

		targetDate, err := time.Parse("2006-01-02", order.Shelf_life)
		if err != nil {
			return fmt.Errorf("некорректный формат даты срока хранения: %v", err)
		}

		currentDate := time.Now().Truncate(24 * time.Hour)
		if targetDate.Before(currentDate) {
			return fmt.Errorf("срок хранения заказа %d истёк", id)
		}
	}

	return m.Storage.IssueOrder(ordersIds)
}

func (m Module) AcceptReturn(order models.Order) error {
	existing_orders, err := m.Storage.ListOrders(order)
	if err != nil {
		return err
	}

	for _, e := range existing_orders {
		if e.Order_id == order.Order_id {
			if !e.Issued {
				return fmt.Errorf("заказ %d не был выдан клиенту %d", e.Order_id, e.Customer_id)
			}

			if e.Returned {
				return fmt.Errorf("заказ %d уже был возвращен", e.Order_id)
			}

			issuedDate, err := time.Parse("2006-01-02", e.Issued_date)
			if err != nil {
				return fmt.Errorf("некорректный формат даты выдачи: %v", err)
			}

			if time.Since(issuedDate) > 48*time.Hour {
				return fmt.Errorf("с момента выдачи заказа %d прошло больше двух дней", order.Order_id)
			}

			return m.Storage.AcceptReturn(order)
		}
	}

	return fmt.Errorf("заказ %d не найден", order.Order_id)
}

func (m Module) ListReturns(customer_id int, page, pageSize int) ([]models.Order, error) {
	return m.Storage.ListReturns(customer_id, page, pageSize)
}
