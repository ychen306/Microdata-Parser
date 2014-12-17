package mcrdata

// an item either has a list of properties, or a value (content), not both
type Item struct {
	properties map[string][]interface{}
	itemType   string
}

// given a property name, get all the values
func (item *Item) Get(prop string) []interface{} {
	return item.properties[prop]
}

// get list of properties of an item
func (item *Item) Properties() (properties []string) {
	properties = make([]string, len(item.properties))
	i := 0
	for p, _ := range item.properties {
		properties[i] = p
		i++
	}
	return
}

// add property
func (item *Item) addProp(prop string, val interface{}) {
	propVals := item.properties[prop]
	propVals = append(propVals, val)
	item.properties[prop] = propVals
}
