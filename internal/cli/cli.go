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
				description: "выдать заказ клиенту: использование issueOrder --orders_id=[id1, id2, id3, ...]",
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
				description: "получить список возвратов: использование listReturns --customer_id=someId",
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

// func (c CLI) deleteContact(args []string) error {
// 	var name string

// 	fs := flag.NewFlagSet(deleteContact, flag.ContinueOnError)
// 	fs.StringVar(&name, "name", "", "use --name=SomeName")
// 	if err := fs.Parse(args); err != nil {
// 		return err
// 	}

// 	if len(name) == 0 {
// 		return errors.New("name is empty")
// 	}

// 	return c.Module.DeleteContact(models.Telephone{
// 		Name: models.Name(name),
// 	})
// }

// func (c CLI) findContact(args []string) error {
// 	var name string

// 	fs := flag.NewFlagSet(findContact, flag.ContinueOnError)
// 	fs.StringVar(&name, "name", "", "use --name=SomeName")

// 	if err := fs.Parse(args); err != nil {
// 		return err
// 	}

// 	if len(name) == 0 {
// 		return errors.New("name is empty")
// 	}

// 	telephone, err := c.Module.FindContact(models.Telephone{
// 		Name: models.Name(name),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("Имя: %7s\nТелефон: %s\n", telephone.Name, telephone.Telephone)
// 	return nil
// }

func (c CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
}
