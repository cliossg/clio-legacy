package ssg

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

type mockPublisherWithTracking struct {
	publishCalled  bool
	publishErr     error
	lastCfg        PublisherConfig
	lastSourceDir  string
	commitURL      string
	planCalled     bool
	planErr        error
	planReport     PlanReport
}

func (m *mockPublisherWithTracking) Validate(cfg PublisherConfig) error {
	return nil
}

func (m *mockPublisherWithTracking) Publish(ctx context.Context, cfg PublisherConfig, sourceDir string) (string, error) {
	m.publishCalled = true
	m.lastCfg = cfg
	m.lastSourceDir = sourceDir
	if m.publishErr != nil {
		return "", m.publishErr
	}
	return m.commitURL, nil
}

func (m *mockPublisherWithTracking) Plan(ctx context.Context, cfg PublisherConfig, sourceDir string) (PlanReport, error) {
	m.planCalled = true
	m.lastCfg = cfg
	m.lastSourceDir = sourceDir
	if m.planErr != nil {
		return PlanReport{}, m.planErr
	}
	return m.planReport, nil
}

func TestServiceGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func() context.Context
		setupRepo func(*mockServiceRepo)
		wantErr   bool
	}{
		{
			name: "fails without site slug in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			setupRepo: func(m *mockServiceRepo) {},
			wantErr:   true,
		},
		{
			name: "fails when GetAllContentWithMeta returns error",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), siteSlugKey, "test-site")
			},
			setupRepo: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)

			cfg := hm.NewConfig()
			gen := NewGenerator(hm.XParams{Cfg: cfg})

			svc := &BaseService{
				Service: hm.NewService("test-service", hm.XParams{Cfg: cfg}),
				repo:    repo,
				gen:     gen,
			}
			ctx := tt.setupCtx()

			err := svc.GenerateMarkdown(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePublish(t *testing.T) {
	tests := []struct {
		name          string
		setupRepo     func(*mockServiceRepo)
		setupPub      func(*mockPublisherWithTracking)
		commitMsg     string
		wantErr       bool
		checkPub      func(*testing.T, *mockPublisherWithTracking)
		wantCommitURL string
	}{
		{
			name: "publishes successfully with all params",
			setupRepo: func(m *mockServiceRepo) {
				m.paramsByRef[SSGKey.PublishRepoURL] = Param{ID: uuid.New(), RefKey: SSGKey.PublishRepoURL, Value: "https://github.com/user/repo"}
				m.paramsByRef[SSGKey.PublishBranch] = Param{ID: uuid.New(), RefKey: SSGKey.PublishBranch, Value: "gh-pages"}
				m.paramsByRef[SSGKey.PublishAuthToken] = Param{ID: uuid.New(), RefKey: SSGKey.PublishAuthToken, Value: "token123"}
				m.paramsByRef[SSGKey.PublishCommitUserName] = Param{ID: uuid.New(), RefKey: SSGKey.PublishCommitUserName, Value: "Test User"}
				m.paramsByRef[SSGKey.PublishCommitUserEmail] = Param{ID: uuid.New(), RefKey: SSGKey.PublishCommitUserEmail, Value: "test@example.com"}
				m.paramsByRef[SSGKey.PublishCommitMessage] = Param{ID: uuid.New(), RefKey: SSGKey.PublishCommitMessage, Value: "Default commit message"}
			},
			setupPub: func(m *mockPublisherWithTracking) {
				m.commitURL = "https://github.com/user/repo/commit/abc123"
			},
			commitMsg:     "",
			wantErr:       false,
			wantCommitURL: "https://github.com/user/repo/commit/abc123",
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if !m.publishCalled {
					t.Error("Expected Publish to be called")
				}
				if m.lastCfg.RepoURL != "https://github.com/user/repo" {
					t.Errorf("Expected RepoURL 'https://github.com/user/repo', got %q", m.lastCfg.RepoURL)
				}
				if m.lastCfg.Branch != "gh-pages" {
					t.Errorf("Expected Branch 'gh-pages', got %q", m.lastCfg.Branch)
				}
				if m.lastCfg.Auth.Token != "token123" {
					t.Errorf("Expected Token 'token123', got %q", m.lastCfg.Auth.Token)
				}
				if m.lastCfg.CommitAuthor.UserName != "Test User" {
					t.Errorf("Expected UserName 'Test User', got %q", m.lastCfg.CommitAuthor.UserName)
				}
				if m.lastCfg.CommitAuthor.UserEmail != "test@example.com" {
					t.Errorf("Expected UserEmail 'test@example.com', got %q", m.lastCfg.CommitAuthor.UserEmail)
				}
				if m.lastCfg.CommitAuthor.Message != "Default commit message" {
					t.Errorf("Expected commit message 'Default commit message', got %q", m.lastCfg.CommitAuthor.Message)
				}
			},
		},
		{
			name: "overrides commit message when provided",
			setupRepo: func(m *mockServiceRepo) {
				m.paramsByRef[SSGKey.PublishCommitMessage] = Param{ID: uuid.New(), RefKey: SSGKey.PublishCommitMessage, Value: "Default message"}
			},
			setupPub: func(m *mockPublisherWithTracking) {
				m.commitURL = "https://github.com/user/repo/commit/def456"
			},
			commitMsg:     "Custom commit message",
			wantErr:       false,
			wantCommitURL: "https://github.com/user/repo/commit/def456",
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if m.lastCfg.CommitAuthor.Message != "Custom commit message" {
					t.Errorf("Expected commit message 'Custom commit message', got %q", m.lastCfg.CommitAuthor.Message)
				}
			},
		},
		{
			name:      "uses default values when params not found",
			setupRepo: func(m *mockServiceRepo) {},
			setupPub: func(m *mockPublisherWithTracking) {
				m.commitURL = ""
			},
			commitMsg: "",
			wantErr:   false,
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if m.lastCfg.RepoURL != "" {
					t.Errorf("Expected empty RepoURL, got %q", m.lastCfg.RepoURL)
				}
				if m.lastCfg.Branch != "" {
					t.Errorf("Expected empty Branch, got %q", m.lastCfg.Branch)
				}
			},
		},
		{
			name:      "fails when publisher returns error",
			setupRepo: func(m *mockServiceRepo) {},
			setupPub: func(m *mockPublisherWithTracking) {
				m.publishErr = fmt.Errorf("publish failed")
			},
			commitMsg: "",
			wantErr:   true,
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if !m.publishCalled {
					t.Error("Expected Publish to be called")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			pub := &mockPublisherWithTracking{}
			tt.setupPub(pub)

			cfg := hm.NewConfig()
			cfg.Set(SSGKey.HTMLPath, "_workspace/documents/html")
			pm := NewParamManager(repo, hm.XParams{Cfg: cfg})

			svc := &BaseService{
				Service: hm.NewService("test-service", hm.XParams{Cfg: cfg}),
				pub:     pub,
				pm:      pm,
			}

			commitURL, err := svc.Publish(context.Background(), tt.commitMsg)

			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && commitURL != tt.wantCommitURL {
				t.Errorf("Publish() commitURL = %q, want %q", commitURL, tt.wantCommitURL)
			}

			if tt.checkPub != nil {
				tt.checkPub(t, pub)
			}
		})
	}
}

