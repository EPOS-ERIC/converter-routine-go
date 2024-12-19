package orms

type SoftwareApplicationOperation struct {
	tableName                     struct{} `gorm:"softwareapplication_operation,alias:softwareapplication_operation"`
	InstanceOperationID           string   `gorm:"column:instance_operation_id"`
	InstanceSoftwareApplicationID string   `gorm:"column:instance_softwareapplication_id"`
}

func (s *SoftwareApplicationOperation) GetInstanceOperationID() string {
	return s.InstanceOperationID
}

func (s *SoftwareApplicationOperation) SetInstanceOperationID(instanceOperationID string) {
	s.InstanceOperationID = instanceOperationID
}

func (s *SoftwareApplicationOperation) GetInstanceSoftwareApplicationID() string {
	return s.InstanceSoftwareApplicationID
}

func (s *SoftwareApplicationOperation) SetInstanceSoftwareApplicationID(instanceSoftwareApplicationID string) {
	s.InstanceSoftwareApplicationID = instanceSoftwareApplicationID
}

func (SoftwareApplicationOperation) TableName() string {
	return "softwareapplication_operation" // Replace this with your actual table name
}
