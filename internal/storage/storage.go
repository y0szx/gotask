package storage

import (
	"encoding/json"
	"errors"
	"os"
	"phone_book/internal/models"
)

type Storage struct {
	fileName string
}

func NewStorage(fileName string) Storage {
	return Storage{fileName: fileName}
}

func (s Storage) AddContact(telephone models.Telephone) error {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		if errCreateFile := s.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records []telephoneRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}

	records = append(records, transform(telephone))

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s Storage) ReWrite(telephones []models.Telephone) error {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		if errCreateFile := s.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}

	var records []telephoneRecord
	for _, telephone := range telephones {
		records = append(records, transform(telephone))
	}

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)

}

func (s Storage) ListContact() ([]models.Telephone, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	var records []telephoneRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	result := make([]models.Telephone, 0, len(records))
	for _, record := range records {
		result = append(result, record.toDomain())
	}
	return result, nil
}

func (s Storage) createFile() error {
	f, err := os.Create(s.fileName)
	os.WriteFile(s.fileName, []byte("[]"), 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
