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
items := node.FindAll(doc)
```
APIs:
To prepare a web page before searching for items, use `func Parse`.
```go
func Parse(page []byte) *Node {...}
```
To search items given within a Node, use `func Findall` for all items or `func Find` for item of a specific type. 
```go
func (node *Node) FindAll() []*Item {...}
func (node *Node) Find(itemtype string) []*Item {...}
```
To get property of an item, use `func Get`, which returns a slice of values, type of which can be either `string` or `Item`.
```go
func (item *Item) Get(property string) []interface{} {...}
```
To get type of an item, use `func Type`.
```go
func (item *Item) Type() string {...}
```
To get available properties of an item, use `func Properties`.
```go
func (item *Item) Properties() []string {...}
```
Note that `Node` is a type within package `mcrdata`, which internally uses a `Node` from gokogiri(https://github.com/moovweb/gokogiri/).
To search items within an existing `gokogiri node`, wrap it like this
```go
mcrdataNode := mcrdata.Node{Data: gokogiriNode}
```
