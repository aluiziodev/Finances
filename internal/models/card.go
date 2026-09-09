package models

import "fmt"

type Card struct {
	Id     string
	Name   string
	Bank   string
	Limit  float64
	DueDay int
}

func (c *Card) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("Name is required")
	}

	if c.Bank == "" {
		return fmt.Errorf("Bank is required")
	}

	if c.Limit < 0 {
		return fmt.Errorf("Limit cannot be negative")
	}

	if c.DueDay < 1 || c.DueDay > 28 {
		return fmt.Errorf("Due day must be between 1 and 28")
	}
	return nil
}
