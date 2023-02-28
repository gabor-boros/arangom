package arangom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	OperationKindAQLExecute            OperationKind = iota + 1 // operation to execute an AQL query
	OperationKindCollectionCreate                               // operation to create a collection
	OperationKindCollectionUpdate                               // operation to update a collection
	OperationKindCollectionDelete                               // operation to delete a collection
	OperationKindGraphCreate                                    // operation to create a graph
	OperationKindGraphAddVertex                                 // operation to add a vertex collection to a graph
	OperationKindGraphRemoveVertex                              // operation to remove a vertex collection from a graph
	OperationKindGraphAddEdge                                   // operation to add an edge definition to a graph
	OperationKindGraphRemoveEdge                                // operation to remove an edge definition from a graph
	OperationKindGraphDelete                                    // operation to delete a graph
	OperationKindViewCreate                                     // operation to create a view
	OperationKindViewUpdate                                     // operation to update a view
	OperationKindViewDelete                                     // operation to delete a view
	OperationKindFulltextIndexCreate                            // operation to create a fulltext index
	OperationKindGeoSpatialIndexCreate                          // operation to create a geospatial index
	OperationKindHashIndexCreate                                // operation to create a hash index
	OperationKindInvertedIndexCreate                            // operation to create an inverted index
	OperationKindPersistentIndexCreate                          // operation to create a persistent index
	OperationKindSkipListIndexCreate                            // operation to create a skiplist index
	OperationKindTTLIndexCreate                                 // operation to create a TTL index
	OperationKindZKDIndexCreate                                 // operation to create a ZKD index
	OperationKindIndexDelete                                    // operation to delete an index
)

var (
	// ErrInvalidOperationKind is returned when an invalid operation kind is
	// specified.
	ErrInvalidOperationKind = fmt.Errorf("invalid operation kind")

	// operationMap is a map of operation names to operation kinds.
	operationMap = map[string]OperationKind{
		"executeAQL":            OperationKindAQLExecute,
		"createCollection":      OperationKindCollectionCreate,
		"updateCollection":      OperationKindCollectionUpdate,
		"deleteCollection":      OperationKindCollectionDelete,
		"createGraph":           OperationKindGraphCreate,
		"addVertexToGraph":      OperationKindGraphAddVertex,
		"removeVertexFromGraph": OperationKindGraphRemoveVertex,
		"addEdgeToGraph":        OperationKindGraphAddEdge,
		"removeEdgeFromGraph":   OperationKindGraphRemoveEdge,
		"deleteGraph":           OperationKindGraphDelete,
		"createView":            OperationKindViewCreate,
		"updateView":            OperationKindViewUpdate,
		"deleteView":            OperationKindViewDelete,
		"createFulltextIndex":   OperationKindFulltextIndexCreate,
		"createGeoSpatialIndex": OperationKindGeoSpatialIndexCreate,
		"createHashIndex":       OperationKindHashIndexCreate,
		"createInvertedIndex":   OperationKindInvertedIndexCreate,
		"createPersistentIndex": OperationKindPersistentIndexCreate,
		"createSkipListIndex":   OperationKindSkipListIndexCreate,
		"createTTLIndex":        OperationKindTTLIndexCreate,
		"createZKDIndex":        OperationKindZKDIndexCreate,
		"deleteIndex":           OperationKindIndexDelete,
	}

	// operationKindMap is a map of operation kinds to operations.
	operationKindMap = map[OperationKind]func(o *Operation) OperationFn{
		OperationKindAQLExecute:            ExecuteAQLOperation,
		OperationKindCollectionCreate:      CreateCollectionOperation,
		OperationKindCollectionUpdate:      UpdateCollectionOperation,
		OperationKindCollectionDelete:      DeleteCollectionOperation,
		OperationKindGraphCreate:           CreateGraphOperation,
		OperationKindGraphAddVertex:        AddVertexOperation,
		OperationKindGraphRemoveVertex:     RemoveVertexOperation,
		OperationKindGraphAddEdge:          AddEdgeOperation,
		OperationKindGraphRemoveEdge:       RemoveEdgeOperation,
		OperationKindGraphDelete:           DeleteGraphOperation,
		OperationKindViewCreate:            CreateViewOperation,
		OperationKindViewUpdate:            UpdateViewOperation,
		OperationKindViewDelete:            DeleteViewOperation,
		OperationKindFulltextIndexCreate:   CreateFulltextIndexOperation,
		OperationKindGeoSpatialIndexCreate: CreateGeoSpatialIndexOperation,
		OperationKindHashIndexCreate:       CreateHashIndexOperation,
		OperationKindInvertedIndexCreate:   CreateInvertedIndexOperation,
		OperationKindPersistentIndexCreate: CreatePersistentIndexOperation,
		OperationKindSkipListIndexCreate:   CreateSkipListIndexOperation,
		OperationKindTTLIndexCreate:        CreateTTLIndexOperation,
		OperationKindZKDIndexCreate:        CreateZKDIndexOperation,
		OperationKindIndexDelete:           DeleteIndexOperation,
	}
)

