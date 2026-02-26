package mongorm

import "go.mongodb.org/mongo-driver/v2/bson"

// BeforeFindHook is called before executing a find operation.
// It allows you to modify the query or perform any necessary setup.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) BeforeFind(m *MongORM[ToDo], query *bson.M) error {
//	    // Modify the query or perform setup here
//	    return nil
//	}
type BeforeFindHook[T any] interface {
	BeforeFind(*MongORM[T], *bson.M) error
}

// AfterFindHook is called after executing a find operation.
// It allows you to perform any necessary cleanup or post-processing.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) AfterFind(m *MongORM[ToDo]) error {
//	    // Perform cleanup or post-processing here
//	    return nil
//	}
type AfterFindHook[T any] interface {
	AfterFind(*MongORM[T]) error
}

// BeforeSaveHook is called before executing a save operation.
// It allows you to modify the document or perform any necessary setup.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) BeforeSave(m *MongORM[ToDo], doc *bson.M) error {
//	    // Modify the document or perform setup here
//	    return nil
//	}
type BeforeSaveHook[T any] interface {
	BeforeSave(*MongORM[T], *bson.M) error
}

// AfterSaveHook is called after executing a save operation.
// It allows you to perform any necessary cleanup or post-processing.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) AfterSave(m *MongORM[ToDo]) error {
//	    // Perform cleanup or post-processing here
//	    return nil
//	}
type AfterSaveHook[T any] interface {
	AfterSave(*MongORM[T]) error
}

// BeforeCreateHook is called before executing a create operation.
// It allows you to modify the document or perform any necessary setup.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) BeforeCreate(m *MongORM[ToDo]) error {
//	    // Modify the document or perform setup here
//	    return nil
//	}
type BeforeCreateHook[T any] interface {
	BeforeCreate(*MongORM[T]) error
}

// AfterCreateHook is called after executing a create operation.
// It allows you to perform any necessary cleanup or post-processing.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) AfterCreate(m *MongORM[ToDo]) error {
//	    // Perform cleanup or post-processing here
//	    return nil
//	}
type AfterCreateHook[T any] interface {
	AfterCreate(*MongORM[T]) error
}

// BeforeUpdateHook is called before executing an update operation.
// It allows you to modify the query, update document, or perform any necessary setup.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) BeforeUpdate(m *MongORM[ToDo], query *bson.M, update *bson.M) error {
//	    // Modify the query, update document, or perform setup here
//	    return nil
//	}
type BeforeUpdateHook[T any] interface {
	BeforeUpdate(*MongORM[T], *bson.M, *bson.M) error
}

// AfterUpdateHook is called after executing an update operation.
// It allows you to perform any necessary cleanup or post-processing.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) AfterUpdate(m *MongORM[ToDo]) error {
//	    // Perform cleanup or post-processing here
//	    return nil
//	}
type AfterUpdateHook[T any] interface {
	AfterUpdate(*MongORM[T]) error
}

// BeforeDeleteHook is called before executing a delete operation.
// It allows you to modify the query or perform any necessary setup.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) BeforeDelete(m *MongORM[ToDo], query *bson.M) error {
//	    // Modify the query or perform setup here
//	    return nil
//	}
type BeforeDeleteHook[T any] interface {
	BeforeDelete(*MongORM[T], *bson.M) error
}

// AfterDeleteHook is called after executing a delete operation.
// It allows you to perform any necessary cleanup or post-processing.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) AfterDelete(m *MongORM[ToDo]) error {
//	    // Perform cleanup or post-processing here
//	    return nil
//	}
type AfterDeleteHook[T any] interface {
	AfterDelete(*MongORM[T]) error
}

// BeforeFinalizeHook is called before finalizing the MongORM instance.
// It allows you to perform any necessary setup or modifications before the instance is finalized.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) BeforeFinalize(m *MongORM[ToDo]) error {
//	    // Perform setup or modifications here
//	    return nil
//	}
type BeforeFinalizeHook[T any] interface {
	BeforeFinalize(*MongORM[T]) error
}

// AfterFinalizeHook is called after finalizing the MongORM instance.
// It allows you to perform any necessary cleanup or post-processing after the instance is finalized.
//
// Example usage:
//
//	type ToDo struct {
//	    Text *string `bson:"text"`
//	}
//
//	func (u *ToDo) AfterFinalize(m *MongORM[ToDo]) error {
//	    // Perform cleanup or post-processing here
//	    return nil
//	}
type AfterFinalizeHook[T any] interface {
	AfterFinalize(*MongORM[T]) error
}
