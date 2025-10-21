package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"phone_book/internal/models"
	"time"
)

type Storage struct {
	fileName string
}

func NewStorage(fileName string) Storage {
	return Storage{fileName: fileName}
}

func (s Storage) AcceptOrder(order models.Order) error {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		if errCreateFile := s.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}

	records = append(records, transform(order))

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s Storage) ReWrite(orders []models.Order) error {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		if errCreateFile := s.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}

	var records []orderRecord
	for _, order := range orders {
		records = append(records, transform(order))
	}

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)

}

func (s Storage) ListOrders(order models.Order) ([]models.Order, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	result := make([]models.Order, 0, len(records))
	for _, record := range records {
		if order.Customer_id == 0 || record.Customer_id == order.Customer_id {
			result = append(result, record.toDomain())
		}
	}
	return result, nil
}

func (s Storage) ReturnOrder(order models.Order) error {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}

	found := false
	for i, record := range records {
		if record.Order_id == order.Order_id {
			records[i].Deleted = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("заказ %d не найден", order.Order_id)
	}

	updatedData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.fileName, updatedData, 0644)
}

func (s Storage) IssueOrder(ordersIds []int) error {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}

	idSet := make(map[int]struct{}, len(ordersIds))
	for _, id := range ordersIds {
		idSet[id] = struct{}{}
	}

	marked := make(map[int]bool, len(ordersIds))
	for i := range records {
		if _, ok := idSet[records[i].Order_id]; ok {
			records[i].Issued = true
			records[i].Issued_date = time.Now().Format("2006-01-02")
			marked[records[i].Order_id] = true
		}
	}

	for _, id := range ordersIds {
		if !marked[id] {
			return fmt.Errorf("заказ %d не найден", id)
		}
	}

	updatedData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.fileName, updatedData, 0644)
}

func (s Storage) AcceptReturn(order models.Order) error {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}

	found := false
	for i, record := range records {
		if record.Order_id == order.Order_id {
			records[i].Returned = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("заказ %d не найден", order.Order_id)
	}

	updatedData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.fileName, updatedData, 0644)
}

func (s Storage) ListReturns(customer_id int, page, pageSize int) ([]models.Order, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	var returnedOrders []models.Order
	for _, record := range records {
		if record.Returned && (customer_id == 0 || record.Customer_id == customer_id) {
			returnedOrders = append(returnedOrders, record.toDomain())
		}
	}

	start := page * pageSize
	end := start + pageSize

	if start >= len(returnedOrders) {
		return []models.Order{}, nil
	}

	if end > len(returnedOrders) {
		end = len(returnedOrders)
	}

	return returnedOrders[start:end], nil
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
