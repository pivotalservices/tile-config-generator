package generator

type Template struct {
	NetworkProperties *NetworkProperties     `yaml:"network-properties"`
	ProductProperties map[string]interface{} `yaml:"product-properties"`
	ResourceConfig    map[string]Resource    `yaml:"resource-config"`
}

type FormType struct {
	Description string     `yaml:"description"`
	Label       string     `yaml:"label"`
	Name        string     `yaml:"name"`
	Properties  []Property `yaml:"property_inputs"`
}

type Property struct {
	Description string             `yaml:"description"`
	Label       string             `yaml:"label"`
	Placeholder string             `yaml:"placeholder"`
	Reference   string             `yaml:"reference"`
	Selectors   []SelectorProperty `yaml:"selector_property_inputs"`
}

type SelectorProperty struct {
	Label      string     `yaml:"label"`
	Reference  string     `yaml:"reference"`
	Properties []Property `yaml:"property_inputs"`
}

type Option struct {
	Label string      `json:"label"`
	Name  interface{} `json:"name"`
}

type Ops struct {
	Type  string      `yaml:"type"`
	Path  string      `yaml:"path"`
	Value interface{} `yaml:"value,omitempty"`
}
