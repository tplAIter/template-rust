package checks

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"
)

type target struct {
	snippet, path string
	numbered      bool
}
type anchor struct{ file, mark, insert string }
type generator struct {
	kind     string
	targets  []target
	anchors  []anchor
	postgres bool
}

func source(t *testing.T, p string) []byte {
	t.Helper()
	b, e := os.ReadFile(filepath.Join("..", p))
	if e != nil {
		t.Fatal(e)
	}
	return b
}

// Read this repository's deliberately small declared YAML shape; core's author
// checker separately performs actual manifest/schema validation. No general YAML
// parser or installed generation authority is supplied by this test fixture.
func declared(t *testing.T) map[string]*generator {
	t.Helper()
	result := map[string]*generator{}
	var g *generator
	section := ""
	active := false
	for _, line := range strings.Split(string(source(t, "template.manifest.yaml")), "\n") {
		if line == "generators:" {
			active = true
			continue
		}
		if !active {
			continue
		}
		v := strings.TrimSpace(line)
		value := func(prefix string) string {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(v, prefix)), "\"")
		}
		if strings.HasPrefix(line, "  - kind:") {
			g = &generator{kind: value("- kind:")}
			result[g.kind] = g
			section = ""
			continue
		}
		if g == nil {
			continue
		}
		switch {
		case strings.HasPrefix(v, "when:"):
			g.postgres = strings.Contains(v, "database=postgres")
		case v == "targets:":
			section = "targets"
		case v == "anchors:":
			section = "anchors"
		case strings.HasPrefix(v, "snippet:") && section == "":
			g.targets = append(g.targets, target{snippet: value("snippet:")})
		case strings.HasPrefix(v, "target:"):
			g.targets[len(g.targets)-1].path = value("target:")
		case strings.HasPrefix(v, "- snippet:"):
			g.targets = append(g.targets, target{snippet: value("- snippet:")})
		case strings.HasPrefix(v, "numbered:"):
			g.targets[len(g.targets)-1].numbered = value("numbered:") == "goose"
		case strings.HasPrefix(v, "- file:"):
			g.anchors = append(g.anchors, anchor{file: value("- file:")})
		case strings.HasPrefix(v, "anchor:"):
			g.anchors[len(g.anchors)-1].mark = value("anchor:")
		case strings.HasPrefix(v, "insert:"):
			g.anchors[len(g.anchors)-1].insert = value("insert:")
		}
	}
	return result
}
func render(t *testing.T, b []byte, data any, db, layout string) []byte {
	t.Helper()
	f := template.FuncMap{"is": func(key, value string) bool {
		if key == "database" {
			return db == value
		}
		return key == "layout" && layout == value
	}}
	tmpl, e := template.New("public").Option("missingkey=error").Funcs(f).Parse(string(b))
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e = tmpl.Execute(&out, data); e != nil {
		t.Fatal(e)
	}
	return out.Bytes()
}
func fixture(t *testing.T, db, layout string) string {
	t.Helper()
	dir := t.TempDir()
	if base := os.Getenv("RUST_FIXTURE_ROOT"); base != "" && t.Name() == "TestMaterializeAllProfiles" {
		dir = filepath.Join(base, db+"-"+layout)
		if _, e := os.Stat(dir); !os.IsNotExist(e) {
			t.Fatal("fixture destination must be fresh")
		}
		if e := os.MkdirAll(dir, 0700); e != nil {
			t.Fatal(e)
		}
	}
	data := map[string]any{"Project": map[string]any{"Slug": "ci_fixture"}, "Runtime": map[string]any{"Port": 8080}, "Settings": map[string]any{"postgres_host_port": 5432}}
	e := filepath.WalkDir("../rust-files", func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel("../rust-files", p)
		if db == "none" && (rel == "src/db.rs.tmpl" || strings.HasPrefix(rel, "migrations/") || rel == "docker-compose.yml.tmpl" || rel == ".env.example.tmpl") {
			return nil
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		name := strings.TrimSuffix(rel, ".tmpl")
		if strings.HasSuffix(rel, ".tmpl") {
			b = render(t, b, data, db, layout)
		}
		dest := filepath.Join(dir, name)
		if e = os.MkdirAll(filepath.Dir(dest), 0755); e != nil {
			return e
		}
		return os.WriteFile(dest, b, 0644)
	})
	if e != nil {
		t.Fatal(e)
	}
	return dir
}

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
var tablePattern = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

