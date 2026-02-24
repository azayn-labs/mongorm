package orm

type ModelTags string

const (
	ModelTagPrimary  ModelTags = "primary"
	ModelTagReadonly ModelTags = "readonly"
)
