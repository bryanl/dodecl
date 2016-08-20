package dodecl

import "gopkg.in/yaml.v2"

// Dodecl is a set of declarative rules for infrastructure.
type Dodecl struct {
	Resources []Resource `yaml:"resources"`
}

// Resource defines an infrastructure resource that needs to be constructed.
type Resource struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	Properties Properties `yaml:"properties"`
}

// Properties are properties for a Resource.
type Properties map[string]interface{}

// ReadFromYAML generates a Dodecl instance from a string.
func ReadFromYAML(yamlData []byte) (*Dodecl, error) {
	d := &Dodecl{}

	err := yaml.Unmarshal(yamlData, d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
