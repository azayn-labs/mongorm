package mongorm

import (
	"reflect"
)

func (m *Model[T]) Set(value *T) *Model[T] {
	if value == nil {
		return m
	}

	src := reflect.ValueOf(value).Elem()
	dst := reflect.ValueOf(m.clone).Elem()
	t := dst.Type()

	for i := 0; i < src.NumField(); i++ {
		df := dst.Field(i)
		if !df.CanSet() {
			continue
		}

		if hasModelFlag(t.Field(i).Tag, string(ModelTagPrimary)) ||
			hasModelFlag(t.Field(i).Tag, string(ModelTagReadonly)) {
			// These fields cannot be updated
			continue
		}

		sf := src.Field(i)
		if sf.Kind() == reflect.Pointer {
			if !sf.IsNil() {
				df.Set(sf)
			}
			continue
		}

		if !sf.IsZero() {
			df.Set(sf)
		}
	}

	return m
}

func (m *Model[T]) Unset(value *T) *Model[T] {
	if value == nil {
		return m
	}

	src := reflect.ValueOf(value).Elem()
	dst := reflect.ValueOf(m.clone).Elem()
	t := dst.Type()

	for i := 0; i < src.NumField(); i++ {
		df := dst.Field(i)
		if !df.CanSet() {
			continue
		}
		if hasModelFlag(t.Field(i).Tag, string(ModelTagPrimary)) ||
			hasModelFlag(t.Field(i).Tag, string(ModelTagReadonly)) {
			// These fields cannot be unset
			continue
		}

		sf := src.Field(i)
		if sf.Kind() == reflect.Pointer {
			if !sf.IsNil() {
				df.Set(reflect.Zero(df.Type()))
			}
			continue
		}

		if !sf.IsZero() {
			df.Set(reflect.Zero(df.Type()))
		}
	}

	return m
}