func TestServicePlan(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockServiceRepo)
		setupPub  func(*mockPublisherWithTracking)
		wantErr   bool
		checkPub  func(*testing.T, *mockPublisherWithTracking)
	}{
		{
			name: "plans successfully with all params",
			setupRepo: func(m *mockServiceRepo) {
				m.paramsByRef[SSGKey.PublishRepoURL] = Param{ID: uuid.New(), RefKey: SSGKey.PublishRepoURL, Value: "https://github.com/user/repo"}
				m.paramsByRef[SSGKey.PublishBranch] = Param{ID: uuid.New(), RefKey: SSGKey.PublishBranch, Value: "gh-pages"}
				m.paramsByRef[SSGKey.PublishAuthToken] = Param{ID: uuid.New(), RefKey: SSGKey.PublishAuthToken, Value: "token123"}
			},
			setupPub: func(m *mockPublisherWithTracking) {
				m.planReport = PlanReport{
					Summary: "Would publish 5 files",
					Added:   []string{"file1.html", "file2.html", "file3.html"},
				}
			},
			wantErr: false,
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if !m.planCalled {
					t.Error("Expected Plan to be called")
				}
				if m.publishCalled {
					t.Error("Expected Publish not to be called")
				}
				if m.lastCfg.RepoURL != "https://github.com/user/repo" {
					t.Errorf("Expected RepoURL 'https://github.com/user/repo', got %q", m.lastCfg.RepoURL)
				}
			},
		},
		{
			name:      "uses default values when params not found",
			setupRepo: func(m *mockServiceRepo) {},
			setupPub: func(m *mockPublisherWithTracking) {
				m.planReport = PlanReport{Summary: "No changes"}
			},
			wantErr: false,
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if !m.planCalled {
					t.Error("Expected Plan to be called")
				}
			},
		},
		{
			name:      "fails when publisher returns error",
			setupRepo: func(m *mockServiceRepo) {},
			setupPub: func(m *mockPublisherWithTracking) {
				m.planErr = fmt.Errorf("plan failed")
			},
			wantErr: true,
			checkPub: func(t *testing.T, m *mockPublisherWithTracking) {
				if !m.planCalled {
					t.Error("Expected Plan to be called")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			pub := &mockPublisherWithTracking{}
			tt.setupPub(pub)

			cfg := hm.NewConfig()
			cfg.Set(SSGKey.HTMLPath, "_workspace/documents/html")
			pm := NewParamManager(repo, hm.XParams{Cfg: cfg})

			svc := &BaseService{
				Service: hm.NewService("test-service", hm.XParams{Cfg: cfg}),
				pub:     pub,
				pm:      pm,
			}

			report, err := svc.Plan(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Plan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && report.Summary != pub.planReport.Summary {
				t.Errorf("Plan() report.Summary = %q, want %q", report.Summary, pub.planReport.Summary)
			}

			if tt.checkPub != nil {
				tt.checkPub(t, pub)
			}
		})
	}
}
