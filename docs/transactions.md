# Transactions

MongORM provides a clean transaction wrapper for atomic, multi-step workflows. It is designed for production-critical flows where data consistency matters across multiple writes.

## WithTransaction

Use `WithTransaction(ctx, fn)` on any initialized ORM instance.

The callback receives a transaction-bound context (`txCtx`). Any MongORM operation using `txCtx` participates in the same transaction.

```go
err := mongorm.New(&ToDo{}).WithTransaction(ctx, func(txCtx context.Context) error {
    a := &ToDo{Text: mongorm.String("step-1")}
    if err := mongorm.New(a).Save(txCtx); err != nil {
        return err
    }

    b := &ToDo{Text: mongorm.String("step-2")}
    if err := mongorm.New(b).Save(txCtx); err != nil {
        return err
    }

    return nil // commit
})
if err != nil {
    panic(err)
}
```

Return an error from the callback to abort and roll back:

```go
rollbackErr := errors.New("abort transaction")

err := mongorm.New(&ToDo{}).WithTransaction(ctx, func(txCtx context.Context) error {
    toDo := &ToDo{Text: mongorm.String("temporary")}
    if err := mongorm.New(toDo).Save(txCtx); err != nil {
        return err
    }

    return rollbackErr
})
```

## Transaction Options

Pass MongoDB transaction options through `WithTransaction`:

```go
err := mongorm.New(&ToDo{}).WithTransaction(
    ctx,
    func(txCtx context.Context) error {
        // operations...
        return nil
    },
    options.Transaction().SetReadConcern(readconcern.Majority()),
)
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
