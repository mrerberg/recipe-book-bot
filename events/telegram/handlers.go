package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	tgClient "recipe-book-bot/clients/telegram"
	lib "recipe-book-bot/lib/error"
	"recipe-book-bot/lib/utils"
	"recipe-book-bot/storage"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

func (p *Processor) sendHelp(_ context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, helpMsg)
}

func (p *Processor) sendHello(_ context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, helloMsg)
}

func (p *Processor) startRecipeSave(_ context.Context, chatID int, username string) error {
	data, _ := json.Marshal(struct{}{})
	err := p.cache.Add(&memcache.Item{Key: username, Value: data})
	if err != nil {
		return err
	}

	return p.tg.SendMessage(chatID, recipeNameMsg)
}

func (p *Processor) saveRecipe(ctx context.Context, chatID int, rawRecipe string, username string) error {
	cachedRecipe, err := p.cache.Get(username)
	if err != nil {
		return lib.WrapErr("saveRecipe", err)
	}

	var res storage.Recipe
	err = json.Unmarshal(cachedRecipe.Value, &res)
	if err != nil {
		return lib.WrapErr("saveRecipe", err)
	}

	// TODO: handle duplicates
	if res.Name == "" {
		data, _ := json.Marshal(storage.Recipe{
			Name:     rawRecipe,
			Username: username,
		})
		err = p.cache.Set(&memcache.Item{Key: username, Value: data})
		if err != nil {
			return lib.WrapErr("saveRecipe", err)
		}
		return p.tg.SendMessage(chatID, recipeDescriptionMsg)
	}

	if res.Name != "" && res.Description == "" {
		data, _ := json.Marshal(storage.Recipe{
			Name:        res.Name,
			Username:    res.Username,
			Description: rawRecipe,
		})
		err = p.cache.Set(&memcache.Item{Key: username, Value: data})
		if err != nil {
			return lib.WrapErr("saveRecipe", err)
		}
		return p.tg.SendMessage(chatID, ingredientsMsg)
	}

	if res.Name != "" && res.Description != "" && len(res.Ingredients) == 0 {
		data, _ := json.Marshal(storage.Recipe{
			Name:        res.Name,
			Username:    res.Username,
			Description: res.Description,
			Ingredients: utils.ExtractIngredients(rawRecipe),
		})
		err = p.cache.Set(&memcache.Item{Key: username, Value: data})
		if err != nil {
			return lib.WrapErr("saveRecipe", err)
		}
		return p.tg.SendMessage(chatID, recipeProcessMsg)
	}

	if res.Name != "" && res.Description != "" && len(res.Ingredients) > 0 && res.Instructions == "" {
		_, err = p.storage.Save(ctx, &storage.Recipe{
			Name:         res.Name,
			Username:     res.Username,
			Description:  res.Description,
			Ingredients:  res.Ingredients,
			Instructions: rawRecipe,
		})
		if err != nil {
			// TODO: log
			return p.tg.SendMessage(chatID, notSavedRecipeMsg)
		}

		err = p.cache.Delete(username)
		if err != nil {
			return lib.WrapErr("saveRecipe", err)
		}

		return p.tg.SendMessage(chatID, fmt.Sprintf(recipeSavedMsg, res.Name))
	}

	return p.tg.SendMessage(chatID, "Ошибка. Попробуй добавить рецепт позже")
}

func (p *Processor) sendAll(ctx context.Context, chatID int, username string) error {
	recipes, err := p.storage.GetAllByUserName(ctx, username)
	if err != nil {
		return err
	}

	if len(recipes) == 0 {
		return p.tg.SendMessage(chatID, noRecipesMsg)
	}

	var str string

	for idx, recipe := range recipes {
		str = fmt.Sprintf("%s\n%d. %s", str, idx+1, recipe.Name)
	}

	return p.tg.SendMessage(chatID, str)
}

func (p *Processor) deleteRecipe(ctx context.Context, chatID int, recipeName string, username string) error {
	err := p.storage.Delete(ctx, recipeName, username)
	if err != nil {
		log.Printf("[ERR] could not delete recipe %v from user %v. Error: %v", recipeName, username, err)

		return err
	}

	msg := fmt.Sprintf(recipeDeletedMsg, recipeName)
	return p.tg.SendMessage(chatID, msg)
}

func (p *Processor) getRecipe(ctx context.Context, chatID int, recipeName string, username string) error {
	recipe, err := p.storage.FindByName(ctx, recipeName, username)
	if err != nil {
		log.Printf("[ERR] could not find recipe %v from user %v", recipeName, username)

		return p.tg.SendMessage(chatID, fmt.Sprintf(notFoundRecipeMsg, recipeName))
	}

	inlineKeyboard := tgClient.InlineKeyboard{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: "Показать рецепт", CallbackData: fmt.Sprintf("cb:get:%s", recipeName)},
				{Text: "Удалить рецепт", CallbackData: fmt.Sprintf("cb:delete:%s", recipeName)},
			},
		},
	}

	return p.tg.SendMessageWithMarkup(
		chatID,
		fmt.Sprintf("Блюдо: %v", recipe.Name),
		&inlineKeyboard,
	)
}

func (p *Processor) getFullRecipe(ctx context.Context, chatID int, recipeName string, username string) error {
	recipe, err := p.storage.FindByName(ctx, recipeName, username)
	if err != nil {
		log.Printf("[ERR] could not find recipe %v from user %v. Error: %v", recipeName, username, err)

		return p.tg.SendMessage(chatID, fmt.Sprintf(notFoundRecipeMsg, recipeName))
	}

	var result []string
	for _, ingredient := range recipe.Ingredients {
		result = append(result, "- "+ingredient)
	}

	return p.tg.SendMessage(
		chatID,
		fmt.Sprintf(templateMsg, recipe.Name, strings.Join(result, "\n"),
			recipe.Instructions),
	)
}
