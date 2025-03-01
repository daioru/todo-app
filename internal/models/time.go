package models

import (
	"fmt"
	"time"
)

type JSONTime time.Time

// Метод для маршалинга JSON
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}
