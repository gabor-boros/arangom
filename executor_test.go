package arangom

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/arangodb/go-driver/v2/arangodb"
	"go.uber.org/mock/gomock"
)

func expectMigrationDocument(ctx context.Context, coll *MockArangoCollection, key string, item *MigrationItem) {
	coll.EXPECT().ReadDocument(ctx, key, gomock.Any()).DoAndReturn(
		func(_ context.Context, _ string, result any) (arangodb.DocumentMeta, error) {
			encoded, err := json.Marshal(item)
			if err != nil {
				return arangodb.DocumentMeta{}, err
			}
			if err := json.Unmarshal(encoded, result); err != nil {
				return arangodb.DocumentMeta{}, err
			}

			return arangodb.DocumentMeta{}, nil
		},
	)
}

func TestWithDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := NewMockArangoDB(ctrl)
	type fields struct {
		executor *Executor
	}
	type args struct {
		db arangodb.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Executor
		wantErr bool
	}{
		{
			name: "with database",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				db: db,
			},
			want: &Executor{
				db: db,
			},
		},
		{
			name: "with no database",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				db: nil,
			},
			want:    new(Executor),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := WithDatabase(tt.args.db)(tt.fields.executor); (err != nil) != tt.wantErr {
				t.Errorf("WithDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.fields.executor, tt.want) {
				t.Errorf("WithDatabase() = %v, want %v", tt.fields.executor, tt.want)
			}
		})
	}
}

func TestWithCollection(t *testing.T) {
	type fields struct {
		executor *Executor
	}
	type args struct {
		collection string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Executor
		wantErr bool
	}{
		{
			name: "with collection",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				collection: "test",
			},
			want: &Executor{
				collection: "test",
			},
		},
		{
			name: "with no collection",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				collection: "",
			},
			want:    new(Executor),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := WithCollection(tt.args.collection)(tt.fields.executor); (err != nil) != tt.wantErr {
				t.Errorf("WithCollection() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.fields.executor, tt.want) {
				t.Errorf("WithCollection() = %v, want %v", tt.fields.executor, tt.want)
			}
		})
	}
}

func TestWithMigrations(t *testing.T) {
	type fields struct {
		executor *Executor
	}
	type args struct {
		migrations []*Migration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Executor
		wantErr bool
	}{
		{
			name: "with migrations",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
			},
			want: &Executor{
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
			},
		},
		{
			name: "with empty migrations",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				migrations: []*Migration{},
			},
			want:    new(Executor),
			wantErr: true,
		},
		{
			name: "with no migrations",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				migrations: nil,
			},
			want:    new(Executor),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := WithMigrations(tt.args.migrations)(tt.fields.executor); (err != nil) != tt.wantErr {
				t.Errorf("WithMigrations() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.fields.executor, tt.want) {
				t.Errorf("WithMigrations() = %v, want %v", tt.fields.executor, tt.want)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := NewMockLogger(ctrl)
	type fields struct {
		executor *Executor
	}
	type args struct {
		logger Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Executor
		wantErr bool
	}{
		{
			name: "with logger",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				logger: logger,
			},
			want: &Executor{
				logger: logger,
			},
		},
		{
			name: "with no logger",
			fields: fields{
				executor: new(Executor),
			},
			args: args{
				logger: nil,
			},
			want:    new(Executor),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := WithLogger(tt.args.logger)(tt.fields.executor); (err != nil) != tt.wantErr {
				t.Errorf("WithLogger() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.fields.executor, tt.want) {
				t.Errorf("WithLogger() = %v, want %v", tt.fields.executor, tt.want)
			}
		})
	}
}

func TestNewExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := NewMockArangoDB(ctrl)
	logger := NewMockLogger(ctrl)
	type args struct {
		opts []ExecutorOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Executor
		wantErr bool
	}{
		{
			name: "with options",
			args: args{
				opts: []ExecutorOption{
					WithDatabase(db),
					WithCollection("test"),
					WithMigrations([]*Migration{
						{
							ID: 123,
						},
					}),
					WithLogger(logger),
				},
			},
			want: &Executor{
				db:         db,
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: logger,
			},
		},
		{
			name: "with invalid options",
			args: args{
				opts: []ExecutorOption{
					WithDatabase(nil),
					WithCollection("test"),
					WithMigrations([]*Migration{
						{
							ID: 123,
						},
					}),
					WithLogger(logger),
				},
			},
			wantErr: true,
		},
		{
			name: "with no database",
			args: args{
				opts: []ExecutorOption{
					WithCollection("test"),
					WithMigrations([]*Migration{
						{
							ID: 123,
						},
					}),
					WithLogger(logger),
				},
			},
			wantErr: true,
		},
		{
			name: "with no collection",
			args: args{
				opts: []ExecutorOption{
					WithDatabase(db),
					WithMigrations([]*Migration{
						{
							ID: 123,
						},
					}),
					WithLogger(logger),
				},
			},
			want: &Executor{
				db:         db,
				collection: DefaultMigrationCollection,
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: logger,
			},
		},
		{
			name: "with no migrations",
			args: args{
				opts: []ExecutorOption{
					WithDatabase(db),
					WithCollection("test"),
					WithLogger(logger),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewExecutor(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecutor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExecutor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutor_Execute(t *testing.T) {
	type fields struct {
		db         func(*gomock.Controller) arangodb.Database
		collection string
		migrations []*Migration
		logger     func(*gomock.Controller) Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "execute migrations",
			fields: fields{
				db: func(ctrl *gomock.Controller) arangodb.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().DocumentExists(ctx, checksum).Return(false, nil).MinTimes(1)
					coll.EXPECT().CreateDocument(ctx, gomock.Any()).Return(arangodb.CollectionDocumentCreateResponse{}, nil).MinTimes(1)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(ctx, "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return db
				},
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func(ctrl *gomock.Controller) Logger {
					logger := NewMockLogger(ctrl)
					logger.EXPECT().Infof("connecting to the migration collection \"%s\"", "test").Return()
					logger.EXPECT().Infof("[%d] fetching migration status", 123).Return()
					logger.EXPECT().Infof("[%d] executing migration", 123).Return()
					logger.EXPECT().Infof("[%d] migration executed successfully", 123).Return()
					logger.EXPECT().Info("all migrations executed successfully").Return()
					return logger
				},
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "execute on invalid collection",
			fields: fields{
				db: func(ctrl *gomock.Controller) arangodb.Database {
					ctx := context.Background()
					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(ctx, "test", (*arangodb.GetCollectionOptions)(nil)).Return(nil, fmt.Errorf("error"))
					return db
				},
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func(ctrl *gomock.Controller) Logger {
					logger := NewMockLogger(ctrl)
					logger.EXPECT().Infof("connecting to the migration collection \"%s\"", "test").Return()
					return logger
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "execute without migrations",
			fields: fields{
				db: func(ctrl *gomock.Controller) arangodb.Database {
					ctx := context.Background()

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(ctx, "test", (*arangodb.GetCollectionOptions)(nil)).Return(NewMockArangoCollection(ctrl), nil)

					return db
				},
				collection: "test",
				migrations: make([]*Migration, 0),
				logger: func(ctrl *gomock.Controller) Logger {
					logger := NewMockLogger(ctrl)
					logger.EXPECT().Infof("connecting to the migration collection \"%s\"", "test").Return()
					logger.EXPECT().Info("all migrations executed successfully").Return()
					return logger
				},
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "execute with migration error",
			fields: fields{
				db: func(ctrl *gomock.Controller) arangodb.Database {
					ctx := context.Background()
					checksum := "99a2534f50213b0508bddd9e0b1d908d381adfc8e9c1897a2df9a7cd9ec1bb4f"

					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().DocumentExists(ctx, checksum).Return(false, nil).MinTimes(1)
					coll.EXPECT().CreateDocument(ctx, gomock.Any()).Return(arangodb.CollectionDocumentCreateResponse{}, nil).MinTimes(1)

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(ctx, "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)
					db.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

					return db
				},
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
						Operations: []*Operation{
							{
								Kind: OperationKindAQLExecute,
								Options: map[string]any{
									"query": "FOR doc IN @@collection RETURN doc",
									"bindVars": map[string]any{
										"@collection": "test",
									},
								},
							},
						},
					},
				},
				logger: func(ctrl *gomock.Controller) Logger {
					logger := NewMockLogger(ctrl)
					logger.EXPECT().Infof("connecting to the migration collection \"%s\"", "test").Return()
					logger.EXPECT().Infof("[%d] fetching migration status", 123).Return()
					logger.EXPECT().Infof("[%d] executing migration", 123).Return()
					logger.EXPECT().Errorf("[%d] migration failed; err=%s", 123, "migration failed: error").Return()
					return logger
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "execute applied migration",
			fields: fields{
				db: func(ctrl *gomock.Controller) arangodb.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().DocumentExists(ctx, checksum).Return(true, nil)
					expectMigrationDocument(ctx, coll, checksum, &MigrationItem{
						Key:    "123",
						Status: MigrationStatusDone,
					})

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(ctx, "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return db
				},
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func(ctrl *gomock.Controller) Logger {
					logger := NewMockLogger(ctrl)
					logger.EXPECT().Infof("connecting to the migration collection \"%s\"", "test").Return()
					logger.EXPECT().Infof("[%d] fetching migration status", 123).Return()
					logger.EXPECT().Infof("[%d] migration already executed", 123).Return()
					logger.EXPECT().Info("all migrations executed successfully").Return()
					return logger
				},
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "execute failed migration",
			fields: fields{
				db: func(ctrl *gomock.Controller) arangodb.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := NewMockArangoCollection(ctrl)
					coll.EXPECT().DocumentExists(ctx, checksum).Return(true, nil)
					expectMigrationDocument(ctx, coll, checksum, &MigrationItem{
						Key:    "123",
						Status: MigrationStatusFailed,
					})

					db := NewMockArangoDB(ctrl)
					db.EXPECT().GetCollection(ctx, "test", (*arangodb.GetCollectionOptions)(nil)).Return(coll, nil)

					return db
				},
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func(ctrl *gomock.Controller) Logger {
					logger := NewMockLogger(ctrl)
					logger.EXPECT().Infof("connecting to the migration collection \"%s\"", "test").Return()
					logger.EXPECT().Infof("[%d] fetching migration status", 123).Return()
					logger.EXPECT().Errorf("[%d] migration failed", 123).Return()
					return logger
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			e := &Executor{
				db:         tt.fields.db(ctrl),
				collection: tt.fields.collection,
				migrations: tt.fields.migrations,
				logger:     tt.fields.logger(ctrl),
			}

			if err := e.Execute(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Executor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
