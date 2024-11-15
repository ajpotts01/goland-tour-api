package todo

import (
	"context"
	"errors"
	"fmt"
	"github.com/ajpotts01/goland-tour-api/internal/db"
	"strings"
)

type Item struct {
	Task   string
	Status string
}

type Manager interface {
	InsertItem(ctx context.Context, item db.Item) error
	GetAllItems(ctx context.Context) ([]db.Item, error)
}

type Service struct {
	db Manager
}

func NewService(db Manager) *Service {
	return &Service{
		db: db,
	}
}

func (svc *Service) Add(todo string) error {
	existingItems, err := svc.GetAll()

	if err != nil {
		return fmt.Errorf("failed to get existing items from db: %w", err)
	}

	for _, t := range existingItems {
		if strings.ToLower(t.Task) == strings.ToLower(todo) {
			return errors.New("todo is not unique")
		}
	}

	err = svc.db.InsertItem(context.Background(), db.Item{
		Task:   todo,
		Status: "TO_BE_STARTED",
	})

	if err != nil {
		return fmt.Errorf("failed to insert item into db: %w", err)
	}

	return nil
}

func (svc *Service) Search(query string) ([]string, error) {
	var results []string

	existingItems, err := svc.GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to get existing items from db: %w", err)
	}

	for _, t := range existingItems {
		if strings.Contains(
			strings.ToLower(t.Task),
			strings.ToLower(query),
		) {
			results = append(results, t.Task)
		}
	}

	return results, nil
}

func (svc *Service) GetAll() ([]Item, error) {
	var results []Item

	items, err := svc.db.GetAllItems(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to read from db: %w", err)
	}

	for _, i := range items {
		results = append(results, Item{
			Task:   i.Task,
			Status: i.Status,
		})
	}
	return results, nil
}
