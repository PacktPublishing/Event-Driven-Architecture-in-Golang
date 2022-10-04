package ddd

type Metadata map[string]any

func (m Metadata) Set(key string, value any) {
	m[key] = value
}

func (m Metadata) Get(key string) any {
	return m[key]
}

func (m Metadata) Del(key string) {
	delete(m, key)
}

func (m Metadata) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m Metadata) configureEvent(e *event) {
	for key, value := range m {
		e.metadata[key] = value
	}
}

func (m Metadata) configureCommand(c *command) {
	for key, value := range m {
		c.metadata[key] = value
	}
}

func (m Metadata) configureReply(r *reply) {
	for key, value := range m {
		r.metadata[key] = value
	}
}
