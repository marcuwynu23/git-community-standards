package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// baseURL is the raw content root on the GitHub repository. File specs are
// resolved relative to this base.
const baseURL = "https://raw.githubusercontent.com/marcuwynu23/git-community-standards/refs/heads/main/docs/templates"

// FileSpec describes a single template to download.
type FileSpec struct {
	RemotePath string
	LocalPath  string
	Fallbacks  []string
	Optional   bool
}

// Fetcher downloads the content for a remote path and returns it.
type Fetcher func(remotePath string) ([]byte, error)

// DefaultFetcher fetches templates over HTTP from baseURL.
func DefaultFetcher(remotePath string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", baseURL, remotePath)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

// App is the runtime configuration for the community-standards commands.
type App struct {
	Version string
	Fetcher Fetcher
	Stdout  io.Writer
	Stderr  io.Writer
	// BaseDir, when set, is prepended to every LocalPath. It allows tests to
	// sandbox filesystem operations in a temporary directory.
	BaseDir  string
	useColor bool
}

// NewApp returns an App wired to the process stdout/stderr with the HTTP fetcher.
func NewApp(version string) *App {
	return &App{
		Version:  version,
		Fetcher:  DefaultFetcher,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		useColor: detectTTY(),
	}
}

// resolve returns a spec's LocalPath relative to BaseDir (when set).
func (a *App) resolve(localPath string) string {
	if a.BaseDir == "" {
		return localPath
	}
	return filepath.Join(a.BaseDir, localPath)
}

// Run executes the parsed arguments. It returns the process exit code.
func (a *App) Run(args []string) int {
	if len(args) > 0 && args[0] == "community-standards" {
		args = args[1:]
	}

	if len(args) < 1 {
		a.printUsage()
		return 1
	}

	switch args[0] {
	case "--version", "-v":
		a.printSuccess("git-community-standards %s", a.Version)
		return 0
	case "list":
		a.listCategories()
		return 0
	case "apply":
		if err := a.applyCommand(args[1:]); err != nil {
			a.printError("%v", err)
			return 1
		}
		return 0
	default:
		a.printUsage()
		return 1
	}
}

func (a *App) applyCommand(args []string) error {
	platform, override, err := ParseApplyArgs(args)
	if err != nil {
		return err
	}

	a.printPlain("Applying community standards...\n")
	if !override {
		a.printInfo("Override mode: OFF (existing files will be skipped)")
	} else {
		a.printWarn("Override mode: ON (existing files will be replaced)")
	}

	if err := a.applySpecs(generalFiles, override); err != nil {
		return err
	}

	if platform == "" || platform == "none" || platform == "general" {
		a.printSuccess("General community docs applied successfully.")
		return nil
	}

	specs := platforms[platform]
	if err := a.applySpecs(specs, override); err != nil {
		return err
	}
	a.printSuccess("Platform %q templates applied successfully.", platform)
	return nil
}

func (a *App) applySpecs(specs []FileSpec, override bool) error {
	for _, spec := range specs {
		local := a.resolve(spec.LocalPath)
		if !override {
			if _, err := os.Stat(local); err == nil {
				a.printWarn("Skipping %s (already exists). Use `apply override` to replace it.", spec.LocalPath)
				continue
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed checking %s: %w", spec.LocalPath, err)
			}
		}

		a.printStep("Downloading %s...", spec.LocalPath)
		if err := a.downloadAndWrite(spec, local); err != nil {
			if spec.Optional {
				a.printWarn("Skipping optional file %s: %v", spec.LocalPath, err)
				continue
			}
			return fmt.Errorf("failed downloading %s: %w", spec.RemotePath, err)
		}
	}
	return nil
}

func (a *App) downloadAndWrite(spec FileSpec, local string) error {
	paths := append([]string{spec.RemotePath}, spec.Fallbacks...)
	var body []byte
	var usedRemotePath string
	var err error

	for _, remotePath := range paths {
		body, err = a.Fetcher(remotePath)
		if err == nil {
			usedRemotePath = remotePath
			break
		}
	}
	if err != nil {
		return err
	}

	if usedRemotePath != spec.RemotePath {
		a.printInfo("Using fallback source %s for %s.", usedRemotePath, spec.LocalPath)
	}

	parentDir := filepath.Dir(local)
	if parentDir != "." {
		if err := os.MkdirAll(parentDir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(local, body, 0o644)
}

// SupportedPlatforms returns the available platforms in a stable order.
func SupportedPlatforms() []string {
	out := make([]string, len(platformOrder))
	copy(out, platformOrder)
	return out
}

// IsPlatform reports whether name names a known platform template set.
func IsPlatform(name string) bool {
	_, ok := platforms[name]
	return ok
}

// ParseApplyArgs parses the arguments of an `apply` command into a platform
// token (which may be "none"/"general" for general-only) and the override
// flag. It returns an error if the arguments are invalid.
func ParseApplyArgs(args []string) (platform string, override bool, err error) {
	if len(args) > 2 {
		return "", false, errors.New("usage: git community-standards apply [github|gitlab|bitbucket|none] [override]")
	}
	for _, arg := range args {
		if arg == "override" {
			override = true
			continue
		}
		if platform != "" {
			return "", false, errors.New("usage: git community-standards apply [github|gitlab|bitbucket|none] [override]")
		}
		platform = arg
	}
	if platform != "" && platform != "none" && platform != "general" {
		if !IsPlatform(platform) {
			return "", false, fmt.Errorf("unknown platform %q. Run `git-community-standards list` to see valid options", platform)
		}
	}
	return platform, override, nil
}

func (a *App) listCategories() {
	a.printHeader("General docs (always applied):")
	a.printPlain("  general: README, LICENSE, CONTRIBUTING, CODE_OF_CONDUCT, RELEASE-NOTES, SECURITY\n\n")
	a.printHeader("Available platforms (applied on top of general docs):")
	for _, platform := range platformOrder {
		description := platformDescriptions[platform]
		if description == "" {
			description = "No description available."
		}
		a.printStep("%s: %s", platform, description)
	}
}

func (a *App) printUsage() {
	a.printHeader("Usage:")
	a.printPlain("  git community-standards list\n")
	a.printPlain("  git community-standards apply\n")
	a.printPlain("  git community-standards apply <github|gitlab|bitbucket|none>\n")
	a.printPlain("  git community-standards apply override\n")
	a.printPlain("  git community-standards apply <github|gitlab|bitbucket> override\n")
	a.printPlain("  git community-standards --version\n")
	a.printPlain("  git community-standards -v\n\n")
	a.printInfo("If no platform (or `none`) is given, only the general community docs are applied.")
}
