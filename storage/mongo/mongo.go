package mongo

import (
	"context"
	"fmt"
	"log"
	lib "recipe-book-bot/lib/error"
	"recipe-book-bot/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	recipes Recipes
}

type Recipes struct {
	*mongo.Collection
}

type Recipe struct {
	Name         string   `bjson:"name"`
	Description  string   `bjson:"description"`
	Ingredients  []string `bjson:"ingredients"`
	Instructions string   `bjson:"instructions"`
	Username     string   `bjson:"username"`
}

func New(connectSting string) Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(connectSting)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Panic("[DB] Error withing connection", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Panic("[DB] Error withing connection", err)
	}

	recipes := Recipes{
		client.Database("recipe-book").Collection("recipes"),
	}

	return Storage{
		recipes: recipes,
	}
}

func (s Storage) Save(ctx context.Context, recipe *storage.Recipe) (interface{}, error) {
	res, err := s.recipes.InsertOne(ctx, Recipe{
		Name:         recipe.Name,
		Description:  recipe.Description,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		Username:     recipe.Username,
	})
	if err != nil {
		return nil, lib.WrapErr("can't save recipe", err)
	}

	return res.InsertedID, nil
}

func (s Storage) Delete(ctx context.Context, name string, username string) error {
	res, err := s.recipes.DeleteOne(ctx, bson.M{
		"name":     name,
		"username": username,
	})
	if err != nil {
		return lib.WrapErr("can't remove recipe", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("can't find recipe with name=%s and username=%s", name, username)
	}

	return nil
}

func (s Storage) GetAllByUserName(ctx context.Context, username string) ([]*storage.Recipe, error) {
	cursor, err := s.recipes.Find(ctx, bson.M{
		"username": username,
	})
	if err != nil {
		return nil, lib.WrapErr("can't get recipes", err)
	}

	var results []Recipe
	if err = cursor.All(ctx, &results); err != nil {
		return nil, lib.WrapErr("can't decode recipes", err)
	}

	var r []*storage.Recipe

	for _, result := range results {
		r = append(r, &storage.Recipe{
			Name:         result.Name,
			Description:  result.Description,
			Ingredients:  result.Ingredients,
			Instructions: result.Instructions,
		})
	}

	return r, nil
}

func (s Storage) GetAllByUserNameTest(ctx context.Context, username string, page int64, recipesPerPage int64) ([]*storage.Recipe, error) {
	normalizedPage := page - 1
	if normalizedPage <= 0 {
		normalizedPage = 0
	}
	opts := options.Find().SetSkip(normalizedPage * recipesPerPage).SetLimit(recipesPerPage)
	cursor, err := s.recipes.Find(ctx, bson.M{
		"username": username,
	}, opts)
	if err != nil {
		return nil, lib.WrapErr("can't get recipes", err)
	}

	var results []Recipe
	if err = cursor.All(ctx, &results); err != nil {
		return nil, lib.WrapErr("can't decode recipes", err)
	}

	var r []*storage.Recipe

	for _, result := range results {
		r = append(r, &storage.Recipe{
			Name:         result.Name,
			Description:  result.Description,
			Ingredients:  result.Ingredients,
			Instructions: result.Instructions,
		})
	}

	return r, nil
}

func (s Storage) Exists(ctx context.Context, name string, username string) (bool, error) {
	var result Recipe
	err := s.recipes.FindOne(ctx, bson.M{
		"username": username,
		"name":     name,
	}).Decode(&result)
	if err != nil {
		return false, lib.WrapErr("can't find recipe", err)
	}

	return true, nil
}

func (s Storage) FindByName(ctx context.Context, name string, username string) (*storage.Recipe, error) {
	var result Recipe
	err := s.recipes.FindOne(ctx, bson.M{
		"username": username,
		"name":     name,
	}).Decode(&result)
	if err != nil {
		return nil, lib.WrapErr("can't find recipe", err)
	}

	return &storage.Recipe{
		Name:         result.Name,
		Description:  result.Description,
		Ingredients:  result.Ingredients,
		Instructions: result.Instructions,
	}, nil
}

func (r Recipe) Filter() bson.M {
	return bson.M{
		"name":     r.Name,
		"username": r.Username,
	}
}
