package db

import (
	"database/sql"
	"strings"
)

type MigrationHistory struct {
	Version   string
	Statement string
	CreatedTs int64
}

type MigrationHistoryCreate struct {
	Version   string
	Statement string
}

type MigrationHistoryFind struct {
	Version *string
}

func findMigrationHistoryList(db *sql.DB, find *MigrationHistoryFind) ([]*MigrationHistory, error) {
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := find.Version; v != nil {
		where, args = append(where, "version = ?"), append(args, *v)
	}

	rows, err := db.Query(`
		SELECT 
			version,
			created_ts
		FROM
			migration_history
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_ts DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	migrationHistoryList := make([]*MigrationHistory, 0)
	for rows.Next() {
		var migrationHistory MigrationHistory
		if err := rows.Scan(
			&migrationHistory.Version,
			&migrationHistory.CreatedTs,
		); err != nil {
			return nil, err
		}

		migrationHistoryList = append(migrationHistoryList, &migrationHistory)
	}

	return migrationHistoryList, nil
}

func findMigrationHistory(db *sql.DB, find *MigrationHistoryFind) (*MigrationHistory, error) {
	list, err := findMigrationHistoryList(db, find)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	} else {
		return list[0], nil
	}
}

func upsertMigrationHistory(db *sql.DB, create *MigrationHistoryCreate) (*MigrationHistory, error) {
	row, err := db.Query(`
		INSERT INTO migration_history (
			version, 
			statement
		)
		VALUES (?, ?)
		ON CONFLICT(version) DO UPDATE
		SET
			version=EXCLUDED.version,
			statement=EXCLUDED.statement
		RETURNING version, statement, created_ts
	`,
		create.Version,
		create.Statement,
	)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	row.Next()
	var migrationHistory MigrationHistory
	if err := row.Scan(
		&migrationHistory.Version,
		&migrationHistory.Statement,
		&migrationHistory.CreatedTs,
	); err != nil {
		return nil, err
	}

	return &migrationHistory, nil
}
