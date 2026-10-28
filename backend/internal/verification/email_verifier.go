package verification

import (
	"context"
	"errors"
	"net"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidSyntax     = errors.New("invalid email syntax")
	ErrDisposableDomain  = errors.New("disposable and temporary email domains are not allowed")
	ErrNoActiveMXRecords = errors.New("domain has no active MX records or valid mail exchangers")
)

// Standard RFC 5322 compliant regex for basic syntax validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

// Comprehensive blacklist of known disposable and temporary email domains
var disposableDomains = map[string]struct{}{
	"mailinator.com":      {},
	"tempmail.com":        {},
	"10minutemail.com":    {},
	"guerrillamail.com":   {},
	"guerrillamail.net":   {},
	"guerrillamail.org":   {},
	"throwawaymail.com":   {},
	"yopmail.com":         {},
	"yopmail.net":         {},
	"trashmail.com":       {},
	"trashmail.net":       {},
	"trashmail.org":       {},
	"sharklasers.com":     {},
	"getairmail.com":      {},
	"dispostable.com":     {},
	"fakeinbox.com":       {},
	"mytempemail.com":     {},
	"mohmal.com":          {},
	"temp-mail.org":       {},
	"temp-mail.io":        {},
	"nada.ltd":            {},
	"inboxkitten.com":     {},
	"burnermail.io":       {},
	"generator.email":     {},
	"crazymailing.com":    {},
	"dropmail.me":         {},
	"crazymail.com":       {},
	"tempm.com":           {},
	"minutemail.com":      {},
	"emailondeck.com":     {},
	"crazymailing.org":    {},
}

type EmailVerifier struct {
	resolver *net.Resolver
	timeout  time.Duration
}

func NewEmailVerifier(timeout time.Duration) *EmailVerifier {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &EmailVerifier{
		resolver: net.DefaultResolver,
		timeout:  timeout,
	}
}

// Verify performs:
// 1. Syntax regex validation
// 2. Disposable domain blacklist check
// 3. DNS net.LookupMX verification
func (v *EmailVerifier) Verify(ctx context.Context, email string) error {
	trimmed := strings.TrimSpace(strings.ToLower(email))
	if trimmed == "" || len(trimmed) > 254 {
		return ErrInvalidSyntax
	}

	// 1. Syntax Regex Check
	if !emailRegex.MatchString(trimmed) {
		return ErrInvalidSyntax
	}

	parts := strings.Split(trimmed, "@")
	if len(parts) != 2 {
		return ErrInvalidSyntax
	}
	domain := parts[1]

	// 2. Disposable Domain Blacklist Check
	if _, isDisposable := disposableDomains[domain]; isDisposable {
		return ErrDisposableDomain
	}

	// 3. DNS LookupMX Check with context timeout
	lookupCtx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	mxRecords, err := v.resolver.LookupMX(lookupCtx, domain)
	if err == nil && len(mxRecords) > 0 {
		return nil
	}

	// Fallback per RFC 5321: If no MX records are published, check for direct host A/AAAA records
	ips, ipErr := v.resolver.LookupIPAddr(lookupCtx, domain)
	if ipErr == nil && len(ips) > 0 {
		return nil
	}

	return ErrNoActiveMXRecords
}
