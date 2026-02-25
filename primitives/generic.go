package primitives

type GenericField struct {
	name string
}

func GenericType(name string) *GenericField {
	return &GenericField{name: name}
}

func (f *GenericField) BSONName() string {
	return f.name
}
