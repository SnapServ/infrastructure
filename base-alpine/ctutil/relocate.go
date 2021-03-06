package main

import (
	"fmt"
	"os"
	"regexp"
	"syscall"
)

var reEnvName = regexp.MustCompile(`[^A-Z0-9_]+`)

type RelocateCmd struct {
	OldPath string `kong:"required,arg,help='Path to old directory'"`
	NewPath string `kong:"required,arg,help='Path for new directory (must not exist)'"`
}

func (c *RelocateCmd) Run() error {
	// Ensure old path is symlink or exists and gather information
	oldStat, err := os.Stat(c.OldPath)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("old path [%s] does not exist: %w", c.OldPath, err)
	}
	if !oldStat.IsDir() {
		return fmt.Errorf("old path [%s] must point to a directory", c.OldPath)
	}

	// Check if new path still needs to be created
	if _, err := os.Stat(c.NewPath); os.IsNotExist(err) {
		// Create new path with same permission as old path
		if err := os.MkdirAll(c.NewPath, oldStat.Mode()); err != nil {
			return fmt.Errorf("could not create new path [%s]: %w", c.NewPath, err)
		}

		// Adjust ownership if available
		if sys, ok := oldStat.Sys().(*syscall.Stat_t); ok {
			if err := syscall.Chown(c.NewPath, int(sys.Uid), int(sys.Gid)); err != nil {
				return fmt.Errorf("could not set ownership of new path [%s]: %w", c.NewPath, err)
			}
		}
	}

	// Attempt to remove old path
	if err := os.Remove(c.OldPath); err != nil {
		return fmt.Errorf("could not remove old path [%s]: %w", c.OldPath, err)
	}

	// Symlink old path to new path
	if err := os.Symlink(c.NewPath, c.OldPath); err != nil {
		return fmt.Errorf("could not create symlink from [%s] to [%s]: %w", c.OldPath, c.NewPath, err)
	}

	return nil
}
