package physical

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"

	_ "github.com/lib/pq"
)

func TestCockroachDBBackend(t *testing.T) {
	connURL := os.Getenv("CRURL")
	if connURL == "" {
		t.SkipNow()
	}

	table := os.Getenv("CRTABLE")
	if table == "" {
		table = "vault_kv_store"
	}

	// Run vault tests
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("cockroachdb", logger, map[string]string{
		"connection_url": connURL,
		"table":          table,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		crdb := b.(*CockroachDBBackend)
		_, err := crdb.client.Exec("TRUNCATE TABLE " + crdb.table)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}