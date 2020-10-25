package cvm

import (
	"context"
	"log"

	"tencentcloud/connection"
	"tencentcloud/neo4j"
	"tencentcloud/product/base"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/neo4j"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type CvmService struct {
	Client       *connection.TencentCloudClient
	Neo4jSession neo4jDriver.Session
}

func (me *CvmService) ServiceInit(client *connection.TencentCloudClient, neo4jSession neo4jDriver.Session) error {
	me.Client = client
	me.Neo4jSession = neo4jSession
	return nil
}

func (me *CvmService) DescribeInstances() (instances []*cvm.Instance, errRet error) {
	request := cvm.NewDescribeInstancesRequest()
	params := "{}"
	errRet = request.FromJsonString(params)
	if errRet != nil {
		return
	}
	response, errRet := me.Client.UseCvmClient().DescribeInstances(request)
	if errRet != nil {
		return
	}
	instances = response.Response.InstanceSet
	return
}

func (me *CvmService) transfer2Neo4j(ins *cvm.Instance) (neo4jRet *neo4j.NeoNodeOps) {
	var attrMap = make(map[string]interface{})
	attrMap["Id"] = *ins.InstanceId
	attrMap["ImageId"] = *ins.ImageId
	attrMap["RegionId"] = me.Client.Region
	neo4jRet = &neo4j.NeoNodeOps{
		Operation: "CREATE",
		DeferOp:   "RETURN",
		Node:      "CVM",
		Attr:      attrMap,
	}
	// create relationship
	base.RelationGen(neo4jRet)
	return
}

func (me *CvmService) read2Neo4j() (neo4jList []*neo4j.NeoNodeOps, errRet error) {
	ins, errRet := me.DescribeInstances()
	if errRet != nil {
		return
	}
	for _, item := range ins {
		neo4jList = append(neo4jList, me.transfer2Neo4j(item))
	}
	return
}

func (me *CvmService) Create2Neo4j(ctx context.Context) (errRet error) {
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	neoList, errRet := me.read2Neo4j()
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
