package neo4jhelper

type Neo4j interface {
	connectNeo4j()

	CreateRelation(relationType string, fromLabel, toLabel string, fromID, toID string) error
	CreateNode(label string, id string, obj any) error
	//CreateUpdate(typeDB int, key string, value any) (string, error)
	//Delete(typeDB int, key string) error
}
