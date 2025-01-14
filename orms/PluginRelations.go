package orms

type PluginRelations struct {
	tableName    struct{} `gorm:"plugin_relations,alias:plugin_relations"`
	Id           string   `gorm:"primaryKey"`
	PluginID     string   `gorm:"column:plugin_id"`
	RelationID   string   `gorm:"column:relation_id"`
	RelationType string   `gorm:"column:relation_type"`
	InputFormat  string   `gorm:"column:input_format"`
	OutputFormat string   `gorm:"column:output_format"`
}

func (p *PluginRelations) GetId() string {
	return p.Id
}

func (p *PluginRelations) SetId(id string) {
	p.Id = id
}

func (p *PluginRelations) GetPluginID() string {
	return p.PluginID
}

func (p *PluginRelations) SetPluginID(pluginID string) {
	p.PluginID = pluginID
}

func (p *PluginRelations) GetRelationID() string {
	return p.RelationID
}

func (p *PluginRelations) SetRelationID(relationID string) {
	p.RelationID = relationID
}

func (p *PluginRelations) GetRelationType() string {
	return p.RelationType
}

func (p *PluginRelations) SetRelationType(relationType string) {
	p.RelationType = relationType
}

func (p *PluginRelations) GetInputFormat() string {
	return p.InputFormat
}

func (p *PluginRelations) SetInputFormat(inputFormat string) {
	p.InputFormat = inputFormat
}

func (p *PluginRelations) GetOutputFormat() string {
	return p.OutputFormat
}

func (p *PluginRelations) SetOutputFormat(outputFormat string) {
	p.OutputFormat = outputFormat
}

func (PluginRelations) TableName() string {
	return "plugin_relations" // Replace this with your actual table name
}
