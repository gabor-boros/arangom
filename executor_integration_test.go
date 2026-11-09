//go:build integration

package arangom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"github.com/arangodb/go-driver/v2/connection"
	"github.com/testcontainers/testcontainers-go"
	tcarangodb "github.com/testcontainers/testcontainers-go/modules/arangodb"
	"gopkg.in/yaml.v3"
)

const (
	integrationImage       = "arangodb:3.12.13"
	integrationScenarioDir = "testdata/integration"
	maxDatabaseNameLength  = 64

	// integrationRootPassword is only used by the disposable test container.
	integrationRootPassword = "arangom-integration" //nolint:gosec // disposable local container credential
)

var integrationClient arangodb.Client

type quietLogger struct{}

var _ Logger = quietLogger{}

func (quietLogger) Info(...any)           {}
func (quietLogger) Infof(string, ...any)  {}
func (quietLogger) Error(...any)          {}
func (quietLogger) Errorf(string, ...any) {}

func (quietLogger) Fatal(args ...any) {
	panic(fmt.Sprint(args...))
}

func (quietLogger) Fatalf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

type scenario struct {
	File          string               `yaml:"-"`
	Name          string               `yaml:"name"`
	ExecuteCount  int                  `yaml:"execute_count"`
	WantError     bool                 `yaml:"want_error"`
	ErrorContains string               `yaml:"error_contains"`
	Migrations    []scenarioMigration  `yaml:"migrations"`
	Expect        scenarioExpectations `yaml:"expect"`
}

type scenarioMigration struct {
	ID         int          `yaml:"id"`
	Path       string       `yaml:"path"`
	Operations []*Operation `yaml:"operations"`
}

type scenarioExpectations struct {
	Collections      []namedExistence       `yaml:"collections"`
	Graphs           []namedExistence       `yaml:"graphs"`
	Views            []namedExistence       `yaml:"views"`
	Analyzers        []typedExistence       `yaml:"analyzers"`
	Indexes          []indexExpectation     `yaml:"indexes"`
	Documents        []documentExpectation  `yaml:"documents"`
	MigrationRecords []migrationExpectation `yaml:"migration_records"`
}

type namedExistence struct {
	Name   string `yaml:"name"`
	Exists bool   `yaml:"exists"`
}

type typedExistence struct {
	Name   string `yaml:"name"`
	Exists bool   `yaml:"exists"`
	Type   string `yaml:"type"`
}

type indexExpectation struct {
	Collection string `yaml:"collection"`
	Name       string `yaml:"name"`
	Exists     bool   `yaml:"exists"`
	Type       string `yaml:"type"`
}

type documentExpectation struct {
	Collection string           `yaml:"collection"`
	Count      *int             `yaml:"count"`
	Contains   []map[string]any `yaml:"contains"`
}

type migrationExpectation struct {
	Name   string `yaml:"name"`
	Status string `yaml:"status"`
}

// TestMain starts one ArangoDB container for every integration scenario.
func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := startArango(ctx)
	if err != nil {
		log.Fatalf("start arangodb: %v", err)
	}

	code := runIntegration(ctx, m, container)
	if termErr := testcontainers.TerminateContainer(container); termErr != nil {
		log.Printf("terminate arangodb: %v", termErr)

		if code == 0 {
			code = 1
		}
	}

	os.Exit(code)
}

func startArango(ctx context.Context) (*tcarangodb.Container, error) {
	// The image entrypoint treats "arangod" as the server command and keeps
	// the following arguments. Vector indexes are disabled without this flag.
	container, err := tcarangodb.Run(
		ctx,
		integrationImage,
		tcarangodb.WithRootPassword(integrationRootPassword),
		testcontainers.WithCmd("arangod", "--vector-index"),
	)
	if err != nil {
		if container != nil {
			err = errors.Join(err, testcontainers.TerminateContainer(container))
		}

		return nil, err
	}

	return container, nil
}

