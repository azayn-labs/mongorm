package todostruct

import (
	"testing"

	"github.com/CdTgr/mongorm"
)

func UpdateAllToDo(t *testing.T, update *ToDo) {
	logger(t, "[TODO] Trying to update all\n")

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Set(update)

	if _, err := todoModel.SaveMulti(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, "[TODO] Updated all\n")
}
