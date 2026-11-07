package util

import (
	"testing"

	"github.com/rtech91/justjump/pkg/config/global"
)

func TestDetermineJumpRootIgnoresSiblingPaths(t *testing.T) {
	jumpRoots := global.JumpRoots{
		"project": {Name: "project", Root: "/tmp/project"},
	}

	exists, root := DetermineJumpRoot("/tmp/project2/subdir", jumpRoots)
	if exists {
		t.Fatalf("expected sibling path to not match jump root, got root %q", root)
	}
	if root != "" {
		t.Fatalf("expected empty root for no match, got %q", root)
	}
}

func TestDetermineJumpRootMatchesSubdirectory(t *testing.T) {
	jumpRoots := global.JumpRoots{
		"project": {Name: "project", Root: "/tmp/project"},
	}

	exists, root := DetermineJumpRoot("/tmp/project/src/app", jumpRoots)
	if !exists {
		t.Fatalf("expected subdirectory to match jump root")
	}
	if root != "/tmp/project" {
		t.Fatalf("expected jump root %q, got %q", "/tmp/project", root)
	}
}

func TestParseGitWorktrees(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected []map[string]string
	}{
		{
			name: "filters prunable record and preserves neighboring records",
			output: "worktree /repo/main\nHEAD abc123\nbranch refs/heads/main\n\n" +
				"worktree /repo/stale\nHEAD def456\nbranch refs/heads/stale\n" +
				"prunable gitdir file points to a non-existent location\n\n" +
				"worktree /repo/active\nHEAD ghi789\nbranch refs/heads/active\n",
			expected: []map[string]string{
				{"jumpPoint": "main", "fullPath": "/repo/main"},
				{"jumpPoint": "active", "fullPath": "/repo/active"},
			},
		},
		{
			name: "does not infer status from explanatory text",
			output: "worktree /repo/main\nHEAD abc123\nbranch refs/heads/main\n\n" +
				"worktree /repo/reason-only\nHEAD def456\nbranch refs/heads/reason-only\n" +
				"reason: prunable because a path is missing\nprunable-ish is not the status token\n",
			expected: []map[string]string{
				{"jumpPoint": "main", "fullPath": "/repo/main"},
				{"jumpPoint": "reason-only", "fullPath": "/repo/reason-only"},
			},
		},
		{
			name: "retains current worktree when every linked worktree is prunable",
			output: "worktree /repo/main\nHEAD abc123\nbranch refs/heads/main\n\n" +
				"worktree /repo/stale-one\nHEAD def456\nprunable reason one\n\n" +
				"worktree /repo/stale-two\nHEAD ghi789\nprunable reason two\n",
			expected: []map[string]string{
				{"jumpPoint": "main", "fullPath": "/repo/main"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGitWorktrees(tt.output)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d worktrees, got %d: %#v", len(tt.expected), len(got), got)
			}
			for i := range tt.expected {
				if got[i]["jumpPoint"] != tt.expected[i]["jumpPoint"] ||
					got[i]["fullPath"] != tt.expected[i]["fullPath"] {
					t.Errorf("worktree %d: expected %#v, got %#v", i, tt.expected[i], got[i])
				}
			}
		})
	}
}