func runIntegration(ctx context.Context, m *testing.M, container *tcarangodb.Container) int {
	endpoint, err := container.HTTPEndpoint(ctx)
	if err != nil {
		log.Printf("arangodb endpoint: %v", err)
		return 1
	}

	user, password := container.Credentials()
	conn := connection.NewHttpConnection(connection.HttpConfiguration{
		Endpoint:       connection.NewRoundRobinEndpoints([]string{endpoint}),
		Authentication: connection.NewBasicAuth(user, password),
	})
	integrationClient = arangodb.NewClient(conn)

	if _, err := integrationClient.Version(ctx); err != nil {
		log.Printf("arangodb version: %v", err)
		return 1
	}

	return m.Run()
}

// TestMigrationScenarios runs each YAML scenario in testdata/integration
// against a new database.
func TestMigrationScenarios(t *testing.T) {
	scenarios := loadScenarios(t)
	for i := range scenarios {
		sc := &scenarios[i]
		t.Run(sc.Name, func(t *testing.T) {
			runScenario(t, sc)
		})
	}
}

func loadScenarios(t *testing.T) []scenario {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(integrationScenarioDir, "*.yaml"))
	if err != nil {
		t.Fatalf("list scenarios: %v", err)
	}

	if len(matches) == 0 {
		t.Fatalf("no scenarios in %s", integrationScenarioDir)
	}

	sort.Strings(matches)

	scenarios := make([]scenario, 0, len(matches))
	names := make(map[string]string, len(matches))

	for _, path := range matches {
		sc, err := loadScenario(path)
		if err != nil {
			t.Fatalf("load %s: %v", path, err)
		}

		if previous, ok := names[sc.Name]; ok {
			t.Fatalf("duplicate scenario %q in %s and %s", sc.Name, previous, path)
		}

		names[sc.Name] = path
		scenarios = append(scenarios, sc)
	}

	return scenarios
}

func loadScenario(path string) (scenario, error) {
	body, err := os.ReadFile(path) //nolint:gosec // fixtures are read from the module testdata directory
	if err != nil {
		return scenario{}, err
	}

	var sc scenario
	if err := yaml.Unmarshal(body, &sc); err != nil {
		return scenario{}, err
	}

	sc.File = path
	if err := validateScenario(sc); err != nil {
		return scenario{}, err
	}

	return sc, nil
}

func validateScenario(sc scenario) error {
	if strings.TrimSpace(sc.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if sc.ExecuteCount < 0 {
		return fmt.Errorf("execute_count must be zero or greater")
	}

	if len(sc.Migrations) == 0 {
		return fmt.Errorf("migrations are required")
	}

	for i, migration := range sc.Migrations {
		if strings.TrimSpace(migration.Path) == "" {
			return fmt.Errorf("migrations[%d].path is required", i)
		}

		if len(migration.Operations) == 0 {
			return fmt.Errorf("migrations[%d].operations are required", i)
		}

		for j, operation := range migration.Operations {
			if operation == nil {
				return fmt.Errorf("migrations[%d].operations[%d] is empty", i, j)
			}
		}
	}

	if err := validateExpectations(sc.Expect); err != nil {
		return err
	}

	return nil
}

func validateExpectations(expect scenarioExpectations) error {
	for i, item := range expect.Collections {
		if item.Name == "" {
			return fmt.Errorf("expect.collections[%d].name is required", i)
		}
	}

	for i, item := range expect.Graphs {
		if item.Name == "" {
			return fmt.Errorf("expect.graphs[%d].name is required", i)
		}
	}

	for i, item := range expect.Views {
		if item.Name == "" {
			return fmt.Errorf("expect.views[%d].name is required", i)
		}
	}

	for i, item := range expect.Analyzers {
		if item.Name == "" {
			return fmt.Errorf("expect.analyzers[%d].name is required", i)
		}
	}

	for i, item := range expect.Indexes {
		if item.Collection == "" || item.Name == "" {
			return fmt.Errorf("expect.indexes[%d].collection and name are required", i)
		}
	}

	for i, item := range expect.Documents {
		if item.Collection == "" || item.Count == nil {
			return fmt.Errorf("expect.documents[%d].collection and count are required", i)
		}
	}

	for i, item := range expect.MigrationRecords {
		if item.Name == "" || item.Status == "" {
			return fmt.Errorf("expect.migration_records[%d].name and status are required", i)
		}
	}

	return nil
}

func runScenario(t *testing.T, sc *scenario) {
	t.Helper()

	ctx := t.Context()
	db := newScenarioDatabase(t, sc.File)
	migrations := make([]*Migration, 0, len(sc.Migrations))

	for _, item := range sc.Migrations {
		migrations = append(migrations, &Migration{
			ID:         item.ID,
			Path:       item.Path,
			Operations: item.Operations,
		})
	}

	executor, err := NewExecutor(
		WithDatabase(db),
		WithCollection(DefaultMigrationCollection),
		WithMigrations(migrations),
		WithLogger(quietLogger{}),
	)
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}

	executeCount := sc.ExecuteCount
	if executeCount == 0 {
		executeCount = 1
	}

	for attempt := 1; attempt <= executeCount; attempt++ {
		err := executor.Execute(ctx)
		if sc.WantError {
			assertExecutionError(t, attempt, err, sc.ErrorContains)
			continue
		}

		if err != nil {
			t.Fatalf("execution %d: %v", attempt, err)
		}
	}

	assertExpectations(t, db, sc.Expect)
}

