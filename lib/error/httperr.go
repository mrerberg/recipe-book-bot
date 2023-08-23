package lib

import "fmt"

func WrapErr(msg string, err error) error {
	return fmt.Errorf("%s %w", msg, err)
}
