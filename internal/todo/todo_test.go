package todo_test

import (
	"context"
	"errors"
	"github.com/ajpotts01/goland-tour-api/internal/db"
	"github.com/ajpotts01/goland-tour-api/internal/todo"
	"reflect"
	"testing"
)

type MockDB struct {
	items []db.Item
}

func (m *MockDB) InsertItem(_ context.Context, item db.Item) error {
	m.items = append(m.items, item)
	return nil
}

func (m *MockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	return m.items, nil
}

func TestService_Search(t *testing.T) {
	tests := []struct {
		name       string
		todosToAdd []string
		query      string
		want       []string
	}{
		{
			name:       "given a todo of shop and a search of sh, i should get shop back",
			todosToAdd: []string{"shop"},
			query:      "sh",
			want:       []string{"shop"},
		},
		{
			name:       "still returns shop, even if the case doesn't match",
			todosToAdd: []string{"Shopping"},
			query:      "sh",
			want:       []string{"Shopping"},
		},
		{
			name:       "spaces",
			todosToAdd: []string{"go Shopping"},
			query:      "go",
			want:       []string{"go Shopping"},
		},
		{
			name:       "space at start of word",
			todosToAdd: []string{" Space at beginning"},
			query:      "space",
			want:       []string{" Space at beginning"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			for _, td := range tt.todosToAdd {
				err := svc.Add(td)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := svc.Search(tt.query)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Add_Duplicates(t *testing.T) {
	tests := []struct {
		name       string
		todosToAdd []string
		want       error
	}{
		{
			name:       "can't add duplicates",
			todosToAdd: []string{"shop", "shop"},
			want:       errors.New("todo is not unique"),
		},
		{
			name:       "can't add duplicates - casing",
			todosToAdd: []string{"shop", "Shop"},
			want:       errors.New("todo is not unique"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			var got error
			for _, td := range tt.todosToAdd {
				err := svc.Add(td)
				if err != nil {
					got = err
				}
			}

			if got == nil || got.Error() != tt.want.Error() {
				t.Errorf("Add() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}
