package checks

import (
	"strings"
	"testing"
)

func TestFiveDeclaredGeneratorsAndSixCRUDTargets(t *testing.T) {
	g := declared(t)
	if len(g) != 5 {
		t.Fatal("five generators required")
	}
	for _, kind := range []string{"domain", "http-handler", "worker", "migration", "crud"} {
		if g[kind] == nil {
			t.Fatal(kind)
		}
	}
	if len(g["crud"].targets) != 6 || len(g["crud"].anchors) != 6 {
		t.Fatal("complete six-target/six-insertion vertical required")
	}
	if !g["crud"].postgres || !g["migration"].postgres {
		t.Fatal("database gating absent")
	}
	if !strings.Contains(string(source(t, "template.manifest.yaml")), "required: true") {
		t.Fatal("explicit table contract absent")
	}
}
