package plugins

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Plugin struct {
	ID          string
	Name        string //Name to display in the UI
	Description string
	Version     string
	Author      string
	Category    string
	Entrypoint  string // The run command
	Dir         string // Work directory
	Parameters  []*PluginParameter
}

type PluginParameter struct {
	ID                 string // Same as key for the value, ie: "-t", "--ports=", ...
	Name               string // The field label for the UI.
	MinValue, MaxValue int    // Not valid for string/bool Type.
	Value              any    //
	ValueType          string // Using typeof() in utils to get type
}

func NewPlugin() *Plugin {

	return &Plugin{
		Parameters: []*PluginParameter{}, // Initialize as an empty slice
	}
}

func NewPluginFromJson(jsonContent []byte) (*Plugin, error) {

	result := NewPlugin()

	// Define a temporary struct for parsing the JSON
	var jsonData struct {
		Plugin struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Version     string `json:"version"`
			Author      string `json:"author"`
			Category    string `json:"category"`
			Entrypoint  string `json:"entrypoint"`
		} `json:"plugin"`
		Parameters map[string]struct {
			Name   string      `json:"name"`
			Limits []int       `json:"limits,omitempty"`
			Value  interface{} `json:"value"`
		} `json:"parameters"`
	}
	// Parse the JSON content
	if err := json.Unmarshal(jsonContent, &jsonData); err != nil {
		return result, fmt.Errorf("failed to parse JSON content: %w", err)
	}

	// Fill the Plugin
	result = &Plugin{
		ID:          jsonData.Plugin.ID,
		Name:        jsonData.Plugin.Name,
		Description: jsonData.Plugin.Description,
		Version:     jsonData.Plugin.Version,
		Author:      jsonData.Plugin.Author,
		Category:    jsonData.Plugin.Category,
		Entrypoint:  jsonData.Plugin.Entrypoint,
	}

	// Convert Parameters map to slice
	for key, param := range jsonData.Parameters {
		pluginParam := &PluginParameter{
			ID:        key,
			Name:      param.Name,
			Value:     param.Value,
			ValueType: typeof(param.Value),
		}
		if len(param.Limits) == 2 {
			pluginParam.MinValue = param.Limits[0]
			pluginParam.MaxValue = param.Limits[1]
		}
		result.Parameters = append(result.Parameters, pluginParam) // Append to slice
	}
	return result, nil
}

func (p *Plugin) SetParameterValue(paramID string, value any) error { //TODO: Check error?
	for _, param := range p.Parameters {
		if paramID == param.ID {
			param.Value = value
			return nil
		}
	}
	return fmt.Errorf("parameter %s in plugin %s does not exists", paramID, p.ID)
}

func (p *Plugin) GetParametersAsSlice() []string {
    result := []string{}
    for _, param := range p.Parameters {
        pStr := param.ID
        if ! strings.HasSuffix(param.ID, "=") {
            pStr = pStr + " "
        }
		if slice, ok := param.Value.([]interface{}); ok {
			pStr = pStr + SliceToCommaSepString(slice)
		} else {
			pStr = pStr + fmt.Sprintf("%v", param.Value)
		}
        result = append(result, pStr)        
    }
    return result
}

func (p *Plugin) GetParametersAsString() string {
	result := ""
	for _, param := range p.Parameters {
		result = result + param.ID
		if !strings.HasSuffix(param.ID, "=") {
			result = result + " "
		}
		// Slice must be a comma separated string
		if slice, ok := param.Value.([]interface{}); ok {
			result = result + SliceToCommaSepString(slice)
		} else {
			result = result + fmt.Sprintf("%v", param.Value)
		}
		result = result + " " //SPACE at the end!!
	}
	return result
}

// Implement String method to use Printf "%s" for debug purposes.
func (p *Plugin) String() string {

	// Convert PluginParameters to "string"
	paramsString := "\n"
	for _, param := range p.Parameters {
		paramsString = paramsString + "    " + param.ID
		if !strings.HasSuffix(param.ID, "=") {
			paramsString = paramsString + " "
		}
		paramsString = paramsString + fmt.Sprintf("%v (%s)\n", param.Value, param.ValueType)
	}

	template := `ID         : %s
Name       : %s
Description: %s
Version    : %s
Author     : %s
Category   : %s
Entrypoint : %s
Dir        : %s
Parameters : %s`
	return fmt.Sprintf(template, p.ID, p.Name, p.Description, p.Version, p.Author, p.Category, p.Entrypoint, p.Dir, paramsString)
}
