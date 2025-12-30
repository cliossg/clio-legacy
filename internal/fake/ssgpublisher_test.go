package fake_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/ssg"
)

func TestSSGPublisherPublish(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SSGPublisher)
		ctx         context.Context
		cfg         ssg.PublisherConfig
		sourceDir   string
		expectedURL string
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful publish",
			setupFake:   func(f *fake.SSGPublisher) {},
			ctx:         context.Background(),
			cfg:         ssg.PublisherConfig{RepoURL: "http://example.com/repo"},
			sourceDir:   "/tmp/src",
			expectedURL: "fake-commit-url",
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "publish returns error",
			setupFake: func(f *fake.SSGPublisher) {
				f.PublishFn = func(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (string, error) {
					return "", errors.New("publish failed")
				}
			},
			ctx:         context.Background(),
			cfg:         ssg.PublisherConfig{RepoURL: "http://example.com/repo"},
			sourceDir:   "/tmp/src",
			expectedURL: "",
			expectedErr: errors.New("publish failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSSGPublisher()
			tt.setupFake(f)

			url, err := f.Publish(tt.ctx, tt.cfg, tt.sourceDir)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if url != tt.expectedURL {
				t.Errorf("expected URL %q, got %q", tt.expectedURL, url)
			}

			if len(f.PublishCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.PublishCalls))
			}
			if tt.expectCalls > 0 {
				call := f.PublishCalls[0]
				if call.Cfg.RepoURL != tt.cfg.RepoURL || call.SourceDir != tt.sourceDir {
					t.Errorf("captured call arguments mismatch")
				}
			}
		})
	}
}

func TestSSGPublisherValidate(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SSGPublisher)
		cfg         ssg.PublisherConfig
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful validation",
			setupFake:   func(f *fake.SSGPublisher) {},
			cfg:         ssg.PublisherConfig{RepoURL: "http://example.com/repo"},
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "validation returns error",
			setupFake: func(f *fake.SSGPublisher) {
				f.ValidateFn = func(cfg ssg.PublisherConfig) error {
					return errors.New("validation failed")
				}
			},
			cfg:         ssg.PublisherConfig{RepoURL: ""},
			expectedErr: errors.New("validation failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSSGPublisher()
			tt.setupFake(f)

			err := f.Validate(tt.cfg)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.ValidateCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.ValidateCalls))
			}
			if tt.expectCalls > 0 {
				call := f.ValidateCalls[0]
				if call.Cfg.RepoURL != tt.cfg.RepoURL {
					t.Errorf("captured call arguments mismatch")
				}
			}
		})
	}
}

func TestSSGPublisherPlan(t *testing.T) {
	tests := []struct {
		name            string
		setupFake       func(f *fake.SSGPublisher)
		ctx             context.Context
		cfg             ssg.PublisherConfig
		sourceDir       string
		expectedSummary string
		expectedErr     error
		expectCalls     int
	}{
		{
			name:            "successful plan",
			setupFake:       func(f *fake.SSGPublisher) {},
			ctx:             context.Background(),
			cfg:             ssg.PublisherConfig{RepoURL: "http://example.com/repo"},
			sourceDir:       "/tmp/src",
			expectedSummary: "fake plan",
			expectedErr:     nil,
			expectCalls:     1,
		},
		{
			name: "plan returns error",
			setupFake: func(f *fake.SSGPublisher) {
				f.PlanFn = func(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (ssg.PlanReport, error) {
					return ssg.PlanReport{}, errors.New("plan failed")
				}
			},
			ctx:             context.Background(),
			cfg:             ssg.PublisherConfig{RepoURL: "http://example.com/repo"},
			sourceDir:       "/tmp/src",
			expectedSummary: "",
			expectedErr:     errors.New("plan failed"),
			expectCalls:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSSGPublisher()
			tt.setupFake(f)

			report, err := f.Plan(tt.ctx, tt.cfg, tt.sourceDir)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if report.Summary != tt.expectedSummary {
				t.Errorf("expected summary %q, got %q", tt.expectedSummary, report.Summary)
			}

			if len(f.PlanCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.PlanCalls))
			}
			if tt.expectCalls > 0 {
				call := f.PlanCalls[0]
				if call.Cfg.RepoURL != tt.cfg.RepoURL || call.SourceDir != tt.sourceDir {
					t.Errorf("captured call arguments mismatch")
				}
			}
		})
	}
}
