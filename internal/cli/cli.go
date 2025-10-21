package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"phone_book/internal/models"
	"strconv"
	"strings"
)

type Module interface {
	AcceptOrder(order models.Order) error
	ReturnOrder(order models.Order) error
	ListOrders(order models.Order) ([]models.Order, error)
	IssueOrder(ordersId []int) error
	AcceptReturn(order models.Order) error
	ListReturns(customer_id int, page, pageSize int) ([]models.Order, error)
}

type Deps struct {
	Module Module
}

type CLI struct {
	Deps
	commandList []command
}

func NewCli(d Deps) CLI {
	return CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				description: "справка",
			},
			{
				name:        acceptOrder,
				description: "принять заказ от курьера: использование acceptOrder --order_id=someId --customer_id=someId --shelf_life=someDate",
			},
			{
				name:        returnOrder,
				description: "вернуть заказ курьеру: использование returnOrder --order_id=someId",
			},
			{
				name:        issueOrder,
				description: "выдать заказ клиенту: использование issueOrder --orders=1,2,3",
			},
			{
				name:        listOrders,
				description: "получить список заказов: использование listOrders --customer_id=someId (--last_n=n)",
			},
			{
				name:        acceptReturn,
				description: "принять возврат от клиента: использование acceptReturn --customer_id=someId --order_id=someId",
			},
			{
				name:        listReturns,
				description: "получить список возвратов: использование listReturns --customer_id=someId --page=0 --page_size=10",
			},
		},
	}
}

func (c CLI) Run() error {
	args := os.Args[1:]
	if len(args) == 0 {
		return fmt.Errorf("command isn't set")
	}

	commandName := args[0]
	switch commandName {
	case help:
		c.help()
		return nil
	case acceptOrder:
		return c.acceptOrder(args[1:])
	case returnOrder:
		return c.returnOrder(args[1:])
	case listOrders:
		return c.listOrders(args[1:])
	case issueOrder:
		return c.issueOrder(args[1:])
	case acceptReturn:
		return c.acceptReturn(args[1:])
	case listReturns:
		return c.listReturns(args[1:])
	}
	return fmt.Errorf("command isn't set")
}

func (c CLI) acceptOrder(args []string) error {
	var order_id, customer_id int
	var shelf_life string

	fs := flag.NewFlagSet(acceptOrder, flag.ContinueOnError)
	fs.IntVar(&order_id, "order_id", 0, "use --order_id=SomeId")
	fs.IntVar(&customer_id, "customer_id", 0, "use --customer_id=SomeId")
	fs.StringVar(&shelf_life, "shelf_life", "", "use --shelf_life=SomeDate")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if order_id == 0 {
		return errors.New("order_id is empty")
	}

	if customer_id == 0 {
		return errors.New("customer_id is empty")
	}

	if len(shelf_life) == 0 {
		return errors.New("shelf life is empty")
	}
	return c.Module.AcceptOrder(models.Order{
		Order_id:    order_id,
		Customer_id: customer_id,
		Shelf_life:  shelf_life,
		Issued:      false,
		Issued_date: "",
		Deleted:     false,
		Returned:    false,
	})
}

func (c CLI) returnOrder(args []string) error {
	var order_id int

	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.IntVar(&order_id, "order_id", 0, "use --order_id=SomeId")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if order_id == 0 {
		return errors.New("order_id is empty")
	}

	return c.Module.ReturnOrder(models.Order{
		Order_id: order_id,
	})
}

func (c CLI) listOrders(args []string) error {
	var customer_id int

	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.IntVar(&customer_id, "customer_id", 0, "use --customer_id=SomeId")

	if err := fs.Parse(args); err != nil {
		return err
	}

	list, err := c.Module.ListOrders(models.Order{
		Customer_id: customer_id,
	})
	if err != nil {
		return err
	}

	if customer_id == 0 {
		fmt.Println("Список всех заказов:")
	} else {
		fmt.Printf("Список заказов пользователя %d:\n", customer_id)
	}

	for _, order := range list {
		fmt.Printf("Заказ №%d, клиент %d\n", order.Order_id, order.Customer_id)
	}
	return nil
}

func (c CLI) issueOrder(args []string) error {
	var orders string

	fs := flag.NewFlagSet(issueOrder, flag.ContinueOnError)
	fs.StringVar(&orders, "orders", "", "use --orders=1,2,3")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(orders) == 0 {
		return errors.New("orders list is empty")
	}

	raw := strings.Split(orders, ",")
	if len(raw) == 0 {
		return errors.New("order list is empty")
	}

	ids := make([]int, 0, len(raw))
	for _, part := range raw {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.Atoi(part)
		if err != nil {
			return err
		}
		ids = append(ids, id)

	}

	if len(ids) == 0 {
		return errors.New("order list is empty")
	}

	return c.Module.IssueOrder(ids)
}

func (c CLI) acceptReturn(args []string) error {
	var customer_id, order_id int
	fs := flag.NewFlagSet(acceptReturn, flag.ContinueOnError)
	fs.IntVar(&customer_id, "customer_id", 0, "use --customer_id=someId")
	fs.IntVar(&order_id, "order_id", 0, "use --order_id=someId")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if order_id == 0 {
		return errors.New("order_id is empty")
	}

	if customer_id == 0 {
		return errors.New("customer_id is empty")
	}

	return c.Module.AcceptReturn(models.Order{
		Order_id:    order_id,
		Customer_id: customer_id,
	})
}

func (c CLI) listReturns(args []string) error {
	var customer_id int
	var page, pageSize int

	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.IntVar(&customer_id, "customer_id", 0, "use --customer_id=someId")
	fs.IntVar(&page, "page", 0, "use --page=0 (номер страницы, начиная с 0)")
	fs.IntVar(&pageSize, "page_size", 10, "use --page_size=10 (размер страницы)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if page < 0 {
		return errors.New("page must be >= 0")
	}

	if pageSize <= 0 {
		return errors.New("page_size must be > 0")
	}

	list, err := c.Module.ListReturns(customer_id, page, pageSize)
	if err != nil {
		return err
	}

	if customer_id == 0 {
		fmt.Printf("Список возвратов (страница %d, размер %d):\n", page, pageSize)
	} else {
		fmt.Printf("Список возвратов клиента %d (страница %d, размер %d):\n", customer_id, page, pageSize)
	}

	if len(list) == 0 {
		fmt.Println("Возвратов не найдено")
		return nil
	}

	for _, order := range list {
		fmt.Printf("Заказ №%d, клиент %d, дата выдачи: %s\n", order.Order_id, order.Customer_id, order.Issued_date)
	}

	return nil
}

func (c CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
}
