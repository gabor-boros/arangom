package arangom

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
)

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

		if err := migration.Migrate(ctx, e.db); err != nil {
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
