package telegram

const helpMsg = `I can save and keep you recipes.

In order to save the recipe, just send me /add command.

To get list of all recipes use /all command
`

const helloMsg = "Hi ✌️\n\n" + helpMsg

const templateMsg = `
------
Meal: %s

Ingredients: 
%s

Process:
%s
------
`

const ingredientsMsg = `
What ingredients?

List them in order:
- Ingredient 1
- Ingredient 2
...
`

const (
	recipeNameMsg        = "What is the name of dish or recipe?"
	recipeDescriptionMsg = "Description?"
	recipeProcessMsg     = "What is the process?"
)

const (
	noRecipesMsg           = "You don't have saved recipes yet 😢"
	existingRecipeMsg      = "Cannot save this recipe as the recipe has already been saved previously ✋"
	recipeSavedMsg         = "I saved your recipe «%s» 👌"
	recipeDeletedMsg       = "I deleted your recipe «%s» 👌"
	canNotDeleteRecipesMsg = "Something went wrong while deleting recipe 😢"
	notFoundRecipeMsg      = "I can't find recipe «%s» 😢"
	notSavedRecipeMsg      = "Could not save the recipe 😢. Try again later"
)

const (
	unknownCommandMsg = "I don't recognize command"
	unknownCbMsg      = "Sorry. I could not do it right now 😢. Please, try again later"
)
