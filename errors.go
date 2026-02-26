package mongorm

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrNotFound               = errors.New("mongorm: not found")
	ErrDuplicateKey           = errors.New("mongorm: duplicate key")
	ErrInvalidConfig          = errors.New("mongorm: invalid configuration")
	ErrTransactionUnsupported = errors.New("mongorm: transaction unsupported")
	ErrOptimisticLockConflict = errors.New("mongorm: optimistic lock conflict")
)

func normalizeError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.Join(ErrNotFound, err)
	}

	if mongo.IsDuplicateKeyError(err) {
		return errors.Join(ErrDuplicateKey, err)
	}

	if IsTransactionUnsupported(err) {
		return errors.Join(ErrTransactionUnsupported, err)
	}

	return err
}

func configErrorf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidConfig, fmt.Sprintf(format, args...))
}

func IsTransactionUnsupported(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrTransactionUnsupported) {
		return true
	}

	return isTransactionUnsupportedError(err)
}

func isTransactionUnsupportedError(err error) bool {
	message := strings.ToLower(err.Error())

	if strings.Contains(message, "transaction numbers are only allowed") {
		return true
	}

	if strings.Contains(message, "transactions are not supported") {
		return true
	}

	if strings.Contains(message, "replica set") && strings.Contains(message, "transaction") {
		return true
	}

	return false
}

func mapUpdateOneError(err error, optimisticLockEnabled bool) error {
	if optimisticLockEnabled && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.Join(ErrOptimisticLockConflict, err)
	}

	return normalizeError(err)
}
