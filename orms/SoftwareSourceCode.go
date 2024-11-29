package orms

type SoftwareSourceCode struct {
	tableName        struct{} `pg:"softwaresourcecode,alias:softwaresourcecode"`
	Instance_id      string
	Meta_id          string
	Uid              string
	Name             string
	Description      string
	Licenseurl       string
	Downloadurl      string
	Runtimeplatform  string
	Softwareversion  string
	Keywords         string
	Coderepository   string
	Mainentityofpage string
	Operation        string
	State            string
}

func (s *SoftwareSourceCode) GetInstance_id() string {
	return s.Instance_id
}

func (s *SoftwareSourceCode) SetInstance_id(Instance_id string) {
	s.Instance_id = Instance_id
}

func (s *SoftwareSourceCode) GetMeta_id() string {
	return s.Meta_id
}

func (s *SoftwareSourceCode) SetMeta_id(Meta_id string) {
	s.Meta_id = Meta_id
}

func (s *SoftwareSourceCode) GetUid() string {
	return s.Uid
}

func (s *SoftwareSourceCode) SetUid(Uid string) {
	s.Uid = Uid
}

func (s *SoftwareSourceCode) GetName() string {
	return s.Name
}

func (s *SoftwareSourceCode) SetName(Name string) {
	s.Name = Name
}

func (s *SoftwareSourceCode) GetDescription() string {
	return s.Description
}

func (s *SoftwareSourceCode) SetDescription(Description string) {
	s.Description = Description
}

func (s *SoftwareSourceCode) GetLicenseurl() string {
	return s.Licenseurl
}

func (s *SoftwareSourceCode) SetLicenseurl(Licenseurl string) {
	s.Licenseurl = Licenseurl
}

func (s *SoftwareSourceCode) GetDownloadurl() string {
	return s.Downloadurl
}

func (s *SoftwareSourceCode) SetDownloadurl(Downloadurl string) {
	s.Downloadurl = Downloadurl
}

func (s *SoftwareSourceCode) GetRuntimeplatform() string {
	return s.Runtimeplatform
}

func (s *SoftwareSourceCode) SetRuntimeplatform(Runtimeplatform string) {
	s.Runtimeplatform = Runtimeplatform
}

func (s *SoftwareSourceCode) GetSoftwareversion() string {
	return s.Softwareversion
}

func (s *SoftwareSourceCode) SetSoftwareversion(Softwareversion string) {
	s.Softwareversion = Softwareversion
}

func (s *SoftwareSourceCode) GetKeywords() string {
	return s.Keywords
}

func (s *SoftwareSourceCode) SetKeywords(Keywords string) {
	s.Keywords = Keywords
}

func (s *SoftwareSourceCode) GetCoderepository() string {
	return s.Coderepository
}

func (s *SoftwareSourceCode) SetCoderepository(Coderepository string) {
	s.Coderepository = Coderepository
}

func (s *SoftwareSourceCode) GetMainentityofpage() string {
	return s.Mainentityofpage
}

func (s *SoftwareSourceCode) SetMainentityofpage(Mainentityofpage string) {
	s.Mainentityofpage = Mainentityofpage
}

func (s *SoftwareSourceCode) GetOperation() string {
	return s.Operation
}

func (s *SoftwareSourceCode) SetOperation(Operation string) {
	s.Operation = Operation
}

func (s *SoftwareSourceCode) GetState() string {
	return s.State
}

func (s *SoftwareSourceCode) SetState(State string) {
	s.State = State
}
