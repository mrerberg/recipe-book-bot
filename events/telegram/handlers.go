package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

func (p *Processor) sendAll(ctx context.Context, chatID int, username string, page int64, messageID int) error {
	var recipesPerPage int64 = 6

	allRecipes, err := p.storage.GetAllByUserName(ctx, username)
	if err != nil {
		return err
	}

	recipes, err := p.storage.GetAllByUserNameTest(ctx, username, page, recipesPerPage)
	if err != nil {
		return err
	}

	if len(allRecipes) == 0 {
		return p.tg.SendMessage(chatID, noRecipesMsg)
	}

	recipeButtons := make([][]tgClient.InlineKeyboardButton, 0, len(recipes))
	chunkSize := 2

	for i := 0; i < len(recipes); i += chunkSize {
		end := i + chunkSize
		if end > len(recipes) {
			end = len(recipes)
		}
		chunk := recipes[i:end]
		items := make([]tgClient.InlineKeyboardButton, len(chunk))
		for j, r := range chunk {
			items[j] = tgClient.InlineKeyboardButton{Text: r.Name, CallbackData: fmt.Sprintf("cb:get:%v", r.Name)}
		}
		recipeButtons = append(recipeButtons, items)
	}

	totalPagesNum := math.Ceil(float64(len(allRecipes)) / float64(recipesPerPage))
	totalCount := fmt.Sprintf("%v/%v", page, totalPagesNum)

	nexPageValue := page
	if page >= int64(totalPagesNum) {
		nexPageValue = 1
	} else {
		nexPageValue++
	}

	prevPageValue := page
	if page-1 <= 0 {
		prevPageValue = int64(totalPagesNum)
	} else {
		prevPageValue--
	}

	inlineKeyboard := tgClient.InlineKeyboard{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{},
	}

	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, recipeButtons...)
	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []tgClient.InlineKeyboardButton{
		{Text: "←", CallbackData: fmt.Sprintf("cb:getall:%v", prevPageValue)},
		{Text: totalCount, CallbackData: "_"},
		{Text: "→", CallbackData: fmt.Sprintf("cb:getall:%v", nexPageValue)},
	})

	err = p.tg.SendMessageWithMarkup(chatID, "Ваши рецепты:", &inlineKeyboard)
	if err != nil {
		return err
	}

	return p.tg.DeleteMessage(chatID, messageID)
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
