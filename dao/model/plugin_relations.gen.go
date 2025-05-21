package model

import (
	"fmt"

	"github.com/google/uuid"
)

const TableNamePluginRelation = "converter_catalogue.plugin_relations"

// PluginRelation mapped from table <plugin_relations>
type PluginRelation struct {
	// the id of the relation (generated when the relation is created)
	ID string `gorm:"column:id;primaryKey" json:"id"`
	// the id of the plugin (from the plugin table)
	PluginID string `gorm:"column:plugin_id;not null" json:"plugin_id"`
	// the instanceId of the distribution
	RelationID string `gorm:"column:relation_id;not null" json:"relation_id"`
	// the file format expected by the plugin for the input
	InputFormat string `gorm:"column:input_format;not null" json:"input_format"`
	// the file format expected as the output from the plugin execution
	OutputFormat string `gorm:"column:output_format;not null" json:"output_format"`
}

// TableName PluginRelation's table name
func (*PluginRelation) TableName() string {
	return TableNamePluginRelation
}

func (r *PluginRelation) Validate() error {
	// TODO: better validation
	if r.ID == "" || uuid.Validate(r.ID) != nil {
		return fmt.Errorf("invalid Id in relation: %+v", r)
	}
	if r.InputFormat == "" {
		return fmt.Errorf("invalid InputFormat in relation: %+v", r)
	}
	if r.OutputFormat == "" {
		return fmt.Errorf("invalid OutputFormat in relation: %+v", r)
	}

	if r.PluginID == "" || uuid.Validate(r.PluginID) != nil {
		return fmt.Errorf("invalid PluginID in relation: %+v", r)
	}

	if r.RelationID == "" || uuid.Validate(r.RelationID) != nil {
		return fmt.Errorf("invalid RelationID in relation: %+v", r)
	}

	return nil
}
