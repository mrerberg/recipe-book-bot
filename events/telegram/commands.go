package telegram

import (
	"context"
	"fmt"
	"log"
	"recipe-book-bot/events"
	"recipe-book-bot/storage"
	"regexp"
	"strings"
)

const (
	StartCmd    = "/start"
	HelpCmd     = "/help"
	AllCmd      = "/all"
	TemplateCmd = "/template"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	// log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.saveRecipe(ctx, chatID, text, username)
	}

	if isDeleteCmd(text) {
		return p.deleteRecipe(ctx, chatID, text, username)
	}

	if isGetCmd(text) {
		return p.getRecipe(ctx, chatID, text, username)
	}

	fmt.Printf("CONTEXT %v", ctx)

	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	case AllCmd:
		return p.sendAll(ctx, chatID, username)
	case TemplateCmd:
		return p.tg.SendMessage(chatID, TemplateMessage)
	default:
		return p.tg.SendMessage(chatID, UnknownCommand)
	}
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, MsgHelp)
}

func isAddCmd(text string) bool {
	return strings.Contains(strings.ToLower(text), "блюдо")
}

func isDeleteCmd(text string) bool {
	one := strings.Contains(strings.ToLower(text), "удалить")
	two := strings.Contains(strings.ToLower(text), "удали")

	return one || two
}

func isGetCmd(text string) bool {
	return !events.IsCommand(text) && !isAddCmd(text) && !isDeleteCmd(text)
}

func (p *Processor) saveRecipe(ctx context.Context, chatID int, rawRecipe string, username string) error {
	recipe := parseRecipe(rawRecipe)
	recipe.Username = username
	fmt.Printf("%+v\n", recipe)

	exists, err := p.storage.Exists(ctx, recipe.Name, username)
	if err != nil {
		// TODO: log
	}
	if exists == true {
		return p.tg.SendMessage(chatID, "Данный рецепт был сохранен ранее ✋")
	}

	kek, err := p.storage.Save(ctx, recipe)
	if err != nil {
		// TODO: log
		return p.tg.SendMessage(chatID, "Не удалось сохранить рецепт!")
	}

	log.Printf("Saved new recipe %v", kek)

	return p.tg.SendMessage(chatID, fmt.Sprintf(RecipeSaved, recipe.Name))
}

func parseRecipe(text string) *storage.Recipe {
	recipe := storage.Recipe{}

	// Extract recipe name
	nameRegexp := regexp.MustCompile(`Блюдо:\s+(.*)`)
	nameMatch := nameRegexp.FindStringSubmatch(text)
	if len(nameMatch) > 1 {
		recipe.Name = nameMatch[1]
	}

	// Extract ingredients
	ingredients := make([]string, 0)
	regex := regexp.MustCompile(`(?m)^\- (.+?)$`)
	matches := regex.FindAllStringSubmatch(text, -1)
	if len(matches) >= 1 {
		for _, match := range matches {
			ingredients = append(ingredients, match[1])
		}

		recipe.Ingredients = ingredients
	}

	// Extract instructions
	instructionsRegexp := regexp.MustCompile(`Процесс:\n(.+)`)
	instructionsMatch := instructionsRegexp.FindStringSubmatch(text)
	fmt.Printf("instructionsMatch %v \n", instructionsMatch)
	if len(instructionsMatch) > 1 {
		recipe.Instructions = instructionsMatch[1]
	}

	return &recipe
}

func (p *Processor) sendAll(ctx context.Context, chatID int, username string) error {
	recipes, err := p.storage.GetAllByUserName(ctx, username)
	if err != nil {
		// TODO: ...
	}

	if len(recipes) == 0 {
		p.tg.SendMessage(chatID, NoRecipes)
	}

	var str string

	for idx, recipe := range recipes {
		str = fmt.Sprintf("%s\n%d. %s", str, idx+1, recipe.Name)
	}

	return p.tg.SendMessage(chatID, str)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, MsgHello)
}

func (p *Processor) deleteRecipe(ctx context.Context, chatID int, text string, username string) error {
	regex := regexp.MustCompile(`\s+(.*)`)
	match := regex.FindStringSubmatch(text)

	if len(match) == 0 {
		log.Printf("[ERR] could not delete recipe %v from user %v", text, username)

		return p.tg.SendMessage(chatID, CanNotDeleteRecipes)
	}

	recipeName := match[1]

	err := p.storage.Delete(ctx, recipeName, username)
	if err != nil {
		log.Printf("[ERR] could not delete recipe %v from user %v", recipeName, username)

		return err
	}

	msg := fmt.Sprintf(RecipeDeleted, recipeName)
	return p.tg.SendMessage(chatID, msg)
}

func (p *Processor) getRecipe(ctx context.Context, chatID int, recipeName string, username string) error {
	recipe, err := p.storage.FindByName(ctx, recipeName, username)
	if err != nil {
		log.Printf("[ERR] could not find recipe %v from user %v", recipeName, username)

		return p.tg.SendMessage(chatID, fmt.Sprintf(NotFoundRecipe, recipeName))
	}

	var result []string
	for _, ingredient := range recipe.Ingredients {
		result = append(result, "- "+ingredient)
	}

	// msg := fmt.Sprintf(RecipeDeleted, recipeName)
	return p.tg.SendMessage(
		chatID,
		fmt.Sprintf(TemplateMessage2, recipe.Name, strings.Join(result, "\n"),
			recipe.Instructions),
	)
}
