package neo4j

import (
	"fmt"
	"log"
	"strings"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/neo4j"
)

type NeoNodeOps struct {
	Operation string
	Node      string
	Attr      map[string]interface{}
	DeferOp   string
	DeferAttr map[string]interface{}
}

func (me *NeoNodeOps) Create(session neo4jDriver.Session) (errRet error) {
	Neo4jLockObj.Lock.Lock()
	defer Neo4jLockObj.Lock.Unlock()

	me.Operation = "create"
	var mapStr string
	var id string
	for k, v := range me.Attr {
		if k == "Id" {
			id = v.(string)
		}
		mapStr += "n." + k + "=$" + k + ", "
	}
	mapStr = strings.Trim(mapStr, ", ")
	cypher := fmt.Sprintf("merge (n: %s {Id: \"%s\"}) on %s set %s RETURN n", me.Node, id, me.Operation, mapStr)
	var cypherMap = make(map[string]interface{})
	for k, v := range me.Attr {
		cypherMap[k] = v
	}
	_, errRet = session.Run(cypher, cypherMap)
	if errRet != nil {
		log.Printf("create fail, %+v", errRet)
		return
	}

	return
}

func (me *NeoNodeOps) Read(session neo4jDriver.Session) (result neo4jDriver.Result, errRet error) {
	me.Operation = "match"
	var mapStr string
	for k := range me.Attr {
		if k == "Id" {
			mapStr += "\"" + k + "\"" + ": $" + k + ","
			break
		}
	}
	mapStr = strings.Trim(mapStr, ",")
	cypher := fmt.Sprintf("%s (n: %s {%s}) RETURN n", me.Operation, me.Node, mapStr)
	var cypherMap = make(map[string]interface{})
	for k, v := range me.Attr {
		if k == "Id" {
			cypherMap[k] = v
			break
		}
	}
	result, errRet = session.Run(cypher, cypherMap)
	log.Printf("hello here: %+v", result)
	if errRet != nil {
		log.Printf("read fail, %+v", errRet)
		return
	}

	return
}

func (me *NeoNodeOps) Update(session neo4jDriver.Session) (errRet error) {
	me.Operation = "match"
	queryStr := "Id: $id"
	var setStr string
	for k := range me.Attr {
		setStr += "set n." + k + "= $" + k + ","
	}
	setStr = strings.Trim(setStr, ",")
	cypher := fmt.Sprintf("%s (n: %s {%s}) %s RETURN n", me.Operation, me.Node, queryStr, setStr)
	var cypherMap = make(map[string]interface{})
	for k, v := range me.Attr {
		cypherMap[k] = v
		break
	}
	_, errRet = session.Run(cypher, cypherMap)
	if errRet != nil {
		log.Printf("update fail, %+v", errRet)
		return
	}

	return
}

func (me *NeoNodeOps) Delete(session neo4jDriver.Session) (errRet error) {
	me.Operation = "match"
	me.DeferOp = "delete"

	cypher := fmt.Sprintf("%s (n: %s) %s n", me.Operation, me.Node, me.DeferOp)
	_, errRet = session.Run(cypher, nil)
	if errRet != nil {
		log.Printf("delete fail, %+v", errRet)
		return
	}

	return
}
