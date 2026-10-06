package checks

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	h := sha256.New()
	e := filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		b, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		h.Write([]byte(rel + "\x00"))
		h.Write(b)
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	return hex.EncodeToString(h.Sum(nil))
}
func TestArchitectureAndConcreteBridges(t *testing.T) {
	for _, p := range []string{"rust-files/src/domain/entity.rs", "rust-files/src/domain/error.rs", "rust-generators/crud_entity.rs.tmpl"} {
		b := string(source(t, p))
		if strings.Contains(b, "sqlx") || strings.Contains(b, "axum") {
			t.Fatal("framework import leaked into domain", p)
		}
	}
	for _, p := range []string{"rust-generators/crud_repository.rs.tmpl", "rust-generators/crud_application.rs.tmpl", "rust-generators/crud_service.rs.tmpl"} {
		b := string(source(t, p))
		for _, op := range []string{"create", "find", "update", "delete"} {
			if !strings.Contains(b, "fn "+op) {
				t.Fatal("missing CRUD operation", p, op)
			}
		}
	}
	if strings.Contains(string(source(t, "rust-generators/crud_controller.rs.tmpl")), "NOT_IMPLEMENTED") {
		t.Fatal("CRUD endpoint remains a placeholder")
	}
}
