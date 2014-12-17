Microdata Parser in Go. 
Usage:
```go
import (
  "mcrdata"
  ...
)

resp, _ := http.Get(url)
doc, _ := ioutil.ReadAll(resp.Body)
items := mcrdata.FindAll(doc)
```
