package lib

import "fmt"

func WrapErr(msg string, err error) error {
	if err == nil {
		return fmt.Errorf("%s", msg)
	}

	return fmt.Errorf("%s %w", msg, err)
}
