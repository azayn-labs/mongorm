package todostruct

import (
	"fmt"
	"testing"

	"github.com/CdTgr/mongorm"
)

func CreateTodo(t *testing.T, toDo *ToDo) {
	logger(t, "Creating TODO")

	todoModel := mongorm.New(toDo)
	if err := todoModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("Created TODO: %+v\n", toDo))
}
