package helpers

import (
	"fmt"
)

var allowedFields = map[string]bool{
	"id":          true, //Для указания какой Task обновлять
	"user_id":     true, //Для проверки пользователя, который обновляет
	"title":       true,
	"description": true,
	"status":      true,
}

func ValidateUpdates(updates map[string]interface{}) (map[string]interface{}, error) {
	_, ok := updates["id"]
	if !ok {
		return nil, fmt.Errorf("validation failed: %w", NewSpecificValidationError("id", "cannot be empty"))
	}

	_, ok = updates["user_id"]
	if !ok {
		return nil, fmt.Errorf("validation failed: %w", NewSpecificValidationError("user_id", "cannot be empty"))
	}

	if len(updates) <= 2 {
		return nil, fmt.Errorf("validation failed: %w", NewSpecificValidationError("", "no fields to update"))
	}

	return FilterAllowedFields(updates)
}

// Функция для фильтрации разрешенных полей
func FilterAllowedFields(updates map[string]interface{}) (map[string]interface{}, error) {
	validUpdates := make(map[string]interface{})

	for key, value := range updates {
		if _, ok := allowedFields[key]; ok {
			validUpdates[key] = value
		} else {
			return nil, fmt.Errorf("validation failed: %w", NewSpecificValidationError(key, "field not allowed"))
		}
	}

	return validUpdates, nil
}
