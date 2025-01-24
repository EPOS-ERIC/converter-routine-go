package connection

import (
	"fmt"
	"strings"

	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/google/uuid"
)

func GetSoftwareSourceCodes() ([]model.Softwaresourcecode, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	var listOfSoftwareSourceCodes []model.Softwaresourcecode
	err = db.Where(&model.Versioningstatus{Status: "PUBLISHED"}).
		Where("softwaresourcecode.uid ilike ?", "%plugin%").
		Joins("versioningstatus").
		Find(&listOfSoftwareSourceCodes).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareSourceCodes, nil
}

func GetSoftwareApplications() ([]model.Softwareapplication, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareApplications []model.Softwareapplication
	err = db.Where(&model.Versioningstatus{Status: "PUBLISHED"}).
		Where("softwareapplication.uid ilike ?", "%plugin%").
		Joins("versioningstatus").
		Find(&listOfSoftwareApplications).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplications, nil
}

func GetSoftwareApplicationsOperations() ([]model.SoftwareapplicationOperation, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareApplicationsOperations []model.SoftwareapplicationOperation
	err = db.Find(&listOfSoftwareApplicationsOperations).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplicationsOperations, nil
}

func GetPlugins() ([]model.Plugin, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfPlugins []model.Plugin
	err = db.Find(&listOfPlugins).Error
	if err != nil {
		return nil, err
	}
	return listOfPlugins, nil
}

func GetPluginById(id string) (model.Plugin, error) {
	db, err := Connect()
	if err != nil {
		return model.Plugin{}, err
	}
	var plugin model.Plugin
	err = db.Find(&plugin, "id = ?", id).Error
	if err != nil {
		return model.Plugin{}, err
	}
	return plugin, nil
}

func GetPluginRelations() ([]model.PluginRelation, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfPluginRelations []model.PluginRelation
	err = db.Find(&listOfPluginRelations).Error
	if err != nil {
		return nil, err
	}
	return listOfPluginRelations, nil
}

func SetPlugins(ph []model.Plugin) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	// truncate
	err = db.Exec("TRUNCATE plugin CASCADE").Error // c("TRUNCATE plugin CASCADE", nil)
	if err != nil {
		return err
	}
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}

	return nil
}

func SetPluginsRelations(ph []model.PluginRelation) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	// truncate
	err = db.Exec("TRUNCATE plugin_relations CASCADE").Error
	if err != nil {
		return err
	}
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}

	return nil
}

// getNewSoftwareSourceCode returns the (new) software source codes that are not in the plugins table
func getNewSoftwareSourceCode() ([]model.Softwaresourcecode, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	var softwareSourceCode []model.Softwaresourcecode
	// Select all plugins that are not in the plugins table (new plugin)
	err = db.Model(&softwareSourceCode).
		Joins("LEFT JOIN plugin ON softwaresourcecode.instance_id = plugin.software_source_code_id").
		Where("plugin.software_source_code_id IS NULL").
		Where("softwaresourcecode.state = ?", "PUBLISHED").
		Where("softwaresourcecode.uid ILIKE '%' || ? || '%'", "plugin").Find(&softwareSourceCode).Error
	if err != nil {
		return nil, err
	}

	return softwareSourceCode, nil
}

func InsertPlugins(plugins []model.Plugin) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	err = db.Create(&plugins).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertPluginsRelations(pluginRelations []model.PluginRelation) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	err = db.Create(&pluginRelations).Error
	if err != nil {
		return err
	}
	return nil
}

func GeneratePluginsRelations() ([]model.PluginRelation, error) {
	newApplicationsOperations, err := getNewApplicationOperations()
	if err != nil {
		return nil, fmt.Errorf("failed to get new application operations: %w", err)
	}

	var listOfPluginsRelations []model.PluginRelation

	for _, newOperation := range newApplicationsOperations {
		plugin, err := getPluginFromSoftwareApplicationInstanceId(newOperation.SoftwareapplicationInstanceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin for software application instance ID %s: %w", newOperation.SoftwareapplicationInstanceID, err)
		}

		softwareApplicationParameters, err := getSoftwareApplicationParameters(newOperation.SoftwareapplicationInstanceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get software application parameters for instance ID %s: %w", newOperation.SoftwareapplicationInstanceID, err)
		}

		if len(softwareApplicationParameters) != 2 {
			return nil, fmt.Errorf("unexpected number of software application parameters (%d) for instance ID %s", len(softwareApplicationParameters), newOperation.SoftwareapplicationInstanceID)
		}

		var inputFormat, outputFormat string
		for _, sap := range softwareApplicationParameters {
			switch sap.Action {
			case "OBJECT":
				inputFormat = sap.Encodingformat
			case "RESULT":
				outputFormat = sap.Encodingformat
			default:
				return nil, fmt.Errorf("unknown action type '%s' in software application parameters for instance ID %s", sap.Action, newOperation.SoftwareapplicationInstanceID)
			}
		}

		pluginRelation := model.PluginRelation{
			ID:           uuid.New().String(),
			PluginID:     plugin.ID,
			RelationID:   newOperation.SoftwareapplicationInstanceID,
			RelationType: "Operation",
			InputFormat:  inputFormat,
			OutputFormat: outputFormat,
		}

		listOfPluginsRelations = append(listOfPluginsRelations, pluginRelation)
	}

	return listOfPluginsRelations, nil
}