func generate(t *testing.T, dir, db, layout, kind, name, table string) error {
	t.Helper()
	g := declared(t)[kind]
	if g == nil {
		return fmt.Errorf("unknown generator")
	}
	if g.postgres && db != "postgres" {
		return fmt.Errorf("unavailable generator")
	}
	if !namePattern.MatchString(name) {
		return fmt.Errorf("invalid name")
	}
	if (kind == "migration" || kind == "crud") && !tablePattern.MatchString(table) {
		return fmt.Errorf("explicit table required")
	}
	pascal := ""
	for _, part := range strings.Split(name, "_") {
		pascal += strings.ToUpper(part[:1]) + part[1:]
	}
	marker := "// gen:" + kind + ":" + name
	seq := 1
	files, _ := filepath.Glob(filepath.Join(dir, "migrations", "*.sql"))
	for _, p := range files {
		n, _ := strconv.Atoi(strings.Split(filepath.Base(p), "_")[0])
		if n >= seq {
			seq = n + 1
		}
	}
	data := map[string]any{"Name": map[string]any{"Snake": name, "Pascal": pascal, "Kebab": strings.ReplaceAll(name, "_", "-")}, "Params": map[string]any{"table": table}, "MigrationSeq": fmt.Sprintf("%04d", seq), "Marker": marker}
	planned := map[string][]byte{}
	for _, spec := range g.targets {
		path := string(render(t, []byte(spec.path), data, db, layout))
		if filepath.IsAbs(path) || strings.Contains(path, "..") {
			return fmt.Errorf("unconfined target")
		}
		dest := filepath.Join(dir, path)
		if _, e := os.Stat(dest); !os.IsNotExist(e) {
			return fmt.Errorf("occupied target")
		}
		planned[dest] = render(t, source(t, spec.snippet), data, db, layout)
	}
	for _, a := range g.anchors {
		path := filepath.Join(dir, a.file)
		b, ok := planned[path]
		if !ok {
			var e error
			b, e = os.ReadFile(path)
			if e != nil {
				return e
			}
		}
		needle := "// " + a.mark
		if bytes.Count(b, []byte(needle)) != 1 {
			return fmt.Errorf("missing or ambiguous anchor")
		}
		if bytes.Contains(b, []byte(marker)) {
			return fmt.Errorf("duplicate generation")
		}
		insert := render(t, source(t, a.insert), data, db, layout)
		if bytes.Count(insert, []byte(marker)) != 1 {
			return fmt.Errorf("marker not preserved")
		}
		planned[path] = bytes.Replace(b, []byte(needle), append(append(insert, '\n'), []byte(needle)...), 1)
	}
	for p, b := range planned {
		if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
			return e
		}
		if e := os.WriteFile(p, b, 0644); e != nil {
			return e
		}
	}
	return nil
}
func TestMaterializeAllProfiles(t *testing.T) {
	for _, db := range []string{"none", "postgres"} {
		for _, layout := range []string{"layer_files", "per_entity"} {
			dir := fixture(t, db, layout)
			for _, op := range [][3]string{{"domain", "location", ""}, {"http-handler", "status", ""}, {"worker", "cleanup", ""}} {
				if e := generate(t, dir, db, layout, op[0], op[1], op[2]); e != nil {
					t.Fatal(e)
				}
			}
			if db == "postgres" {
				for _, op := range [][3]string{{"crud", "ride", "rides"}, {"crud", "driver", "drivers"}, {"migration", "create_events", "events"}} {
					if e := generate(t, dir, db, layout, op[0], op[1], op[2]); e != nil {
						t.Fatal(e)
					}
				}
			}
			t.Logf("materialized %s %s at %s; source-only, no installed authority", db, layout, dir)
		}
	}
}
func TestGeneratorRefusalsHaveNoEffects(t *testing.T) {
	dir := fixture(t, "none", "layer_files")
	before := treeDigest(t, dir)
	for _, op := range [][3]string{{"crud", "ride", "rides"}, {"migration", "events", "events"}, {"domain", "../outside", ""}} {
		if generate(t, dir, "none", "layer_files", op[0], op[1], op[2]) == nil {
			t.Fatal("accepted invalid operation")
		}
		if treeDigest(t, dir) != before {
			t.Fatal("refusal changed fixture")
		}
	}
	dir = fixture(t, "postgres", "per_entity")
	for _, tab := range []string{"", "rides;DROP TABLE users", "Ride", "a/b"} {
		before = treeDigest(t, dir)
		if generate(t, dir, "postgres", "per_entity", "crud", "ride", tab) == nil {
			t.Fatal("bad table accepted")
		}
		if treeDigest(t, dir) != before {
			t.Fatal("invalid table effects")
		}
	}
	if e := generate(t, dir, "postgres", "per_entity", "crud", "ride", "rides"); e != nil {
		t.Fatal(e)
	}
	before = treeDigest(t, dir)
	if generate(t, dir, "postgres", "per_entity", "crud", "ride", "rides") == nil {
		t.Fatal("duplicate accepted")
	}
	if treeDigest(t, dir) != before {
		t.Fatal("duplicate effects")
	}
}

func TestAmbiguousAnchorRefusesBeforeCreatingTargets(t *testing.T) {
	for _, variant := range []string{"missing", "duplicate"} {
		dir := fixture(t, "postgres", "layer_files")
		p := filepath.Join(dir, "src/domain/mod.rs")
		b, e := os.ReadFile(p)
		if e != nil {
			t.Fatal(e)
		}
		if variant == "missing" {
			b = bytes.ReplaceAll(b, []byte("// CODEGEN:CRUD_DOMAINS"), []byte("// removed"))
		} else {
			b = append(b, []byte("\n// CODEGEN:CRUD_DOMAINS\n")...)
		}
		if e = os.WriteFile(p, b, 0644); e != nil {
			t.Fatal(e)
		}
		before := treeDigest(t, dir)
		if generate(t, dir, "postgres", "layer_files", "crud", "ride", "rides") == nil {
			t.Fatal("bad anchor accepted")
		}
		if treeDigest(t, dir) != before {
			t.Fatal("anchor refusal created partial files")
		}
	}
}
