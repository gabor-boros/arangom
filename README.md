# arangom

[![GoDoc](https://godoc.org/github.com/gabor-boros/arangom?status.svg)](https://godoc.org/github.com/gabor-boros/arangom)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabor-boros/arangom)](https://goreportcard.com/report/github.com/gabor-boros/arangom)
[![Maintainability](https://api.codeclimate.com/v1/badges/322e5839b4b4d1710351/maintainability)](https://codeclimate.com/github/gabor-boros/arangom/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/322e5839b4b4d1710351/test_coverage)](https://codeclimate.com/github/gabor-boros/arangom/test_coverage)

arangom is an [ArangoDB] migration tool written in Go. It is heavily inspired
by [ArangoMiGo] which is already a great tool.

[ArangoDB]: https://www.arangodb.com/
[ArangoMiGo]: https://github.com/deusdat/ArangoMiGo

## Installation

### As a CLI tool

**Using `brew`**

``` shell
$ brew tap gabor-boros/brew
$ brew install arangom
```

**Manual install**

To install `arangom`, use one of the [release artifacts]. If you have `go`
installed, you can build from source as well.

### As a package

```bash
go get github.com/gabor-boros/arangom
```

[release artifacts]: https://github.com/gabor-boros/arangom/releases

## Migrations

Migrations are a set of operations that are executed in a given order. The
operations are executed in the order they are defined in the migration.

When a migration is executed, depending on the status of the operations, the
migration is marked as `missing`, `running`, `done` or `failed`. If a migration
is marked as `missing` or `running`, that means that the migration is not
executed yet, therefore those migrations will be executed. If a migration is
marked as `done`, it is skipped. In the case of a `failed` migration, the
migration is skipped along with all the migrations that follow it and an error
is returned.

The lifecycle of a migration is the following:

1. The migration is marked as `missing`.
2. When executing the migration, the migration is marked as `running`.
3. When the migration is marked as `done` or `failed` based on the result of
   the operations.

Every migration has a checksum that is calculated based on a unique ID and the
operations of the migration. The checksum is used to determine if a migration
has changed since the last execution. If the checksum of a migration has changed
since the last execution, the migration is marked as `missing` and will be
executed again.

## Migration files

Migrations are stored in YAML files and are executed in alphabetical order. The
files are expected to be in the following schema:

```yaml
id: <numerical ID> # unique identifier of the migration, only used for checksum
operations: # list of operations to execute in the given order
  - kind: <operation kind> # the kind of the operation
    collection: <collection name> # the name of the collection to operate on
    options: <options> # the options of the operation
  # ...
```

Example:

```yaml
id: 1677564649
operations:
  - kind: createCollection
    collection: mycollection
    options:
      waitForSync: true
      cacheEnabled: true
      schema:
        level: moderate
        message: The document does not match the schema.
        rule:
          properties:
            kind:
              type: string
              maximum: 12
          additionalProperties:
            type: string
          required:
            - kind
```

The files are expected to be in a `migrations` directory in the current
working directory, but this can be changed with the `-migration-dir`
command line flag.

The migration directory can be structured in subdirectories as desired, arangom
will recursively search for migration files.

## Example usage

### As a package

```go
package main

import (
	"context"
	"os"

	arangoDriver "github.com/arangodb/go-driver"
	arangoHTTP "github.com/arangodb/go-driver/http"

	arangoMigrate "github.com/gabor-boros/arangom"
)

func main() {
	conn, _ := arangoHTTP.NewConnection(arangoHTTP.ConnectionConfig{
		Endpoints: []string{os.Getenv("ARANGO_URL")},
	})

	client, _ := arangoDriver.NewClient(arangoDriver.ClientConfig{
		Connection:     conn,
		Authentication: arangoDriver.BasicAuthentication(os.Getenv("ARANGO_USER"), os.Getenv("ARANGO_PASSWORD")),
	})

	db, _ := client.Database(context.Background(), os.Getenv("ARANGO_DB"))

	migrations := []*arangom.Migration{
		{
			Path:       "0001.initial", // any unique identifier
			Operations: []arangom.Operation{
				{
					Kind:       arangom.OperationKindCollectionCreate,
					Collection: os.Getenv("ARANGO_MIGRATION_COLLECTION"),
					Options:    map[string]any{
						"waitForSync": true,
					},
				},
			},
		},
	}

	executor, _ := arangom.NewExecutor(
		arangom.WithDatabase(db),
		arangom.WithCollection(os.Getenv("ARANGO_MIGRATION_COLLECTION")),
		arangom.WithMigrations(migrations),
	)

	if err := executor.Execute(context.Background()); err != nil {
		panic(err)
	}
}
```

### As a CLI tool

```bash
$ arangom -username "root" -password "openSesame" -database "mydb" -migration-dir "migrations"
2023/02/28 06:47:43 [INFO] connecting to the migration collection "migrations"
2023/02/28 06:47:43 [INFO] [1677564649] fetching migration status
2023/02/28 06:47:43 [INFO] [1677564649] executing migration
2023/02/28 06:47:43 [INFO] [1677564649] migration executed successfully
2023/02/28 06:47:43 [INFO] [1677564650] fetching migration status
2023/02/28 06:47:43 [INFO] [1677564650] executing migration
2023/02/28 06:47:43 [INFO] [1677564650] migration executed successfully
2023/02/28 06:47:43 [INFO] all migrations executed successfully
```

#### Flags

```bash
Usage of arangom:
  -collection string
        Migration collection (default "migrations")
  -create-collection
        Create migration collection if it does not exist (default true)
  -database string
        Database name
  -endpoints string
        Comma-separated list of database endpoints (default "http://localhost:8529")
  -migration-dir string
        Migration directory (default "migrations")
  -password string
        Database password
  -username string
        Database user
  -version
        Print version and exit
```

## Supported operations

The options of the operations are the same as the request body options of the
corresponding ArangoDB operations. However, some operations have additional
options that are not supported by ArangoDB. These options are documented in the
table below with the caveat column. For more information about operation
options, please refer to the [ArangoDB documentation].

| Operation kind          | Description                               | Caveats                 |
|-------------------------|-------------------------------------------|-------------------------|
| `executeAQL`            | Executes an AQL query.                    | [executeAQL]            |
| `createCollection`      | Creates a new collection.                 | -                       |
| `updateCollection`      | Updates an existing collection.           | -                       |
| `deleteCollection`      | Deletes an existing collection.           | -                       |
| `createGraph`           | Creates a new graph.                      | [createGraph]           |
| `addVertexToGraph`      | Adds a vertex collection to a graph.      | [addVertexToGraph]      |
| `removeVertexFromGraph` | Removes a vertex collection from a graph. | [removeVertexFromGraph] |
| `addEdgeToGraph`        | Adds an edge definition to a graph.       | [addEdgeToGraph]        |
| `removeEdgeFromGraph`   | Removes an edge definition from a graph.  | [removeEdgeFromGraph]   |
| `deleteGraph`           | Deletes an existing graph.                | [deleteGraph]           |
| `createView`            | Creates a new search view.                | [createView]            |
| `updateView`            | Updates an existing search view.          | [updateView]            |
| `deleteView`            | Deletes an existing search view.          | [deleteView]            |
| `createFulltextIndex`   | Creates a fulltext index.                 | [createFulltextIndex]   |
| `createGeoSpatialIndex` | Creates a Geo-spatial index.              | [createGeoSpatialIndex] |
| `createHashIndex`       | Creates a hash index.                     | [createHashIndex]       |
| `createInvertedIndex`   | Creates an inverted index.                | -                       |
| `createPersistentIndex` | Creates a persistent index.               | [createPersistentIndex] |
| `createSkipListIndex`   | Creates a skiplist index.                 | [createSkipListIndex]   |
| `createTTLIndex`        | Creates a TTL index.                      | [createTTLIndex]        |
| `createZKDIndex`        | Creates a ZKD index.                      | [createZKDIndex]        |
| `deleteIndex`           | Deletes an index.                         | [deleteIndex]           |
| `createAnalyzer`        | Creates an analyzer.                      | [createAnalyzer]        |
| `deleteAnalyzer`        | Deletes an analyzer.                      | [deleteIndex]           |

Database operations are intentionally not supported. To apply migrations on a
database, create the database first. **Feature requests and/or pull requests
are welcomed for other missing operations.**

[ArangoDB documentation]: https://www.arangodb.com/docs/stable/http/
[executeAQL]: #executeaql-options
[createGraph]: #creategraph-options
[addVertexToGraph]: #addvertextograph-options
[removeVertexFromGraph]: #removevertexfromgraph-options
[addEdgeToGraph]: #addedgetograph-options
[removeEdgeFromGraph]: #removeedgefromgraph-options
[deleteGraph]: #deletegraph-options
[createView]: #createview-options
[updateView]: #updateview-options
[deleteView]: #deleteview-options
[createFulltextIndex]: #createfulltextindex-options
[createGeoSpatialIndex]: #creategeospatialindex-options
[createHashIndex]: #createhashindex-options
[createPersistentIndex]: #createpersistentindex-options
[createSkipListIndex]: #createskiplistindex-options
[createTTLIndex]: #createttlindex-options
[createZKDIndex]: #createzkdindex-options
[deleteIndex]: #deleteindex-options
[createAnalyzer]: #createanalyzer-options
[deleteAnalyzer]: #deleteanalyzer-options

### Operation option caveats

#### `executeAQL` options

Executing AQL queries has no options defined by ArangoDB, though it is possible
to bind variables to the query. The `executeAQL` operation therefore has the
following options:

| Option name | Description                      |
|-------------|----------------------------------|
| `query`     | The AQL query to execute.        |
| `bindVars`  | The bind variables of the query. |

#### `createGraph` options

The graph name is defined in the `collection` field of the operation.

#### `addVertexToGraph` options

The graph name is defined in the `collection` field of the operation. The
options of the operation is extended with the following options:

| Option name  | Description                                            |
|--------------|--------------------------------------------------------|
| `collection` | The name of the vertex collection to add to the graph. |

#### `removeVertexFromGraph` options

The graph name is defined in the `collection` field of the operation. The
options of the operation is extended with the following options:

| Option name  | Description                                                 |
|--------------|-------------------------------------------------------------|
| `collection` | The name of the vertex collection to remove from the graph. |

#### `addEdgeToGraph` options

The graph name is defined in the `collection` field of the operation. The
options of the operation is extended with the following options:

| Option name   | Description                                          |
|---------------|------------------------------------------------------|
| `collection`  | The name of the edge definition to add to the graph. |
| `constraints` | The edge constraints of the edge definition.         |

The edge `constraints` are defined as follows:

| Option name | Description                                                                      |
|-------------|----------------------------------------------------------------------------------|
| `from`      | The list of name of the vertex collections that the edge definition starts from. |
| `to`        | The list of name of the vertex collections that the edge definition ends to.     |

#### `removeEdgeFromGraph` options

The graph name is defined in the `collection` field of the operation. The
options of the operation is extended with the following options:

| Option name  | Description                                          |
|--------------|------------------------------------------------------|
| `collection` | The name of the edge definition to add to the graph. |

#### `deleteGraph` options

The graph name is defined in the `collection` field of the operation.

#### `createView` options

The view name is defined in the `collection` field of the operation.

#### `updateView` options

The view name is defined in the `collection` field of the operation.

#### `deleteView` options

#### `createFulltextIndex` options

The options of the operation is extended with the following options:

| Option name | Description                  |
|-------------|------------------------------|
| `fields`    | The list of fields to index. |

#### `createGeoSpatialIndex` options

The options of the operation is extended with the following options:

| Option name | Description                  |
|-------------|------------------------------|
| `fields`    | The list of fields to index. |

#### `createHashIndex` options

The options of the operation is extended with the following options:

| Option name | Description                  |
|-------------|------------------------------|
| `fields`    | The list of fields to index. |

#### `createPersistentIndex` options

The options of the operation is extended with the following options:

| Option name | Description                  |
|-------------|------------------------------|
| `fields`    | The list of fields to index. |

#### `createSkipListIndex` options

The options of the operation is extended with the following options:

| Option name | Description                  |
|-------------|------------------------------|
| `fields`    | The list of fields to index. |

#### `createTTLIndex` options

The options of the operation is extended with the following options:

| Option name   | Description                     |
|---------------|---------------------------------|
| `field`       | The field to index.             |
| `expireAfter` | The expiration time in seconds. |

#### `createZKDIndex` options

The options of the operation is extended with the following options:

| Option name | Description                  |
|-------------|------------------------------|
| `fields`    | The list of fields to index. |

#### `deleteIndex` options

The options of the operation is extended with the following options:

| Option name | Description                      |
|-------------|----------------------------------|
| `name`      | The name of the index to delete. |

#### `createAnalyzer` options

In case analyzer with the same name already exists, the operation will fail.

| Option name  | Description                                      |
|--------------|--------------------------------------------------|
| `name`       | The name of the analyzer to create.              |
| `type`       | The type of the analyzer to create.              |
| `properties` | The properties of the analyzer to create.        |
| `features`   | The array of features of the analyzer to create. |

#### `deleteAnalyzer` options

| Option name  | Description                         |
|--------------|-------------------------------------|
| `name`       | The name of the analyzer to delete. |
| `force`      | Delete even if in use.              |

## Compatibility

The compatibility of arangom is equal to the compatibility of the ArangoDB
Go driver. For more information, please refer to the [ArangoDB Go driver
documentation].

[ArangoDB Go driver documentation]: https://github.com/arangodb/go-driver#supported-versions

## Q&A

### Why not contributing to ArangoMiGo instead?

ArangoMiGo is a great tool, but it is highly opinionated and has design
decisions that are not flexible enough. Refactoring ArangoMiGo to be more
flexible would be a huge effort, almost a rewrite. Therefore, I created a new
project instead.

### Why using `collection` for the name of the graph and view?

The naming is indeed confusing, and I am open to suggestions.

## Contributing

Contributions are welcomed! Please open an issue or a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
