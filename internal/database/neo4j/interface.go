package neo4jhelper

type Neo4j interface {
	connectNeo4j()

	Read(typeDB int, key string) (any, error)
	CreateUpdate(typeDB int, key string, value any) (string, error)
	Delete(typeDB int, key string) error
}
