package arangom

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/arangodb/go-driver/v2/arangodb"
	"go.uber.org/mock/gomock"
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
			name:  "unmarshal createGeoSpatialIndex",
			value: []byte(`createGeoSpatialIndex`),
			want:  OperationKindGeoSpatialIndexCreate,
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
			name:  "unmarshal createTTLIndex",
			value: []byte(`createTTLIndex`),
			want:  OperationKindTTLIndexCreate,
		},
		{
			name:  "unmarshal createMDIIndex",
			value: []byte(`createMDIIndex`),
			want:  OperationKindMDIIndexCreate,
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
			name:  "unmarshal createMDIPrefixedIndex",
			value: []byte(`createMDIPrefixedIndex`),
			want:  OperationKindMDIPrefixedIndexCreate,
		},
		{
			name:  "unmarshal createVectorIndex",
			value: []byte(`createVectorIndex`),
			want:  OperationKindVectorIndexCreate,
		},
		{
			name:  "unmarshal createSearchAliasView",
			value: []byte(`createSearchAliasView`),
			want:  OperationKindSearchAliasViewCreate,
		},
		{
			name:  "unmarshal updateSearchAliasView",
			value: []byte(`updateSearchAliasView`),
			want:  OperationKindSearchAliasViewUpdate,
		},
		{
			name:    "unmarshal invalid operation kind",
			value:   []byte(`invalid`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
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
			name: "get createGeoSpatialIndex operation",
			operation: Operation{
				Kind: OperationKindGeoSpatialIndexCreate,
			},
			want: CreateGeoSpatialIndexOperation,
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
			name: "get createTTLIndex operation",
			operation: Operation{
				Kind: OperationKindTTLIndexCreate,
			},
			want: CreateTTLIndexOperation,
		},
		{
			name: "get createMDIIndex operation",
			operation: Operation{
				Kind: OperationKindMDIIndexCreate,
			},
			want: CreateMDIIndexOperation,
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
			name: "get createMDIPrefixedIndex operation",
			operation: Operation{
				Kind: OperationKindMDIPrefixedIndexCreate,
			},
			want: CreateMDIPrefixedIndexOperation,
		},
		{
			name: "get createVectorIndex operation",
			operation: Operation{
				Kind: OperationKindVectorIndexCreate,
			},
			want: CreateVectorIndexOperation,
		},
		{
			name: "get createSearchAliasView operation",
			operation: Operation{
				Kind: OperationKindSearchAliasViewCreate,
			},
			want: CreateSearchAliasViewOperation,
		},
		{
			name: "get updateSearchAliasView operation",
			operation: Operation{
				Kind: OperationKindSearchAliasViewUpdate,
			},
			want: UpdateSearchAliasViewOperation,
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Query(context.Background(), "FOR i IN @range RETURN i", &arangodb.QueryOptions{BindVars: map[string]any{
						"range": "1..10",
					}}).Return(nil, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Query(context.Background(), "FOR i IN @range RETURN i", &arangodb.QueryOptions{BindVars: map[string]any{
						"range": "1..10",
					}}).Return(nil, fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := ExecuteAQLOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					wfs := true
					db.EXPECT().CreateCollectionV2(context.Background(), "test", &arangodb.CreateCollectionPropertiesV2{
						WaitForSync: &wfs,
					}).Return(NewMockArangoCollection(ctrl), nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					wfs := true
					db.EXPECT().CreateCollectionV2(context.Background(), "test", &arangodb.CreateCollectionPropertiesV2{
						WaitForSync: &wfs,
					}).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateCollectionOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateCollectionOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCollectionOperation(t *testing.T) {
	keepNull := false
	computed := []arangodb.ComputedValue{
		{
			Name:       "svid_pending",
			Expression: "RETURN @doc.matched_at == null ? true : null",
			ComputeOn: []arangodb.ComputeOn{
				"insert",
				"update",
				"replace",
			},
			Overwrite: true,
			KeepNull:  &keepNull,
		},
	}
	emptyComputed := []arangodb.ComputedValue{}

	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					wfs := true

					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().SetPropertiesV2(context.Background(), arangodb.SetCollectionPropertiesOptionsV2{
						WaitForSync: &wfs,
					}).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "update collection operation sends computed values",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionUpdate,
					Collection: "test",
					Options: map[string]any{
						"computedValues": []any{
							map[string]any{
								"name":       "svid_pending",
								"expression": "RETURN @doc.matched_at == null ? true : null",
								"computeOn":  []any{"insert", "update", "replace"},
								"overwrite":  true,
								"keepNull":   false,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().SetPropertiesV2(context.Background(), arangodb.SetCollectionPropertiesOptionsV2{
						ComputedValues: &computed,
					}).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "update collection operation sends an empty computed values list",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindCollectionUpdate,
					Collection: "test",
					Options: map[string]any{
						"computedValues": []any{},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().SetPropertiesV2(context.Background(), arangodb.SetCollectionPropertiesOptionsV2{
						ComputedValues: &emptyComputed,
					}).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					wfs := true

					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().SetPropertiesV2(context.Background(), arangodb.SetCollectionPropertiesOptionsV2{
						WaitForSync: &wfs,
					}).Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(nil, fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := UpdateCollectionOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().Remove(context.Background()).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().Remove(context.Background()).Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "test", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))
					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := DeleteCollectionOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().CreateGraph(context.Background(), "test", &arangodb.GraphDefinition{
						EdgeDefinitions: []arangodb.EdgeDefinition{
							{
								Collection: "test",
								From:       []string{"test"},
								To:         []string{"test"},
							},
						},
					}, (*arangodb.CreateGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().CreateGraph(context.Background(), "test", &arangodb.GraphDefinition{
						EdgeDefinitions: []arangodb.EdgeDefinition{
							{
								Collection: "test",
								From:       []string{"test"},
								To:         []string{"test"},
							},
						},
					}, (*arangodb.CreateGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateGraphOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().CreateVertexCollection(context.Background(), "collection", &arangodb.CreateVertexCollectionOptions{
						Satellites: []string{
							"sat1",
							"sat2",
						},
					}).Return(arangodb.CreateVertexCollectionResponse{}, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().CreateVertexCollection(context.Background(), "collection", &arangodb.CreateVertexCollectionOptions{
						Satellites: []string{
							"sat1",
							"sat2",
						},
					}).Return(arangodb.CreateVertexCollectionResponse{}, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := AddVertexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().DeleteVertexCollection(context.Background(), "collection", (*arangodb.DeleteVertexCollectionOptions)(nil)).Return(arangodb.DeleteVertexCollectionResponse{}, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().DeleteVertexCollection(context.Background(), "collection", (*arangodb.DeleteVertexCollectionOptions)(nil)).Return(arangodb.DeleteVertexCollectionResponse{}, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().DeleteVertexCollection(context.Background(), "collection", (*arangodb.DeleteVertexCollectionOptions)(nil)).Return(arangodb.DeleteVertexCollectionResponse{}, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := RemoveVertexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().CreateEdgeDefinition(context.Background(),
						"collection",
						[]string{"from"},
						[]string{"to"},
						&arangodb.CreateEdgeDefinitionOptions{
							Satellites: []string{"sat1", "sat2"},
						},
					).Return(arangodb.CreateEdgeDefinitionResponse{}, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().CreateEdgeDefinition(context.Background(),
						"collection",
						[]string{"from"},
						[]string{"to"},
						&arangodb.CreateEdgeDefinitionOptions{
							Satellites: []string{"sat1", "sat2"},
						},
					).Return(arangodb.CreateEdgeDefinitionResponse{}, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := AddEdgeOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().DeleteEdgeDefinition(context.Background(), "collection", (*arangodb.DeleteEdgeDefinitionOptions)(nil)).Return(arangodb.DeleteEdgeDefinitionResponse{}, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().DeleteEdgeDefinition(context.Background(), "collection", (*arangodb.DeleteEdgeDefinitionOptions)(nil)).Return(arangodb.DeleteEdgeDefinitionResponse{}, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().DeleteEdgeDefinition(context.Background(), "collection", (*arangodb.DeleteEdgeDefinitionOptions)(nil)).Return(arangodb.DeleteEdgeDefinitionResponse{}, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := RemoveEdgeOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().Remove(context.Background(), (*arangodb.RemoveGraphOptions)(nil)).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					graph := NewMockArangoGraph(ctrl)
					graph.EXPECT().Remove(context.Background(), (*arangodb.RemoveGraphOptions)(nil)).Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(graph, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().Graph(context.Background(), "graph", (*arangodb.GetGraphOptions)(nil)).Return(NewMockArangoGraph(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := DeleteGraphOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("DeleteGraphOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateViewOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					cis := int64(1)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().CreateArangoSearchView(context.Background(), "view", &arangodb.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(NewMockArangoSearchView(ctrl), nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					cis := int64(1)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().CreateArangoSearchView(context.Background(), "view", &arangodb.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(NewMockArangoSearchView(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateViewOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateViewOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateViewOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					cis := int64(1)

					searchView := NewMockArangoSearchView(ctrl)
					searchView.EXPECT().SetProperties(context.Background(), arangodb.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(nil)

					view := NewMockArangoView(ctrl)
					view.EXPECT().ArangoSearchView().Return(searchView, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(view, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					cis := int64(1)

					searchView := NewMockArangoSearchView(ctrl)
					searchView.EXPECT().SetProperties(context.Background(), arangodb.ArangoSearchViewProperties{
						CleanupIntervalStep: &cis,
					}).Return(fmt.Errorf("error"))

					view := NewMockArangoView(ctrl)
					view.EXPECT().ArangoSearchView().Return(searchView, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(view, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					view := NewMockArangoView(ctrl)
					view.EXPECT().ArangoSearchView().Return(NewMockArangoSearchView(ctrl), fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(view, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(NewMockArangoView(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := UpdateViewOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("UpdateViewOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteViewOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					view := NewMockArangoView(ctrl)
					view.EXPECT().Remove(context.Background()).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(view, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					view := NewMockArangoView(ctrl)
					view.EXPECT().Remove(context.Background()).Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(view, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "view").Return(NewMockArangoView(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := DeleteViewOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("DeleteViewOperation() error = %v, wantErr %v", err, tt.wantErr)
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureGeoIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreateGeoIndexOptions{
							Name: "geo",
						},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureGeoIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreateGeoIndexOptions{
							Name: "geo",
						},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateGeoSpatialIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateGeoSpatialIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureInvertedIndex(context.Background(),
						&arangodb.InvertedIndexOptions{
							Name: "inverted",
						},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureInvertedIndex(context.Background(),
						&arangodb.InvertedIndexOptions{
							Name: "inverted",
						},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateInvertedIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsurePersistentIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreatePersistentIndexOptions{
							Name: "persistent",
						},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsurePersistentIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreatePersistentIndexOptions{
							Name: "persistent",
						},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreatePersistentIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreatePersistentIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureTTLIndex(context.Background(), []string{"field1"}, 3, &arangodb.CreateTTLIndexOptions{
						Name: "TTL",
					},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureTTLIndex(context.Background(), []string{"field1"}, 3, &arangodb.CreateTTLIndexOptions{
						Name: "TTL",
					},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateTTLIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateTTLIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateMDIIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create MDI index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindMDIIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "mdi",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureMDIIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreateMDIIndexOptions{
							Name:            "mdi",
							FieldValueTypes: arangodb.MDIDoubleFieldType,
						},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "create MDI index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindMDIIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "mdi",
						"fields": []string{
							"field1",
							"field2",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureMDIIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreateMDIIndexOptions{
							Name:            "mdi",
							FieldValueTypes: arangodb.MDIDoubleFieldType,
						},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
		{
			name: "create MDI index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindMDIIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateMDIIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateMDIIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().DeleteIndex(context.Background(), "index").Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().DeleteIndex(context.Background(), "index").Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().DeleteIndex(context.Background(), "index").Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := DeleteIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().EnsureCreatedAnalyzer(context.Background(), &arangodb.AnalyzerDefinition{
						Name: "text_en",
						Type: "text",
					}).Return(nil, true, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "create analyzer operation with properties and features",
			fields: fields{
				operation: &Operation{
					Kind: OperationKindAnalyzerCreate,
					Options: map[string]any{
						"name": "ngram_lower",
						"type": "pipeline",
						"properties": map[string]any{
							"pipeline": []any{
								map[string]any{"type": "norm", "properties": map[string]any{"locale": "en"}},
							},
						},
						"features": []any{"frequency", "norm", "position"},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().EnsureCreatedAnalyzer(context.Background(), &arangodb.AnalyzerDefinition{
						Name: "ngram_lower",
						Type: "pipeline",
						Properties: arangodb.ArangoSearchAnalyzerProperties{
							Pipeline: []arangodb.ArangoSearchAnalyzerPipeline{
								{
									Type: "norm",
									Properties: arangodb.ArangoSearchAnalyzerProperties{
										Locale: "en",
									},
								},
							},
						},
						Features: []arangodb.ArangoSearchFeature{
							arangodb.ArangoSearchFeatureFrequency,
							arangodb.ArangoSearchFeatureNorm,
							arangodb.ArangoSearchFeaturePosition,
						},
					}).Return(nil, false, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().EnsureCreatedAnalyzer(context.Background(), &arangodb.AnalyzerDefinition{
						Name: "text_en",
						Type: "text",
					}).Return(nil, false, fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateAnalyzerOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
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
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					analyzer := NewMockArangoSearchAnalyzer(ctrl)
					analyzer.EXPECT().Remove(context.Background(), true).Return(nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Analyzer(context.Background(), "text_en").Return(analyzer, nil)

					return context.Background(), db
				},
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
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					analyzer := NewMockArangoSearchAnalyzer(ctrl)
					analyzer.EXPECT().Remove(context.Background(), true).Return(fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().Analyzer(context.Background(), "text_en").Return(analyzer, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := DeleteAnalyzerOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAnalyzerOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateMDIPrefixedIndexOperation(t *testing.T) {
	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create MDI prefixed index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindMDIPrefixedIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "mdi-prefixed",
						"fields": []string{
							"field1",
							"field2",
						},
						"prefixFields": []string{"tenant"},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureMDIPrefixedIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreateMDIPrefixedIndexOptions{
							CreateMDIIndexOptions: arangodb.CreateMDIIndexOptions{
								Name:            "mdi-prefixed",
								FieldValueTypes: arangodb.MDIDoubleFieldType,
							},
							PrefixFields: []string{"tenant"},
						},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "create MDI prefixed index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindMDIPrefixedIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name": "mdi-prefixed",
						"fields": []string{
							"field1",
							"field2",
						},
						"prefixFields": []string{"tenant"},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureMDIPrefixedIndex(context.Background(),
						[]string{"field1", "field2"},
						&arangodb.CreateMDIPrefixedIndexOptions{
							CreateMDIIndexOptions: arangodb.CreateMDIIndexOptions{
								Name:            "mdi-prefixed",
								FieldValueTypes: arangodb.MDIDoubleFieldType,
							},
							PrefixFields: []string{"tenant"},
						},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
		{
			name: "create MDI prefixed index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindMDIPrefixedIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateMDIPrefixedIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateMDIPrefixedIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateVectorIndexOperation(t *testing.T) {
	dimension := 8
	metric := arangodb.VectorMetricCosine
	nLists := 4
	strategy := arangodb.VectorNListsStrategyAutoSqrt
	multiplier := 2
	minNLists := 1
	name := "embeddings"

	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create vector index operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindVectorIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"name":   name,
						"fields": []string{"embedding"},
						"params": map[string]any{
							"dimension": dimension,
							"metric":    string(metric),
							"nLists":    nLists,
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureVectorIndex(context.Background(),
						[]string{"embedding"},
						&arangodb.VectorParams{
							Dimension: &dimension,
							Metric:    &metric,
							NLists:    arangodb.NewVectorNLists(nLists),
						},
						&arangodb.CreateVectorIndexOptions{Name: &name},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "create vector index operation with scaling nLists",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindVectorIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"fields": []string{"embedding"},
						"params": map[string]any{
							"dimension": dimension,
							"metric":    string(metric),
							"nLists": map[string]any{
								"strategy":   string(strategy),
								"multiplier": multiplier,
								"minNLists":  minNLists,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureVectorIndex(context.Background(),
						[]string{"embedding"},
						&arangodb.VectorParams{
							Dimension: &dimension,
							Metric:    &metric,
							NLists: &arangodb.VectorNLists{
								Scaling: &arangodb.VectorNListsScaling{
									Strategy:   &strategy,
									Multiplier: &multiplier,
									MinNLists:  &minNLists,
								},
							},
						},
						&arangodb.CreateVectorIndexOptions{},
					).Return(arangodb.IndexResponse{}, false, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "create vector index operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindVectorIndexCreate,
					Collection: "collection",
					Options: map[string]any{
						"fields": []string{"embedding"},
						"params": map[string]any{
							"dimension": dimension,
							"metric":    string(metric),
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					collection := NewMockArangoCollection(ctrl)
					collection.EXPECT().EnsureVectorIndex(context.Background(),
						[]string{"embedding"},
						&arangodb.VectorParams{
							Dimension: &dimension,
							Metric:    &metric,
						},
						&arangodb.CreateVectorIndexOptions{},
					).Return(arangodb.IndexResponse{}, false, fmt.Errorf("error"))

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(collection, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
		{
			name: "create vector index operation with collection open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindVectorIndexCreate,
					Collection: "collection",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(context.Background(), "collection", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateVectorIndexOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateVectorIndexOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateSearchAliasViewOperation(t *testing.T) {
	aliasIndexes := []arangodb.ArangoSearchAliasIndex{
		{Collection: "docs", Index: "inv"},
	}

	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create search alias view operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSearchAliasViewCreate,
					Collection: "search",
					Options: map[string]any{
						"indexes": []any{
							map[string]any{"collection": "docs", "index": "inv"},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().CreateArangoSearchAliasView(context.Background(), "search", &arangodb.ArangoSearchAliasViewProperties{
						Indexes: aliasIndexes,
					}).Return(nil, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "create search alias view operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSearchAliasViewCreate,
					Collection: "search",
					Options: map[string]any{
						"indexes": []any{
							map[string]any{"collection": "docs", "index": "inv"},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().CreateArangoSearchAliasView(context.Background(), "search", &arangodb.ArangoSearchAliasViewProperties{
						Indexes: aliasIndexes,
					}).Return(nil, fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := CreateSearchAliasViewOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("CreateSearchAliasViewOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateSearchAliasViewOperation(t *testing.T) {
	properties := arangodb.ArangoSearchAliasViewProperties{
		Indexes: []arangodb.ArangoSearchAliasIndex{
			{Collection: "docs", Index: "inv"},
		},
	}

	type fields struct {
		operation *Operation
	}
	type args struct {
		ctx context.Context
		db  func(*gomock.Controller) (context.Context, arangodb.Database)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update search alias view operation",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSearchAliasViewUpdate,
					Collection: "search",
					Options: map[string]any{
						"indexes": []any{
							map[string]any{"collection": "docs", "index": "inv"},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					alias := NewMockArangoSearchViewAlias(ctrl)
					alias.EXPECT().SetProperties(context.Background(), properties).Return(nil)

					view := NewMockArangoView(ctrl)
					view.EXPECT().ArangoSearchViewAlias().Return(alias, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "search").Return(view, nil)

					return context.Background(), db
				},
			},
		},
		{
			name: "update search alias view operation with error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSearchAliasViewUpdate,
					Collection: "search",
					Options: map[string]any{
						"indexes": []any{
							map[string]any{"collection": "docs", "index": "inv"},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					alias := NewMockArangoSearchViewAlias(ctrl)
					alias.EXPECT().SetProperties(context.Background(), properties).Return(fmt.Errorf("error"))

					view := NewMockArangoView(ctrl)
					view.EXPECT().ArangoSearchViewAlias().Return(alias, nil)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "search").Return(view, nil)

					return context.Background(), db
				},
			},
			wantErr: true,
		},
		{
			name: "update search alias view operation with view open error",
			fields: fields{
				operation: &Operation{
					Kind:       OperationKindSearchAliasViewUpdate,
					Collection: "search",
				},
			},
			args: args{
				ctx: context.Background(),
				db: func(ctrl *gomock.Controller) (context.Context, arangodb.Database) {
					db := NewMockArangoDB(ctrl)
					db.EXPECT().View(context.Background(), "search").Return(NewMockArangoView(ctrl), fmt.Errorf("error"))

					return context.Background(), db
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx, db := tt.args.db(ctrl)
			if err := UpdateSearchAliasViewOperation(tt.fields.operation)(ctx, db); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSearchAliasViewOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
