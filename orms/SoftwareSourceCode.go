package orms

type SoftwareSourceCode struct {
	tableName        struct{} `gorm:"softwaresourcecode,alias:softwaresourcecode"`
	InstanceID       string   `gorm:"column:instance_id"`
	MetaID           string   `gorm:"column:meta_id"`
	UID              string   `gorm:"column:uid"`
	Name             string
	Description      string
	LicenseURL       string `gorm:"column:licenseurl"`
	DownloadURL      string `gorm:"column:downloadurl"`
	RuntimePlatform  string `gorm:"column:runtimeplatform"`
	SoftwareVersion  string `gorm:"column:softwareversion"`
	Keywords         string
	CodeRepository   string `gorm:"column:coderepository"`
	MainEntityOfPage string `gorm:"column:mainentityofpage"`
	Operation        string
	State            string
}

func (SoftwareSourceCode) TableName() string {
	return "softwaresourcecode" // Replace this with your actual table name
}

func (s *SoftwareSourceCode) GetInstanceID() string {
	return s.InstanceID
}

func (s *SoftwareSourceCode) SetInstanceID(instanceID string) {
	s.InstanceID = instanceID
}

func (s *SoftwareSourceCode) GetMetaID() string {
	return s.MetaID
}

func (s *SoftwareSourceCode) SetMetaID(metaID string) {
	s.MetaID = metaID
}

func (s *SoftwareSourceCode) GetUID() string {
	return s.UID
}

func (s *SoftwareSourceCode) SetUID(uid string) {
	s.UID = uid
}

func (s *SoftwareSourceCode) GetName() string {
	return s.Name
}

func (s *SoftwareSourceCode) SetName(name string) {
	s.Name = name
}

func (s *SoftwareSourceCode) GetDescription() string {
	return s.Description
}

func (s *SoftwareSourceCode) SetDescription(description string) {
	s.Description = description
}

func (s *SoftwareSourceCode) GetLicenseURL() string {
	return s.LicenseURL
}

func (s *SoftwareSourceCode) SetLicenseURL(licenseURL string) {
	s.LicenseURL = licenseURL
}

func (s *SoftwareSourceCode) GetDownloadURL() string {
	return s.DownloadURL
}

func (s *SoftwareSourceCode) SetDownloadURL(downloadURL string) {
	s.DownloadURL = downloadURL
}

func (s *SoftwareSourceCode) GetRuntimePlatform() string {
	return s.RuntimePlatform
}

func (s *SoftwareSourceCode) SetRuntimePlatform(runtimePlatform string) {
	s.RuntimePlatform = runtimePlatform
}

func (s *SoftwareSourceCode) GetSoftwareVersion() string {
	return s.SoftwareVersion
}

func (s *SoftwareSourceCode) SetSoftwareVersion(softwareVersion string) {
	s.SoftwareVersion = softwareVersion
}

func (s *SoftwareSourceCode) GetKeywords() string {
	return s.Keywords
}

func (s *SoftwareSourceCode) SetKeywords(keywords string) {
	s.Keywords = keywords
}

func (s *SoftwareSourceCode) GetCodeRepository() string {
	return s.CodeRepository
}

func (s *SoftwareSourceCode) SetCodeRepository(codeRepository string) {
	s.CodeRepository = codeRepository
}

func (s *SoftwareSourceCode) GetMainEntityOfPage() string {
	return s.MainEntityOfPage
}

func (s *SoftwareSourceCode) SetMainEntityOfPage(mainEntityOfPage string) {
	s.MainEntityOfPage = mainEntityOfPage
}

func (s *SoftwareSourceCode) GetOperation() string {
	return s.Operation
}

func (s *SoftwareSourceCode) SetOperation(operation string) {
	s.Operation = operation
}

func (s *SoftwareSourceCode) GetState() string {
	return s.State
}

func (s *SoftwareSourceCode) SetState(state string) {
	s.State = state
}
