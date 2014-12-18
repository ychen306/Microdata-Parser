package mcrdata

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

type Node struct {
	root     xml.Node
	basePath string
}

var (
	propSearchTmpl  string            = "%s//*[@itemprop]"
	allScopePath    *xpath.Expression = xpath.Compile("ancestor::*[@itemscope and @itemtype]")
	scopeSearchTmpl string            = "ancestor::*[@itemscope and @itemtype=\"%s\"]"
)

func getAttr(attrs map[string]*xml.AttributeNode, attr string) (val string, ok bool) {
	if attrNode := attrs[attr]; attrNode != nil {
		val = attrNode.Content()
		ok = true
	}
	return
}

func getPropVal(node xml.Node) (prop string, val Property, ok bool) {
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
	var content string
	switch node.Name() {
	case "img":
		content, ok = getAttr(attrs, "src")
	case "a":
		content, ok = getAttr(attrs, "href")
	default:
		attrContent, hasContent := getAttr(attrs, "content")
		if !hasContent {
			content = node.Content()
		} else {
			content = attrContent
		}
	}
	val = PlainData(content)
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
		properties: make(map[string][]Property),
		itemType:   itype}
	return
}

func (node *Node) find(scopeSearchPath *xpath.Expression, itype string) (found []*Item, err error) {
	searchPath := xpath.Compile(fmt.Sprintf(propSearchTmpl, node.basePath))
	propNodes, err := node.root.Search(searchPath)
	if err != nil {
		return
	}

	found = make([]*Item, 0)
	items := make(map[string]*Item)

	for _, propNode := range propNodes {
		prop, val, ok := getPropVal(propNode)
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
		if val.Type() == itype {
			item := val.(*Item)
			items[propNode.Path()] = item
			found = append(found, item)
		}
	}
	return
}

func Parse(page []byte) (*Node, error) {
	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		return nil, err
	}
	return ParseXmlNode(doc.Root()), nil
}

func ParseXmlNode(node xml.Node) *Node {
	return &Node{
		root:     node,
		basePath: node.Path()}
}

func (node *Node) FindAll() ([]*Item, error) {
	return node.find(allScopePath, "")
}

func (node *Node) Find(itype string) ([]*Item, error) {
	scopeSearchPath := xpath.Compile(fmt.Sprintf(scopeSearchTmpl, itype))
	return node.find(scopeSearchPath, itype)
}