// OperationFn runs an operation on a database.
type OperationFn func(ctx context.Context, db driver.Database) error

// OperationKind is the kind of operation to run.
type OperationKind int

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (o *OperationKind) UnmarshalYAML(value *yaml.Node) error {
	if kind, ok := operationMap[value.Value]; ok {
		*o = kind
		return nil
	}

	return errors.Wrap(ErrInvalidOperationKind, value.Value)
}

// Operation is an operation to run in a migration.
type Operation struct {
	Kind       OperationKind  `yaml:"kind"`
	Collection string         `yaml:"collection"`
	Options    map[string]any `yaml:"options"`
}

// GetOperationFn returns the operation function for the operation kind.
func (o *Operation) GetOperationFn() (OperationFn, error) {
	if op, ok := operationKindMap[o.Kind]; ok && op != nil {
		return op(o), nil
	}

	return nil, ErrInvalidOperationKind
}

// convertToOperationOptions converts a map of options to a struct using JSON.
// This is used to convert the options from the YAML file to the options for
// the ArangoDB driver. Since the ArangoDB driver uses a JSON tag for the
// options, we can use JSON to convert the options from the YAML file to the
// options for the ArangoDB driver.
func convertToOperationOptions(opts map[string]any, dst any) error {
	b, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dst)
}

// ExecuteAQLOperation executes an AQL query.
func ExecuteAQLOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type aqlOpts struct {
			Query    string         `json:"query"`
			BindVars map[string]any `json:"bindVars"`
		}

		opts := aqlOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.Query(ctx, opts.Query, opts.BindVars)

		return err
	}
}

// CreateCollectionOperation creates a collection.
func CreateCollectionOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		opts := driver.CreateCollectionOptions{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.CreateCollection(ctx, o.Collection, &opts)
		return err
	}
}

// UpdateCollectionOperation updates a collection.
func UpdateCollectionOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		opts := driver.SetCollectionPropertiesOptions{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		return coll.SetProperties(ctx, opts)
	}
}

// DeleteCollectionOperation removes a collection.
func DeleteCollectionOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		return coll.Remove(ctx)
	}
}

// CreateGraphOperation creates a graph.
func CreateGraphOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		opts := driver.CreateGraphOptions{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.CreateGraphV2(ctx, o.Collection, &opts)
		return err
	}
}

// AddVertexOperation adds a vertex collection to a graph.
func AddVertexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type addVertexOpts struct {
			driver.CreateVertexCollectionOptions
			Collection string `json:"collection"`
		}

		opts := addVertexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, err = graph.CreateVertexCollectionWithOptions(ctx, opts.Collection, opts.CreateVertexCollectionOptions)

		return err
	}
}

// RemoveVertexOperation removes a vertex collection from a graph.
func RemoveVertexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type removeVertexOpts struct {
			Collection string `json:"collection"`
		}

		opts := removeVertexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection)
		if err != nil {
			return err
		}

		vertex, err := graph.VertexCollection(ctx, opts.Collection)
		if err != nil {
			return err
		}

		return vertex.Remove(ctx)
	}
}

// AddEdgeOperation adds an edge definition to a graph.
func AddEdgeOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type addEdgeOpts struct {
			driver.CreateEdgeCollectionOptions
			Collection  string                   `json:"collection"`
			Constraints driver.VertexConstraints `json:"constraints"`
		}

		opts := addEdgeOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, err = graph.CreateEdgeCollectionWithOptions(ctx, opts.Collection, opts.Constraints, opts.CreateEdgeCollectionOptions)

		return err
	}
}

