package ssg

import (
	"context"

	"github.com/adrianpk/clio/internal/am"
)

type ParamManager struct {
	am.Core
	repo Repo
}


// NewParamManagerWithParams creates a ParamManager with XParams.
func NewParamManager(repo Repo, params am.XParams) *ParamManager {
	core := am.NewCoreWithParams("param-manager", params)
	return &ParamManager{
		Core: core,
		repo: repo,
	}
}

func (pm *ParamManager) FindParam(ctx context.Context, ref string) (Param, error) {
	return pm.findParamByName(ctx, ref)

}

func (pm *ParamManager) findParamByName(ctx context.Context, name string) (Param, error) {
	// TODO: This is a simple wrapper for now but a simple caching strategy can be added here.
	return pm.repo.GetParamByName(ctx, name)
}

func (pm *ParamManager) FindParamByRef(ctx context.Context, refKey string) (Param, error) {
	// TODO: This is a simple wrapper for now but a simple caching strategy can be added here.
	return pm.repo.GetParamByRefKey(ctx, refKey)
}

func (pm *ParamManager) Get(ctx context.Context, refKey string, defVal string) string {
	param, err := pm.repo.GetParamByRefKey(ctx, refKey)
	if err == nil && !param.IsZero() {
		return param.Value
	}

	// Fallback to configuration
	return pm.Cfg().StrValOrDef(refKey, defVal)
}
