package mongorm

import "time"

// Utility functions for working with pointers to basic types. These functions are used to
// create pointers to values of basic types (e.g., string, int, bool) and to retrieve the
// values from pointers, providing a convenient way to work with optional fields in the
// schema. These functions are commonly used when defining the schema struct for a MongORM
// instance, allowing you to easily create pointers to values and retrieve their values when
// needed.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	}
//	todo := &ToDo{Text: mongorm.String("Buy milk")}
func String(s string) *string {
	str := string(s)
	return &str
}

// StringVal retrieves the value from a pointer to a string. If the pointer is nil, it returns
// an empty string. This function is useful for safely accessing the value of optional string
// fields in the schema without having to check for nil pointers every time.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	}
//	todo := &ToDo{Text: mongorm.String("Buy milk")}
//	textValue := mongorm.StringVal(todo.Text) // textValue will be "Buy milk"
func StringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Similar utility functions for other basic types (e.g., int, bool) can be defined in the
// same way.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	  Done *bool   `bson:"done"`
//	}
//	todo := &ToDo{
//	    Text: mongorm.String("Buy milk"),
//	    Done: mongorm.Bool(false),
//	}
func Bool(b bool) *bool {
	return &b
}

// BoolVal retrieves the value from a pointer to a bool. If the pointer is nil, it returns
// false. This function is useful for safely accessing the value of optional bool fields in the
// schema without having to check for nil pointers every time.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	  Done *bool   `bson:"done"`
//	}
//	todo := &ToDo{
//	    Text: mongorm.String("Buy milk"),
//	    Done: mongorm.Bool(false),
//	}
//	doneValue := mongorm.BoolVal(todo.Done) // doneValue will be false
func BoolVal(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Similar utility functions for other basic types (e.g., int, uint) can be defined in the
// same way.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	  Done *bool   `bson:"done"`
//	  Priority *int64 `bson:"priority"`
//	}
//	todo := &ToDo{
//	    Text: mongorm.String("Buy milk"),
//	    Done: mongorm.Bool(false),
//	    Priority: mongorm.Int64(1),
//	}
func Int64(i int64) *int64 {
	return &i
}

// Int64Val retrieves the value from a pointer to an int64. If the pointer is nil, it returns
// 0. This function is useful for safely accessing the value of optional int64 fields in the
// schema without having to check for nil pointers every time.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	  Done *bool   `bson:"done"`
//	  Priority *int64 `bson:"priority"`
//	}
//	todo := &ToDo{
//	    Text: mongorm.String("Buy milk"),
//	    Done: mongorm.Bool(false),
//	    Priority: mongorm.Int64(1),
//	}
//	priorityValue := mongorm.Int64Val(todo.Priority) // priorityValue will be 1
func Int64Val(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Float64 creates a pointer to a float64 value. This function is useful for creating
// pointers to float values when defining schema structs for a MongORM instance.
//
// Example usage:
//
//	type Product struct {
//	  Price *float64 `bson:"price"`
//	}
//	product := &Product{Price: mongorm.Float64(19.99)}
func Float64(f float64) *float64 {
	return &f
}

// Float64Val retrieves the value from a pointer to a float64. If the pointer is nil, it
// returns 0. This function is useful for safely accessing optional float64 fields.
//
// Example usage:
//
//	type Product struct {
//	  Price *float64 `bson:"price"`
//	}
//	product := &Product{Price: mongorm.Float64(19.99)}
//	priceValue := mongorm.Float64Val(product.Price) // priceValue will be 19.99
func Float64Val(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// Timestamp creates a pointer to a time.Time value. This function is useful for creating
// pointers to time values when defining the schema struct for a MongORM instance, allowing you
// to easily manage timestamp fields in your documents.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	  TaskTime *time.Time `bson:"taskTime"`
//	}
//	todo := &ToDo{
//	    Text: mongorm.String("Buy milk"),
//	    TaskTime: mongorm.Timestamp(time.Now()),
//	}
func Timestamp(t time.Time) *time.Time {
	return &t
}

// TimestampVal retrieves the value from a pointer to a time.Time. If the pointer is nil, it
// returns the zero value of time.Time. This function is useful for safely accessing the value of
// optional time fields in the schema without having to check for nil pointers every time.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	  TaskTime *time.Time `bson:"taskTime"`
//	}
//	todo := &ToDo{
//	    Text: mongorm.String("Buy milk"),
//	    TaskTime: mongorm.Timestamp(time.Now()),
//	}
//	taskTimeValue := mongorm.TimestampVal(todo.TaskTime) // taskTimeValue will be the time value of TaskTime
func TimestampVal(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
