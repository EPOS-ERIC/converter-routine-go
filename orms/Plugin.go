package orms

import "fmt"

type Plugin struct {
	tableName             struct{} `gorm:"plugin,alias:plugin"`
	Id                    string   `gorm:"primaryKey"`
	SoftwareSourceCodeID  string   `gorm:"column:software_source_code_id"`
	SoftwareApplicationID string   `gorm:"column:software_application_id"`
	Version               string
	ProxyType             string `gorm:"column:proxy_type"`
	Runtime               string
	Execution             string
	Installed             bool
	Enabled               bool
}

func (p *Plugin) GetId() string {
	return p.Id
}

func (p *Plugin) SetId(id string) {
	p.Id = id
}

func (p *Plugin) GetSoftwareSourceCodeID() string {
	return p.SoftwareSourceCodeID
}

func (p *Plugin) SetSoftwareSourceCodeID(softwareSourceCodeID string) {
	p.SoftwareSourceCodeID = softwareSourceCodeID
}

func (p *Plugin) GetSoftwareApplicationID() string {
	return p.SoftwareApplicationID
}

func (p *Plugin) SetSoftwareApplicationID(softwareApplicationID string) {
	p.SoftwareApplicationID = softwareApplicationID
}

func (p *Plugin) GetVersion() string {
	return p.Version
}

func (p *Plugin) SetVersion(version string) {
	p.Version = version
}

func (p *Plugin) GetProxyType() string {
	return p.ProxyType
}

func (p *Plugin) SetProxyType(proxyType string) {
	p.ProxyType = proxyType
}

func (p *Plugin) GetRuntime() string {
	return p.Runtime
}

func (p *Plugin) SetRuntime(runtime string) {
	p.Runtime = runtime
}

func (p *Plugin) GetExecution() string {
	return p.Execution
}

func (p *Plugin) SetExecution(execution string) {
	p.Execution = execution
}

func (p *Plugin) GetInstalled() bool {
	return p.Installed
}

func (p *Plugin) SetInstalled(installed bool) {
	p.Installed = installed
}

func (p *Plugin) GetEnabled() bool {
	return p.Enabled
}

func (p *Plugin) SetEnabled(enabled bool) {
	p.Enabled = enabled
}
func (Plugin) TableName() string {
	return "plugin" // Replace this with your actual table name
}


func (p *Plugin) IsValid() error {
	if p.Execution == "" {
		return fmt.Errorf("invalid Execution in plugin: %+v", p)
	}
	if p.Id == "" {
		return fmt.Errorf("invalid Id in plugin: %+v", p)
	}
	if p.ProxyType == "" {
		return fmt.Errorf("invalid ProxyType in plugin: %+v", p)
	}
	if p.Runtime == "" {
		return fmt.Errorf("invalid Runtime in plugin: %+v", p)
	}
	if p.SoftwareApplicationID == "" {
		return fmt.Errorf("invalid SoftwareApplicationID in plugin: %+v", p)
	}
	if p.SoftwareSourceCodeID == "" {
		return fmt.Errorf("invalid SoftwareSourceCodeID in plugin: %+v", p)
	}
	if p.Version == "" {
		return fmt.Errorf("invalid Version in plugin: %+v", p)
	}

	return nil
}
