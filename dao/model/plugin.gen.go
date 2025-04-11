package model

import (
	"fmt"

	"github.com/google/uuid"
)

const TableNamePlugin = "plugin"

// ENUM(branch, tag)
type VersionType string

// ENUM(binary, java, python)
type SupportedRuntimes string

// Plugin mapped from table <plugin>
type Plugin struct {
	// the id of the plugin (generated when the plugin is created)
	ID string `gorm:"column:id;primaryKey" json:"id"`
	// the name of the plugin
	Name string `gorm:"column:name;not null" json:"name"`
	// a description of the plugin
	Description string `gorm:"column:description;not null" json:"description"`
	// the name of the branch if version_type is branch or the tag number if it is tag
	Version string `gorm:"column:version;not null" json:"version"`
	// either 'branch' or 'tag'
	VersionType VersionType `gorm:"column:version_type;not null" json:"version_type"`
	// the url from which to clone the repository
	Repository string `gorm:"column:repository;not null" json:"repository"`
	// the runtime (binary, java, python, ...)
	Runtime SupportedRuntimes `gorm:"column:runtime;not null" json:"runtime"`
	// the path for the executable
	Executable string `gorm:"column:executable;not null" json:"executable"`
	// arguments for the execution (if needed (like the main java class name))
	Arguments string `gorm:"column:arguments;not null" json:"arguments"`
	// if the plugin is currently installed
	Installed bool `gorm:"column:installed;not null" json:"installed"`
	// if the plugin is enabled aka if it can be used
	Enabled bool `gorm:"column:enabled;not null" json:"enabled"`
}

// TableName Plugin's table name
func (*Plugin) TableName() string {
	return TableNamePlugin
}

func (p *Plugin) Validate() error {
	// TODO: better validation
	if p.ID == "" || uuid.Validate(p.ID) != nil {
		return fmt.Errorf("invalid Id in plugin: %+v", p)
	}
	if p.Name == "" {
		return fmt.Errorf("invalid Name in plugin: %+v", p)
	}
	if p.Version == "" {
		return fmt.Errorf("invalid Version in plugin: %+v", p)
	}
	if p.VersionType == "" {
		return fmt.Errorf("invalid VersionType in plugin: %+v", p)
	}
	if !p.VersionType.IsValid() {
		return fmt.Errorf("invalid VersionType in plugin: %s is not in any of %+v", p.VersionType, VersionTypeValues())
	}
	if p.Repository == "" {
		return fmt.Errorf("invalid Repository in plugin: %+v", p)
	}
	if p.Runtime == "" {
		return fmt.Errorf("invalid Runtime in plugin: %+v", p)
	}
	if !p.Runtime.IsValid() {
		return fmt.Errorf("invalid Runtime in plugin: %s is not in any of %+v", p.Runtime, SupportedRuntimesValues())
	}
	if p.Executable == "" {
		return fmt.Errorf("invalid Executable in plugin: %+v", p)
	}

	return nil
}
