package storage

import (
	"phone_book/internal/models"
	"time"
)

type telephoneRecord struct {
	Name      string    `json:"name"`
	Telephone string    `json:"telephone"`
	CreatedAt time.Time `json:"created_at"`
}

func (t telephoneRecord) toDomain() models.Telephone {
	return models.Telephone{
		Name:      models.Name(t.Name),
		Telephone: models.Number(t.Telephone),
	}
}

func transform(telephone models.Telephone) telephoneRecord {
	return telephoneRecord{
		Name:      string(telephone.Name),
		Telephone: string(telephone.Telephone),
		CreatedAt: time.Now(),
	}
}
