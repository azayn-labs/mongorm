# Utility Types

MongORM includes helper functions for creating pointers to basic Go types. Because model fields use pointer types (e.g., `*string`, `*bool`), these helpers reduce boilerplate when setting field values.

## Import

```go
import "github.com/CdTgr/mongorm"
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
| `Timestamp(t time.Time) *time.Time` | `TimestampVal(t *time.Time) time.Time` | `time.Time{}` |

---

[Back to Documentation Index](./index.md) | [README](../README.md)
