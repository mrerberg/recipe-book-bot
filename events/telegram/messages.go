package telegram

const helpMsg = `I can save and keep you recipes.

In order to save the recipe, just send me it using template from /template command.

To delete recipe send me "Удалить <Название рецепта>"

To get list of all recipes use /all command
`

const helloMsg = "Hi ✌️\n\n" + helpMsg

const templateMsg = `
------
Блюдо: %s

Ингредиенты: 
%s

Процесс:
%s
------
`

const ingredientsMsg = `
Какие ингридиенты?

Перечислите по порядку:
- Ингридиент 1
- Ингридиент 2
...
`

const tryAgainLaterMsg = "Sorry. I could not do it right now 😢. Please, try again later"

const (
	recipeNameMsg        = "Название блюда или рецепта?"
	recipeDescriptionMsg = "Какое описание блюда?"
	recipeProcessMsg     = "Каков процесс приготовления?"
)

const (
	UnknownCommand      = "I don't recognize command"
	NoRecipes           = "You don't have saved recipes yet 😢"
	existingRecipeMsg   = "Не могу сохранить данный рецепт, так как рецепт уже был сохранен ранее ✋"
	RecipeSaved         = "I saved your recipe «%s» 👌"
	RecipeDeleted       = "I deleted your recipe «%s» 👌"
	CanNotDeleteRecipes = "Something went wrong while deleting recipe 😢"
	NotFoundRecipe      = "I can't find recipe «%s» 😢"
	notSavedRecipe      = "Не удалось сохранить рецепт 😢. Попробуйте еще раз позже"
)

const (
	UnknownCb = "Sorry. I could not do it right now 😢. Please, try again later"
)
