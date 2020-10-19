package base

import (
	"strings"

	"tencentcloud/neo4j"
)

var Resources = []string{"CVM", "Region", "Image"}
var Relationships = make([]*neo4j.NeoRelaOps, 0)

func RelationGen(nodeOps *neo4j.NeoNodeOps) {
	for k, v := range nodeOps.Attr {
		for _, item := range Resources {
			if strings.Contains(strings.ToLower(k), strings.ToLower(item)) {
				Relationships = append(Relationships, &neo4j.NeoRelaOps{
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
