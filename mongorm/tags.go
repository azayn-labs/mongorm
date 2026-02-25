package mongorm

type ModelTags string

const (
	ModelTagReadonly ModelTags = "readonly"
	ModelTagPrimary  ModelTags = "primary"
)
