package orms

type SoftwareApplicationParameters struct {
	tableName                     struct{} `gorm:"softwareapplication_parameters,alias:softwareapplication_parameters"`
	ID                            string   `gorm:"primaryKey"`
	EncodingFormat                string   `gorm:"column:encodingformat"`
	ConformsTo                    string   `gorm:"column:conformsto"`
	Action                        string
	InstanceSoftwareApplicationID string `gorm:"column:instance_softwareapplication_id"`
}

func (SoftwareApplicationParameters) TableName() string {
	return "softwareapplication_parameters" // Replace this with your actual table name
}
