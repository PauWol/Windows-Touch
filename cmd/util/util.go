package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"moul.io/banner"
)

func Banner(text string) string {
	return banner.Inline(text)
}

func Intro() string {
	var intro string

	intro += Banner("touch")
	intro += "\n\n"
	intro += "An advanced touch command with extra features and flexibility.\n\n"
	intro += "Supports custom timestamps, overwrite protection, and file attribute control.\n\n"
	intro += "https://github.com/pauwol/touch\n\n"
	intro += "Usage:\n"
	intro += "touch [flags] [file(s)]\n\n"
	intro += "Flags:\n"
	intro += "-f, --force\t\tOverwrite existing files\n"
	intro += "-d, --directory\t\tCreate directories instead of files\n"
	intro += "-t, --timestamp <timestamp>\tSet the creation timestamp for the file (YYYY-MM-DD HH:MM:SS)\n"
	intro += "-p, --permissions <level>\tSet file permissions (USER or ADMIN)\n"
	intro += "-u, --update\t\tUpdate timestamps and permissions without creating new files\n"
	intro += "-h, --help\t\tShows this help message\n"

	return intro
}

// Determines if the path looks like a file (has an extension)
func LooksLikeFile(path string) bool {
	return filepath.Ext(path) != ""
}

type Path struct {
	Path string
}

func (f *Path) Exists() bool {
	_, err := os.Stat(f.Path)
	return err == nil
}

func (f *Path) IsFile() bool {
	info, err := os.Stat(f.Path)
	return err == nil && !info.IsDir()
}

func (f *Path) IsDir() bool {
	info, err := os.Stat(f.Path)
	return err == nil && info.IsDir()
}

