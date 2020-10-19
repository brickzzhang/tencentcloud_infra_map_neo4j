package neo4j

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type NeoNodeInRela struct {
	Node string
	Id   string
	Attr map[string]interface{}
}

type NeoRelaOps struct {
	PreOp     string
	Operation string
	DeferOp   string
	LeftNode  NeoNodeInRela
	RightNode NeoNodeInRela
	RelaAttr  map[string]interface{}
}

func (me *NeoRelaOps) Create(session neo4j.Session) (errRet error) {
	me.PreOp = "match"
	me.Operation = "create"
	me.DeferOp = "return"

	var (
		matchStr  string
		createStr string
		cypher    string
	)
	matchStr = fmt.Sprintf("match (a: %s),(b: %s) where a.Id = \"%s\" and b.Id = \"%s\"",
		me.LeftNode.Node, me.RightNode.Node, me.LeftNode.Attr["Id"], me.RightNode.Attr["Id"])
	createStr = fmt.Sprintf("create (a)-[r:related]->(b) return r")
	cypher = matchStr + " " + createStr

	_, errRet = session.Run(cypher, nil)
	if errRet != nil {
		log.Printf("relationship create fail, %+v", errRet)
		return
	}

	return
}
