package arangom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	OperationKindAQLExecute             OperationKind = 1  // executeAQL
	OperationKindCollectionCreate       OperationKind = 2  // createCollection
	OperationKindCollectionUpdate       OperationKind = 3  // updateCollection
	OperationKindCollectionDelete       OperationKind = 4  // deleteCollection
	OperationKindGraphCreate            OperationKind = 5  // createGraph
	OperationKindGraphAddVertex         OperationKind = 6  // addVertexToGraph
	OperationKindGraphRemoveVertex      OperationKind = 7  // removeVertexFromGraph
	OperationKindGraphAddEdge           OperationKind = 8  // addEdgeToGraph
	OperationKindGraphRemoveEdge        OperationKind = 9  // removeEdgeFromGraph
	OperationKindGraphDelete            OperationKind = 10 // deleteGraph
	OperationKindViewCreate             OperationKind = 11 // createView
	OperationKindViewUpdate             OperationKind = 12 // updateView
	OperationKindViewDelete             OperationKind = 13 // deleteView
	OperationKindGeoSpatialIndexCreate  OperationKind = 15 // createGeoSpatialIndex
	OperationKindInvertedIndexCreate    OperationKind = 17 // createInvertedIndex
	OperationKindPersistentIndexCreate  OperationKind = 18 // createPersistentIndex
	OperationKindTTLIndexCreate         OperationKind = 20 // createTTLIndex
	OperationKindMDIIndexCreate         OperationKind = 21 // createMDIIndex
	OperationKindIndexDelete            OperationKind = 22 // deleteIndex
	OperationKindAnalyzerCreate         OperationKind = 23 // createAnalyzer
	OperationKindAnalyzerDelete         OperationKind = 24 // deleteAnalyzer
	OperationKindMDIPrefixedIndexCreate OperationKind = 25 // createMDIPrefixedIndex
	OperationKindVectorIndexCreate      OperationKind = 26 // createVectorIndex
	OperationKindSearchAliasViewCreate  OperationKind = 27 // createSearchAliasView
	OperationKindSearchAliasViewUpdate  OperationKind = 28 // updateSearchAliasView
)

var (
	// ErrInvalidOperationKind is returned when an invalid operation kind is
	// specified.
	ErrInvalidOperationKind = fmt.Errorf("invalid operation kind")

	// operationKindMap is a map of operation kinds to operations.
	operationKindMap = map[OperationKind]func(o *Operation) OperationFn{
		OperationKindAQLExecute:             ExecuteAQLOperation,
		OperationKindCollectionCreate:       CreateCollectionOperation,
		OperationKindCollectionUpdate:       UpdateCollectionOperation,
		OperationKindCollectionDelete:       DeleteCollectionOperation,
		OperationKindGraphCreate:            CreateGraphOperation,
		OperationKindGraphAddVertex:         AddVertexOperation,
		OperationKindGraphRemoveVertex:      RemoveVertexOperation,
		OperationKindGraphAddEdge:           AddEdgeOperation,
		OperationKindGraphRemoveEdge:        RemoveEdgeOperation,
		OperationKindGraphDelete:            DeleteGraphOperation,
		OperationKindViewCreate:             CreateViewOperation,
		OperationKindViewUpdate:             UpdateViewOperation,
		OperationKindViewDelete:             DeleteViewOperation,
		OperationKindGeoSpatialIndexCreate:  CreateGeoSpatialIndexOperation,
		OperationKindInvertedIndexCreate:    CreateInvertedIndexOperation,
		OperationKindPersistentIndexCreate:  CreatePersistentIndexOperation,
		OperationKindTTLIndexCreate:         CreateTTLIndexOperation,
		OperationKindMDIIndexCreate:         CreateMDIIndexOperation,
		OperationKindIndexDelete:            DeleteIndexOperation,
		OperationKindAnalyzerCreate:         CreateAnalyzerOperation,
		OperationKindAnalyzerDelete:         DeleteAnalyzerOperation,
		OperationKindMDIPrefixedIndexCreate: CreateMDIPrefixedIndexOperation,
		OperationKindVectorIndexCreate:      CreateVectorIndexOperation,
		OperationKindSearchAliasViewCreate:  CreateSearchAliasViewOperation,
		OperationKindSearchAliasViewUpdate:  UpdateSearchAliasViewOperation,
	}
)

