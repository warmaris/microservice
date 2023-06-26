package exibillia

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// model is struct representing database relation.
// It is used for sqlx Scan and distinguished from main entity. So, we can separate business logic and
// outgoing port's implementation.
type model struct {
	ID          uint64 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Tags        tags   `db:"tags"`

	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type tags []string

func (t tags) Value() (driver.Value, error) {
	v, err := json.Marshal(t)
	if err != nil {
		return []byte{}, err
	}

	return v, nil
}

func (t *tags) Scan(src interface{}) error {
	var source []byte
	switch s := src.(type) {
	case string:
		source = []byte(s)
	case []byte:
		if len(s) == 0 {
			source = []byte("[]")
		} else {
			source = s
		}
	case nil:
		*t = []string{}
	default:
		return fmt.Errorf("incompatible type for tags")
	}

	return json.Unmarshal(source, t)
}
