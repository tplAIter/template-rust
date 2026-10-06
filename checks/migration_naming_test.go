package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrationTitleTableSequenceAndForwardOnly(t *testing.T) {
	dir := fixture(t, "postgres", "layer_files")
	for _, op := range [][2]string{{"ride", "rides"}, {"create_driver_locations", "driver_locations"}} {
		if e := generate(t, dir, "postgres", "layer_files", "migration", op[0], op[1]); e != nil {
			t.Fatal(e)
		}
	}
	for _, name := range []string{"0002_ride.sql", "0003_create_driver_locations.sql"} {
		b, e := os.ReadFile(filepath.Join(dir, "migrations", name))
		if e != nil {
			t.Fatal(e)
		}
		if strings.Contains(string(b), "DROP TABLE") {
			t.Fatal("forward migration destroys table")
		}
	}
	b := string(source(t, "rust-generators/crud_migration.sql.tmpl"))
	if strings.Contains(b, "goose Down") || strings.Contains(b, "DROP TABLE") {
		t.Fatal("SQLx executes destructive down body")
	}
}

// PostgreSQL permits the declared lowercase identifier grammar even when the
// identifier is a reserved word; quote it instead of narrowing the contract.
func TestQuotedTableIdentifiersAndNoEffectRefusals(t *testing.T) {
	for _, table := range []string{"order", "purchases"} {
		for _, kind := range []string{"migration", "crud"} {
			t.Run(kind+"/"+table, func(t *testing.T) {
				dir := fixture(t, "postgres", "layer_files")
				if err := generate(t, dir, "postgres", "layer_files", kind, "purchase", table); err != nil {
					t.Fatal(err)
				}
				suffix := ""
				if kind == "crud" {
					suffix = "_crud"
				}
				sql, err := os.ReadFile(filepath.Join(dir, "migrations", "0002_purchase"+suffix+".sql"))
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(sql), `CREATE TABLE "`+table+`" (`) {
					t.Fatalf("unquoted migration: %s", sql)
				}
				if kind != "crud" {
					return
				}
				b, err := os.ReadFile(filepath.Join(dir, "src/repository/purchase_repository.rs"))
				if err != nil {
					t.Fatal(err)
				}
				body := string(b)
				for _, want := range []string{
					`INSERT INTO \"` + table + `\" (id, name)`,
					`SELECT id::text AS id, name FROM \"` + table + `\" WHERE`,
					`UPDATE \"` + table + `\" SET`,
					`DELETE FROM \"` + table + `\" WHERE`,
				} {
					if !strings.Contains(body, want) {
						t.Fatalf("missing Rust-escaped quoted query %q", want)
					}
				}
			})
		}
	}
	for _, kind := range []string{"migration", "crud"} {
		for _, table := range []string{"", "Order", "a/b", "a\\b", `a"b`, "orders;DROP TABLE users"} {
			t.Run(kind+"/refuse/"+table, func(t *testing.T) {
				dir := fixture(t, "postgres", "layer_files")
				before := treeDigest(t, dir)
				if generate(t, dir, "postgres", "layer_files", kind, "purchase", table) == nil {
					t.Fatal("invalid/missing table accepted")
				}
				if treeDigest(t, dir) != before {
					t.Fatal("refusal changed fixture before effects")
				}
			})
		}
	}
}
