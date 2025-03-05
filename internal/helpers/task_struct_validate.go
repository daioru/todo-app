package helpers

import (
	"fmt"

	"github.com/daioru/todo-app/internal/models"
)

func ValidateTaskFields(task *models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("title", "cannot be blank"))
	}

	if len(task.Title) > 100 {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("title", "field too long"))
	}

	if task.Status == "" {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("status", "cannot be blank"))
	}

	if len(task.Status) > 100 {
		return fmt.Errorf("validation failed: %w", NewSpecificValidationError("title", "status too long"))
	}

	return nil
}
