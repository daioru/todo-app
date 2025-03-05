package helpers

import "errors"

var allowedFields = map[string]bool{
	"id":          true, //Для указания какой Task обновлять
	"user_id":     true, //Для проверки пользователя, который обновляет
	"title":       true,
	"description": true,
	"status":      true,
}

func Validate(updates map[string]interface{}) (map[string]interface{}, error) {
	_, ok := updates["id"]
	if !ok {
		return nil, errors.New("no id for task provided")
	}

	_, ok = updates["user_id"]
	if !ok {
		return nil, errors.New("no user_id provided")
	}

	if len(updates) <= 2 {
		return nil, errors.New("no fields to update")
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
			return nil, errors.New("field not allowed")
		}
	}

	return validUpdates, nil
}
