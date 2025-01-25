package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnToDbPostgres() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_URL"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	var now time.Time
	db.Raw("SELECT NOW()").Scan(&now)

	fmt.Println(now)
}

func conf(c *neo4j.Config) {
	c.Log = nil
}

func ConnToNeoDb() neo4j.DriverWithContext {
	ctx := context.Background()
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := os.Getenv("DB_URL")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""), conf)

	if err != nil {
		panic(err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	return driver
}
