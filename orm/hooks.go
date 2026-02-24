package orm

type BeforeSaveHook interface {
	BeforeSave() error
}

type AfterSaveHook interface {
	AfterSave() error
}

type BeforeCreateHook interface {
	BeforeCreate() error
}

type BeforeUpdateHook interface {
	BeforeUpdate() error
}
