package mcrdata

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

type Node struct {
	Data xml.Node
}

var (
	allPropPath  *xpath.Expression = xpath.Compile(".//*[@itemprop]")
	allScopePath *xpath.Expression = xpath.Compile("ancestor::*[@itemscope and @itemtype]")
	scopeTmpl    string            = "ancestor::*[@itemscope and @itemtype=\"%s\"]"
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

func getScope(scopeSearchPath *xpath.Expression, node xml.Node) (scopePath, scopeType string, ok bool) {
	scopes, _ := node.Search(scopeSearchPath)
	scopeCount := len(scopes)
	if scopeCount > 0 {
		ok = true
		scope := scopes[scopeCount-1]
		scopePath = scope.Path()
		scopeType, _ = getAttr(scope.Attributes(), "itemtype")
	}
	return
}

func makeItem(itype string) (item *Item) {
	item = &Item{
		properties: make(map[string][]interface{}),
		itemType:   itype}
	return
}

func (node *Node) find(scopeSearchPath *xpath.Expression, itype string) (found []*Item, err error) {

	propNodes, err := node.Data.Search(allPropPath)
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

		scopePath, scopeType, ok := getScope(scopeSearchPath, propNode)
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
			if itype == "" || itype == item.itemType {
				items[propNode.Path()] = item
				found = append(found, item)
			}
		}
	}

	return
}

func Parse(page []byte) (*Node, error) {
	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		return nil, err
	}
	return &Node{Data: doc.Root()}, nil
}

func (node *Node) FindAll() ([]*Item, error) {
	return node.find(allScopePath, "")
}

func (node *Node) Find(itype string) ([]*Item, error) {
	scopeSearchPath := xpath.Compile(fmt.Sprintf(scopeTmpl, itype))
	return node.find(scopeSearchPath, itype)
}
