package module

import (
	"errors"
	"phone_book/internal/models"
)

type Storage interface {
	AddContact(telephone models.Telephone) error
	ListContact() ([]models.Telephone, error)
	ReWrite(telephones []models.Telephone) error
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

func (m Module) AddContact(telephone models.Telephone) error {
	return m.Storage.AddContact(telephone)
}

func (m Module) ListContact() ([]models.Telephone, error) {
	return m.Storage.ListContact()
}

func (m Module) DeleteContact(telephone models.Telephone) error {
	contacts, err := m.Storage.ListContact()
	if err != nil {
		return err
	}
	set := make(map[models.Name]models.Telephone, len(contacts))
	for _, contact := range contacts {
		set[contact.Name] = contact
	}

	_, ok := set[telephone.Name]
	if !ok {
		return nil
	}
	delete(set, telephone.Name)

	newContacts := make([]models.Telephone, 0, len(set))
	for _, value := range set {
		newContacts = append(newContacts, value)
	}
	return m.Storage.ReWrite(newContacts)
}

func (m Module) FindContact(telephone models.Telephone) (models.Telephone, error) {
	contacts, err := m.Storage.ListContact()
	if err != nil {
		return models.Telephone{}, err
	}

	set := make(map[models.Name]models.Telephone, len(contacts))
	for _, contact := range contacts {
		set[contact.Name] = contact
	}

	contact, ok := set[telephone.Name]
	if !ok {
		return models.Telephone{}, errors.New("contact not found")
	}

	return contact, nil
}
