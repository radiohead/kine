package logstructured

import "context"

// TTLExpirer
type TTLExpirer struct{}

// NewTTLExpirer
func NewTTLExpirer() *TTLExpirer {
	return &TTLExpirer{}
}

// Start
func (t *TTLExpirer) Start(ctx context.Context) error {
	return nil
}
