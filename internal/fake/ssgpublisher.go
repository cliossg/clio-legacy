package fake

import (
	"context"

	"github.com/hermesgen/clio/internal/feat/ssg"
)

// SSGPublisher is a fake implementation of ssg.Publisher for testing.
type SSGPublisher struct {
	// Expected results
	ValidateFn func(cfg ssg.PublisherConfig) error
	PublishFn  func(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (string, error)
	PlanFn     func(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (ssg.PlanReport, error)

	// Captured arguments
	ValidateCalls []struct{ Cfg ssg.PublisherConfig }
	PublishCalls  []struct {
		Ctx       context.Context
		Cfg       ssg.PublisherConfig
		SourceDir string
	}
	PlanCalls []struct {
		Ctx       context.Context
		Cfg       ssg.PublisherConfig
		SourceDir string
	}
}

// NewSSGPublisher creates a new fake SSGPublisher.
func NewSSGPublisher() *SSGPublisher {
	return &SSGPublisher{}
}

func (f *SSGPublisher) Validate(cfg ssg.PublisherConfig) error {
	f.ValidateCalls = append(f.ValidateCalls, struct{ Cfg ssg.PublisherConfig }{Cfg: cfg})
	if f.ValidateFn != nil {
		return f.ValidateFn(cfg)
	}
	return nil
}

func (f *SSGPublisher) Publish(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (string, error) {
	f.PublishCalls = append(f.PublishCalls, struct {
		Ctx       context.Context
		Cfg       ssg.PublisherConfig
		SourceDir string
	}{Ctx: ctx, Cfg: cfg, SourceDir: sourceDir})
	if f.PublishFn != nil {
		return f.PublishFn(ctx, cfg, sourceDir)
	}
	return "fake-commit-url", nil
}

func (f *SSGPublisher) Plan(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (ssg.PlanReport, error) {
	f.PlanCalls = append(f.PlanCalls, struct {
		Ctx       context.Context
		Cfg       ssg.PublisherConfig
		SourceDir string
	}{Ctx: ctx, Cfg: cfg, SourceDir: sourceDir})
	if f.PlanFn != nil {
		return f.PlanFn(ctx, cfg, sourceDir)
	}
	return ssg.PlanReport{Summary: "fake plan"}, nil
}
