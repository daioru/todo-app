package models

import (
	"fmt"
	"time"
)

type JSONTime time.Time


func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, time.Time(t).UTC().Format(time.RFC3339))
	return []byte(formatted), nil
}
