package orms

type PluginRelations struct {
	tableName     struct{} `pg:"plugin_relations,alias:plugin_relations"`
	Id            string
	Plugin_id     string
	Relation_id   string
	Relation_type string
	Input_format  string
	Output_format string
}

func (p *PluginRelations) GetId() string {
	return p.Id
}

func (p *PluginRelations) SetId(Id string) {
	p.Id = Id
}

func (p *PluginRelations) GetPlugin_id() string {
	return p.Plugin_id
}

func (p *PluginRelations) SetPlugin_id(Plugin_id string) {
	p.Plugin_id = Plugin_id
}

func (p *PluginRelations) GetRelation_id() string {
	return p.Relation_id
}

func (p *PluginRelations) SetRelation_id(Relation_id string) {
	p.Relation_id = Relation_id
}

func (p *PluginRelations) GetRelation_type() string {
	return p.Relation_type
}

func (p *PluginRelations) SetRelation_type(Relation_type string) {
	p.Relation_type = Relation_type
}