// Create() only creates if the path doesn't exist
func (f *Path) Create() error {
	if f.Exists() {
		fmt.Println("Path already exists.")
		return errors.New("path already exists")
	}

	if LooksLikeFile(f.Path) {
		// Ensure parent directory exists before creating file
		parentDir := filepath.Dir(f.Path)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Create the file
		if err := os.WriteFile(f.Path, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	} else {
		// Create directory
		if err := os.MkdirAll(f.Path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return nil
}

// ForceCreate() ensures the path is created, replacing any existing file/directory
func (f *Path) ForceCreate() error {
	if f.Exists() {
		// Remove existing file/directory before creation
		if err := os.RemoveAll(f.Path); err != nil {
			return fmt.Errorf("failed to remove existing path: %w", err)
		}
	}

	if LooksLikeFile(f.Path) {
		// Ensure parent directory exists before creating file
		parentDir := filepath.Dir(f.Path)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Create the file
		if err := os.WriteFile(f.Path, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	} else {
		// Create directory
		if err := os.MkdirAll(f.Path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return nil
}

// CreateDir() only creates directories, even if the name looks like a file
func (f *Path) CreateDir() error {
	if f.Exists() {
		fmt.Println("Path already exists.")
		return errors.New("path already exists")
	}

	// Force path to be treated as a directory
	if err := os.MkdirAll(f.Path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ModifyTimestamps updates the file's timestamp.
// If only a time is provided, the file's original date is retained.
// If only a date is provided, the file's original time is retained.
// Otherwise, if both are provided, both are updated.
func (f *Path) ModifyTimestamps(timestamp string) error {
	// Ensure the file exists.
	if !f.Exists() {
		return fmt.Errorf("path does not exist: %s", f.Path)
	}

	// Get the file's current modification time.
	info, err := os.Stat(f.Path)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	origTime := info.ModTime()

	// Define possible layouts.
	formats := []string{
		"2006-01-02 15:04:05", // full date and time
		"02.01.2006 15:04",    // date and time (dot-separated)
		"02-01-2006 15:04",    // date and time (dash-separated date)
		"2006-01-02",          // date only (ISO)
		"02.01.2006",          // date only (dot-separated)
		"02-01-2006",          // date only (dash-separated)
		"15:04",               // time only (colon-separated)
		"15-04",               // time only (dash-separated)
	}

	var parsedTime time.Time
	var usedFormat string
	var parseErr error

	// Try to parse the input timestamp using each format.
	for _, format := range formats {
		parsedTime, parseErr = time.Parse(format, timestamp)
		if parseErr == nil {
			usedFormat = format
			break
		}
	}
	if parseErr != nil {
		return fmt.Errorf("invalid timestamp format: tried formats %v", formats)
	}

	var newModTime time.Time
	switch usedFormat {
	// For time-only input: keep the file's original date.
	case "15:04", "15-04":
		newModTime = time.Date(
			origTime.Year(), origTime.Month(), origTime.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
			origTime.Location(),
		)
	// For date-only input: keep the file's original time.
	case "2006-01-02", "02.01.2006", "02-01-2006":
		newModTime = time.Date(
			parsedTime.Year(), parsedTime.Month(), parsedTime.Day(),
			origTime.Hour(), origTime.Minute(), origTime.Second(), origTime.Nanosecond(),
			origTime.Location(),
		)
	// Otherwise, use the fully provided date and time.
	default:
		newModTime = parsedTime
	}

	// Windows-specific update using syscall.
	if err := updateTimestampsWindows(f.Path, newModTime); err != nil {
		return fmt.Errorf("failed to update timestamps: %w", err)
	}

	fmt.Printf("Updated timestamps for %s to %v\n", f.Path, newModTime)
	return nil
}

// updateTimestampsWindows updates the file timestamps on Windows using syscall.SetFileTime.
func updateTimestampsWindows(path string, modTime time.Time) error {
	// Open file with write access.
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Convert modTime to Windows file time format.
	nanoTime := modTime.UnixNano()
	winTime := syscall.NsecToFiletime(nanoTime)

	// Update creation, access, and modification times.
	err = syscall.SetFileTime(syscall.Handle(file.Fd()), &winTime, &winTime, &winTime)
	if err != nil {
		return fmt.Errorf("failed to set Windows file time: %w", err)
	}

	return nil
}

// ModifyPermissions sets file permissions based on "USER" or "ADMIN".
// Works on both Windows & Linux/macOS.
func (f *Path) ModifyPermissions(level string) error {
	if !f.Exists() {
		return fmt.Errorf("path does not exist: %s", f.Path)
	}

	// Linux/macOS: Use chmod
	if runtime.GOOS != "windows" {
		var perm os.FileMode

		switch level {
		case "USER":
			perm = 0644 // Owner: read/write, Others: read
		case "ADMIN":
			perm = 0600 // Owner: read/write, Others: no access
		default:
			return fmt.Errorf("invalid permission level: %s", level)
		}

		// Apply permissions
		if err := os.Chmod(f.Path, perm); err != nil {
			return fmt.Errorf("failed to change permissions: %w", err)
		}
		fmt.Printf("Updated permissions for: %s (Level: %s)\n", f.Path, level)
		return nil
	}

	// Windows: Use `icacls` to set permissions
	return setWindowsPermissions(f.Path, level)
}

// setWindowsPermissions applies Windows ACLs for "USER" or "ADMIN"
func setWindowsPermissions(path string, level string) error {
	var cmd *exec.Cmd

	switch level {
	case "USER":
		// Grant full control to the user, readable for everyone
		cmd = exec.Command("icacls", path, "/grant", "Everyone:R", "/grant", "%USERNAME%:F")
	case "ADMIN":
		// Only allow the user, remove everyone else's access
		cmd = exec.Command("icacls", path, "/grant", "%USERNAME%:F", "/inheritance:r")
	default:
		return fmt.Errorf("invalid permission level: %s", level)
	}

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Windows permissions: %s\n%s", err, output)
	}

	fmt.Printf("Updated permissions for: %s (Level: %s)\n", path, level)
	return nil
}
