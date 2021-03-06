// config.go - set pseudoCtx from a config file with a "config":"pseudo" JSON object.

package pseudo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clbanning/checkjson"
)

var pseudoCtx context

// ConfigJSON returns the runtime context settings that were configured 
// as a JSON object.
func ConfigJSON() string {
	j, _ := json.Marshal(pseudoCtx)
	return string(j)
}

// Config parses a JSON file with runtime settings.
// This is called by an init() perhaps with a default file name of "./pseudo.json"
// or "pseudo.config".  The init() can also handle CLI flags to override default
// settings.
func Config(file string) error {
	// read file into an array of JSON objects
	objs, err := checkjson.ReadJSONFile(file)
	if err != nil {
		return fmt.Errorf("config file: %s - %s", file, err.Error())
	}

	// get a JSON object that has "config":"pseudo" key:value pair
	type config struct {
		Config string
	}
	var ctxset bool // make sure just one pseudo config entry
	for n, obj := range objs {
		c := new(config)
		// unmarshal the object - and try and retrule a meaningful error
		if err := json.Unmarshal(obj, c); err != nil {
			return fmt.Errorf("parsing config file: %s entry: %d - %s",
				file, n+1, checkjson.ResolveJSONError(obj, err).Error())
		}
		switch strings.ToLower(c.Config) {
		case "pseudo":
			if ctxset {
				return fmt.Errorf("duplicate 'pseudo' entry in config file: %s entry: %d", file, n)
			}
			if err := checkjson.Validate(obj, pseudoCtx); err != nil {
				return fmt.Errorf("checking pseudo config JSON object: %s", err)
			}
			if err := json.Unmarshal(obj, &pseudoCtx); err != nil {
				return fmt.Errorf("config file: %s - %s", file, err)
			}
			ctxset = true
		default:
			// return fmt.Errorf("unknown config option in config file: %s entry: %d", file, n+1)
			// for now, just ignore stuff we're not interested in
		}
	}
	if !ctxset {
		return fmt.Errorf("no pseudo config object in %s", file)
	}
	return nil
}