// RemoveEdgeOperation removes an edge definition from a graph.
func RemoveEdgeOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type removeEdgeOpts struct {
			Collection string `json:"collection"`
		}

		opts := removeEdgeOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection)
		if err != nil {
			return err
		}

		edge, _, err := graph.EdgeCollection(ctx, opts.Collection)
		if err != nil {
			return err
		}

		return edge.Remove(ctx)
	}
}

// DeleteGraphOperation deletes a graph.
func DeleteGraphOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		graph, err := db.Graph(ctx, o.Collection)
		if err != nil {
			return err
		}

		return graph.Remove(ctx)
	}
}

// CreateViewOperation creates a search view.
func CreateViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		opts := driver.ArangoSearchViewProperties{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.CreateArangoSearchView(ctx, o.Collection, &opts)
		return err
	}
}

// UpdateViewOperation updates an existing search view.
func UpdateViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		opts := driver.ArangoSearchViewProperties{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		view, err := db.View(ctx, o.Collection)
		if err != nil {
			return err
		}

		searchView, err := view.ArangoSearchView()
		if err != nil {
			return err
		}

		return searchView.SetProperties(ctx, opts)
	}
}

// DeleteViewOperation removes a search view.
func DeleteViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		view, err := db.View(ctx, o.Collection)
		if err != nil {
			return err
		}

		return view.Remove(ctx)
	}
}

func CreateFulltextIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type fulltextIndexOpts struct {
			driver.EnsureFullTextIndexOptions
			Fields []string `json:"fields"`
		}

		opts := fulltextIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureFullTextIndex(ctx, opts.Fields, &opts.EnsureFullTextIndexOptions) //nolint:staticcheck

		return err
	}
}

func CreateGeoSpatialIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type geoSpatialIndexOpts struct {
			driver.EnsureGeoIndexOptions
			Fields []string `json:"fields"`
		}

		opts := geoSpatialIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureGeoIndex(ctx, opts.Fields, &opts.EnsureGeoIndexOptions)

		return err
	}
}

func CreateHashIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type hashIndexOpts struct {
			driver.EnsureHashIndexOptions
			Fields []string `json:"fields"`
		}

		opts := hashIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureHashIndex(ctx, opts.Fields, &opts.EnsureHashIndexOptions)

		return err
	}
}

func CreateInvertedIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		opts := driver.InvertedIndexOptions{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureInvertedIndex(ctx, &opts)

		return err
	}
}

func CreatePersistentIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type persistentIndexOpts struct {
			driver.EnsurePersistentIndexOptions
			Fields []string `json:"fields"`
		}

		opts := persistentIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsurePersistentIndex(ctx, opts.Fields, &opts.EnsurePersistentIndexOptions)

		return err
	}
}

func CreateSkipListIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type skipListIndexOpts struct {
			driver.EnsureSkipListIndexOptions
			Fields []string `json:"fields"`
		}

		opts := skipListIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureSkipListIndex(ctx, opts.Fields, &opts.EnsureSkipListIndexOptions)

		return err
	}
}

func CreateTTLIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type ttlIndexOpts struct {
			driver.EnsureTTLIndexOptions
			Field       string `json:"field"`
			ExpireAfter int    `json:"expireAfter"`
		}

		opts := ttlIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureTTLIndex(ctx, opts.Field, opts.ExpireAfter, &opts.EnsureTTLIndexOptions)

		return err
	}
}

func CreateZKDIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type zkdIndexOpts struct {
			driver.EnsureZKDIndexOptions
			Fields []string `json:"fields"`
		}

		opts := zkdIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureZKDIndex(ctx, opts.Fields, &opts.EnsureZKDIndexOptions)

		return err
	}
}

func DeleteIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db driver.Database) error {
		type deleteIndexOpts struct {
			Name string `json:"name"`
		}

		opts := deleteIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.Collection(ctx, o.Collection)
		if err != nil {
			return err
		}

		index, err := coll.Index(ctx, opts.Name)
		if err != nil {
			return err
		}

		return index.Remove(ctx)
	}
}
