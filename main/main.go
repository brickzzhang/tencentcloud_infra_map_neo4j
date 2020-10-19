package main

import (
	"log"
	"os"

	"tencentcloud/connection"
	"tencentcloud/neo4j"
	"tencentcloud/product/base"
	"tencentcloud/product/cvm"
	"tencentcloud/product/image"
	"tencentcloud/product/region"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func newTCloudClient(region string) (client *connection.TencentCloudClient) {
	ak := os.Getenv("TENCENTCLOUD_SECRET_ID")
	sk := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	if ak == "" || sk == "" {
		log.Fatalf("ak or sk can't be nil")
	}
	client = new(connection.TencentCloudClient)
	client.Credential = common.NewCredential(ak, sk)
	client.Region = region

	return
}

func main() {
	rClient := newTCloudClient("ap-guangzhou")
	neo4jSess, err := neo4j.StartNeo4j()
	if err != nil {
		return
	}
	regionClient := region.RegionService{
		Client:       rClient,
		Neo4jSession: neo4jSess,
	}
	regions, err := regionClient.DescribeRegions()
	if err != nil {
		log.Fatalf("describe regions error: %+v", err)
	}
	// clean existing resources
	if err = regionClient.ResourceClean(base.Resources); err != nil {
		log.Fatalf("resource clean error: %+v", err)
	}

	// generate region nodes
	if err := regionClient.Create2Neo4j(); err != nil {
		log.Fatalf("region node create2neo4j error: %+v", err)
	}

	for _, region := range regions {
		client := newTCloudClient(*region)
		cvmClient := cvm.CvmService{
			Client:       client,
			Neo4jSession: neo4jSess,
		}
		if err := cvmClient.Create2Neo4j(); err != nil {
			log.Fatalf("cvm create2neo4j error: %+v", err)
		}
		imgClient := image.ImageService{
			Client:       client,
			Neo4jSession: neo4jSess,
		}
		if err := imgClient.Create2Neo4j(); err != nil {
			log.Fatalf("image create2neo4j error: %+v", err)
		}
	}

	for _, item := range base.Relationships {
		if err := item.Create(neo4jSess); err != nil {
			log.Fatalf("create relationship error: %+v", err)
		}
	}
}
