package todostruct

import (
	"fmt"
	"testing"

	"github.com/CdTgr/mongorm"
)

func FindByRegexToDo(t *testing.T, pattern string) {
	logger(t, fmt.Sprintf("Using pattern %s\n", pattern))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.Text.Reg(pattern))

	if err := todoModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("Found with regex %s: %+v\n", pattern, toDo))
}
