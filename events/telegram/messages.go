package telegram

const MsgHelp = `I can save and keep you recipes.

In order to save the recipe, just send me it using template from /template command.

To delete recipe send me "Удалить <Название рецепта>"

To get list of all recipes use /all command
`

const MsgHello = "Hi ✌️\n\n" + MsgHelp

const TemplateMessage = `
Используйте данный шаблон, чтобы добавить ваше блюдо:

Блюдо: <Название блюда>

Ингредиенты: 
- <Ингредиент 1>
- <Ингредиент 2>
- ...

Процесс:
<Процесс готовки>
`

const TemplateMessage2 = `
Блюдо: %s

Ингредиенты: 
%s

Процесс:
%s
`

const (
	UnknownCommand      = "I don't recognize command"
	NoRecipes           = "You don't have saved recipes yet 😢"
	RecipeSaved         = "I saved your recipe «%s» 👌"
	RecipeDeleted       = "I deleted your recipe «%s» 👌"
	CanNotDeleteRecipes = "Something went wrong while deleting recipe 😢"
	NotFoundRecipe      = "I can't find recipe «%s» 😢"
)
