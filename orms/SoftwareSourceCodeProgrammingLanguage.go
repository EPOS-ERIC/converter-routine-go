package orms

type SoftwareSourceCodeProgrammingLanguage struct {
	tableName                      struct{} `pg:"softwaresourcecode_programminglanguage,alias:softwaresourcecode_programminglanguage"`
	Id                             string
	Language                       string
	Instance_softwaresourcecode_id string
}
