package tnt

type registry map[string]interface{}
type registryRactory func() interface{}

func (r *registry) Get(k string, f registryRactory) interface{} {
	entry, ok := (*r)[k]

	if !ok {
		entry = f()
		(*r)[k] = entry
	}

	return entry
}
