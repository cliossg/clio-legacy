package ssg

import (
	"context"
	"embed"
	"fmt"
	"testing"

	"github.com/hermesgen/hm"
)

func TestServiceGenerateHTMLFromContent(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockServiceRepo)
		setupCtx  func() context.Context
		setupSvc  func(*BaseService)
		wantErr   bool
	}{
		{
			name: "fails when GetAllContentWithMeta returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("db error")
			},
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), siteSlugKey, "test-site")
			},
			setupSvc: func(svc *BaseService) {},
			wantErr:  true,
		},
		{
			name: "fails when GetSections returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getSectionErr = fmt.Errorf("db error")
			},
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), siteSlugKey, "test-site")
			},
			setupSvc: func(svc *BaseService) {},
			wantErr:  true,
		},
		{
			name:      "fails without site slug in context",
			setupRepo: func(m *mockServiceRepo) {},
			setupCtx: func() context.Context {
				return context.Background()
			},
			setupSvc: func(svc *BaseService) {},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)

			cfg := hm.NewConfig()
			cfg.Set(SSGKey.SitesBasePath, t.TempDir())
			params := hm.XParams{Cfg: cfg}

			pm := NewParamManager(repo, params)

			svc := &BaseService{
				Service:  hm.NewService("test-service", params),
				assetsFS: embed.FS{},
				repo:     repo,
				pm:       pm,
			}
			tt.setupSvc(svc)

			ctx := tt.setupCtx()
			err := svc.GenerateHTMLFromContent(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateHTMLFromContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
