Microdata Parser in Go. 
Usage:
```go
import (
  "mcrdata"
  ...
)

resp, _ := http.Get(url)
page, _ := ioutil.ReadAll(resp.Body)
node := mcrdata.Parse(page)
items := node.FindAll()
```
APIs: 
To prepare a document before searching items, use `func Parse`. 
```go
func Parse(page []byte) *Node {...}
```
To search items within a Node, use `func Findall` for all items or `func Find` for item of a specific type. 
```go
func (node *Node) FindAll() []*Item {...}
func (node *Node) Find(itemtype string) []*Item {...}
```
To get property of an item, use `func Get`, which returns a slice of `Property`. 
```go
func (item *Item) Get(property string) []Property {...}
```
Here is the list of functions implemented by `Property`.
```go
func Value() string {...} // returns its value if the property is a plain data
func Properties() []string {...} // returns a slice of properties if the property is an Item
func Type() string {...} // returns type of the property
func Get(property string) []Property {...} // returns values of given property if the property is an Item
```
Note that `Node` is a type within package `mcrdata`, which internally uses `Node` from [gokogiri](https://github.com/moovweb/gokogiri/). 

To search items within an existing gokogiri node, wrap it using `func ParseNode`
```go
func ParseNode(node *xml.Node) *Node {...}
```
