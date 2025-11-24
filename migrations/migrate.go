package migrations

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

type Migration struct {
    Version string
    Name    string
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
    // Create migrations table if it doesn't exist
    if err := createMigrationsTable(db); err != nil {
        return fmt.Errorf("error creating migrations table: %v", err)
    }

    // Get applied migrations
    applied, err := getAppliedMigrations(db)
    if err != nil {
        return fmt.Errorf("error getting applied migrations: %v", err)
    }

    // Get available migrations
    available, err := getAvailableMigrations(migrationsPath)
    if err != nil {
        return fmt.Errorf("error getting available migrations: %v", err)
    }

    // Run pending migrations
    for _, migration := range available {
        if _, exists := applied[migration.Version]; !exists {
            if err := runMigration(db, migration, migrationsPath); err != nil {
                return fmt.Errorf("error running migration %s: %v", migration.Version, err)
            }
        }
    }

    log.Printf("Migrations completed. %d migrations applied.", len(available))
    return nil
}

func createMigrationsTable(db *sql.DB) error {
    query := `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
    _, err := db.Exec(query)
    return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
    rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    applied := make(map[string]bool)
    for rows.Next() {
        var version string
        if err := rows.Scan(&version); err != nil {
            return nil, err
        }
        applied[version] = true
    }
    return applied, nil
}

func getAvailableMigrations(migrationsPath string) ([]Migration, error) {
    files, err := os.ReadDir(migrationsPath)
    if err != nil {
        return nil, err
    }

    var migrations []Migration
    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
            parts := strings.Split(file.Name(), "_")
            if len(parts) > 0 {
                migration := Migration{
                    Version: parts[0],
                    Name:    file.Name(),
                }
                migrations = append(migrations, migration)
            }
        }
    }

    // Sort by version
    sort.Slice(migrations, func(i, j int) bool {
        return migrations[i].Version < migrations[j].Version
    })

    return migrations, nil
}

func runMigration(db *sql.DB, migration Migration, migrationsPath string) error {
    filePath := filepath.Join(migrationsPath, migration.Name)
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    // Start transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Execute migration
    if _, err := tx.Exec(string(content)); err != nil {
        return fmt.Errorf("error executing migration %s: %v", migration.Name, err)
    }

    // Record migration
    if _, err := tx.Exec(
        "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
        migration.Version, migration.Name,
    ); err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    log.Printf("Applied migration: %s", migration.Name)
    return nil
}