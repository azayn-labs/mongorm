package todostruct

import (
	"fmt"
	"testing"
	"time"

	"github.com/CdTgr/mongorm"
)

func logger(t *testing.T, message string) {
	t.Logf("TODO [struct options] %s\n", message)
}

func TestMain(t *testing.T) {
	var toDo = &ToDo{
		Text: mongorm.String("This is an example todo created with struct options"),
	}

	t.Run("Create TODO", func(t *testing.T) {
		CreateTodo(t, toDo)
	})

	t.Run("Find TODO by ID", func(t *testing.T) {
		FindByIDToDo(t, toDo.ID)
	})

	t.Run("Find TODO by regex", func(t *testing.T) {
		FindByRegexToDo(t, "struct options$")
	})

	t.Run("Update TODO by ID", func(t *testing.T) {
		update := &ToDo{
			Text: mongorm.String("This is an updated todo text"),
		}
		UpdateToDoByID(t, toDo.ID, update)
	})

	t.Run("Update All TODOs", func(t *testing.T) {
		update := &ToDo{
			Text: mongorm.String(
				fmt.Sprintf("This is an updated todo text for all todos at %s", time.Now()),
			),
		}
		UpdateAllToDo(t, update)
	})
}