func getPluginFromSoftwareApplicationInstanceId(softwareApplicationInstanceId string) (model.Plugin, error) {
	db, err := Connect()
	if err != nil {
		return model.Plugin{}, err
	}

	plugin := model.Plugin{
		SoftwareApplicationID: softwareApplicationInstanceId,
	}
	err = db.Find(&plugin).Error
	if err != nil {
		return model.Plugin{}, err
	}
	return plugin, nil
}

func getSoftwareApplicationParameters(softwareApplicationInstanceId string) ([]model.Parameter, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	var sap []model.SoftwareapplicationParameter
	err = db.Model(&sap).
		Where("softwareapplication_instance_id = ?", softwareApplicationInstanceId).Find(&sap).Error
	if err != nil {
		return nil, err
	}

	var p []model.Parameter
	err = db.Model(&p).Joins("join softwareapplication_parameters on parameter.instance_id = software_application_parameters.parameter_instance_id").Find(&p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetSoftwareSourceCodeProgrammingLanguage(ssc string) (string, error) {
	db, err := Connect()
	if err != nil {
		return "", err
	}

	type Result struct {
		result string
	}
	var result Result
	err = db.Table("softwaresourcecode").
		Select(`"element".value as value`).
		Joins(`JOIN softwaresourcecode_element ON softwaresourcecode_element.softwaresourcecode_instance_id = softwaresourcecode.instance_id`).
		Joins(`JOIN "element" ON "element".instance_id = softwaresourcecode_element.element_instance_id`).
		Where(`"element".type = ?`, "PROGRAMMINGLANGUAGE").
		Where("softwaresourcecode.instance_id = ?", ssc).
		Take(&result).Error
	if err != nil {
		return "", err
	}
	return result.result, nil
}

func GeneratePlugins(installedRepos []model.Softwaresourcecode) ([]model.Plugin, error) {
	// Retrieve new software source codes
	listOfSoftwareSourceCodes, err := getNewSoftwareSourceCode()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve new software source codes: %w", err)
	}

	// Retrieve software applications
	listOfSoftwareApplications, err := GetSoftwareApplications()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve software applications: %w", err)
	}

	var listOfPlugins []model.Plugin

	// For each software source code (that is a plugin)
	for _, objSoftwareSourceCode := range listOfSoftwareSourceCodes {
		// If the objSoftwareSourceCode is not in the installedRepos, skip it
		found := false
		for _, repo := range installedRepos {
			if objSoftwareSourceCode.UID == repo.UID {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// Initialize a new plugin
		plugin := model.Plugin{
			ID:                   uuid.New().String(),
			SoftwareSourceCodeID: objSoftwareSourceCode.InstanceID,
			Version:              objSoftwareSourceCode.Softwareversion,
			Installed:            true,
			Enabled:              true,
		}

		// For each software application
		for _, objSoftwareApplication := range listOfSoftwareApplications {
			// If the software source code and the software application don't match, continue
			if strings.Replace(objSoftwareSourceCode.UID, "SoftwareSourceCode/", "", -1) != strings.Replace(objSoftwareApplication.UID, "SoftwareApplication/", "", -1) {
				continue
			}

			lang, err := GetSoftwareSourceCodeProgrammingLanguage(objSoftwareSourceCode.InstanceID)
			if err != nil {
				return nil, err
			}
			// Set the plugin properties
			plugin.ProxyType = lang
			plugin.SoftwareApplicationID = objSoftwareApplication.InstanceID
			plugin.Runtime = lang
			plugin.Execution = objSoftwareApplication.Requirements
		}

		// Add the plugin to the list
		listOfPlugins = append(listOfPlugins, plugin)
	}

	return listOfPlugins, nil
}

func getNewApplicationOperations() ([]model.SoftwareapplicationOperation, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareApplicationsOperations []model.SoftwareapplicationOperation
	// Select all the software application operations that are not in the plugin relations table
	err = db.Model(&listOfSoftwareApplicationsOperations).
		Joins("LEFT JOIN plugin_relations ON softwareapplication_operation.instance_operation_id = plugin_relations.relation_id").
		Where("plugin_relations.relation_id IS NULL").
		Find(&listOfSoftwareApplicationsOperations).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplicationsOperations, nil
}
