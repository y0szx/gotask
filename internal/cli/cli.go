package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"phone_book/internal/models"
)

type Module interface {
	AddContact(telephone models.Telephone) error
	ListContact() ([]models.Telephone, error)
	DeleteContact(telephone models.Telephone) error
	FindContact(telephone models.Telephone) (models.Telephone, error)
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
				name:        addContact,
				description: "добавить контакт: использование add --name=SomeName --telephone=+79191111111",
			},
			{
				name:        listContact,
				description: "удалить контакт: использование list",
			},
			{
				name:        deleteContact,
				description: "удалить контакт: использование delete --name=SomeName",
			},
			{
				name:        findContact,
				description: "найти контакт: использование find --name=SomeName",
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
	case addContact:
		return c.addContact(args[1:])
	case listContact:
		return c.listContact()
	case deleteContact:
		return c.deleteContact(args[1:])
	case findContact:
		return c.findContact(args[1:])
	}
	return fmt.Errorf("command isn't set")
}

func (c CLI) addContact(args []string) error {
	var name, telephone string

	fs := flag.NewFlagSet(addContact, flag.ContinueOnError)
	fs.StringVar(&name, "name", "", "use --name=SomeName")
	fs.StringVar(&telephone, "telephone", "", "use --telephone=+71111111111")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(name) == 0 {
		return errors.New("name is empty")
	}

	if len(telephone) == 0 {
		return errors.New("telephone is empty")
	}
	return c.Module.AddContact(models.Telephone{
		Name:      models.Name(name),
		Telephone: models.Number(telephone),
	})
}

func (c CLI) listContact() error {
	list, err := c.Module.ListContact()
	if err != nil {
		return err
	}
	for _, telephone := range list {
		fmt.Printf("Имя: %7s\nТелефон: %s\n", telephone.Name, telephone.Telephone)
	}
	return nil
}

func (c CLI) deleteContact(args []string) error {
	var name string

	fs := flag.NewFlagSet(deleteContact, flag.ContinueOnError)
	fs.StringVar(&name, "name", "", "use --name=SomeName")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(name) == 0 {
		return errors.New("name is empty")
	}

	return c.Module.DeleteContact(models.Telephone{
		Name: models.Name(name),
	})
}

func (c CLI) findContact(args []string) error {
	var name string

	fs := flag.NewFlagSet(findContact, flag.ContinueOnError)
	fs.StringVar(&name, "name", "", "use --name=SomeName")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(name) == 0 {
		return errors.New("name is empty")
	}

	telephone, err := c.Module.FindContact(models.Telephone{
		Name: models.Name(name),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Имя: %7s\nТелефон: %s\n", telephone.Name, telephone.Telephone)
	return nil
}

func (c CLI) help() {
	fmt.Println("command list:")
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
}
