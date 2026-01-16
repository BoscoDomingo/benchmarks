package main

import "strconv"

// PrefDef defines a preference with display name, description, and default value.
type PrefDef struct {
	DisplayName  string
	Description  string
	DefaultValue string
}

// PreferenceDefinitions maps snake_case keys to their definitions.
var PreferenceDefinitions = map[string]PrefDef{
	"min_runs": {
		DisplayName:  "Minimum Runs",
		Description:  "Minimum number of benchmark iterations",
		DefaultValue: "300",
	},
	"save_results": {
		DisplayName:  "Save Results",
		Description:  "Export to results.md (yes/no)",
		DefaultValue: "yes",
	},
}

// PreferenceKeys defines the order in which preferences are displayed.
var PreferenceKeys = []string{"min_runs", "save_results"}

// Preferences holds user-configurable run settings as string values.
type Preferences struct {
	values map[string]string
}

func newPreferences() *Preferences {
	values := make(map[string]string)
	for key, def := range PreferenceDefinitions {
		values[key] = def.DefaultValue
	}
	return &Preferences{values: values}
}

func (p *Preferences) Get(key string) string {
	return p.values[key]
}

func (p *Preferences) GetInt(key string) int {
	n, _ := strconv.Atoi(p.values[key])
	return n
}

func (p *Preferences) GetBool(key string) bool {
	v := p.values[key]
	return v == "yes" || v == "true" || v == "1"
}

func (p *Preferences) Set(key, value string) {
	p.values[key] = value
}
