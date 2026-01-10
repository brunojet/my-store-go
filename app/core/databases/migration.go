package databases

type MigrationRunner interface {
	MigrateUp() error
	MigrateDown() error
}
