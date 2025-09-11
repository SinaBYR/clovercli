package internal

type Field struct {
	Name  string
	Value string
}

type Collection struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name string
	Func func(args ...interface{})
}

type Object struct {
	Name        string
	Fields      []Field
	Collections []Collection
	Methods     []Method
}

func NewObject(name string) *Object {
	return &Object{
		Name: name,
	}
}

func (o *Object) RegisterField(field Field) {
	o.Fields = append(o.Fields, field)
}

func (o *Object) RegisterMethod(method Method) {
	o.Methods = append(o.Methods, method)
}

func (o *Object) RegisterCollection(collName string) {
	coll := Collection{
		Name: collName,
		Methods: []Method{
			Method{Name: "find", Func: EvaluateFind},
		},
	}
	o.Collections = append(o.Collections, coll)
}
