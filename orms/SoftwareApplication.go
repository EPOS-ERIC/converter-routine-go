package orms

type SoftwareApplication struct {
	tableName       struct{} `pg:"softwareapplication,alias:softwareapplication"`
	Instance_id     string
	Meta_id         string
	Uid             string
	Name            string
	Description     string
	Licenseurl      string
	Downloadurl     string
	Softwareversion string
	Keywords        string
	Requirements    string
	State           string
}

func (s *SoftwareApplication) GetInstance_id() string {
	return s.Instance_id
}

func (s *SoftwareApplication) SetInstance_id(Instance_id string) {
	s.Instance_id = Instance_id
}

func (s *SoftwareApplication) GetMeta_id() string {
	return s.Meta_id
}

func (s *SoftwareApplication) SetMeta_id(Meta_id string) {
	s.Meta_id = Meta_id
}

func (s *SoftwareApplication) GetUid() string {
	return s.Uid
}

func (s *SoftwareApplication) SetUid(Uid string) {
	s.Uid = Uid
}

func (s *SoftwareApplication) GetName() string {
	return s.Name
}

func (s *SoftwareApplication) SetName(Name string) {
	s.Name = Name
}

func (s *SoftwareApplication) GetDescription() string {
	return s.Description
}

func (s *SoftwareApplication) SetDescription(Description string) {
	s.Description = Description
}

func (s *SoftwareApplication) GetLicenseurl() string {
	return s.Licenseurl
}

func (s *SoftwareApplication) SetLicenseurl(Licenseurl string) {
	s.Licenseurl = Licenseurl
}

func (s *SoftwareApplication) GetDownloadurl() string {
	return s.Downloadurl
}

func (s *SoftwareApplication) SetDownloadurl(Downloadurl string) {
	s.Downloadurl = Downloadurl
}

func (s *SoftwareApplication) GetSoftwareversion() string {
	return s.Softwareversion
}

func (s *SoftwareApplication) SetSoftwareversion(Softwareversion string) {
	s.Softwareversion = Softwareversion
}

func (s *SoftwareApplication) GetKeywords() string {
	return s.Keywords
}

func (s *SoftwareApplication) SetKeywords(Keywords string) {
	s.Keywords = Keywords
}

func (s *SoftwareApplication) GetRequirements() string {
	return s.Requirements
}

func (s *SoftwareApplication) SetRequirements(Requirements string) {
	s.Requirements = Requirements
}

func (s *SoftwareApplication) GetState() string {
	return s.State
}

func (s *SoftwareApplication) SetState(State string) {
	s.State = State
}
