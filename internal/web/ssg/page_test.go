package ssg

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestNewParamPage(t *testing.T) {
	param := Param{
		ID:   uuid.New(),
		Name: "Test Param",
	}

	req := httptest.NewRequest("GET", "/", nil)
	page := NewParamPage(req, param)

	if page == nil {
		t.Fatal("NewParamPage() returned nil")
	}

	if page.Param.ID != param.ID {
		t.Errorf("NewParamPage() Param.ID = %v, want %v", page.Param.ID, param.ID)
	}

	if page.Param.Name != param.Name {
		t.Errorf("NewParamPage() Param.Name = %v, want %v", page.Param.Name, param.Name)
	}
}
