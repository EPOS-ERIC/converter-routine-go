package orms

type SoftwareApplicationParameters struct {
	tableName                     struct{} `pg:"softwareapplication_parameters,alias:softwareapplication_parameters"`
	Id                            string
	Encodingformat                string
	Conformsto                    string
	Action                        string
	InstanceSoftwareapplicationId string
}
