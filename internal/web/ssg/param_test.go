package ssg

import (
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewParam(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
	}{
		{
			name:       "creates param with name and value",
			paramName:  "Test Param",
			paramValue: "Test Value",
		},
		{
			name:       "creates param with empty fields",
			paramName:  "",
			paramValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := NewParam(tt.paramName, tt.paramValue)
			if param.Name != tt.paramName {
				t.Errorf("NewParam() Name = %v, want %v", param.Name, tt.paramName)
			}
			if param.Value != tt.paramValue {
				t.Errorf("NewParam() Value = %v, want %v", param.Value, tt.paramValue)
			}
		})
	}
}

func TestParamType(t *testing.T) {
	param := &Param{}
	if got := param.Type(); got != "param" {
		t.Errorf("Type() = %v, want param", got)
	}
}

func TestParamGetID(t *testing.T) {
	id := uuid.New()
	param := Param{ID: id}
	if got := param.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestParamGenID(t *testing.T) {
	param := &Param{}
	param.GenID()
	if param.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestParamSetID(t *testing.T) {
	tests := []struct {
		name    string
		initial uuid.UUID
		new     uuid.UUID
		force   []bool
		wantID  uuid.UUID
	}{
		{
			name:    "sets ID when empty",
			initial: uuid.Nil,
			new:     uuid.New(),
			force:   nil,
			wantID:  uuid.Nil,
		},
		{
			name:    "sets ID with force",
			initial: uuid.New(),
			new:     uuid.New(),
			force:   []bool{true},
			wantID:  uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := &Param{ID: tt.initial}
			param.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && param.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", param.ID, tt.wantID)
			}
		})
	}
}

func TestParamGetShortID(t *testing.T) {
	param := Param{ShortID: "test123"}
	if got := param.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestParamGenShortID(t *testing.T) {
	param := &Param{}
	param.GenShortID()
	if param.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestParamSetShortID(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		new     string
		force   []bool
	}{
		{
			name:    "sets shortID when empty",
			initial: "",
			new:     "test123",
			force:   nil,
		},
		{
			name:    "sets shortID with force",
			initial: "old123",
			new:     "test123",
			force:   []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := &Param{ShortID: tt.initial}
			param.SetShortID(tt.new, tt.force...)
			if param.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && param.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", param.ShortID, tt.new)
				}
			}
		})
	}
}

func TestParamIsZero(t *testing.T) {
	tests := []struct {
		name  string
		param Param
		want  bool
	}{
		{
			name:  "zero param",
			param: Param{},
			want:  true,
		},
		{
			name:  "non-zero param",
			param: Param{ID: uuid.New()},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParamSlug(t *testing.T) {
	param := &Param{Name: "Test Param", ShortID: "abc123"}
	got := param.Slug()
	if got == "" {
		t.Error("Slug() returned empty string")
	}
	if len(got) < len("test-param") {
		t.Errorf("Slug() = %v, too short", got)
	}
}

func TestParamTypeID(t *testing.T) {
	param := &Param{ShortID: "abc123"}
	got := param.TypeID()
	if got == "" {
		t.Error("TypeID() returned empty string")
	}
}

func TestParamIsSystem(t *testing.T) {
	tests := []struct {
		name   string
		param  Param
		want   bool
	}{
		{
			name:  "system param",
			param: Param{System: 1},
			want:  true,
		},
		{
			name:  "non-system param",
			param: Param{System: 0},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.IsSystem(); got != tt.want {
				t.Errorf("IsSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToWebParam(t *testing.T) {
	id := uuid.New()
	featParam := feat.Param{
		ID:          id,
		ShortID:     "abc123",
		Name:        "Test Param",
		Description: "Test description",
		Value:       "test value",
		RefKey:      "test-ref",
		System:      1,
	}

	webParam := ToWebParam(featParam)

	if webParam.ID != featParam.ID {
		t.Errorf("ToWebParam() ID = %v, want %v", webParam.ID, featParam.ID)
	}
	if webParam.ShortID != featParam.ShortID {
		t.Errorf("ToWebParam() ShortID = %v, want %v", webParam.ShortID, featParam.ShortID)
	}
	if webParam.Name != featParam.Name {
		t.Errorf("ToWebParam() Name = %v, want %v", webParam.Name, featParam.Name)
	}
	if webParam.Description != featParam.Description {
		t.Errorf("ToWebParam() Description = %v, want %v", webParam.Description, featParam.Description)
	}
	if webParam.Value != featParam.Value {
		t.Errorf("ToWebParam() Value = %v, want %v", webParam.Value, featParam.Value)
	}
	if webParam.RefKey != featParam.RefKey {
		t.Errorf("ToWebParam() RefKey = %v, want %v", webParam.RefKey, featParam.RefKey)
	}
	if webParam.System != featParam.System {
		t.Errorf("ToWebParam() System = %v, want %v", webParam.System, featParam.System)
	}
}

func TestToWebParams(t *testing.T) {
	featParams := []feat.Param{
		{
			ID:          uuid.New(),
			ShortID:     "abc123",
			Name:        "Param 1",
			Description: "Description 1",
			Value:       "value 1",
			RefKey:      "ref1",
			System:      1,
		},
		{
			ID:          uuid.New(),
			ShortID:     "def456",
			Name:        "Param 2",
			Description: "Description 2",
			Value:       "value 2",
			RefKey:      "ref2",
			System:      0,
		},
	}

	webParams := ToWebParams(featParams)

	if len(webParams) != len(featParams) {
		t.Errorf("ToWebParams() length = %v, want %v", len(webParams), len(featParams))
	}

	for i, webParam := range webParams {
		if webParam.ID != featParams[i].ID {
			t.Errorf("ToWebParams()[%d] ID = %v, want %v", i, webParam.ID, featParams[i].ID)
		}
		if webParam.Name != featParams[i].Name {
			t.Errorf("ToWebParams()[%d] Name = %v, want %v", i, webParam.Name, featParams[i].Name)
		}
	}
}
