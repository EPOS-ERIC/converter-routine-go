package connection

import (
	"fmt"
	"os"
	"strings"

	"github.com/epos-eu/converter-routine/orms"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

func Connect() (*pg.DB, error) {
	conn := ""

	if val, res := os.LookupEnv("POSTGRESQL_CONNECTION_STRING"); res == true {
		conn = val
	} else {
		return nil, fmt.Errorf("POSTGRESQL_CONNECTION_STRING is not set")
	}

	opt, err := pg.ParseURL(conn)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)

	return db, nil
}

func GetSoftwareSourceCodes() ([]orms.SoftwareSourceCode, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// Select all users.
	var listOfSoftwareSourceCodes []orms.SoftwareSourceCode
	err = db.Model(&listOfSoftwareSourceCodes).Where("state = ?", "PUBLISHED").Where("uid ILIKE '%' || ? || '%'", "plugin").Select()
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
	defer db.Close()
	// Select all users.
	var listOfSoftwareApplications []orms.SoftwareApplication
	err = db.Model(&listOfSoftwareApplications).Where("state = ?", "PUBLISHED").Where("uid ILIKE '%' || ? || '%'", "plugin").Select()
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
	defer db.Close()
	// Select all users.
	var listOfSoftwareApplicationsOperations []orms.SoftwareApplicationOperation
	err = db.Model(&listOfSoftwareApplicationsOperations).Select()
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
	defer db.Close()
	// Select all users.
	var listOfPlugins []orms.Plugin
	err = db.Model(&listOfPlugins).Select()
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
	defer db.Close()
	// Select all users.
	var listOfPluginRelations []orms.PluginRelations
	err = db.Model(&listOfPluginRelations).Select()
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
	defer db.Close()

	//truncate
	_, err = db.Exec("TRUNCATE plugin CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Model(&ph).Insert()
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
	defer db.Close()

	//truncate
	_, err = db.Exec("TRUNCATE plugin_relations CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Model(&ph).Insert()
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
	defer db.Close()

	var softwareSourceCode []orms.SoftwareSourceCode
	// Select all plugins that are not in the plugins table (new plugin)
	err = db.Model(&softwareSourceCode).
		Join("LEFT JOIN plugin ON softwaresourcecode.instance_id = plugin.software_source_code_id").
		Where("plugin.software_source_code_id IS NULL").
		Where("softwaresourcecode.state = ?", "PUBLISHED").
		Where("softwaresourcecode.uid ILIKE '%' || ? || '%'", "plugin").
		Select()
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
	defer db.Close()

	_, err = db.Model(&plugins).Insert()
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
	defer db.Close()

	_, err = db.Model(&pluginRelations).Insert()
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
		plugin, err := getPluginFromSoftwareApplicationInstanceId(newOperation.Instance_softwareapplication_id)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin for software application instance ID %s: %w", newOperation.Instance_softwareapplication_id, err)
		}

		softwareApplicationParameters, err := getSoftwareApplicationParameters(newOperation.Instance_softwareapplication_id)
		if err != nil {
			return nil, fmt.Errorf("failed to get software application parameters for instance ID %s: %w", newOperation.Instance_softwareapplication_id, err)
		}

		if len(softwareApplicationParameters) != 2 {
			return nil, fmt.Errorf("unexpected number of software application parameters (%d) for instance ID %s", len(softwareApplicationParameters), newOperation.Instance_softwareapplication_id)
		}

		var inputFormat, outputFormat string
		for _, sap := range softwareApplicationParameters {
			switch sap.Action {
			case "object":
				inputFormat = sap.Encodingformat
			case "result":
				outputFormat = sap.Encodingformat
			default:
				return nil, fmt.Errorf("unknown action type '%s' in software application parameters for instance ID %s", sap.Action, newOperation.Instance_softwareapplication_id)
			}
		}

		pluginRelation := orms.PluginRelations{
			Id:            uuid.New().String(),
			Plugin_id:     plugin.Id,
			Relation_id:   newOperation.Instance_operation_id,
			Relation_type: "Operation",
			Input_format:  inputFormat,
			Output_format: outputFormat,
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
	defer db.Close()

	var plugin orms.Plugin
	err = db.Model(&plugin).
		Where("software_application_id = ?", softwareApplicationInstanceId).
		Select()
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
	defer db.Close()

	var sa []orms.SoftwareApplicationParameters
	err = db.Model(&sa).
		Where("instance_softwareapplication_id = ?", softwareApplicationInstanceId).
		Select()
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
	defer db.Close()

	var sscpg orms.SoftwareSourceCodeProgrammingLanguage
	err = db.Model(&sscpg).Where("instance_softwaresourcecode_id = ?", ssc).Select()
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
			if objSoftwareSourceCode.Uid == repo.Uid {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// Initialize a new plugin
		plugin := orms.Plugin{
			Id:                      uuid.New().String(),
			Software_source_code_id: objSoftwareSourceCode.Instance_id,
			Version:                 objSoftwareSourceCode.Softwareversion,
			Installed:               true,
			Enabled:                 true,
		}

		// For each software application
		for _, objSoftwareApplication := range listOfSoftwareApplications {
			// If the software source code and the software application don't match, continue
			if strings.Replace(objSoftwareSourceCode.Uid, "SoftwareSourceCode/", "", -1) != strings.Replace(objSoftwareApplication.Uid, "SoftwareApplication/", "", -1) {
				continue
			}

			lang, err := GetSoftwareSourceCodeProgrammingLanguage(objSoftwareSourceCode.Instance_id)
			if err != nil {
				return nil, err
			}
			// Set the plugin properties
			plugin.Proxy_type = lang
			plugin.Software_application_id = objSoftwareApplication.Instance_id
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
	defer db.Close()
	// Select all users.
	var listOfSoftwareApplicationsOperations []orms.SoftwareApplicationOperation
	// Select all the software application operations that are not in the plugin relations table
	err = db.Model(&listOfSoftwareApplicationsOperations).
		Join("LEFT JOIN plugin_relations ON softwareapplication_operation.instance_operation_id = plugin_relations.relation_id").
		Where("plugin_relations.relation_id IS NULL").
		Select()
	if err != nil {
		return nil, err
	}
	return listOfSoftwareApplicationsOperations, nil
}
