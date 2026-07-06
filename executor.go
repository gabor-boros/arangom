package arangom

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
)

type connCtxKey struct{}

// WithConnectionContext returns a context with the given connection stored in it.
func WithConnectionContext(ctx context.Context, conn driver.Connection) context.Context {
	return context.WithValue(ctx, connCtxKey{}, conn)
}

// ConnectionFromContext retrieves the connection from the context, or nil if not set.
func ConnectionFromContext(ctx context.Context) driver.Connection {
	conn, _ := ctx.Value(connCtxKey{}).(driver.Connection)
	return conn
}

// ExecutorOption is a function that sets configuration options on Executor.
type ExecutorOption func(*Executor) error

// WithDatabase sets the database of the executor.
func WithDatabase(db driver.Database) ExecutorOption {
	return func(e *Executor) error {
		if db == nil {
			return ErrNoDatabase
		}

		e.db = db
		return nil
	}
}

// WithCollection sets the collection of the executor.
func WithCollection(collection string) ExecutorOption {
	return func(e *Executor) error {
		if collection == "" {
			return ErrNoCollection
		}

		e.collection = collection
		return nil
	}
}

// WithMigrations sets the migrations of the executor.
func WithMigrations(migrations []*Migration) ExecutorOption {
	return func(e *Executor) error {
		if len(migrations) == 0 {
			return ErrNoMigrations
		}

		e.migrations = migrations
		return nil
	}
}

// WithConnection sets the connection of the executor. This is required for
// operations that need raw HTTP access (e.g. createAnalyzer).
func WithConnection(conn driver.Connection) ExecutorOption {
	return func(e *Executor) error {
		e.conn = conn
		return nil
	}
}

// WithLogger sets the logger of the executor.
func WithLogger(logger Logger) ExecutorOption {
	return func(e *Executor) error {
		if logger == nil {
			return ErrNoLogger
		}

		e.logger = logger
		return nil
	}
}

// Executor executes migrations on a database in order.
type Executor struct {
	db         driver.Database
	conn       driver.Connection
	collection string
	migrations []*Migration
	logger     Logger
}

// Execute executes the migrations on the database.
func (e *Executor) Execute(ctx context.Context) error {
	e.logger.Infof("connecting to the migration collection \"%s\"", e.collection)
	coll, err := e.db.Collection(ctx, e.collection)
	if err != nil {
		return err
	}

	for _, migration := range e.migrations {
		e.logger.Infof("[%d] fetching migration status", migration.ID)
		fetchMigrationStatus(ctx, coll, migration)

		// Skip migrations that have already been run.
		if migration.Status != MigrationStatusMissing {
			if migration.Status == MigrationStatusFailed {
				e.logger.Errorf("[%d] migration failed", migration.ID)
				return ErrMigrationFailed
			}

			e.logger.Infof("[%d] migration already executed", migration.ID)
			continue
		}

		e.logger.Infof("[%d] executing migration", migration.ID)
		migration.Status = MigrationStatusRunning
		saveMigration(ctx, coll, migration)

		migrateCtx := ctx
		if e.conn != nil {
			migrateCtx = WithConnectionContext(ctx, e.conn)
		}

		if err := migration.Migrate(migrateCtx, e.db); err != nil {
			migration.Status = MigrationStatusFailed
			saveMigration(ctx, coll, migration)

			err = errors.Wrap(err, ErrMigrationFailed.Error())
			e.logger.Errorf("[%d] migration failed; err=%s", migration.ID, err.Error())
			return err
		}

		e.logger.Infof("[%d] migration executed successfully", migration.ID)
		migration.Status = MigrationStatusDone
		saveMigration(ctx, coll, migration)
	}

	e.logger.Info("all migrations executed successfully")

	return nil
}

// NewExecutor creates a new Executor. If no migrations are provided, an error
// is returned.
func NewExecutor(opts ...ExecutorOption) (*Executor, error) {
	e := &Executor{
		collection: DefaultMigrationCollection,
		logger:     NewDefaultLogger(),
	}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}

	if e.db == nil {
		return nil, ErrNoDatabase
	}

	if e.migrations == nil {
		return nil, ErrNoMigrations
	}

	return e, nil
}
