package checks

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
)

func TestBothLayoutBodiesAndPins(t *testing.T) {
	var ledger struct {
		Bodies map[string]string `json:"bodies"`
	}
	if e := json.Unmarshal(source(t, "block-exports/content-digests.json"), &ledger); e != nil {
		t.Fatal(e)
	}
	if len(ledger.Bodies) != 6 {
		t.Fatal("both layouts missing")
	}
	for p, pin := range ledger.Bodies {
		sum := sha256.Sum256(source(t, p))
		if pin != "sha256:"+hex.EncodeToString(sum[:]) {
			t.Fatal("body pin differs", p)
		}
	}
	b := string(source(t, "block-exports/rust-service-lifecycle.yaml"))
	for _, layer := range []string{"controller", "service", "repository"} {
		for _, suffix := range []string{"base", "entity"} {
			if strings.Count(b, "id: rust."+layer+"."+suffix+",") != 1 {
				t.Fatal("lost or duplicate block identity")
			}
		}
	}
	for _, layout := range []string{"layer_files", "per_entity"} {
		dir := fixture(t, "none", layout)
		if dir == "" {
			t.Fatal("render missing")
		}
	}
}
