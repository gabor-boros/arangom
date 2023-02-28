package arangom

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/stretchr/testify/mock"
)

func TestWithDatabase(t *testing.T) {
	type fields struct {
		executor *Executor
	}
	type args struct {
		db driver.Database
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
				db: new(MockArangoDB),
			},
			want: &Executor{
				db: new(MockArangoDB),
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
		tt := tt
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
		tt := tt
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
		tt := tt
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
				logger: new(MockLogger),
			},
			want: &Executor{
				logger: new(MockLogger),
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
		tt := tt
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
					WithDatabase(new(MockArangoDB)),
					WithCollection("test"),
					WithMigrations([]*Migration{
						{
							ID: 123,
						},
					}),
					WithLogger(new(MockLogger)),
				},
			},
			want: &Executor{
				db:         new(MockArangoDB),
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: new(MockLogger),
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
					WithLogger(new(MockLogger)),
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
					WithLogger(new(MockLogger)),
				},
			},
			wantErr: true,
		},
		{
			name: "with no collection",
			args: args{
				opts: []ExecutorOption{
					WithDatabase(new(MockArangoDB)),
					WithMigrations([]*Migration{
						{
							ID: 123,
						},
					}),
					WithLogger(new(MockLogger)),
				},
			},
			want: &Executor{
				db:         new(MockArangoDB),
				collection: DefaultMigrationCollection,
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: new(MockLogger),
			},
		},
		{
			name: "with no migrations",
			args: args{
				opts: []ExecutorOption{
					WithDatabase(new(MockArangoDB)),
					WithCollection("test"),
					WithLogger(new(MockLogger)),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
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
		db         driver.Database
		collection string
		migrations []*Migration
		logger     Logger
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
				db: func() driver.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := new(MockArangoCollection)
					coll.On("DocumentExists", ctx, checksum).Return(false, nil)
					coll.On("CreateDocument", ctx, mock.Anything).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", ctx, "test").Return(coll, nil)

					return db
				}(),
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func() Logger {
					logger := new(MockLogger)
					logger.On("Infof", "connecting to the migration collection \"%s\"", []any{"test"}).Return()
					logger.On("Infof", "[%d] fetching migration status", []any{123}).Return()
					logger.On("Infof", "[%d] executing migration", []any{123}).Return()
					logger.On("Infof", "[%d] migration executed successfully", []any{123}).Return()
					logger.On("Info", []any{"all migrations executed successfully"}).Return()
					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "execute on invalid collection",
			fields: fields{
				db: func() driver.Database {
					ctx := context.Background()
					db := new(MockArangoDB)
					db.On("Collection", ctx, "test").Return(nil, fmt.Errorf("error"))
					return db
				}(),
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func() Logger {
					logger := new(MockLogger)
					logger.On("Infof", "connecting to the migration collection \"%s\"", []any{"test"}).Return()
					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "execute without migrations",
			fields: fields{
				db: func() driver.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := new(MockArangoCollection)
					coll.On("DocumentExists", ctx, checksum).Return(false, nil)
					coll.On("CreateDocument", ctx, mock.Anything).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", ctx, "test").Return(coll, nil)

					return db
				}(),
				collection: "test",
				migrations: make([]*Migration, 0),
				logger: func() Logger {
					logger := new(MockLogger)
					logger.On("Infof", "connecting to the migration collection \"%s\"", []any{"test"}).Return()
					logger.On("Info", []any{"all migrations executed successfully"}).Return()
					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "execute with migration error",
			fields: fields{
				db: func() driver.Database {
					ctx := context.Background()
					checksum := "99a2534f50213b0508bddd9e0b1d908d381adfc8e9c1897a2df9a7cd9ec1bb4f"

					coll := new(MockArangoCollection)
					coll.On("DocumentExists", ctx, checksum).Return(false, nil)
					coll.On("CreateDocument", ctx, mock.Anything).Return(driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", ctx, "test").Return(coll, nil)
					db.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error"))

					return db
				}(),
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
				logger: func() Logger {
					logger := new(MockLogger)
					logger.On("Infof", "connecting to the migration collection \"%s\"", []any{"test"}).Return()
					logger.On("Infof", "[%d] fetching migration status", []any{123}).Return()
					logger.On("Infof", "[%d] executing migration", []any{123}).Return()
					logger.On("Errorf", "[%d] migration failed; err=%s", []any{123, "migration failed: error"}).Return()
					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "execute applied migration",
			fields: fields{
				db: func() driver.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := new(MockArangoCollection)
					coll.On("DocumentExists", ctx, checksum).Return(true, nil)
					coll.On("ReadDocument", ctx, checksum, mock.Anything).Return(&MigrationItem{
						Key:    "123",
						Status: MigrationStatusDone,
					}, driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", ctx, "test").Return(coll, nil)

					return db
				}(),
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func() Logger {
					logger := new(MockLogger)
					logger.On("Infof", "connecting to the migration collection \"%s\"", []any{"test"}).Return()
					logger.On("Infof", "[%d] fetching migration status", []any{123}).Return()
					logger.On("Infof", "[%d] migration already executed", []any{123}).Return()
					logger.On("Info", []any{"all migrations executed successfully"}).Return()
					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "execute failed migration",
			fields: fields{
				db: func() driver.Database {
					ctx := context.Background()
					checksum := "a4af0f02f3ed717e72f74da096cf4a3e36af3bdae4515fdd664f316163341ef2"

					coll := new(MockArangoCollection)
					coll.On("DocumentExists", ctx, checksum).Return(true, nil)
					coll.On("ReadDocument", ctx, checksum, mock.Anything).Return(&MigrationItem{
						Key:    "123",
						Status: MigrationStatusFailed,
					}, driver.DocumentMeta{}, nil)

					db := new(MockArangoDB)
					db.On("Collection", ctx, "test").Return(coll, nil)

					return db
				}(),
				collection: "test",
				migrations: []*Migration{
					{
						ID: 123,
					},
				},
				logger: func() Logger {
					logger := new(MockLogger)
					logger.On("Infof", "connecting to the migration collection \"%s\"", []any{"test"}).Return()
					logger.On("Infof", "[%d] fetching migration status", []any{123}).Return()
					logger.On("Errorf", "[%d] migration failed", []any{123}).Return()
					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &Executor{
				db:         tt.fields.db,
				collection: tt.fields.collection,
				migrations: tt.fields.migrations,
				logger:     tt.fields.logger,
			}

			if err := e.Execute(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Executor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
