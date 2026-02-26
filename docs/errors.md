# Errors

MongORM exposes sentinel errors to simplify error handling with `errors.Is`.

## Sentinel Errors

- `ErrNotFound`
- `ErrDuplicateKey`
- `ErrInvalidConfig`
- `ErrTransactionUnsupported`
- `ErrOptimisticLockConflict`

## Usage

```go
if err := orm.First(ctx); err != nil {
    switch {
    case errors.Is(err, mongorm.ErrNotFound):
        // handle not found
    case errors.Is(err, mongorm.ErrDuplicateKey):
        // handle duplicate key
    case errors.Is(err, mongorm.ErrInvalidConfig):
        // handle configuration issue
    case errors.Is(err, mongorm.ErrOptimisticLockConflict):
        // handle stale version update
    default:
        // generic error handling
    }
}
```

`ErrNotFound` keeps compatibility with Mongo driver behavior by also matching `mongo.ErrNoDocuments`.

## Transaction Capability Helper

Use `IsTransactionUnsupported(err)` to detect deployments that do not support transactions (for example, standalone servers):

```go
err := mongorm.New(&ToDo{}).WithTransaction(ctx, func(txCtx context.Context) error {
    // transactional operations
    return nil
})

if mongorm.IsTransactionUnsupported(err) {
    // fallback to non-transaction flow
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
