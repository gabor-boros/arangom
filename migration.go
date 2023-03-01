package arangom

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/arangodb/go-driver"
	"gopkg.in/yaml.v3"
)

const (
	// MigrationStatusMissing is the status of a migration that is missing.
	MigrationStatusMissing MigrationStatus = iota
	// MigrationStatusRunning is the status of a migration that is currently running.
	MigrationStatusRunning
	// MigrationStatusDone is the status of a migration that has been run.
	MigrationStatusDone
	// MigrationStatusFailed is the status of a migration that has failed.
	MigrationStatusFailed
)

const (
	// DefaultMigrationCollection is the default collection name for the collection storing migrations.
	DefaultMigrationCollection = "migrations"
)

var (
	// ErrInvalidMigrationStatus is returned when an invalid migration status is provided.
	ErrInvalidMigrationStatus = fmt.Errorf("invalid migration status")
	// ErrNoMigrations is returned when no migrations are provided.
	ErrNoMigrations = fmt.Errorf("no migrations provided")
	// ErrMigrationFailed is returned when a migration has failed.
	ErrMigrationFailed = fmt.Errorf("migration failed")
)

// MigrationStatus is the status of a migration.
type MigrationStatus int

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *MigrationStatus) UnmarshalJSON(status []byte) error {
	switch string(status) {
	case `"missing"`:
		*s = MigrationStatusMissing
		return nil
	case `"running"`:
		*s = MigrationStatusRunning
		return nil
	case `"done"`:
		*s = MigrationStatusDone
		return nil
	case `"failed"`:
		*s = MigrationStatusFailed
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidMigrationStatus, status)
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (s *MigrationStatus) MarshalJSON() ([]byte, error) {
	switch *s {
	case MigrationStatusMissing:
		return []byte(`"missing"`), nil
	case MigrationStatusRunning:
		return []byte(`"running"`), nil
	case MigrationStatusDone:
		return []byte(`"done"`), nil
	case MigrationStatusFailed:
		return []byte(`"failed"`), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrInvalidMigrationStatus, s)
	}
}

// MigrationItem represents a migration in the collection.
type MigrationItem struct {
	Key       string          `json:"_key"`
	Name      string          `json:"name"`
	Status    MigrationStatus `json:"status"`
	AppliedAt time.Time       `json:"appliedAt"`
}

// Migration is a migration that can be run on a database.
type Migration struct {
	ID         int             `yaml:"id"`
	Path       string          `yaml:"-"`
	Status     MigrationStatus `yaml:"-"`
	Operations []*Operation    `yaml:"operations"`
}

// Name returns the name of the migration. The name is the name of the file
// without the extension.
func (m *Migration) Name() string {
	parts := strings.Split(m.Path, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".yaml")
}

// Checksum returns the checksum of the migration. The checksum is the SHA256
// hash of the ID and operations of the migration file. It is used to determine
// if a migration has been run already.
func (m *Migration) Checksum() (string, error) {
	fields := struct {
		ID         int
		Operations []*Operation `yaml:"operations"`
	}{
		ID:         m.ID,
		Operations: m.Operations,
	}

	b, err := yaml.Marshal(fields)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Migrate executes the operations registered to the migration.
func (m *Migration) Migrate(ctx context.Context, db driver.Database) error {
	for _, operation := range m.Operations {
		opFn, err := operation.GetOperationFn()
		if err != nil {
			return err
		}

		if err := opFn(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

// saveMigration saves the migration to the database. The migration is saved
// as a document in the given collection. If the document already exists, the
// migration is updated.
func saveMigration(ctx context.Context, coll driver.Collection, migration *Migration) {
	checksum, err := migration.Checksum()
	if err != nil {
		panic(err)
	}

	exists, err := coll.DocumentExists(ctx, checksum)
	if err != nil {
		panic(err)
	}

	item := &MigrationItem{
		Key:       checksum,
		Name:      migration.Name(),
		Status:    migration.Status,
		AppliedAt: time.Now(),
	}

	if !exists {
		if _, err := coll.CreateDocument(ctx, item); err != nil {
			panic(err)
		}
		return
	}

	if _, err := coll.UpdateDocument(ctx, checksum, item); err != nil {
		panic(err)
	}
}

// fetchMigrationStatus fetches the status of the migration from the database.
func fetchMigrationStatus(ctx context.Context, coll driver.Collection, migration *Migration) {
	checksum, err := migration.Checksum()
	if err != nil {
		panic(err)
	}

	exists, err := coll.DocumentExists(ctx, checksum)
	if err != nil {
		panic(err)
	}

	if !exists {
		migration.Status = MigrationStatusMissing
		return
	}

	item := new(MigrationItem)
	if _, err := coll.ReadDocument(ctx, checksum, item); err != nil {
		panic(err)
	}

	migration.Status = item.Status
}
