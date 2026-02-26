package todostruct

import (
	"fmt"
	"testing"

	"github.com/azayn-labs/mongorm"
)

func CreateTodo(t *testing.T, toDo *ToDo) {
	logger(t, "[TODO] Creating")

	todoModel := mongorm.New(toDo)
	if err := todoModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("[TODO] Created: %+v\n", toDo))
}
