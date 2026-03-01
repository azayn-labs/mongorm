# Utility Types

MongORM includes helper functions for creating pointers to basic Go types. Because model fields use pointer types (e.g., `*string`, `*bool`), these helpers reduce boilerplate when setting field values.

## Import

```go
import "github.com/azayn-labs/mongorm"
```

## Functions

### String / StringVal

```go
func String(s string) *string
func StringVal(s *string) string
```

Create a `*string` or safely read its value (returns `""` for nil).

```go
todo := &ToDo{Text: mongorm.String("Buy groceries")}

text := mongorm.StringVal(todo.Text) // "Buy groceries"
nilText := mongorm.StringVal(nil)    // ""
```

### Bool / BoolVal

```go
func Bool(b bool) *bool
func BoolVal(b *bool) bool
```

Create a `*bool` or safely read its value (returns `false` for nil).

```go
task := &Task{
    Done: mongorm.Bool(false),
}

done := mongorm.BoolVal(task.Done) // false
nilDone := mongorm.BoolVal(nil)    // false
```

### Int64 / Int64Val

```go
func Int64(i int64) *int64
func Int64Val(i *int64) int64
```

Create a `*int64` or safely read its value (returns `0` for nil).

```go
task := &Task{
    Priority: mongorm.Int64(3),
}

priority := mongorm.Int64Val(task.Priority) // 3
nilPriority := mongorm.Int64Val(nil)        // 0
```

### Float64 / Float64Val

```go
func Float64(f float64) *float64
func Float64Val(f *float64) float64
```

Create a `*float64` or safely read its value (returns `0` for nil).

```go
product := &Product{
    Price: mongorm.Float64(19.99),
}

price := mongorm.Float64Val(product.Price) // 19.99
nilPrice := mongorm.Float64Val(nil)        // 0
```

### Decimal128 / Decimal128Val

```go
func Decimal128(d bson.Decimal128) *bson.Decimal128
func Decimal128Val(d *bson.Decimal128) bson.Decimal128
```

Create a `*bson.Decimal128` or safely read its value (returns zero `bson.Decimal128{}` for nil).

```go
amount, err := bson.ParseDecimal128("123.45")
if err != nil {
    panic(err)
}

invoice := &Invoice{Amount: mongorm.Decimal128(amount)}
amountValue := mongorm.Decimal128Val(invoice.Amount)
```

### Timestamp / TimestampVal

```go
func Timestamp(t time.Time) *time.Time
func TimestampVal(t *time.Time) time.Time
```

Create a `*time.Time` or safely read its value (returns `time.Time{}` zero value for nil).

```go
import "time"

task := &Task{
    DueDate: mongorm.Timestamp(time.Now().Add(24 * time.Hour)),
}

due := mongorm.TimestampVal(task.DueDate) // the actual time.Time value
nilDue := mongorm.TimestampVal(nil)       // time.Time{} zero value
```

## Summary Table

| Setter | Reader | Nil-safe default |
| --- | --- | --- |
| `String(s string) *string` | `StringVal(s *string) string` | `""` |
| `Bool(b bool) *bool` | `BoolVal(b *bool) bool` | `false` |
| `Int64(i int64) *int64` | `Int64Val(i *int64) int64` | `0` |
| `Float64(f float64) *float64` | `Float64Val(f *float64) float64` | `0` |
| `Decimal128(d bson.Decimal128) *bson.Decimal128` | `Decimal128Val(d *bson.Decimal128) bson.Decimal128` | `bson.Decimal128{}` |
| `Timestamp(t time.Time) *time.Time` | `TimestampVal(t *time.Time) time.Time` | `time.Time{}` |

## Model Output Helpers

In addition to pointer utilities, MongORM exposes helper methods on `*MongORM[T]` for accessing model data:

### Document

```go
func (m *MongORM[T]) Document() *T
```

Returns the underlying model pointer currently attached to the ORM instance.

```go
orm := mongorm.New(&ToDo{})
doc := orm.Document()
```

### JSON

```go
func (m *MongORM[T]) JSON(doc *T) (map[string]any, error)
```

Converts a model value into `map[string]any`, useful when you need a generic payload for logging or custom processing.

```go
orm := mongorm.New(&ToDo{})
payload, err := orm.JSON(&ToDo{Text: mongorm.String("hello")})
if err != nil {
    panic(err)
}

fmt.Println(payload)
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
