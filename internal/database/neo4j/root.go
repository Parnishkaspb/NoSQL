package neo4jhelper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver   neo4j.DriverWithContext
	initOnce sync.Once
	initErr  error
	ctx      = context.Background()
)

func initNeo4j() {
	_ = godotenv.Load(".env")

	host := os.Getenv("NEO4J_HOST")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	if host == "" || username == "" || password == "" {
		initErr = fmt.Errorf("не заданы обязательные переменные окружения: NEO4J_HOST, NEO4J_USERNAME, NEO4J_PASSWORD")
		log.Println("❌", initErr)
		return
	}

	uri := fmt.Sprintf("bolt://%s", strings.TrimSpace(host))

	for i := 1; i <= 5; i++ {
		driver, initErr = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
			config.MaxConnectionPoolSize = 50
			config.ConnectionAcquisitionTimeout = 10 * time.Second
		})

		if initErr == nil {
			pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := driver.VerifyConnectivity(pingCtx)
			cancel()
			if err == nil {
				log.Println("✅ Успешное подключение к Neo4j")
				return
			}
			initErr = fmt.Errorf("Neo4j не отвечает: %v", err)
		}

		log.Printf("🔄 Попытка %d: %v", i, initErr)
		time.Sleep(2 * time.Second)
	}

	log.Print("❌ Не удалось подключиться к Neo4j")
}

func getDriver() (neo4j.DriverWithContext, error) {
	initOnce.Do(initNeo4j)
	return driver, initErr
}

func CreateRelation(relationType string, fromLabel, toLabel string, fromID, toID string) error {
	drv, err := getDriver()
	if err != nil {
		return err
	}

	session := drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := fmt.Sprintf(`
			MATCH (a:%s {id: $fromID})
			MATCH (b:%s {id: $toID})
			MERGE (a)-[:%s]->(b)
		`, fromLabel, toLabel, relationType)

		params := map[string]any{
			"fromID": fromID,
			"toID":   toID,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func CreateNode[T any](label string, id string, obj T) error {
	drv, err := getDriver()
	if err != nil {
		return err
	}

	props, err := toMap(obj)
	if err != nil {
		return err
	}

	session := drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := fmt.Sprintf(`MERGE (n:%s {id: $id}) SET n += $props`, label)
		params := map[string]any{
			"id":    id,
			"props": props,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}

func Close() {
	if driver != nil {
		_ = driver.Close(ctx)
	}
}
