package arangom

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/arangodb/go-driver"
	"gopkg.in/yaml.v3"
)

func TestOperationKind_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		value   []byte
		want    OperationKind
		wantErr bool
	}{
		{
			name:  "unmarshal executeAQL",
			value: []byte(`executeAQL`),
			want:  OperationKindAQLExecute,
		},
		{
			name:  "unmarshal createCollection",
			value: []byte(`createCollection`),
			want:  OperationKindCollectionCreate,
		},
		{
			name:  "unmarshal updateCollection",
			value: []byte(`updateCollection`),
			want:  OperationKindCollectionUpdate,
		},
		{
			name:  "unmarshal deleteCollection",
			value: []byte(`deleteCollection`),
			want:  OperationKindCollectionDelete,
		},
		{
			name:  "unmarshal createGraph",
			value: []byte(`createGraph`),
			want:  OperationKindGraphCreate,
		},
		{
			name:  "unmarshal addVertexToGraph",
			value: []byte(`addVertexToGraph`),
			want:  OperationKindGraphAddVertex,
		},
		{
			name:  "unmarshal removeVertexFromGraph",
			value: []byte(`removeVertexFromGraph`),
			want:  OperationKindGraphRemoveVertex,
		},
		{
			name:  "unmarshal addEdgeToGraph",
			value: []byte(`addEdgeToGraph`),
			want:  OperationKindGraphAddEdge,
		},
		{
			name:  "unmarshal removeEdgeFromGraph",
			value: []byte(`removeEdgeFromGraph`),
			want:  OperationKindGraphRemoveEdge,
		},
		{
			name:  "unmarshal deleteGraph",
			value: []byte(`deleteGraph`),
			want:  OperationKindGraphDelete,
		},
		{
			name:  "unmarshal createView",
			value: []byte(`createView`),
			want:  OperationKindViewCreate,
		},
		{
			name:  "unmarshal updateView",
			value: []byte(`updateView`),
			want:  OperationKindViewUpdate,
		},
		{
			name:  "unmarshal deleteView",
			value: []byte(`deleteView`),
			want:  OperationKindViewDelete,
		},
		{
			name:  "unmarshal createFulltextIndex",
			value: []byte(`createFulltextIndex`),
			want:  OperationKindFulltextIndexCreate,
		},
		{
			name:  "unmarshal createGeoSpatialIndex",
			value: []byte(`createGeoSpatialIndex`),
			want:  OperationKindGeoSpatialIndexCreate,
		},
		{
			name:  "unmarshal createHashIndex",
			value: []byte(`createHashIndex`),
			want:  OperationKindHashIndexCreate,
		},
		{
			name:  "unmarshal createInvertedIndex",
			value: []byte(`createInvertedIndex`),
			want:  OperationKindInvertedIndexCreate,
		},
		{
			name:  "unmarshal createPersistentIndex",
			value: []byte(`createPersistentIndex`),
			want:  OperationKindPersistentIndexCreate,
		},
		{
			name:  "unmarshal createSkipListIndex",
			value: []byte(`createSkipListIndex`),
			want:  OperationKindSkipListIndexCreate,
		},
		{
			name:  "unmarshal createTTLIndex",
			value: []byte(`createTTLIndex`),
			want:  OperationKindTTLIndexCreate,
		},
		{
			name:  "unmarshal createZKDIndex",
			value: []byte(`createZKDIndex`),
			want:  OperationKindZKDIndexCreate,
		},
		{
			name:  "unmarshal deleteIndex",
			value: []byte(`deleteIndex`),
			want:  OperationKindIndexDelete,
		},
		{
			name:  "unmarshal createAnalyzer",
			value: []byte(`createAnalyzer`),
			want:  OperationKindAnalyzerCreate,
		},
		{
			name:  "unmarshal deleteAnalyzer",
			value: []byte(`deleteAnalyzer`),
			want:  OperationKindAnalyzerDelete,
		},
		{
			name:    "unmarshal invalid operation kind",
			value:   []byte(`invalid`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got OperationKind
			if err := yaml.Unmarshal(tt.value, &got); (err != nil) != tt.wantErr {
				t.Errorf("OperationKind.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("OperationKind.UnmarshalYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_GetOperationFn(t *testing.T) {
	emptyFn := func(_ *Operation) OperationFn {
		return nil
	}

	tests := []struct {
		name      string
		operation Operation
		want      func(o *Operation) OperationFn
		wantErr   bool
	}{
		{
			name: "get executeAQL operation",
			operation: Operation{
				Kind: OperationKindAQLExecute,
			},
			want: ExecuteAQLOperation,
		},
		{
			name: "get createCollection operation",
			operation: Operation{
				Kind: OperationKindCollectionCreate,
			},
			want: CreateCollectionOperation,
		},
		{
			name: "get updateCollection operation",
			operation: Operation{
				Kind: OperationKindCollectionUpdate,
			},
			want: UpdateCollectionOperation,
		},
		{
			name: "get deleteCollection operation",
			operation: Operation{
				Kind: OperationKindCollectionDelete,
			},
			want: DeleteCollectionOperation,
		},
		{
			name: "get createGraph operation",
			operation: Operation{
				Kind: OperationKindGraphCreate,
			},
			want: CreateGraphOperation,
		},
		{
			name: "get addVertexToGraph operation",
			operation: Operation{
				Kind: OperationKindGraphAddVertex,
			},
			want: AddVertexOperation,
		},
		{
			name: "get removeVertexFromGraph operation",
			operation: Operation{
				Kind: OperationKindGraphRemoveVertex,
			},
			want: RemoveVertexOperation,
		},
		{
			name: "get addEdgeToGraph operation",
			operation: Operation{
				Kind: OperationKindGraphAddEdge,
			},
			want: AddEdgeOperation,
		},
		{
			name: "get removeEdgeFromGraph operation",
			operation: Operation{
				Kind: OperationKindGraphRemoveEdge,
			},
			want: RemoveEdgeOperation,
		},
		{
			name: "get deleteGraph operation",
			operation: Operation{
				Kind: OperationKindGraphDelete,
			},
			want: DeleteGraphOperation,
		},
		{
			name: "get createView operation",
			operation: Operation{
				Kind: OperationKindViewCreate,
			},
			want: CreateViewOperation,
		},
		{
			name: "get updateView operation",
			operation: Operation{
				Kind: OperationKindViewUpdate,
			},
			want: UpdateViewOperation,
		},
		{
			name: "get deleteView operation",
			operation: Operation{
				Kind: OperationKindViewDelete,
			},
			want: DeleteViewOperation,
		},
		{
			name: "get createFulltextIndex operation",
			operation: Operation{
				Kind: OperationKindFulltextIndexCreate,
			},
			want: CreateFulltextIndexOperation,
		},
		{
			name: "get createGeoSpatialIndex operation",
			operation: Operation{
				Kind: OperationKindGeoSpatialIndexCreate,
			},
			want: CreateGeoSpatialIndexOperation,
		},
		{
			name: "get createHashIndex operation",
			operation: Operation{
				Kind: OperationKindHashIndexCreate,
			},
			want: CreateHashIndexOperation,
		},
		{
			name: "get createInvertedIndex operation",
			operation: Operation{
				Kind: OperationKindInvertedIndexCreate,
			},
			want: CreateInvertedIndexOperation,
		},
		{
			name: "get createPersistentIndex operation",
			operation: Operation{
				Kind: OperationKindPersistentIndexCreate,
			},
			want: CreatePersistentIndexOperation,
		},
		{
			name: "get createSkipListIndex operation",
			operation: Operation{
				Kind: OperationKindSkipListIndexCreate,
			},
			want: CreateSkipListIndexOperation,
		},
		{
			name: "get createTTLIndex operation",
			operation: Operation{
				Kind: OperationKindTTLIndexCreate,
			},
			want: CreateTTLIndexOperation,
		},
		{
			name: "get createZKDIndex operation",
			operation: Operation{
				Kind: OperationKindZKDIndexCreate,
			},
			want: CreateZKDIndexOperation,
		},
		{
			name: "get deleteIndex operation",
			operation: Operation{
				Kind: OperationKindIndexDelete,
			},
			want: DeleteIndexOperation,
		},
		{
			name: "get createAnalyzer operation",
			operation: Operation{
				Kind: OperationKindAnalyzerCreate,
			},
			want: CreateAnalyzerOperation,
		},
		{
			name: "get deleteAnalyzer operation",
			operation: Operation{
				Kind: OperationKindAnalyzerDelete,
			},
			want: DeleteAnalyzerOperation,
		},
		{
			name: "get unknown operation",
			operation: Operation{
				Kind: OperationKind(0),
			},
			want:    emptyFn,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.operation.GetOperationFn()
			if (err != nil) != tt.wantErr {
				t.Errorf("Operation.GetOperationFn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := tt.want(&tt.operation)
			if reflect.ValueOf(got).Pointer() != reflect.ValueOf(want).Pointer() {
				t.Errorf("Operation.GetOperationFn() = %v, want %v", got, want)
			}
		})
	}
}

func TestExecuteAQLOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "execute aql operation",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAQLExecute,
					Options: map[string]any{
						"query": "FOR i IN @range RETURN i",
						"bindVars": map[string]any{
							"range": "1..10",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Query", context.Background(), "FOR i IN @range RETURN i", map[string]any{
						"range": "1..10",
					}).Return(nil, nil)

					return db
				}(),
			},
		},
		{
			name: "execute aql operation with error",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAQLExecute,
					Options: map[string]any{
						"query": "FOR i IN @range RETURN i",
						"bindVars": map[string]any{
							"range": "1..10",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Query", context.Background(), "FOR i IN @range RETURN i", map[string]any{
						"range": "1..10",
					}).Return(nil, fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := ExecuteAQLOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteAQLOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCollectionOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create collection operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionCreate,
					Collection: "test",
					Options: map[string]any{
						"waitForSync": true,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("CreateCollection", context.Background(), "test", &driver.CreateCollectionOptions{
						WaitForSync: true,
					}).Return(new(MockArangoCollection), nil)

					return db
				}(),
			},
		},
		{
			name: "create collection operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionCreate,
					Collection: "test",
					Options: map[string]any{
						"waitForSync": true,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("CreateCollection", context.Background(), "test", &driver.CreateCollectionOptions{
						WaitForSync: true,
					}).Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateCollectionOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateCollectionOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCollectionOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update collection operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionUpdate,
					Collection: "test",
					Options: map[string]any{
						"waitForSync": true,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					wfs := true

					coll := new(MockArangoCollection)
					coll.On("SetProperties", context.Background(), driver.SetCollectionPropertiesOptions{
						WaitForSync: &wfs,
					}).Return(nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "test").Return(coll, nil)

					return db
				}(),
			},
		},
		{
			name: "update collection operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionUpdate,
					Collection: "test",
					Options: map[string]any{
						"waitForSync": true,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					wfs := true

					coll := new(MockArangoCollection)
					coll.On("SetProperties", context.Background(), driver.SetCollectionPropertiesOptions{
						WaitForSync: &wfs,
					}).Return(fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "test").Return(coll, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "update collection operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionUpdate,
					Collection: "test",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "test").Return(nil, fmt.Errorf("error"))
					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := UpdateCollectionOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCollectionOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestDeleteCollectionOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete collection operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionDelete,
					Collection: "test",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					coll := new(MockArangoCollection)
					coll.On("Remove", context.Background()).Return(nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "test").Return(coll, nil)

					return db
				}(),
			},
		},
		{
			name: "delete collection operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionDelete,
					Collection: "test",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					coll := new(MockArangoCollection)
					coll.On("Remove", context.Background()).Return(fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "test").Return(coll, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "delete collection operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionDelete,
					Collection: "test",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("CollectionExists", context.Background(), "test").Return(true, nil)
					db.On("Collection", context.Background(), "test").Return(new(MockArangoCollection), fmt.Errorf("error"))
					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := DeleteCollectionOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollectionOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateGraphOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create graph operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphCreate,
					Collection: "test",
					Options: map[string]any{
						"edgeDefinitions": []interface{}{
							map[string]interface{}{
								"collection": "test",
								"from":       []string{"test"},
								"to":         []string{"test"},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("CreateGraphV2", context.Background(), "test", &driver.CreateGraphOptions{
						EdgeDefinitions: []driver.EdgeDefinition{
							{
								Collection: "test",
								From:       []string{"test"},
								To:         []string{"test"},
							},
						},
					}).Return(new(MockArangoGraph), nil)

					return db
				}(),
			},
		},
		{
			name: "create graph operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphCreate,
					Collection: "test",
					Options: map[string]any{
						"edgeDefinitions": []interface{}{
							map[string]interface{}{
								"collection": "test",
								"from":       []string{"test"},
								"to":         []string{"test"},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("CreateGraphV2", context.Background(), "test", &driver.CreateGraphOptions{
						EdgeDefinitions: []driver.EdgeDefinition{
							{
								Collection: "test",
								From:       []string{"test"},
								To:         []string{"test"},
							},
						},
					}).Return(new(MockArangoGraph), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateGraphOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateGraphOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddVertexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add vertex operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphAddVertex,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
						"satellites": []string{
							"sat1",
							"sat2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("CreateVertexCollectionWithOptions", context.Background(), "collection", driver.CreateVertexCollectionOptions{
						Satellites: []string{
							"sat1",
							"sat2",
						},
					}).Return(new(MockArangoCollection), nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
		},
		{
			name: "add vertex operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphAddVertex,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
						"satellites": []string{
							"sat1",
							"sat2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("CreateVertexCollectionWithOptions", context.Background(), "collection", driver.CreateVertexCollectionOptions{
						Satellites: []string{
							"sat1",
							"sat2",
						},
					}).Return(new(MockArangoCollection), fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "add vertex operation with graph open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphAddVertex,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(new(MockArangoGraph), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := AddVertexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("AddVertexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveVertexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove vertex operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveVertex,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					coll := new(MockArangoCollection)
					coll.On("Remove", context.Background()).Return(nil)

					graph := new(MockArangoGraph)
					graph.On("VertexCollection", context.Background(), "collection").Return(coll, nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
		},
		{
			name: "remove vertex operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveVertex,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					coll := new(MockArangoCollection)
					coll.On("Remove", context.Background()).Return(fmt.Errorf("error"))

					graph := new(MockArangoGraph)
					graph.On("VertexCollection", context.Background(), "collection").Return(coll, nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "remove vertex operation with graph open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveVertex,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(new(MockArangoGraph), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "remove vertex operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveVertex,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("VertexCollection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := RemoveVertexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("RemoveVertexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestAddEdgeOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add edge operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphAddEdge,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
						"constraints": map[string][]string{
							"from": {"from"},
							"to":   {"to"},
						},
						"satellites": []string{
							"sat1",
							"sat2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("CreateEdgeCollectionWithOptions", context.Background(),
						"collection",
						driver.VertexConstraints{
							From: []string{"from"},
							To:   []string{"to"},
						},
						driver.CreateEdgeCollectionOptions{
							Satellites: []string{"sat1", "sat2"},
						},
					).Return(new(MockArangoCollection), nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
		},
		{
			name: "add edge operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphAddEdge,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
						"constraints": map[string][]string{
							"from": {"from"},
							"to":   {"to"},
						},
						"satellites": []string{
							"sat1",
							"sat2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("CreateEdgeCollectionWithOptions", context.Background(),
						"collection",
						driver.VertexConstraints{
							From: []string{"from"},
							To:   []string{"to"},
						},
						driver.CreateEdgeCollectionOptions{
							Satellites: []string{"sat1", "sat2"},
						},
					).Return(new(MockArangoCollection), fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "add edge operation with graph open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphAddEdge,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(new(MockArangoGraph), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := AddEdgeOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("AddEdgeOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveEdgeOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove edge operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveEdge,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					constraints := driver.VertexConstraints{}

					coll := new(MockArangoCollection)
					coll.On("Remove", context.Background()).Return(nil)

					graph := new(MockArangoGraph)
					graph.On("EdgeCollection", context.Background(), "collection").Return(coll, constraints, nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
		},
		{
			name: "remove edge operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveEdge,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					constraints := driver.VertexConstraints{}

					coll := new(MockArangoCollection)
					coll.On("Remove", context.Background()).Return(fmt.Errorf("error"))

					graph := new(MockArangoGraph)
					graph.On("EdgeCollection", context.Background(), "collection").Return(coll, constraints, nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "remove edge operation with graph open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveEdge,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(new(MockArangoGraph), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "remove edge operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphRemoveEdge,
					Collection: "graph",
					Options: map[string]any{
						"collection": "collection",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("EdgeCollection", context.Background(), "collection").Return(
						new(MockArangoCollection),
						driver.VertexConstraints{},
						fmt.Errorf("error"),
					)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := RemoveEdgeOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("RemoveEdgeOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteGraphOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete graph operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphDelete,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("Remove", context.Background()).Return(nil)

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
		},
		{
			name: "delete graph operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphDelete,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					graph := new(MockArangoGraph)
					graph.On("Remove", context.Background()).Return(fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(graph, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "delete graph operation with graph open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGraphDelete,
					Collection: "graph",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Graph", context.Background(), "graph").Return(new(MockArangoGraph), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})

		if err := DeleteGraphOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
			t.Errorf("DeleteGraphOperation() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestCreateViewOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create view operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewCreate,
					Collection: "view",
					Options: map[string]any{
						"cleanupIntervalStep": 1,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					cis := int64(1)

					db := new(MockArangoDB)
					db.On("CreateArangoSearchView", context.Background(), "view", &driver.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(new(MockArangoSearchView), nil)

					return db
				}(),
			},
		},
		{
			name: "create view operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewCreate,
					Collection: "view",
					Options: map[string]any{
						"cleanupIntervalStep": 1,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					cis := int64(1)

					db := new(MockArangoDB)
					db.On("CreateArangoSearchView", context.Background(), "view", &driver.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(new(MockArangoSearchView), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})

		if err := CreateViewOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
			t.Errorf("CreateViewOperation() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestUpdateViewOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update view operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewUpdate,
					Collection: "view",
					Options: map[string]any{
						"cleanupIntervalStep": 1,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					cis := int64(1)

					searchView := new(MockArangoSearchView)
					searchView.On("SetProperties", context.Background(), driver.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(nil)

					view := new(MockArangoView)
					view.On("ArangoSearchView").Return(searchView, nil)

					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(view, nil)

					return db
				}(),
			},
		},
		{
			name: "update view operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewUpdate,
					Collection: "view",
					Options: map[string]any{
						"cleanupIntervalStep": 1,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					cis := int64(1)

					searchView := new(MockArangoSearchView)
					searchView.On("SetProperties", context.Background(), driver.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(fmt.Errorf("error"))

					view := new(MockArangoView)
					view.On("ArangoSearchView").Return(searchView, nil)

					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(view, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "update view operation with open search view error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewUpdate,
					Collection: "view",
					Options: map[string]any{
						"cleanupIntervalStep": 1,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					view := new(MockArangoView)
					view.On("ArangoSearchView").Return(new(MockArangoSearchView), fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(view, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "update view operation with open view error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewUpdate,
					Collection: "view",
					Options: map[string]any{
						"cleanupIntervalStep": 1,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(new(MockArangoView), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})

		if err := UpdateViewOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
			t.Errorf("UpdateViewOperation() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestDeleteViewOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete view operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewDelete,
					Collection: "view",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					view := new(MockArangoView)
					view.On("Remove", context.Background()).Return(nil)

					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(view, nil)

					return db
				}(),
			},
		},
		{
			name: "delete view operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewDelete,
					Collection: "view",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					view := new(MockArangoView)
					view.On("Remove", context.Background()).Return(fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(view, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "delete view operation with view open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindViewDelete,
					Collection: "view",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("View", context.Background(), "view").Return(new(MockArangoView), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})

		if err := DeleteViewOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
			t.Errorf("DeleteViewOperation() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestCreateFulltextIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create fulltext index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindFulltextIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name":      "fulltext",
						"minLength": 3,
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureFullTextIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureFullTextIndexOptions{ //nolint:staticcheck
							Name:      "fulltext",
							MinLength: 3,
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create fulltext index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindFulltextIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name":      "fulltext",
						"minLength": 3,
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureFullTextIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureFullTextIndexOptions{ //nolint:staticcheck
							Name:      "fulltext",
							MinLength: 3,
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create fulltext index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindFulltextIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateFulltextIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateFulltextIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateGeoSpatialIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create geo index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGeoSpatialIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "geo",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureGeoIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureGeoIndexOptions{
							Name: "geo",
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create geo index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGeoSpatialIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "geo",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureGeoIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureGeoIndexOptions{
							Name: "geo",
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create geo index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindGeoSpatialIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateGeoSpatialIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateGeoSpatialIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateHashIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create hash index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindHashIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "hash",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureHashIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureHashIndexOptions{
							Name: "hash",
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create hash index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindHashIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "hash",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureHashIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureHashIndexOptions{
							Name: "hash",
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create hash index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindHashIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateHashIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateHashIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateInvertedIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create inverted index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindInvertedIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "inverted",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureInvertedIndex", context.Background(),
						&driver.InvertedIndexOptions{
							Name: "inverted",
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create inverted index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindInvertedIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "inverted",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureInvertedIndex", context.Background(),
						&driver.InvertedIndexOptions{
							Name: "inverted",
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create inverted index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindInvertedIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateInvertedIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateInvertedIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatePersistentIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create persistent index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindPersistentIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "persistent",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsurePersistentIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsurePersistentIndexOptions{
							Name: "persistent",
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create persistent index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindPersistentIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "persistent",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsurePersistentIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsurePersistentIndexOptions{
							Name: "persistent",
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create persistent index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindPersistentIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreatePersistentIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreatePersistentIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateSkipListIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create skipList index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSkipListIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "skipList",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureSkipListIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureSkipListIndexOptions{
							Name: "skipList",
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create skipList index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSkipListIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "skipList",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureSkipListIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureSkipListIndexOptions{
							Name: "skipList",
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create skipList index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSkipListIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateSkipListIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateSkipListIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTTLIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create TTL index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindTTLIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name":        "TTL",
						"expireAfter": 3,
						"field":       "field1",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureTTLIndex", context.Background(), "field1", 3, &driver.EnsureTTLIndexOptions{
						Name: "TTL",
					},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create TTL index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindTTLIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name":        "TTL",
						"expireAfter": 3,
						"field":       "field1",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureTTLIndex", context.Background(), "field1", 3, &driver.EnsureTTLIndexOptions{
						Name: "TTL",
					},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create TTL index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindTTLIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateTTLIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateTTLIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateZKDIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create ZKD index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindZKDIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "ZKD",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureZKDIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureZKDIndexOptions{
							Name: "ZKD",
						},
					).Return(new(MockArangoIndex), false, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "create ZKD index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindZKDIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "ZKD",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("EnsureZKDIndex", context.Background(),
						[]string{"field1", "field2"},
						&driver.EnsureZKDIndexOptions{
							Name: "ZKD",
						},
					).Return(new(MockArangoIndex), false, fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "create ZKD index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindZKDIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateZKDIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateZKDIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindIndexDelete,
					Collection: "collection",
					Options: map[string]any{
						"name": "index",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					index := new(MockArangoIndex)
					index.On("Remove", context.Background()).Return(nil)

					collection := new(MockArangoCollection)
					collection.On("Index", context.Background(), "index").Return(index, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
		},
		{
			name: "delete index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindIndexDelete,
					Collection: "collection",
					Options: map[string]any{
						"name": "index",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					index := new(MockArangoIndex)
					index.On("Remove", context.Background()).Return(fmt.Errorf("error"))

					collection := new(MockArangoCollection)
					collection.On("Index", context.Background(), "index").Return(index, nil)

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "delete index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindIndexDelete,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(new(MockArangoCollection), fmt.Errorf("error"))

					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "delete index operation with index open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindIndexDelete,
					Collection: "collection",
					Options: map[string]any{
						"name": "index",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					collection := new(MockArangoCollection)
					collection.On("Index", context.Background(), "index").Return(new(MockArangoIndex), fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Collection", context.Background(), "collection").Return(collection, nil)

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := DeleteIndexOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateAnalyzerOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create analyzer operation",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAnalyzerCreate,
					Options: map[string]any{
						"name": "text_en",
						"type": "text",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("EnsureAnalyzer", context.Background(), driver.ArangoSearchAnalyzerDefinition{
						Name: "text_en",
						Type: "text",
					}).Return(false, new(MockArangoSearchAnalyzer), nil)

					return db
				}(),
			},
		},
		{
			name: "create analyzer operation with error",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAnalyzerCreate,
					Options: map[string]any{
						"name": "text_en",
						"type": "text",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("EnsureAnalyzer", context.Background(), driver.ArangoSearchAnalyzerDefinition{
						Name: "text_en",
						Type: "text",
					}).Return(false, new(MockArangoSearchAnalyzer), fmt.Errorf("error"))
					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "try to recreate existing analyzer",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAnalyzerCreate,
					Options: map[string]any{
						"name": "text_en",
						"type": "text",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("EnsureAnalyzer", context.Background(), driver.ArangoSearchAnalyzerDefinition{
						Name: "text_en",
						Type: "text",
					}).Return(true, new(MockArangoSearchAnalyzer), nil)

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := CreateAnalyzerOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateAnalyzerOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteAnalyzerOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete analyzer operation",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAnalyzerDelete,
					Options: map[string]any{
						"name":  "text_en",
						"force": true,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					analyzer := new(MockArangoSearchAnalyzer)
					analyzer.On("Remove", context.Background(), true).Return(nil)

					db := new(MockArangoDB)
					db.On("Analyzer", context.Background(), "text_en").Return(analyzer, nil)

					return db
				}(),
			},
		},
		{
			name: "delete analyzer operation with error",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAnalyzerDelete,
					Options: map[string]any{
						"name":  "text_en",
						"force": true,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					analyzer := new(MockArangoSearchAnalyzer)
					analyzer.On("Remove", context.Background(), true).Return(fmt.Errorf("error"))

					db := new(MockArangoDB)
					db.On("Analyzer", context.Background(), "text_en").Return(analyzer, nil)

					return db
				}(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := DeleteAnalyzerOperation(tt.fields.operation)(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAnalyzerOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
