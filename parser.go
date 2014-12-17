package mcrdata

import (
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

var (
	allPropPath     *xpath.Expression = xpath.Compile(".//*[@itemprop]")
	scopeSearchPath *xpath.Expression = xpath.Compile("ancestor::*[@itemscope and @itemtype]")
)

func getAttr(attrs map[string]*xml.AttributeNode, attr string) (val string, ok bool) {
	if attrNode := attrs[attr]; attrNode != nil {
		val = attrNode.Content()
		ok = true
	}
	return
}

func getPropAndVal(node xml.Node) (prop string, val interface{}, ok bool) {
	attrs := node.Attributes()
	ok = true

	// get property name
	prop, _ = getAttr(attrs, "itemprop")

	// get property value
	// property is an item
	if _, newScope := getAttr(attrs, "itemscope"); newScope {
		itype, _ := getAttr(attrs, "itemtype")
		val = makeItem(itype)
		return
	}

	// property is a plain datatype
	switch node.Name() {
	case "img":
		val, ok = getAttr(attrs, "src")
	case "a":
		val, ok = getAttr(attrs, "href")
	default:
		content, hasContent := getAttr(attrs, "content")
		if !hasContent {
			val = node.Content()
		} else {
			val = content
		}
	}
	return
}

func getScope(node xml.Node) (scopePath, scopeType string, ok bool) {
	scopes, _ := node.Search(scopeSearchPath)
	scopeCount := len(scopes)
	if scopeCount > 0 {
		ok = true
		scope := scopes[scopeCount-1]
		scopePath = scope.Path()
		scopeType, _ = getAttr(node.Attributes(), "itemtype")
	}
	return
}

func makeItem(itype string) (item *Item) {
	item = &Item{
		properties: make(map[string][]interface{}),
		itemType:   itype}
	return
}

func find(page []byte, searchPath *xpath.Expression) (found []*Item, err error) {
	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		return
	}

	propNodes, err := doc.Root().Search(searchPath)
	if err != nil {
		return
	}

	found = make([]*Item, 0)
	items := make(map[string]*Item)

	for _, propNode := range propNodes {
		prop, val, ok := getPropAndVal(propNode)
		if !ok {
			continue
		}

		scopePath, scopeType, ok := getScope(propNode)
		if !ok {
			continue
		}

		if items[scopePath] == nil {
			item := makeItem(scopeType)
			items[scopePath] = item
			found = append(found, item)
		}

		items[scopePath].addProp(prop, val)
		switch item := val.(type) {
		case *Item:
			items[propNode.Path()] = item
			found = append(found, item)
		}
	}

	return
}

func FindAll(page []byte) ([]*Item, error) {
	return find(page, allPropPath)
}
