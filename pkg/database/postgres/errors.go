package postgres

import "fmt"

func errorPreparingQuery(query string, err error) error {
	return fmt.Errorf("preparing '%s': %s", query, err.Error())
}

func errorExecutingQuery(query string, err error) error {
	return fmt.Errorf("executing '%s': %s", query, err.Error())
}
