package orms

type SoftwareSourceCodeProgrammingLanguage struct {
	tableName                    struct{} `gorm:"softwaresourcecode_programminglanguage,alias:softwaresourcecode_programminglanguage"`
	ID                           string   `gorm:"primaryKey"`
	Language                     string
	InstanceSoftwareSourceCodeID string `gorm:"column:instance_softwaresourcecode_id"`
}

func (SoftwareSourceCodeProgrammingLanguage) TableName() string {
	return "softwaresourcecode_programminglanguage" // Replace this with your actual table name
}
