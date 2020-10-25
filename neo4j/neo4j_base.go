package neo4j

import (
	"log"
	"os"
	"sync"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/neo4j"
)

func init() {
	Neo4jLockObj = &Neo4jLock{}
}

type Neo4jLock struct {
	Lock sync.Mutex
}

var Neo4jLockObj *Neo4jLock

// each node corresponds to one opsList elem
func StartNeo4j() (session neo4jDriver.Session, errRet error) {
	passwd := os.Getenv("NEODB_PASSWD")
	if passwd == "" {
		log.Fatalf("neodb password can't be nil")
	}
	driver, errRet := neo4jDriver.NewDriver("bolt://localhost:7687/", neo4jDriver.BasicAuth("neo4j", passwd, ""))
	if errRet != nil {
		log.Printf("in newdriver err")
		return
	}
	//defer driver.Close()

	sessionConfig := neo4jDriver.SessionConfig{
		AccessMode:   neo4jDriver.AccessModeWrite,
		DatabaseName: "neo4j",
	}
	session, errRet = driver.NewSession(sessionConfig)
	if errRet != nil {
		log.Printf("in new session err")
		return
	}
	//defer session.Close()

	return
}
