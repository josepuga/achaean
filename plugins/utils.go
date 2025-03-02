package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
func typeof(v any) string {
	return fmt.Sprintf("%T", v)
}
*/

func typeof(v any) string {
	switch v.(type) {
	case string:
		return "string"
	case int, int64, int32, int16, int8:
		return "int"
	case float64, float32:
		return "int" //Float is like an int here!!!
	case []interface{}:
		return "list"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

func SliceToCommaSepString(slice []interface{}) string {
    strSlice := make([]string, len(slice))
    for i, v := range slice {
        strSlice[i] = fmt.Sprintf("%v", v)       
    }
    return strings.Join(strSlice, ",")
}

/*
To get the % done...
*/
func GetProgressPipeName(pluginID string) string {
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("plugin_progress_%s", pluginID)
	return filepath.Join(tmpDir, fileName)
}
