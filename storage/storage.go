package storage

import "context"

type Storage interface {
	Save(ctx context.Context, p *Recipe) (interface{}, error)
	Delete(ctx context.Context, name string, username string) error
	Exists(ctx context.Context, name string, username string) (bool, error)
	GetTotalRecipesCount(ctx context.Context, username string) (int64, error)
	GetAllByUserName(ctx context.Context, username string, page int64, recipesPerPage int64) ([]*Recipe, error)
	FindByName(ctx context.Context, name string, username string) (*Recipe, error)
}

type Recipe struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
	Username     string   `json:"username"`
}
