package orms

type Plugin struct {
	tableName               struct{} `pg:"plugin,alias:plugin"`
	Id                      string
	Software_source_code_id string
	Software_application_id string
	Version                 string
	Proxy_type              string
	Runtime                 string
	Execution               string
	Installed               bool
	Enabled                 bool
}

func (p *Plugin) GetId() string {
	return p.Id
}

func (p *Plugin) SetId(Id string) {
	p.Id = Id
}

func (p *Plugin) GetSoftware_source_code_id() string {
	return p.Software_source_code_id
}

func (p *Plugin) SetSoftware_source_code_id(Software_source_code_id string) {
	p.Software_source_code_id = Software_source_code_id
}

func (p *Plugin) GetSoftware_application_id() string {
	return p.Software_application_id
}

func (p *Plugin) SetSoftware_application_id(Software_application_id string) {
	p.Software_application_id = Software_application_id
}

func (p *Plugin) GetVersion() string {
	return p.Version
}

func (p *Plugin) SetVersion(Version string) {
	p.Version = Version
}

func (p *Plugin) GetProxy_type() string {
	return p.Proxy_type
}

func (p *Plugin) SetProxy_type(Proxy_type string) {
	p.Proxy_type = Proxy_type
}

func (p *Plugin) GetRuntime() string {
	return p.Runtime
}

func (p *Plugin) SetRuntime(Runtime string) {
	p.Runtime = Runtime
}

func (p *Plugin) GetExecution() string {
	return p.Execution
}

func (p *Plugin) SetExecution(Execution string) {
	p.Execution = Execution
}

func (p *Plugin) GetInstalled() bool {
	return p.Installed
}

func (p *Plugin) SetInstalled(Installed bool) {
	p.Installed = Installed
}

func (p *Plugin) GetEnabled() bool {
	return p.Enabled
}

func (p *Plugin) SetEnabled(Enabled bool) {
	p.Enabled = Enabled
}
