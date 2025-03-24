package model

import "fmt"

const TableNamePluginRelation = "plugin_relations"

// PluginRelation mapped from table <plugin_relations>
type PluginRelation struct {
	ID           string `gorm:"column:id;primaryKey" json:"id"`
	PluginID     string `gorm:"column:plugin_id;not null" json:"plugin_id"`
	RelationID   string `gorm:"column:relation_id;not null" json:"relation_id"`
	RelationType string `gorm:"column:relation_type;not null" json:"relation_type"`
	InputFormat  string `gorm:"column:input_format;not null" json:"input_format"`
	OutputFormat string `gorm:"column:output_format;not null" json:"output_format"`
}

// TableName PluginRelation's table name
func (*PluginRelation) TableName() string {
	return TableNamePluginRelation
}

func (r *PluginRelation) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("invalid Id in relation: %+v", r)
	}
	if r.InputFormat == "" {
		return fmt.Errorf("invalid InputFormat in relation: %+v", r)
	}
	if r.OutputFormat == "" {
		return fmt.Errorf("invalid OutputFormat in relation: %+v", r)
	}
	if r.PluginID == "" {
		return fmt.Errorf("invalid PluginID in relation: %+v", r)
	}
	if r.RelationID == "" {
		return fmt.Errorf("invalid RelationID in relation: %+v", r)
	}
	if r.RelationType == "" {
		return fmt.Errorf("invalid RelationType in relation: %+v", r)
	}

	return nil
}
