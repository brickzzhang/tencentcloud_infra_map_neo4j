package base

import (
	"context"
	"strings"
	"sync"

	"tencentcloud/connection"
	"tencentcloud/neo4j"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/neo4j"
)

func init() {
	RelationshipsObj = &Relationships{
		RelationshipList: make([]*neo4j.NeoRelaOps, 0),
	}
}

type CloudService interface {
	ServiceInit(*connection.TencentCloudClient, neo4jDriver.Session) error
	Create2Neo4j(context.Context) error
}

type Relationships struct {
	RelationshipList []*neo4j.NeoRelaOps
	lock             sync.Mutex
}

func (me *Relationships) RelationshipAppend(relationship *neo4j.NeoRelaOps) error {
	me.lock.Lock()
	defer me.lock.Unlock()

	me.RelationshipList = append(me.RelationshipList, relationship)
	return nil
}

var RelationshipsObj *Relationships

var Resources = []string{"CVM", "Region", "Image"}

func RelationGen(nodeOps *neo4j.NeoNodeOps) {
	for k, v := range nodeOps.Attr {
		for _, item := range Resources {
			if strings.Contains(strings.ToLower(k), strings.ToLower(item)) {
				_ = RelationshipsObj.RelationshipAppend(&neo4j.NeoRelaOps{
					LeftNode: neo4j.NeoNodeInRela{
						Node: nodeOps.Node,
						Id:   nodeOps.Attr["Id"].(string),
						Attr: nodeOps.Attr,
					},
					RightNode: neo4j.NeoNodeInRela{
						Node: item,
						Id:   v.(string),
						Attr: map[string]interface{}{
							"Id": v,
						},
					},
				})
			}
		}
	}
	return
}
