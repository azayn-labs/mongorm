package mongorm

import (
	"reflect"
	"time"
)

type Timestamps struct {
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (m *Model[T]) applyTimestamps() {
	v := reflect.ValueOf(m.clone).Elem()
	now := time.Now()

	// CreatedAt (set only if zero)
	if f := v.FieldByName("CreatedAt"); f.IsValid() && f.CanSet() {
		if f.Interface().(*time.Time) == nil || f.Interface().(*time.Time).IsZero() {
			f.Set(reflect.ValueOf(&now))
		}
	}

	// UpdatedAt (always update)
	if f := v.FieldByName("UpdatedAt"); f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(&now))
	}
}
