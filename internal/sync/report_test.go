package sync_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/sync"
)

func TestResult_Summary_InSync(t *testing.T) {
	r := &sync.Result{}
	if got := r.Summary(); got != "in sync" {
		t.Errorf("expected 'in sync', got %q", got)
	}
}

func TestResult_Summary_Mixed(t *testing.T) {
	r := &sync.Result{
		Applied: []string{"A"},
		Skipped: []string{"B", "C"},
		Errors:  []error{fmt.Errorf("oops")},
	}
	summary := r.Summary()
	if !strings.Contains(summary, "1 applied") {
		t.Errorf("missing applied count in %q", summary)
	}
	if !strings.Contains(summary, "2 skipped") {
		t.Errorf("missing skipped count in %q", summary)
	}
	if !strings.Contains(summary, "1 error(s)") {
		t.Errorf("missing error count in %q", summary)
	}
}

func TestResult_Print_InSync(t *testing.T) {
	var buf bytes.Buffer
	r := &sync.Result{}
	r.Print(&buf)
	if !strings.Contains(buf.String(), "in sync") {
		t.Errorf("expected in-sync message, got %q", buf.String())
	}
}

func TestResult_Print_Applied(t *testing.T) {
	var buf bytes.Buffer
	r := &sync.Result{Applied: []string{`FOO (default: "bar")`}}
	r.Print(&buf)
	if !strings.Contains(buf.String(), "Applied") {
		t.Errorf("expected Applied header in output")
	}
	if !strings.Contains(buf.String(), "FOO") {
		t.Errorf("expected FOO in output")
	}
}