// OperationFn runs an operation on a database.
type OperationFn func(ctx context.Context, db arangodb.Database) error

// OperationKind is the kind of operation to run.
//
//go:generate go tool enumer -type=OperationKind -linecomment -output=operationkind_enumer.go
type OperationKind int

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (o *OperationKind) UnmarshalYAML(value *yaml.Node) error {
	kind, err := OperationKindString(value.Value)
	if err != nil {
		return errors.Wrap(ErrInvalidOperationKind, value.Value)
	}

	*o = kind
	return nil
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
	return func(ctx context.Context, db arangodb.Database) error {
		type aqlOpts struct {
			Query    string         `json:"query"`
			BindVars map[string]any `json:"bindVars"`
		}

		opts := aqlOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		cursor, err := db.Query(ctx, opts.Query, &arangodb.QueryOptions{BindVars: opts.BindVars})
		if err != nil {
			return err
		}
		if cursor == nil {
			return nil
		}

		return cursor.CloseWithContext(ctx)
	}
}

// CreateCollectionOperation creates a collection.
func CreateCollectionOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.CreateCollectionPropertiesV2{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.CreateCollectionV2(ctx, o.Collection, &opts)
		return err
	}
}

// UpdateCollectionOperation updates a collection's properties.
//
// Computed values are decoded into SetCollectionPropertiesOptionsV2. An empty
// list is kept, so a migration can remove computed values.
func UpdateCollectionOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.SetCollectionPropertiesOptionsV2{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		return coll.SetPropertiesV2(ctx, opts)
	}
}

// DeleteCollectionOperation removes a collection.
func DeleteCollectionOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		return coll.Remove(ctx)
	}
}

// CreateGraphOperation creates a graph.
func CreateGraphOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type graphCreateOptions struct {
			arangodb.GraphDefinition
			Satellites []string `json:"satellites,omitempty"`
		}

		opts := graphCreateOptions{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		var createOpts *arangodb.CreateGraphOptions
		if len(opts.Satellites) > 0 {
			createOpts = &arangodb.CreateGraphOptions{Satellites: opts.Satellites}
		}

		_, err := db.CreateGraph(ctx, o.Collection, &opts.GraphDefinition, createOpts)

		return err
	}
}

// AddVertexOperation adds a vertex collection to a graph.
func AddVertexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type addVertexOpts struct {
			arangodb.CreateVertexCollectionOptions
			Collection string `json:"collection"`
		}

		opts := addVertexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, err = graph.CreateVertexCollection(ctx, opts.Collection, &opts.CreateVertexCollectionOptions)

		return err
	}
}

// RemoveVertexOperation removes a vertex collection from a graph.
func RemoveVertexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type removeVertexOpts struct {
			Collection string `json:"collection"`
		}

		opts := removeVertexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, err = graph.DeleteVertexCollection(ctx, opts.Collection, nil)

		return err
	}
}

// AddEdgeOperation adds an edge definition to a graph.
func AddEdgeOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type addEdgeOpts struct {
			arangodb.CreateEdgeDefinitionOptions
			Collection  string `json:"collection"`
			Constraints struct {
				From []string `json:"from"`
				To   []string `json:"to"`
			} `json:"constraints"`
		}

		opts := addEdgeOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, err = graph.CreateEdgeDefinition(ctx, opts.Collection, opts.Constraints.From, opts.Constraints.To, &opts.CreateEdgeDefinitionOptions)

		return err
	}
}

// RemoveEdgeOperation removes an edge definition from a graph.
func RemoveEdgeOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type removeEdgeOpts struct {
			Collection string `json:"collection"`
		}

		opts := removeEdgeOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		graph, err := db.Graph(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, err = graph.DeleteEdgeDefinition(ctx, opts.Collection, nil)

		return err
	}
}

// DeleteGraphOperation deletes a graph.
func DeleteGraphOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		graph, err := db.Graph(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		return graph.Remove(ctx, nil)
	}
}

// CreateViewOperation creates a search view.
func CreateViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.ArangoSearchViewProperties{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.CreateArangoSearchView(ctx, o.Collection, &opts)
		return err
	}
}

// UpdateViewOperation updates an existing search view.
func UpdateViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.ArangoSearchViewProperties{}
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
	return func(ctx context.Context, db arangodb.Database) error {
		view, err := db.View(ctx, o.Collection)
		if err != nil {
			return err
		}

		return view.Remove(ctx)
	}
}

func CreateGeoSpatialIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type geoSpatialIndexOpts struct {
			arangodb.CreateGeoIndexOptions
			Fields []string `json:"fields"`
		}

		opts := geoSpatialIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureGeoIndex(ctx, opts.Fields, &opts.CreateGeoIndexOptions)

		return err
	}
}

func CreateInvertedIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.InvertedIndexOptions{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureInvertedIndex(ctx, &opts)

		return err
	}
}

func CreatePersistentIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type persistentIndexOpts struct {
			arangodb.CreatePersistentIndexOptions
			Fields []string `json:"fields"`
		}

		opts := persistentIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsurePersistentIndex(ctx, opts.Fields, &opts.CreatePersistentIndexOptions)

		return err
	}
}

func CreateTTLIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type ttlIndexOpts struct {
			arangodb.CreateTTLIndexOptions
			Field       string `json:"field"`
			ExpireAfter int    `json:"expireAfter"`
		}

		opts := ttlIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureTTLIndex(ctx, []string{opts.Field}, opts.ExpireAfter, &opts.CreateTTLIndexOptions)

		return err
	}
}

// defaultMDIFieldType fills the only field value type the driver accepts.
func defaultMDIFieldType(fieldType arangodb.MDIFieldType) arangodb.MDIFieldType {
	if fieldType == "" {
		return arangodb.MDIDoubleFieldType
	}

	return fieldType
}

func CreateMDIIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type mdiIndexOpts struct {
			arangodb.CreateMDIIndexOptions
			Fields []string `json:"fields"`
		}

		opts := mdiIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}
		opts.FieldValueTypes = defaultMDIFieldType(opts.FieldValueTypes)

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureMDIIndex(ctx, opts.Fields, &opts.CreateMDIIndexOptions)

		return err
	}
}

func CreateMDIPrefixedIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type mdiPrefixedIndexOpts struct {
			arangodb.CreateMDIPrefixedIndexOptions
			Fields []string `json:"fields"`
		}

		opts := mdiPrefixedIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}
		opts.FieldValueTypes = defaultMDIFieldType(opts.FieldValueTypes)

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureMDIPrefixedIndex(ctx, opts.Fields, &opts.CreateMDIPrefixedIndexOptions)

		return err
	}
}

func CreateVectorIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type vectorIndexOpts struct {
			arangodb.CreateVectorIndexOptions
			Fields []string              `json:"fields"`
			Params arangodb.VectorParams `json:"params"`
		}

		opts := vectorIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		_, _, err = coll.EnsureVectorIndex(ctx, opts.Fields, &opts.Params, &opts.CreateVectorIndexOptions)

		return err
	}
}

func DeleteIndexOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type deleteIndexOpts struct {
			Name string `json:"name"`
		}

		opts := deleteIndexOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		coll, err := db.GetCollection(ctx, o.Collection, nil)
		if err != nil {
			return err
		}

		return coll.DeleteIndex(ctx, opts.Name)
	}
}

func CreateAnalyzerOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		analyzer := arangodb.AnalyzerDefinition{}
		if err := convertToOperationOptions(o.Options, &analyzer); err != nil {
			return err
		}

		_, _, err := db.EnsureCreatedAnalyzer(ctx, &analyzer)

		return err
	}
}

func CreateSearchAliasViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.ArangoSearchAliasViewProperties{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		_, err := db.CreateArangoSearchAliasView(ctx, o.Collection, &opts)

		return err
	}
}

func UpdateSearchAliasViewOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		opts := arangodb.ArangoSearchAliasViewProperties{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		view, err := db.View(ctx, o.Collection)
		if err != nil {
			return err
		}

		alias, err := view.ArangoSearchViewAlias()
		if err != nil {
			return err
		}

		return alias.SetProperties(ctx, opts)
	}
}

func DeleteAnalyzerOperation(o *Operation) OperationFn {
	return func(ctx context.Context, db arangodb.Database) error {
		type deleteAnalyzerOpts struct {
			Name  string `json:"name"`
			Force bool   `json:"force"`
		}

		opts := deleteAnalyzerOpts{}
		if err := convertToOperationOptions(o.Options, &opts); err != nil {
			return err
		}

		analyzer, err := db.Analyzer(ctx, opts.Name)
		if err != nil {
			return err
		}

		return analyzer.Remove(ctx, opts.Force)
	}
}
