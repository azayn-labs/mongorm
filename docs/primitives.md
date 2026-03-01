# MongORM Primitives

Primitives are type-safe field wrappers used in your schema struct. Each field type maps to a Go type and exposes query-builder methods that return `bson.M` filters, ready to be passed to `Where()`.

## Importing

```go
import "github.com/azayn-labs/mongorm/primitives"
```

## Available Field Types

| Field Type | Go Type | Use for |
| --- | --- | --- |
| `ObjectIDField` | `bson.ObjectID` | MongoDB `_id` and foreign key fields |
| `StringField` | `string` | Text fields |
| `Int64Field` | `int64` | Integer numeric fields (also handles int32, int8, int) |
| `Float64Field` | `float64` | Floating-point fields (also handles float32) |
| `Decimal128Field` | `bson.Decimal128` | High-precision decimal fields |
| `BoolField` | `bool` | Boolean fields |
| `StringArrayField` | `[]string` | Arrays/slices of strings (also handles `*[]string`) |
| `TimestampField` | `time.Time` | Date/time fields |
| `GeoField` | `mongorm.GeoPoint` / `mongorm.GeoLineString` / `mongorm.GeoPolygon` | Geospatial fields |
| `GenericField` | any | Fallback for unmapped types (only `BSONName()` available) |

---

## ObjectIDField

**Package:** `primitives.ObjectIDField`

```go
type ToDoSchema struct {
    ID *primitives.ObjectIDField
}
var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
```

### ObjectID Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v bson.ObjectID)` | `$eq` | Equals |
| `Ne(v bson.ObjectID)` | `$ne` | Not equals |
| `In(v ...bson.ObjectID)` | `$in` | In a list |
| `Nin(v ...bson.ObjectID)` | `$nin` | Not in a list |
| `Gt(v bson.ObjectID)` | `$gt` | Greater than |
| `Gte(v bson.ObjectID)` | `$gte` | Greater than or equal |
| `Lt(v bson.ObjectID)` | `$lt` | Less than |
| `Lte(v bson.ObjectID)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ToDoFields.ID.Eq(someID))
orm.Where(ToDoFields.ID.In(id1, id2, id3))
```

---

## StringField

**Package:** `primitives.StringField`

### String Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v string)` | `$eq` | Equals |
| `Ne(v string)` | `$ne` | Not equals |
| `Reg(pattern string)` | `$regex` | Regex match |
| `In(v ...string)` | `$in` | In a list |
| `Nin(v ...string)` | `$nin` | Not in a list |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ToDoFields.Text.Eq("Buy groceries"))
orm.Where(ToDoFields.Text.Reg("groceries$")) // ends with "groceries"
```

---

## Int64Field

**Package:** `primitives.Int64Field`

Handles `int64`, `int32`, `int8`, and `int` model fields.

### Int64 Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v int64)` | `$eq` | Equals |
| `Ne(v int64)` | `$ne` | Not equals |
| `In(v ...int64)` | `$in` | In a list |
| `Nin(v ...int64)` | `$nin` | Not in a list |
| `Gt(v int64)` | `$gt` | Greater than |
| `Gte(v int64)` | `$gte` | Greater than or equal |
| `Lt(v int64)` | `$lt` | Less than |
| `Lte(v int64)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(TaskFields.Priority.Gte(3))
orm.Where(TaskFields.Count.In(1, 2, 5))
```

---

## Float64Field

**Package:** `primitives.Float64Field`

Handles `float64` and `float32` model fields.

### Float64 Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v float64)` | `$eq` | Equals |
| `Ne(v float64)` | `$ne` | Not equals |
| `In(v ...float64)` | `$in` | In a list |
| `Nin(v ...float64)` | `$nin` | Not in a list |
| `Gt(v float64)` | `$gt` | Greater than |
| `Gte(v float64)` | `$gte` | Greater than or equal |
| `Lt(v float64)` | `$lt` | Less than |
| `Lte(v float64)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ProductFields.Price.Lt(99.99))
```

---

## Decimal128Field

**Package:** `primitives.Decimal128Field`

Handles `bson.Decimal128` model fields.

### Decimal128 Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v bson.Decimal128)` | `$eq` | Equals |
| `Ne(v bson.Decimal128)` | `$ne` | Not equals |
| `In(v []bson.Decimal128)` | `$in` | In a list |
| `Nin(v []bson.Decimal128)` | `$nin` | Not in a list |
| `Gt(v bson.Decimal128)` | `$gt` | Greater than |
| `Gte(v bson.Decimal128)` | `$gte` | Greater than or equal |
| `Lt(v bson.Decimal128)` | `$lt` | Less than |
| `Lte(v bson.Decimal128)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
amount, _ := bson.ParseDecimal128("10.50")
orm.Where(ProductFields.Amount.Gte(amount))
```

---

## BoolField

**Package:** `primitives.BoolField`

