package image

import (
	"log"

	"tencentcloud/connection"
	"tencentcloud/neo4j"
	"tencentcloud/product/base"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/neo4j"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type ImageService struct {
	Client       *connection.TencentCloudClient
	Neo4jSession neo4jDriver.Session
}

func (me *ImageService) DescribeImages() (images []*cvm.Image, errRet error) {
	request := cvm.NewDescribeImagesRequest()
	params := "{}"
	errRet = request.FromJsonString(params)
	if errRet != nil {
		return
	}
	response, errRet := me.Client.UseCvmClient().DescribeImages(request)
	if errRet != nil {
		return
	}
	images = response.Response.ImageSet
	return
}

func (me *ImageService) transfer2Neo4j(img *cvm.Image) (neo4jRet *neo4j.NeoNodeOps) {
	var attrMap = make(map[string]interface{})
	attrMap["Id"] = *img.ImageId
	attrMap["OsName"] = *img.OsName
	//	attrMap["RegionId"] = me.Client.Region
	neo4jRet = &neo4j.NeoNodeOps{
		Operation: "CREATE",
		DeferOp:   "RETURN",
		Node:      "Image",
		Attr:      attrMap,
	}
	// create relationship
	base.RelationGen(neo4jRet)
	return
}

func (me *ImageService) read2Neo4j() (neo4jList []*neo4j.NeoNodeOps, errRet error) {
	img, errRet := me.DescribeImages()
	if errRet != nil {
		return
	}
	for _, item := range img {
		neo4jList = append(neo4jList, me.transfer2Neo4j(item))
	}
	return
}

func (me *ImageService) Create2Neo4j() (errRet error) {
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
