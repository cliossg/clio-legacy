package ssg

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetSiteSlugFromContext(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantSlug  string
		wantFound bool
	}{
		{
			name:      "returns slug when present",
			ctx:       context.WithValue(context.Background(), siteSlugKey, "my-site"),
			wantSlug:  "my-site",
			wantFound: true,
		},
		{
			name:      "returns empty when not present",
			ctx:       context.Background(),
			wantSlug:  "",
			wantFound: false,
		},
		{
			name:      "returns empty when wrong type",
			ctx:       context.WithValue(context.Background(), siteSlugKey, 123),
			wantSlug:  "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlug, gotFound := GetSiteSlugFromContext(tt.ctx)

			if gotSlug != tt.wantSlug {
				t.Errorf("GetSiteSlugFromContext() slug = %v, want %v", gotSlug, tt.wantSlug)
			}

			if gotFound != tt.wantFound {
				t.Errorf("GetSiteSlugFromContext() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestGetSiteIDFromContext(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		wantID    uuid.UUID
		wantFound bool
	}{
		{
			name:      "returns ID when present",
			ctx:       context.WithValue(context.Background(), siteIDKey, testID),
			wantID:    testID,
			wantFound: true,
		},
		{
			name:      "returns nil UUID when not present",
			ctx:       context.Background(),
			wantID:    uuid.Nil,
			wantFound: false,
		},
		{
			name:      "returns nil UUID when wrong type",
			ctx:       context.WithValue(context.Background(), siteIDKey, "not-a-uuid"),
			wantID:    uuid.Nil,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotFound := GetSiteIDFromContext(tt.ctx)

			if gotID != tt.wantID {
				t.Errorf("GetSiteIDFromContext() ID = %v, want %v", gotID, tt.wantID)
			}

			if gotFound != tt.wantFound {
				t.Errorf("GetSiteIDFromContext() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestRequireSiteSlug(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name:    "returns slug when present",
			ctx:     context.WithValue(context.Background(), siteSlugKey, "my-site"),
			want:    "my-site",
			wantErr: false,
		},
		{
			name:    "returns error when not present",
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "returns error when empty string",
			ctx:     context.WithValue(context.Background(), siteSlugKey, ""),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequireSiteSlug(tt.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("RequireSiteSlug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("RequireSiteSlug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequireSiteID(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name    string
		ctx     context.Context
		want    uuid.UUID
		wantErr bool
	}{
		{
			name:    "returns ID when present",
			ctx:     context.WithValue(context.Background(), siteIDKey, testID),
			want:    testID,
			wantErr: false,
		},
		{
			name:    "returns error when not present",
			ctx:     context.Background(),
			want:    uuid.Nil,
			wantErr: true,
		},
		{
			name:    "returns error when nil UUID",
			ctx:     context.WithValue(context.Background(), siteIDKey, uuid.Nil),
			want:    uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequireSiteID(tt.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("RequireSiteID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("RequireSiteID() = %v, want %v", got, tt.want)
			}
		})
	}
}
