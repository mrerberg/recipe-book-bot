package utils

import (
	"regexp"
)

func ExtractIngredients(text string) []string {
	ingredients := make([]string, 0)
	regex := regexp.MustCompile(`(?m)^\- (.+?)$`)
	matches := regex.FindAllStringSubmatch(text, -1)

	if len(matches) >= 1 {
		for _, match := range matches {
			ingredients = append(ingredients, match[1])
		}
	}

	return ingredients
}
