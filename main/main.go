package main

import (
	"context"
	"log"
	"os"

	"tencentcloud/connection"
	"tencentcloud/neo4j"
	"tencentcloud/product/base"
	"tencentcloud/product/region"
	"tencentcloud/product/register"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"golang.org/x/sync/errgroup"
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
	ctx := context.Background()
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
	if err := regionClient.Create2Neo4j(ctx); err != nil {
		log.Fatalf("region node create2neo4j error: %+v", err)
	}

	//	group, _ := errgroup.WithContext(ctx)
	// api frequency limit
	for _, region := range regions {
		//		group.Go(func() error {
		client := newTCloudClient(*region)
		group, groupCtx := errgroup.WithContext(ctx)
		for k, v := range register.ServiceMap {
			innerK := k
			innerV := v
			group.Go(func() error {
				innerV.ServiceInit(client, neo4jSess)
				if err := innerV.Create2Neo4j(groupCtx); err != nil {
					log.Fatalf("%s create2neo4j error: %+v", innerK, err)
					return err
				}
				return nil
			})
		}
		//			return group.Wait()
		_ = group.Wait()
		//		})
	}
	//	_ = group.Wait()

	for _, item := range base.RelationshipsObj.RelationshipList {
		if err := item.Create(neo4jSess); err != nil {
			log.Fatalf("create relationship error: %+v", err)
		}
	}
}
