package checks

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
)

func TestRootV1AndUnavailableProfiles(t *testing.T) {
	var c struct {
		APIVersion   string `json:"apiVersion"`
		Manifest     string `json:"manifestSHA256"`
		Dependencies []any  `json:"dependencies"`
	}
	if e := json.Unmarshal(source(t, "template.contract.json"), &c); e != nil {
		t.Fatal(e)
	}
	if c.APIVersion != "tplaiter.dev/native-template-contract/v1" || len(c.Dependencies) != 0 {
		t.Fatal("root v1 changed")
	}
	sum := sha256.Sum256(source(t, "template.manifest.yaml"))
	if c.Manifest != "sha256:"+hex.EncodeToString(sum[:]) {
		t.Fatal("raw manifest binding differs")
	}
	b := string(source(t, "template.manifest.yaml"))
	for _, token := range []string{"blockExports:", "dependencies:", "commands:", "OpenAPI", "oidc", "otlp"} {
		if strings.Contains(b, token) {
			t.Fatal("unsupported profile/grant", token)
		}
	}
}
