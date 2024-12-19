package orms

type SoftwareApplication struct {
	tableName       struct{} `gorm:"softwareapplication,alias:softwareapplication"`
	InstanceID      string   `gorm:"column:instance_id"`
	MetaID          string   `gorm:"column:meta_id"`
	UID             string   `gorm:"column:uid"`
	Name            string
	Description     string
	LicenseURL      string `gorm:"column:licenseurl"`
	DownloadURL     string `gorm:"column:downloadurl"`
	SoftwareVersion string `gorm:"column:softwareversion"`
	Keywords        string
	Requirements    string
	State           string
}

func (s *SoftwareApplication) GetInstanceID() string {
	return s.InstanceID
}

func (s *SoftwareApplication) SetInstanceID(instanceID string) {
	s.InstanceID = instanceID
}

func (s *SoftwareApplication) GetMetaID() string {
	return s.MetaID
}

func (s *SoftwareApplication) SetMetaID(metaID string) {
	s.MetaID = metaID
}

func (s *SoftwareApplication) GetUID() string {
	return s.UID
}

func (s *SoftwareApplication) SetUID(uid string) {
	s.UID = uid
}

func (s *SoftwareApplication) GetName() string {
	return s.Name
}

func (s *SoftwareApplication) SetName(name string) {
	s.Name = name
}

func (s *SoftwareApplication) GetDescription() string {
	return s.Description
}

func (s *SoftwareApplication) SetDescription(description string) {
	s.Description = description
}

func (s *SoftwareApplication) GetLicenseURL() string {
	return s.LicenseURL
}

func (s *SoftwareApplication) SetLicenseURL(licenseURL string) {
	s.LicenseURL = licenseURL
}

func (s *SoftwareApplication) GetDownloadURL() string {
	return s.DownloadURL
}

func (s *SoftwareApplication) SetDownloadURL(downloadURL string) {
	s.DownloadURL = downloadURL
}

func (s *SoftwareApplication) GetSoftwareVersion() string {
	return s.SoftwareVersion
}

func (s *SoftwareApplication) SetSoftwareVersion(softwareVersion string) {
	s.SoftwareVersion = softwareVersion
}

func (s *SoftwareApplication) GetKeywords() string {
	return s.Keywords
}

func (s *SoftwareApplication) SetKeywords(keywords string) {
	s.Keywords = keywords
}

func (s *SoftwareApplication) GetRequirements() string {
	return s.Requirements
}

func (s *SoftwareApplication) SetRequirements(requirements string) {
	s.Requirements = requirements
}

func (s *SoftwareApplication) GetState() string {
	return s.State
}

func (s *SoftwareApplication) SetState(state string) {
	s.State = state
}

func (SoftwareApplication) TableName() string {
	return "softwareapplication" // Replace this with your actual table name
}
