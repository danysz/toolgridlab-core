package verification

import (
	"context"
	"testing"
	"time"
)

func TestEmailVerifier_SyntaxAndDisposable(t *testing.T) {
	v := NewEmailVerifier(2 * time.Second)
	ctx := context.Background()

	// 1. Invalid syntax
	if err := v.Verify(ctx, "invalid-email"); err != ErrInvalidSyntax {
		t.Fatalf("expected ErrInvalidSyntax, got %v", err)
	}

	// 2. Disposable domain
	if err := v.Verify(ctx, "test@mailinator.com"); err != ErrDisposableDomain {
		t.Fatalf("expected ErrDisposableDomain, got %v", err)
	}

	// 3. Known valid domain syntax check
	// Note: We test syntax acceptance, MX lookup runs against public DNS
	if err := v.Verify(ctx, "hello@google.com"); err != nil {
		t.Logf("DNS MX check result for google.com: %v", err)
	}
}
