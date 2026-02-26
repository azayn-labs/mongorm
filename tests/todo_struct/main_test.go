package todostruct

import (
	"testing"

	"github.com/azayn-labs/mongorm"
)

func logger(t *testing.T, message string) {
	t.Logf("TODO [struct options] %s\n", message)
}

func TestMain(t *testing.T) {
	var toDo = &ToDo{
		Text: mongorm.String("This is an example todo created with struct options"),
	}

	t.Run("Advanced fields generation", func(t *testing.T) {
		ValidateAdvancedFieldsOf(t)
	})

	t.Run("Advanced generic queries", func(t *testing.T) {
		ValidateAdvancedGenericQueries(t)
	})

	t.Run("Create TODO", func(t *testing.T) {
		CreateTodo(t, toDo)
	})

	t.Run("Find TODO by ID", func(t *testing.T) {
		FindByIDToDo(t, toDo.ID)
	})

	t.Run("Update TODO by ID", func(t *testing.T) {
		update := &ToDo{
			Text: mongorm.String("This is an updated todo text"),
		}
		UpdateToDoByID(t, toDo.ID, update)
	})
}
