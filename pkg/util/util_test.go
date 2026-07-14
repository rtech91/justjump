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
