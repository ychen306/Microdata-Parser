package mcrdata

type Item struct {
	properties map[string][]Property
	itemType   string
}

type Property interface {
	Properties() []string
	Get(property string) []Property
	Type() string
	Value() string
}

type PlainData string

// constructor of Item
func makeItem(itype string) (item *Item) {
	item = &Item{
		properties: make(map[string][]Property),
		itemType:   itype}
	return
}

// given a property name, get all the values
func (item *Item) Get(prop string) []Property {
	return item.properties[prop]
}

// get type of the item
func (item *Item) Type() string {
	return item.itemType
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

func (item *Item) Value() string {
	return ""
}

// add property
func (item *Item) addProp(prop string, val Property) {
	propVals := item.properties[prop]
	propVals = append(propVals, val)
	item.properties[prop] = propVals
}

func (data PlainData) Properties() []string {
	return nil
}

func (data PlainData) Get(_ string) []Property {
	return nil
}

func (data PlainData) Type() string {
	return "Text"
}

func (data PlainData) Value() string {
	return string(data)
}
