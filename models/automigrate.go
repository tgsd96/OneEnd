package models

var Models []interface{}

func AddToModel(model interface{}) {
	Models = append(Models, model)
}