func assertExecutionError(t *testing.T, attempt int, err error, contains string) {
	t.Helper()

	if err == nil {
		t.Fatalf("execution %d: expected an error", attempt)
	}

	if contains != "" && !strings.Contains(err.Error(), contains) {
		t.Fatalf("execution %d: error %q does not contain %q", attempt, err.Error(), contains)
	}
}

func newScenarioDatabase(t *testing.T, file string) arangodb.Database {
	t.Helper()

	ctx := t.Context()
	name := scenarioDatabaseName(file)

	db, err := integrationClient.CreateDatabase(ctx, name, nil)
	if err != nil {
		t.Fatalf("create database %s: %v", name, err)
	}

	t.Cleanup(func() {
		if err := db.Remove(context.Background()); err != nil {
			t.Errorf("remove database %s: %v", name, err)
		}
	})

	if _, err := db.CreateCollectionV2(ctx, DefaultMigrationCollection, nil); err != nil {
		t.Fatalf("create migration collection: %v", err)
	}

	return db
}

func scenarioDatabaseName(file string) string {
	base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

	var name strings.Builder
	name.WriteString("it_")

	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			name.WriteRune(r)
		default:
			name.WriteByte('_')
		}
	}

	if name.Len() > maxDatabaseNameLength {
		return name.String()[:maxDatabaseNameLength]
	}

	return name.String()
}

func assertExpectations(t *testing.T, db arangodb.Database, expect scenarioExpectations) {
	t.Helper()

	ctx := t.Context()
	assertNamed(t, "collection", expect.Collections, func(name string) (bool, error) {
		return db.CollectionExists(ctx, name)
	})
	assertNamed(t, "graph", expect.Graphs, func(name string) (bool, error) {
		return db.GraphExists(ctx, name)
	})
	assertNamed(t, "view", expect.Views, func(name string) (bool, error) {
		return db.ViewExists(ctx, name)
	})
	assertAnalyzers(t, db, expect.Analyzers)
	assertIndexes(t, db, expect.Indexes)
	assertDocuments(t, db, expect.Documents)
	assertMigrationRecords(t, db, expect.MigrationRecords)
}

func assertNamed(t *testing.T, kind string, items []namedExistence, exists func(string) (bool, error)) {
	t.Helper()

	for _, item := range items {
		got, err := exists(item.Name)
		if err != nil {
			t.Errorf("%s %s: %v", kind, item.Name, err)
			continue
		}

		if got != item.Exists {
			t.Errorf("%s %s exists = %t, want %t", kind, item.Name, got, item.Exists)
		}
	}
}

func assertAnalyzers(t *testing.T, db arangodb.Database, items []typedExistence) {
	t.Helper()

	ctx := t.Context()

	for _, item := range items {
		analyzer, err := db.Analyzer(ctx, item.Name)
		exists := err == nil

		if err != nil && !shared.IsNotFound(err) {
			t.Errorf("analyzer %s: %v", item.Name, err)
			continue
		}

		if exists != item.Exists {
			t.Errorf("analyzer %s exists = %t, want %t", item.Name, exists, item.Exists)
		}

		if exists && item.Type != "" && string(analyzer.Type()) != item.Type {
			t.Errorf("analyzer %s type = %s, want %s", item.Name, analyzer.Type(), item.Type)
		}
	}
}

