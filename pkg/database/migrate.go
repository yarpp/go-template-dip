package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Migrate(db *sql.DB) error {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		log.Printf("Running migration: %s", f)
		content, err := os.ReadFile(filepath.Join("migrations", f))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", f, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			log.Printf("migration %s failed: %v", f, err)
			return fmt.Errorf("failed to run migration %s: %w", f, err)
		}
		log.Printf("Successfully applied migration: %s", f)
	}
	return nil
}
