package fake_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/hm"
)

func TestGithubClientClone(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		ctx         context.Context
		repoURL     string
		localPath   string
		auth        hm.GitAuth
		env         []string
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful clone",
			setupFake:   func(f *fake.GithubClient) {},
			ctx:         context.Background(),
			repoURL:     "https://github.com/owner/repo.git",
			localPath:   "/tmp/repo",
			auth:        hm.GitAuth{Method: hm.AuthToken, Token: "test-token"},
			env:         []string{"GIT_TERMINAL_PROMPT=0"},
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "clone returns error",
			setupFake: func(f *fake.GithubClient) {
				f.CloneFn = func(ctx context.Context, repoURL, localPath string, auth hm.GitAuth, env []string) error {
					return errors.New("clone failed")
				}
			},
			ctx:         context.Background(),
			repoURL:     "https://github.com/owner/repo.git",
			localPath:   "/tmp/repo",
			auth:        hm.GitAuth{Method: hm.AuthSSH},
			env:         nil,
			expectedErr: errors.New("clone failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			err := f.Clone(tt.ctx, tt.repoURL, tt.localPath, tt.auth, tt.env)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.CloneCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.CloneCalls))
			}
			if tt.expectCalls > 0 {
				call := f.CloneCalls[0]
				if call.RepoURL != tt.repoURL || call.LocalPath != tt.localPath || call.Auth != tt.auth {
					t.Errorf("captured call arguments mismatch")
				}
			}
		})
	}
}

func TestGithubClientCheckout(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		branch      string
		create      bool
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful checkout",
			setupFake:   func(f *fake.GithubClient) {},
			branch:      "main",
			create:      false,
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "checkout returns error",
			setupFake: func(f *fake.GithubClient) {
				f.CheckoutFn = func(ctx context.Context, localRepoPath, branch string, create bool, env []string) error {
					return errors.New("checkout failed")
				}
			},
			branch:      "feature-branch",
			create:      true,
			expectedErr: errors.New("checkout failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			err := f.Checkout(context.Background(), "/tmp/repo", tt.branch, tt.create, nil)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.CheckoutCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.CheckoutCalls))
			}
		})
	}
}

func TestGithubClientAdd(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		pathspec    string
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful add",
			setupFake:   func(f *fake.GithubClient) {},
			pathspec:    ".",
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "add returns error",
			setupFake: func(f *fake.GithubClient) {
				f.AddFn = func(ctx context.Context, localRepoPath, pathspec string, env []string) error {
					return errors.New("add failed")
				}
			},
			pathspec:    "*.txt",
			expectedErr: errors.New("add failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			err := f.Add(context.Background(), "/tmp/repo", tt.pathspec, nil)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.AddCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.AddCalls))
			}
		})
	}
}

func TestGithubClientCommit(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		commit      hm.GitCommit
		expectedErr error
		expectedHash string
		expectCalls int
	}{
		{
			name:        "successful commit",
			setupFake:   func(f *fake.GithubClient) {},
			commit:      hm.GitCommit{Message: "test commit"},
			expectedErr: nil,
			expectedHash: "fake-commit-hash",
			expectCalls: 1,
		},
		{
			name: "commit returns error",
			setupFake: func(f *fake.GithubClient) {
				f.CommitFn = func(ctx context.Context, localRepoPath string, commit hm.GitCommit, env []string) (string, error) {
					return "", errors.New("commit failed")
				}
			},
			commit:      hm.GitCommit{Message: "test commit"},
			expectedErr: errors.New("commit failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			hash, err := f.Commit(context.Background(), "/tmp/repo", tt.commit, nil)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedErr == nil && hash != tt.expectedHash {
				t.Errorf("expected hash %s, got %s", tt.expectedHash, hash)
			}

			if len(f.CommitCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.CommitCalls))
			}
		})
	}
}

func TestGithubClientPush(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		remote      string
		branch      string
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful push",
			setupFake:   func(f *fake.GithubClient) {},
			remote:      "origin",
			branch:      "main",
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "push returns error",
			setupFake: func(f *fake.GithubClient) {
				f.PushFn = func(ctx context.Context, localRepoPath string, auth hm.GitAuth, remote, branch string, env []string) error {
					return errors.New("push failed")
				}
			},
			remote:      "origin",
			branch:      "feature",
			expectedErr: errors.New("push failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			err := f.Push(context.Background(), "/tmp/repo", hm.GitAuth{}, tt.remote, tt.branch, nil)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.PushCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.PushCalls))
			}
		})
	}
}

func TestGithubClientStatus(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		expectedOut string
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful status",
			setupFake:   func(f *fake.GithubClient) {},
			expectedOut: " M somefile.txt\n?? anotherfile.txt",
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "status returns error",
			setupFake: func(f *fake.GithubClient) {
				f.StatusFn = func(ctx context.Context, localRepoPath string, env []string) (string, error) {
					return "", errors.New("status failed")
				}
			},
			expectedErr: errors.New("status failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			out, err := f.Status(context.Background(), "/tmp/repo", nil)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedErr == nil && out != tt.expectedOut {
				t.Errorf("expected output %s, got %s", tt.expectedOut, out)
			}

			if len(f.StatusCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.StatusCalls))
			}
		})
	}
}

func TestGithubClientGitLog(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		args        []string
		expectedOut string
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful git log",
			setupFake:   func(f *fake.GithubClient) {},
			args:        []string{"-1", "--oneline"},
			expectedOut: "fake git log",
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "git log returns error",
			setupFake: func(f *fake.GithubClient) {
				f.GitLogFn = func(ctx context.Context, localRepoPath string, args []string, env []string) (string, error) {
					return "", errors.New("git log failed")
				}
			},
			args:        []string{"-10"},
			expectedErr: errors.New("git log failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			out, err := f.GitLog(context.Background(), "/tmp/repo", tt.args, nil)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedErr == nil && out != tt.expectedOut {
				t.Errorf("expected output %s, got %s", tt.expectedOut, out)
			}

			if len(f.GitLogCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.GitLogCalls))
			}
		})
	}
}
