package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
	"gopkg.in/yaml.v3"

	"github.com/gabor-boros/arangom"
)

const (
	DefaultDatabaseEndpoints = "http://localhost:8529" // default database endpoints
	DefaultMigrationDir      = "migrations"            // default migration directory
)

var (
	dbName      string // database name
	dbUsername  string // database username
	dbPassword  string // database password
	dbEndpoints string // database endpoints

	collection       string // collection in which migrations are stored
	createCollection bool   // create collection if it does not exist

	migrationDir string // directory where migrations are stored

	version = "dev"                           // version of the binary
	commit  = "dirty"                         // git commit hash
	date    = time.Now().Format(time.RFC3339) // build date
)

func init() {
	flag.StringVar(&dbUsername, "username", "", "Database user")
	flag.StringVar(&dbPassword, "password", "", "Database password")
	flag.StringVar(&dbName, "database", "", "Database name")
	flag.StringVar(&dbEndpoints, "endpoints", DefaultDatabaseEndpoints, "Comma-separated list of database endpoints")

	flag.StringVar(&migrationDir, "migration-dir", DefaultMigrationDir, "Migration directory")
	flag.StringVar(&collection, "collection", arangom.DefaultMigrationCollection, "Migration collection")
	flag.BoolVar(&createCollection, "create-collection", true, "Create migration collection if it does not exist")

	printVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if printVersion != nil && *printVersion {
		fmt.Printf("arangom version %s, commit %s (%s)\n", version, commit, date)
		os.Exit(0)
	}

	if err := validateFlags(); err != nil {
		panic(err)
	}
}

func validateFlags() error {
	if dbUsername == "" {
		return fmt.Errorf("database username is required")
	}

	if dbPassword == "" {
		return fmt.Errorf("database password is required")
	}

	if dbEndpoints == "" {
		return fmt.Errorf("database endpoints are required")
	}

	if dbName == "" {
		return fmt.Errorf("database name is required")
	}

	if collection == "" {
		return fmt.Errorf("collection name is required")
	}

	if migrationDir == "" {
		return fmt.Errorf("migration directory is required")
	}

	return nil
}

func initDatabase() (arangodb.Database, error) {
	endpoint := connection.NewRoundRobinEndpoints(strings.Split(dbEndpoints, ","))
	conn := connection.NewHttpConnection(connection.HttpConfiguration{
		Endpoint:       endpoint,
		Authentication: connection.NewBasicAuth(dbUsername, dbPassword),
	})

	client := arangodb.NewClient(conn)
	ctx := context.Background()

	db, err := client.GetDatabase(ctx, dbName, nil)
	if err != nil {
		return nil, err
	}

	collExists, err := db.CollectionExists(ctx, collection)
	if err != nil {
		return nil, err
	}

	if !collExists && createCollection {
		if _, err := db.CreateCollectionV2(ctx, collection, nil); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// loadMigration loads a single migration file.
func loadMigration(path string) (*arangom.Migration, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := new(arangom.Migration)
	if err := yaml.Unmarshal(b, m); err != nil {
		return nil, err
	}

	m.Path = path

	return m, nil
}

// loadMigrations walks the given path and loads all the migration files.
func loadMigrations(path string) ([]*arangom.Migration, error) {
	migrations := make([]*arangom.Migration, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			migration, err := loadMigration(path)
			if err != nil {
				return err
			}

			migrations = append(migrations, migration)
		}

		return nil
	})

	return migrations, err
}

func main() {
	db, err := initDatabase()
	if err != nil {
		panic(err)
	}

	migrations, err := loadMigrations(migrationDir)
	if err != nil {
		panic(err)
	}

	executor, err := arangom.NewExecutor(
		arangom.WithDatabase(db),
		arangom.WithCollection(collection),
		arangom.WithMigrations(migrations),
	)
	if err != nil {
		panic(err)
	}

	if err := executor.Execute(context.Background()); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
