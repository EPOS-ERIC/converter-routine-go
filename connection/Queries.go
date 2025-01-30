package connection

import (
	"fmt"
	"log"
	"strings"

	"github.com/epos-eu/converter-routine/dao/model"
	"github.com/google/uuid"
)

// Initialize the logger
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func GetSoftwareSourceCodes() ([]model.Softwaresourcecode, error) {
	db, err := ConnectMetadata()
	if err != nil {
		return nil, err
	}

	var listOfSoftwareSourceCodes []model.Softwaresourcecode
	err = db.Where("softwaresourcecode.uid ilike ?", "%plugin%").
		Joins("JOIN versioningstatus ON versioningstatus.version_id = softwaresourcecode.version_id").
		Where("versioningstatus.status = ?", "PUBLISHED").
		Find(&listOfSoftwareSourceCodes).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareSourceCodes, nil
}

func GetSoftwareApplications() ([]model.Softwareapplication, error) {
	db, err := ConnectMetadata()
	if err != nil {
		return nil, err
	}

	var listOfSoftwareApplications []model.Softwareapplication
	err = db.Where("softwareapplication.uid ilike ?", "%plugin%").
		Joins("JOIN versioningstatus ON versioningstatus.version_id = softwareapplication.version_id").
		Where("versioningstatus.status = ?", "PUBLISHED").
		Find(&listOfSoftwareApplications).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplications, nil
}

func GetSoftwareApplicationsOperations() ([]model.SoftwareapplicationOperation, error) {
	db, err := ConnectMetadata()
	if err != nil {
		return nil, err
	}

	var listOfSoftwareApplicationsOperations []model.SoftwareapplicationOperation
	err = db.Find(&listOfSoftwareApplicationsOperations).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplicationsOperations, nil
}

func GetPlugins() ([]model.Plugin, error) {
	db, err := ConnectConverter()
	if err != nil {
		return nil, err
	}

	var listOfPlugins []model.Plugin
	err = db.Find(&listOfPlugins).Error
	if err != nil {
		return nil, err
	}
	return listOfPlugins, nil
}

func GetPluginById(id string) (model.Plugin, error) {
	db, err := ConnectConverter()
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
	db, err := ConnectConverter()
	if err != nil {
		return nil, err
	}

	var listOfPluginRelations []model.PluginRelation
	err = db.Find(&listOfPluginRelations).Error
	if err != nil {
		return nil, err
	}
	return listOfPluginRelations, nil
}

func SetPlugins(ph []model.Plugin) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	// Truncate plugins table
	err = db.Exec("TRUNCATE plugin CASCADE").Error
	if err != nil {
		return err
	}

	// Insert new plugins
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}
	return nil
}

func SetPluginsRelations(ph []model.PluginRelation) error {
	db, err := ConnectConverter()
	if err != nil {
		return err
	}

	// Truncate plugin_relations table
	err = db.Exec("TRUNCATE plugin_relations CASCADE").Error
	if err != nil {
		return err
	}

	// Insert new plugin relations
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}
	return nil
}

// getNewSoftwareSourceCode returns the (new) software source codes that are not in the plugins table
func getNewSoftwareSourceCode() ([]model.Softwaresourcecode, error) {
	softwareSourceCodes, err := GetSoftwareSourceCodes()
	if err != nil {
		return nil, err
	}

	plugins, err := GetPlugins()
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of SoftwareSourceCodeIDs in plugins
	pluginSourceCodeMap := make(map[string]struct{}, len(plugins))
	for _, plugin := range plugins {
		pluginSourceCodeMap[plugin.SoftwareSourceCodeID] = struct{}{}
	}

	// Collect software source codes that are not in the plugins map
	var result []model.Softwaresourcecode
	for _, ssc := range softwareSourceCodes {
		if _, exists := pluginSourceCodeMap[ssc.InstanceID]; !exists {
			result = append(result, ssc)
		}
	}
	return result, nil
}

func getNewApplicationOperations() ([]model.SoftwareapplicationOperation, error) {
	softwareApplicationOperations, err := GetSoftwareApplicationsOperations()
	if err != nil {
		return nil, err
	}

	pluginRelations, err := GetPluginRelations()
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of SoftwareSourceCodeIDs in plugins
	pluginApplicationsMap := make(map[string]struct{}, len(pluginRelations))
	for _, plugin := range pluginRelations {
		pluginApplicationsMap[plugin.RelationID] = struct{}{}
	}

	// Collect software source codes that are not in the plugins map
	var result []model.SoftwareapplicationOperation
	for _, sao := range softwareApplicationOperations {
		if _, exists := pluginApplicationsMap[sao.OperationInstanceID]; !exists {
			result = append(result, sao)
		}
	}
	return result, nil
}

func InsertPlugins(plugins []model.Plugin) error {
	db, err := ConnectConverter()
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
	db, err := ConnectConverter()
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

		var inputFormat, outputFormat string
		for _, sap := range softwareApplicationParameters {
			switch sap.Action {
			case "OBJECT":
				inputFormat = sap.Encodingformat
			case "RESULT":
				outputFormat = sap.Encodingformat
			default:
				continue
			}
		}

		if inputFormat == "" || outputFormat == "" {
			return nil, fmt.Errorf("error retrieving input/output formats:\ninputFormat: %s\noutputFormat: %s", inputFormat, outputFormat)
		}

		pluginRelation := model.PluginRelation{
			ID:           uuid.New().String(),
			PluginID:     plugin.ID,
			RelationID:   newOperation.OperationInstanceID,
			RelationType: "Operation",
			InputFormat:  inputFormat,
			OutputFormat: outputFormat,
		}

		listOfPluginsRelations = append(listOfPluginsRelations, pluginRelation)
	}

	return listOfPluginsRelations, nil
}

func getPluginFromSoftwareApplicationInstanceId(softwareApplicationInstanceId string) (model.Plugin, error) {
	db, err := ConnectConverter()
	if err != nil {
		return model.Plugin{}, err
	}

	var plugin model.Plugin
	err = db.Where("software_application_id = ?", softwareApplicationInstanceId).First(&plugin).Error
	if err != nil {
		return model.Plugin{}, err
	}
	return plugin, nil
}

func getSoftwareApplicationParameters(softwareApplicationInstanceId string) ([]model.Parameter, error) {
	db, err := ConnectMetadata()
	if err != nil {
		return nil, err
	}

	var p []model.Parameter
	err = db.Model(&p).
		Joins("JOIN softwareapplication_parameters ON parameter.instance_id = softwareapplication_parameters.parameter_instance_id").
		Where("softwareapplication_parameters.softwareapplication_instance_id = ?", softwareApplicationInstanceId).
		Find(&p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetSoftwareSourceCodeProgrammingLanguage(ssc string) (string, error) {
	db, err := ConnectMetadata()
	if err != nil {
		return "", err
	}

	type Result struct {
		Result string `gorm:"column:value"`
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
	return result.Result, nil
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
		// Check if the software source code is in the installedRepos
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
			// Compare UIDs after removing prefixes
			softwareSourceCodeUID := strings.Replace(objSoftwareSourceCode.UID, "SoftwareSourceCode/", "", -1)
			softwareApplicationUID := strings.Replace(objSoftwareApplication.UID, "SoftwareApplication/", "", -1)
			if softwareSourceCodeUID != softwareApplicationUID {
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
