package arangom

import (
	"context"
	"fmt"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/stretchr/testify/mock"
)

func TestMigrationStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       *MigrationStatus
		status  []byte
		want    MigrationStatus
		wantErr bool
	}{
		{
			name:   "missing",
			s:      new(MigrationStatus),
			status: []byte(`"missing"`),
			want:   MigrationStatusMissing,
		},
		{
			name:   "running",
			s:      new(MigrationStatus),
			status: []byte(`"running"`),
			want:   MigrationStatusRunning,
		},
		{
			name:   "done",
			s:      new(MigrationStatus),
			status: []byte(`"done"`),
			want:   MigrationStatusDone,
		},
		{
			name:   "failed",
			s:      new(MigrationStatus),
			status: []byte(`"failed"`),
			want:   MigrationStatusFailed,
		},
		{
			name:    "invalid",
			s:       new(MigrationStatus),
			status:  []byte(`"invalid"`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.s.UnmarshalJSON(tt.status); (err != nil) != tt.wantErr {
				t.Errorf("MigrationStatus.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if *tt.s != tt.want {
				t.Errorf("MigrationStatus.UnmarshalJSON() = %v, want %v", tt.s, tt.want)
			}
		})
	}
}

func TestMigrationStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       MigrationStatus
		want    []byte
		wantErr bool
	}{
		{
			name: "missing",
			s:    MigrationStatusMissing,
			want: []byte(`"missing"`),
		},
		{
			name: "running",
			s:    MigrationStatusRunning,
			want: []byte(`"running"`),
		},
		{
			name: "done",
			s:    MigrationStatusDone,
			want: []byte(`"done"`),
		},
		{
			name: "failed",
			s:    MigrationStatusFailed,
			want: []byte(`"failed"`),
		},
		{
			name:    "invalid",
			s:       MigrationStatus(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MigrationStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if string(got) != string(tt.want) {
				t.Errorf("MigrationStatus.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestMigration_Name(t *testing.T) {
	tests := []struct {
		name      string
		migration *Migration
		want      string
	}{
		{
			name: "get name from file name",
			migration: &Migration{
				Path: "name.yaml",
			},
			want: "name",
		},
		{
			name: "get name from relative path",
			migration: &Migration{
				Path: "./path/to/name.yaml",
			},
			want: "name",
		},
		{
			name: "get name from absolute path",
			migration: &Migration{
				Path: "/path/to/name.yaml",
			},
			want: "name",
		},
		{
			name: "get name from file name without extension",
			migration: &Migration{
				Path: "name",
			},
			want: "name",
		},
		{
			name: "get name from file name with multiple dots",
			migration: &Migration{
				Path: "file.name.yaml",
			},
			want: "file.name",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.migration.Name(); got != tt.want {
				t.Errorf("Migration.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigration_Checksum(t *testing.T) {
	tests := []struct {
		name      string
		migration *Migration
		want      string
		wantErr   bool
	}{
		{
			name: "get checksum with no operations",
			migration: &Migration{
				Path: "name.yaml",
			},
			want: "3ec9d32b15db2ba18821258a344de1e69111834f2dcddc13c2ca62fd0e2cc66d",
		},
		{
			name: "get checksum with operations",
			migration: &Migration{
				Path: "name.yaml",
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
			want: "78f5c1046a6d4edbedd55a4780de9211e8577592887521e1ed862cb6a501c52b",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.migration.Checksum()
			if (err != nil) != tt.wantErr {
				t.Errorf("Migration.Checksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Migration.Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigration_Migrate(t *testing.T) {
	type args struct {
		ctx context.Context
		db  driver.Database
	}
	tests := []struct {
		name      string
		migration *Migration
		args      args
		wantErr   bool
	}{
		{
			name: "migrate with no operations",
			migration: &Migration{
				Path: "migration.yaml",
			},
			args: args{
				ctx: context.Background(),
				db:  new(MockArangoDB),
			},
		},
		{
			name: "migrate with operations",
			migration: &Migration{
				Path: "migration.yaml",
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
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					query := "FOR doc IN @@collection RETURN doc"
					bindVars := map[string]any{
						"@collection": "test",
					}
					db := new(MockArangoDB)
					db.On("Query", context.Background(), query, bindVars).Return(nil, nil)
					return db
				}(),
			},
		},
		{
			name: "migrate with operation error",
			migration: &Migration{
				Path: "migration.yaml",
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
			args: args{
				ctx: context.Background(),
				db: func() driver.Database {
					db := new(MockArangoDB)
					db.On("Query", context.Background(), mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error"))
					return db
				}(),
			},
			wantErr: true,
		},
		{
			name: "migrate with invalid operation kind",
			migration: &Migration{
				Path: "migration.yaml",
				Operations: []*Operation{
					{
						Kind: OperationKind(0),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				db:  new(MockArangoDB),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.migration.Migrate(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("Migration.Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
