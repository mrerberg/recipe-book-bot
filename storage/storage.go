package storage

import "context"

type Storage interface {
	Save(ctx context.Context, p *Recipe) (interface{}, error)
	Delete(ctx context.Context, name string, username string) error
	// Update(ctx context.Context, p *Recipe) error
	Exists(ctx context.Context, name string, username string) (bool, error)
	GetAllByUserName(ctx context.Context, username string) ([]*Recipe, error)
	FindByName(ctx context.Context, name string, username string) (*Recipe, error)
}

type Recipe struct {
	Name         string
	Description  string
	Ingredients  []string
	Instructions string
	Username     string
}
