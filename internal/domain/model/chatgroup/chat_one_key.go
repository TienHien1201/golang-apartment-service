package chatgroup

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sort"
)

type ChatOneKey struct {
	UserIDs []int64 `json:"user_ids"`
}

func BuildChatOneKey(userIDs []int64) (*ChatOneKey, error) {
	if len(userIDs) != 2 {
		return nil, errors.New("chat one must have exactly 2 users")
	}

	sort.Slice(userIDs, func(i, j int) bool {
		return userIDs[i] < userIDs[j]
	})

	return &ChatOneKey{UserIDs: userIDs}, nil
}

func (k ChatOneKey) Value() (driver.Value, error) {
	return json.Marshal(k)
}

func (k *ChatOneKey) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid chat one key data")
	}

	return json.Unmarshal(bytes, k)
}
