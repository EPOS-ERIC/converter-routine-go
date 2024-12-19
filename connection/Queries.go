package connection

import (
	"fmt"
	"strings"

	"github.com/epos-eu/converter-routine/orms"
	"github.com/google/uuid"
)

func GetSoftwareSourceCodes() ([]orms.SoftwareSourceCode, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareSourceCodes []orms.SoftwareSourceCode
	err = db.Model(orms.SoftwareSourceCode{}).Where("state = ?", "PUBLISHED").Where("uid ILIKE '%' || ? || '%'", "plugin").Find(&listOfSoftwareSourceCodes).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareSourceCodes, nil
}

func GetSoftwareApplications() ([]orms.SoftwareApplication, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareApplications []orms.SoftwareApplication
	err = db.Model(&listOfSoftwareApplications).Where("state = ?", "PUBLISHED").Where("uid ILIKE '%' || ? || '%'", "plugin").Find(&listOfSoftwareApplications).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplications, nil
}

func GetSoftwareApplicationsOperations() ([]orms.SoftwareApplicationOperation, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareApplicationsOperations []orms.SoftwareApplicationOperation
	err = db.Model(&listOfSoftwareApplicationsOperations).Find(&listOfSoftwareApplicationsOperations).Error
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplicationsOperations, nil
}

func GetPlugins() ([]orms.Plugin, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfPlugins []orms.Plugin
	err = db.Model(&listOfPlugins).Find(&listOfPlugins).Error
	if err != nil {
		return nil, err
	}
	return listOfPlugins, nil
}

func GetPluginRelations() ([]orms.PluginRelations, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfPluginRelations []orms.PluginRelations
	err = db.Model(&listOfPluginRelations).Find(&listOfPluginRelations).Error
	if err != nil {
		return nil, err
	}
	return listOfPluginRelations, nil
}

func SetPlugins(ph []orms.Plugin) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	//truncate
	err = db.Exec("TRUNCATE plugin CASCADE").Error //c("TRUNCATE plugin CASCADE", nil)
	if err != nil {
		return err
	}
	err = db.Create(&ph).Error
	if err != nil {
		return err
	}

	return nil
}

func SetPluginsRelations(ph []orms.PluginRelations) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	//truncate
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
func getNewSoftwareSourceCode() ([]orms.SoftwareSourceCode, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	var softwareSourceCode []orms.SoftwareSourceCode
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

func InsertPlugins(plugins []orms.Plugin) error {
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

func InsertPluginsRelations(pluginRelations []orms.PluginRelations) error {
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

func GeneratePluginsRelations() ([]orms.PluginRelations, error) {
	newApplicationsOperations, err := getNewApplicationOperations()
	if err != nil {
		return nil, fmt.Errorf("failed to get new application operations: %w", err)
	}

	var listOfPluginsRelations []orms.PluginRelations

	for _, newOperation := range newApplicationsOperations {
		plugin, err := getPluginFromSoftwareApplicationInstanceId(newOperation.InstanceSoftwareApplicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin for software application instance ID %s: %w", newOperation.InstanceSoftwareApplicationID, err)
		}

		softwareApplicationParameters, err := getSoftwareApplicationParameters(newOperation.InstanceSoftwareApplicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get software application parameters for instance ID %s: %w", newOperation.InstanceSoftwareApplicationID, err)
		}

		if len(softwareApplicationParameters) != 2 {
			return nil, fmt.Errorf("unexpected number of software application parameters (%d) for instance ID %s", len(softwareApplicationParameters), newOperation.InstanceSoftwareApplicationID)
		}

		var inputFormat, outputFormat string
		for _, sap := range softwareApplicationParameters {
			switch sap.Action {
			case "object":
				inputFormat = sap.EncodingFormat
			case "result":
				outputFormat = sap.EncodingFormat
			default:
				return nil, fmt.Errorf("unknown action type '%s' in software application parameters for instance ID %s", sap.Action, newOperation.InstanceSoftwareApplicationID)
			}
		}

		pluginRelation := orms.PluginRelations{
			Id:           uuid.New().String(),
			PluginID:     plugin.Id,
			RelationID:   newOperation.InstanceOperationID,
			RelationType: "Operation",
			InputFormat:  inputFormat,
			OutputFormat: outputFormat,
		}

		listOfPluginsRelations = append(listOfPluginsRelations, pluginRelation)
	}

	return listOfPluginsRelations, nil
}

func getPluginFromSoftwareApplicationInstanceId(softwareApplicationInstanceId string) (orms.Plugin, error) {
	db, err := Connect()
	if err != nil {
		return orms.Plugin{}, err
	}

	var plugin orms.Plugin
	err = db.Model(&plugin).
		Where("software_application_id = ?", softwareApplicationInstanceId).Find(&plugin).Error
	if err != nil {
		return orms.Plugin{}, err
	}
	return plugin, nil
}

func getSoftwareApplicationParameters(softwareApplicationInstanceId string) ([]orms.SoftwareApplicationParameters, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	var sa []orms.SoftwareApplicationParameters
	err = db.Model(&sa).
		Where("instance_softwareapplication_id = ?", softwareApplicationInstanceId).Find(&sa).Error
	if err != nil {
		return nil, err
	}
	return sa, nil
}

func GetSoftwareSourceCodeProgrammingLanguage(ssc string) (string, error) {
	db, err := Connect()
	if err != nil {
		return "", err
	}

	var sscpg orms.SoftwareSourceCodeProgrammingLanguage
	err = db.Model(&sscpg).Where("instance_softwaresourcecode_id = ?", ssc).Find(&sscpg).Error
	if err != nil {
		return "", err
	}
	return sscpg.Language, nil
}

func GeneratePlugins(installedRepos []orms.SoftwareSourceCode) ([]orms.Plugin, error) {
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

	var listOfPlugins []orms.Plugin

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
		plugin := orms.Plugin{
			Id:                   uuid.New().String(),
			SoftwareSourceCodeID: objSoftwareSourceCode.InstanceID,
			Version:              objSoftwareSourceCode.SoftwareVersion,
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

func getNewApplicationOperations() ([]orms.SoftwareApplicationOperation, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// Select all users.
	var listOfSoftwareApplicationsOperations []orms.SoftwareApplicationOperation
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