### Bool Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v bool)` | `$eq` | Equals |
| `Ne(v bool)` | `$ne` | Not equals |
| `In(v ...bool)` | `$in` | In a list |
| `Nin(v ...bool)` | `$nin` | Not in a list |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ToDoFields.Done.Eq(false))
```

---

## StringArrayField

**Package:** `primitives.StringArrayField`

Handles `[]string` and `*[]string` model fields.

### StringArray Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v []string)` | `$eq` | Equals full array |
| `Ne(v []string)` | `$ne` | Not equals full array |
| `In(v []string)` | `$in` | Any value in list |
| `Nin(v []string)` | `$nin` | No value in list |
| `Contains(v string)` | `$in` | Array contains value |
| `ContainsAll(v []string)` | `$all` | Array contains all values |
| `Size(v int)` | `$size` | Array size matches |
| `ElemMatch(v bson.M)` | `$elemMatch` | Element matches filter |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(UserFields.Auth.Scopes.Contains("email"))
orm.Where(UserFields.Auth.Scopes.ContainsAll([]string{"email", "profile"}))
orm.Where(UserFields.Auth.Scopes.Size(2))
```

---

## TimestampField

**Package:** `primitives.TimestampField`

### Timestamp Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v time.Time)` | `$eq` | Equals |
| `Ne(v time.Time)` | `$ne` | Not equals |
| `In(v ...time.Time)` | `$in` | In a list |
| `Nin(v ...time.Time)` | `$nin` | Not in a list |
| `Gt(v time.Time)` | `$gt` | After |
| `Gte(v time.Time)` | `$gte` | After or at |
| `Lt(v time.Time)` | `$lt` | Before |
| `Lte(v time.Time)` | `$lte` | Before or at |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
cutoff := time.Now().Add(-24 * time.Hour)
orm.Where(ToDoFields.CreatedAt.Gte(cutoff))
```

---

## GeoField

**Package:** `primitives.GeoField`

Use `GeoField` for geospatial queries.

### Supported model types

- `*mongorm.GeoPoint`
- `*mongorm.GeoLineString`
- `*mongorm.GeoPolygon`

### Geo Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v any)` | `$eq` | Equals geometry |
| `Ne(v any)` | `$ne` | Not equals geometry |
| `Near(geometry any)` | `$near` | Near geometry |
| `NearWithDistance(geometry, min, max)` | `$near` | Near with min/max distance |
| `NearSphere(geometry any)` | `$nearSphere` | Spherical near |
| `NearSphereWithDistance(geometry, min, max)` | `$nearSphere` | Spherical near with min/max distance |
| `Within(geometry any)` | `$geoWithin` | Within geometry |
| `WithinBox(bottomLeft, upperRight)` | `$geoWithin/$box` | Within box |
| `WithinCenter(center, radius)` | `$geoWithin/$center` | Within flat circle |
| `WithinCenterSphere(center, radius)` | `$geoWithin/$centerSphere` | Within spherical circle |
| `Intersects(geometry any)` | `$geoIntersects` | Geometry intersection |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

{% raw %}

```go
point := mongorm.NewGeoPoint(12.9716, 77.5946)

orm.Where(ToDoFields.Location.Near(point))

cityBounds := mongorm.NewGeoPolygon(
    [][]float64{{77.4, 12.8}, {77.8, 12.8}, {77.8, 13.1}, {77.4, 13.1}, {77.4, 12.8}},
)
orm.Where(ToDoFields.Location.Within(cityBounds))
```

{% endraw %}

---

## GenericField

**Package:** `primitives.GenericField`

Fallback for any model field type that does not map to one of the specific primitives above (for example custom structs, slices, arrays, and maps).

### Generic Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v any)` | `$eq` | Equals |
| `Ne(v any)` | `$ne` | Not equals |
| `In(v []any)` | `$in` | In a list |
| `Nin(v []any)` | `$nin` | Not in a list |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |
| `Contains(v any)` | `$in` | Value exists inside array field |
| `ContainsAll(v []any)` | `$all` | Array contains all values |
| `Size(v int)` | `$size` | Array size matches |
| `ElemMatch(v bson.M)` | `$elemMatch` | Element matches filter |
| `Path(path string)` | dot-notation | Targets nested path under object field |

```go
orm.Where(UserFields.Goth.Path("provider").Eq("google"))
orm.Where(UserFields.Roles.Contains("admin"))
orm.Where(UserFields.Roles.Size(2))
```

### Deep typed fields for custom structs

For nested user-defined struct fields that map to `GenericField`, you can still generate type-safe primitives with `NestedFieldsOf`.

```go
type ToDoMeta struct {
    Source   *string `bson:"source,omitempty"`
    Priority *int64  `bson:"priority,omitempty"`
}

type ToDoMetaSchema struct {
    Source   *primitives.StringField
    Priority *primitives.Int64Field
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
var ToDoMetaFields = mongorm.NestedFieldsOf[ToDoMeta, ToDoMetaSchema](ToDoFields.Meta)

orm.Where(ToDoMetaFields.Source.Eq("import"))
orm.Where(ToDoMetaFields.Priority.Gte(2))
```

### Nested schema pointer support (User inside ToDo)

`FieldsOf` also supports nested schema structs directly, so you can define a schema pointer for a nested model field.

```go
type User struct {
    ID    *bson.ObjectID `bson:"_id,omitempty"`
    Email *string        `bson:"email,omitempty"`
}

type UserSchema struct {
    ID    *primitives.StringField
    Email *primitives.StringField
}

type ToDo struct {
    ID   *bson.ObjectID `bson:"_id,omitempty"`
    User *User          `bson:"user,omitempty"`
}

type ToDoSchema struct {
    ID   *primitives.ObjectIDField
    User *UserSchema
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

orm.Where(ToDoFields.User.Email.Eq("john@example.com"))
// => {"user.email": "john@example.com"}
```

This works for deeper nesting as well (for example `todo.user.profile.provider`).

It also works when the nested model field is an array/slice of structs, including pointer slices such as `Key *[]UserDefinedStruct`, by declaring `Key *UserDefinedStructSchema` in the schema.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
