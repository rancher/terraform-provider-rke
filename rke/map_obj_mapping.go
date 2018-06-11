package rke

func applyMapToObj(mapping *mapObjMapping) {
	mapping.applyMappings()
}

type mapObjMapping struct {
	source         map[string]interface{}
	stringMapping  map[string]*string
	intMapping     map[string]*int
	boolMapping    map[string]*bool
	listStrMapping map[string]*[]string
	mapStrMapping  map[string]*map[string]string
}

func (m *mapObjMapping) applyMappings() {
	m.applyStringMapping()
	m.applyIntMapping()
	m.applyToBoolMapping()
	m.applyToStrListMapping()
	m.applyStrMapMapping()
}

func (m *mapObjMapping) applyStringMapping() {
	for k, d := range m.stringMapping {
		if v, ok := m.source[k]; ok {
			*d = v.(string)
		}
	}
}

func (m *mapObjMapping) applyIntMapping() {
	for k, d := range m.intMapping {
		if v, ok := m.source[k]; ok {
			*d = v.(int)
		}
	}
}

func (m *mapObjMapping) applyToBoolMapping() {
	for k, d := range m.boolMapping {
		if v, ok := m.source[k]; ok {
			*d = v.(bool)
		}
	}
}

func (m *mapObjMapping) applyToStrListMapping() {
	for k, d := range m.listStrMapping {
		if v, ok := m.source[k]; ok {
			var values []string
			for _, e := range v.([]interface{}) {
				values = append(values, e.(string))
			}
			*d = values
		}
	}
}

func (m *mapObjMapping) applyStrMapMapping() {
	for k, d := range m.mapStrMapping {
		if v, ok := m.source[k]; ok {
			values := map[string]string{}
			for k, v := range v.(map[string]interface{}) {
				if v, ok := v.(string); ok {
					values[k] = v
				}
			}
			*d = values
		}
	}
}