func assertIndexes(t *testing.T, db arangodb.Database, items []indexExpectation) {
	t.Helper()

	ctx := t.Context()

	for _, item := range items {
		coll, err := db.GetCollection(ctx, item.Collection, nil)
		if err != nil {
			t.Errorf("collection %s: %v", item.Collection, err)
			continue
		}

		response, err := coll.Index(ctx, item.Name)
		exists := err == nil

		if err != nil && !shared.IsNotFound(err) {
			t.Errorf("index %s/%s: %v", item.Collection, item.Name, err)
			continue
		}

		if exists != item.Exists {
			t.Errorf("index %s/%s exists = %t, want %t", item.Collection, item.Name, exists, item.Exists)
		}

		if exists && item.Type != "" && string(response.Type) != item.Type {
			t.Errorf("index %s/%s type = %s, want %s", item.Collection, item.Name, response.Type, item.Type)
		}
	}
}

func assertDocuments(t *testing.T, db arangodb.Database, items []documentExpectation) {
	t.Helper()

	ctx := t.Context()

	for _, item := range items {
		docs, err := collectionDocuments(ctx, db, item.Collection)
		if err != nil {
			t.Errorf("documents in %s: %v", item.Collection, err)
			continue
		}

		if len(docs) != *item.Count {
			t.Errorf("collection %s document count = %d, want %d", item.Collection, len(docs), *item.Count)
		}

		for _, want := range item.Contains {
			if !documentsContain(docs, want) {
				t.Errorf("collection %s missing document fields %v", item.Collection, want)
			}
		}
	}
}

func collectionDocuments(ctx context.Context, db arangodb.Database, collection string) ([]map[string]any, error) {
	return queryMaps(ctx, db, "FOR doc IN @@collection RETURN doc", map[string]any{
		"@collection": collection,
	})
}

func assertMigrationRecords(t *testing.T, db arangodb.Database, items []migrationExpectation) {
	t.Helper()

	ctx := t.Context()

	for _, item := range items {
		records, err := queryMaps(ctx, db, "FOR item IN @@collection FILTER item.name == @name RETURN item", map[string]any{
			"@collection": DefaultMigrationCollection,
			"name":        item.Name,
		})
		if err != nil {
			t.Errorf("migration %s: %v", item.Name, err)
			continue
		}

		if item.Status == MigrationStatusMissing.String() {
			if len(records) != 0 {
				t.Errorf("migration %s status = present, want missing", item.Name)
			}

			continue
		}

		if len(records) != 1 {
			t.Errorf("migration %s records = %d, want 1", item.Name, len(records))
			continue
		}

		status, _ := records[0]["status"].(string)
		if status != item.Status {
			t.Errorf("migration %s status = %s, want %s", item.Name, status, item.Status)
		}
	}
}

func queryMaps(ctx context.Context, db arangodb.Database, query string, bindVars map[string]any) ([]map[string]any, error) {
	cursor, err := db.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, err
	}

	var docs []map[string]any

	for {
		var doc map[string]any

		_, err := cursor.ReadDocument(ctx, &doc)
		if shared.IsNoMoreDocuments(err) {
			break
		}

		if err != nil {
			return nil, errors.Join(err, cursor.CloseWithContext(ctx))
		}

		docs = append(docs, doc)
	}

	if err := cursor.CloseWithContext(ctx); err != nil {
		return nil, err
	}

	return docs, nil
}

func documentsContain(docs []map[string]any, want map[string]any) bool {
	for _, doc := range docs {
		if documentMatches(doc, want) {
			return true
		}
	}

	return false
}

func documentMatches(doc, want map[string]any) bool {
	for key, value := range want {
		if !jsonEqual(doc[key], value) {
			return false
		}
	}

	return true
}

func jsonEqual(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)

	if leftErr != nil || rightErr != nil {
		return false
	}

	return string(leftJSON) == string(rightJSON)
}
