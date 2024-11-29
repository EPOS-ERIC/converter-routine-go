package orms

type SoftwareApplicationOperation struct {
	tableName                       struct{} `pg:"softwareapplication_operation,alias:softwareapplication_operation"`
	Instance_operation_id           string
	Instance_softwareapplication_id string
}

func (s *SoftwareApplicationOperation) GetInstance_operation_id() string {
	return s.Instance_operation_id
}

func (s *SoftwareApplicationOperation) SetInstance_operation_id(Instance_operation_id string) {
	s.Instance_operation_id = Instance_operation_id
}

func (s *SoftwareApplicationOperation) GetInstance_softwareapplication_id() string {
	return s.Instance_softwareapplication_id
}

func (s *SoftwareApplicationOperation) SetInstance_softwareapplication_id(Instance_softwareapplication_id string) {
	s.Instance_softwareapplication_id = Instance_softwareapplication_id
}
