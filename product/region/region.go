package region

import (
	"fmt"
	"log"

	"tencentcloud/connection"
	"tencentcloud/neo4j"
	"tencentcloud/product/base"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/neo4j"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type RegionService struct {
	Client       *connection.TencentCloudClient
	Neo4jSession neo4jDriver.Session
}

func (me *RegionService) ResourceClean(resources []string) (errRet error) {
	for _, item1 := range resources {
		for _, item := range resources {
			cypher := fmt.Sprintf("match (a: %s)-[r:related]-(b: %s) delete r", item1, item)
			_, errRet = me.Neo4jSession.Run(cypher, nil)
			if errRet != nil {
				log.Printf("clean resource error, %+v", errRet)
				return
			}
		}
	}
	for _, item := range resources {
		cypher := fmt.Sprintf("match (n: %s) delete (n)", item)
		_, errRet = me.Neo4jSession.Run(cypher, nil)
		if errRet != nil {
			log.Printf("clean resource error, %+v", errRet)
			return
		}
	}

	return
}

func (me *RegionService) DescribeRegions() (regions []*string, errRet error) {
	request := cvm.NewDescribeRegionsRequest()
	response, errRet := me.Client.UseCvmClient().DescribeRegions(request)
	if errRet != nil {
		log.Printf("%+v", errRet)
	}
	regions = make([]*string, 0)
	for _, item := range response.Response.RegionSet {
		regions = append(regions, item.Region)
	}
	return
}

func (me *RegionService) transfer2Neo4j(region string) (neo4jRet *neo4j.NeoNodeOps) {
	var attrMap = make(map[string]interface{})
	attrMap["Id"] = region
	neo4jRet = &neo4j.NeoNodeOps{
		Operation: "CREATE",
		DeferOp:   "RETURN",
		Node:      "Region",
		Attr:      attrMap,
	}
	// create relationship
	base.RelationGen(neo4jRet)
	return
}

func (me *RegionService) Read2Neo4j() (neo4jList []*neo4j.NeoNodeOps, errRet error) {
	regions, errRet := me.DescribeRegions()
	if errRet != nil {
		return
	}
	for _, item := range regions {
		neo4jList = append(neo4jList, me.transfer2Neo4j(*item))
	}
	return
}

func (me *RegionService) Create2Neo4j() (errRet error) {
	neoList, errRet := me.Read2Neo4j()
	if errRet != nil {
		log.Printf("%+v", errRet)
		return
	}

	for _, item := range neoList {
		errRet = item.Create(me.Neo4jSession)
		if errRet != nil {
			log.Printf("create error: %+v", errRet)
			return
		}
	}
	return
}
