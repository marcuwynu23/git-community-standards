package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const baseURL = "https://raw.githubusercontent.com/marcuwynu23/git-community-standards/refs/heads/main/docs/templates"

var version = "dev"

type fileSpec struct {
	RemotePath string
	LocalPath  string
	Fallbacks  []string
	Optional   bool
}

var generalFiles = []fileSpec{
	{RemotePath: "general/README.md", LocalPath: "README.md"},
	{
		RemotePath: "general/LICENSE.md",
		LocalPath:  "LICENSE",
		Fallbacks:  []string{"general/LICENSE", "general/LICENCE", "general/LICENCE.md"},
		Optional:   true,
	},
	{RemotePath: "general/CONTRIBUTING.md", LocalPath: "CONTRIBUTING.md"},
	{RemotePath: "general/CODE_OF_CONDUCT.md", LocalPath: "CODE_OF_CONDUCT.md"},
	{RemotePath: "general/RELEASE-NOTES.md", LocalPath: "RELEASE-NOTES.md"},
	{
		RemotePath: "general/SECURITY.md",
		LocalPath:  "SECURITY.md",
		Fallbacks:  []string{"general/SECURITY_POLICY.md"},
		Optional:   true,
	},
}

var platformDescriptions = map[string]string{
	"github":    "GitHub templates under .github (FUNDING, issue templates, pull request template)",
	"gitlab":    "GitLab templates under .gitlab (issue templates, merge request template)",
	"bitbucket": "Bitbucket templates (issue templates and pull request template)",
}

var platforms = map[string][]fileSpec{
	"github": {
		{RemotePath: "github/FUNDING.yml", LocalPath: ".github/FUNDING.yml"},
		{
			RemotePath: "github/bug_report.md",
			LocalPath:  ".github/ISSUE_TEMPLATE/bug_report.md",
			Fallbacks:  []string{"bug_report.md"},
		},
		{
			RemotePath: "github/feature_request.md",
			LocalPath:  ".github/ISSUE_TEMPLATE/feature_request.md",
			Fallbacks:  []string{"feature_request.md"},
		},
		{
			RemotePath: "github/PULL_REQUEST_TEMPLATE.md",
			LocalPath:  ".github/PULL_REQUEST_TEMPLATE.md",
			Fallbacks:  []string{"PULL_REQUEST_TEMPLATE.md"},
		},
	},
	"gitlab": {
		{
			RemotePath: "gitlab/bug_report.md",
			LocalPath:  ".gitlab/issue_templates/bug_report.md",
			Fallbacks:  []string{"bug_report.md"},
		},
		{
			RemotePath: "gitlab/feature_request.md",
			LocalPath:  ".gitlab/issue_templates/feature_request.md",
			Fallbacks:  []string{"feature_request.md"},
		},
		{
			RemotePath: "gitlab/MR_TEMPLATE.md",
			LocalPath:  ".gitlab/merge_request_templates/default.md",
			Fallbacks:  []string{"MR_TEMPLATE.md"},
		},
	},
	"bitbucket": {
		{
			RemotePath: "bitbucket/bug_report.md",
			LocalPath:  ".bitbucket/ISSUE_TEMPLATE/bug_report.md",
			Fallbacks:  []string{"bug_report.md"},
		},
		{
			RemotePath: "bitbucket/feature_request.md",
			LocalPath:  ".bitbucket/ISSUE_TEMPLATE/feature_request.md",
			Fallbacks:  []string{"feature_request.md"},
		},
		{
			RemotePath: "bitbucket/PULL_REQUEST_TEMPLATE.md",
			LocalPath:  ".bitbucket/PULL_REQUEST_TEMPLATE.md",
			Fallbacks:  []string{"PULL_REQUEST_TEMPLATE.md"},
		},
	},
}

var platformOrder = []string{"github", "gitlab", "bitbucket"}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "community-standards" {
		args = args[1:]
	}

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	if args[0] == "--version" || args[0] == "-v" {
		printSuccess("git-community-standards %s", version)
		return
	}

	switch args[0] {
	case "list":
		listCategories()
	case "apply":
		if err := applyCommand(args[1:]); err != nil {
			printError("%v", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func applyCommand(args []string) error {
	if len(args) > 2 {
		return errors.New("usage: git community-standards apply [github|gitlab|bitbucket|none] [override]")
	}

	override := false
	platform := ""

	for _, arg := range args {
		if arg == "override" {
			override = true
			continue
		}
		if platform != "" {
			return errors.New("usage: git community-standards apply [github|gitlab|bitbucket|none] [override]")
		}
		platform = arg
	}

	// Validate the requested platform before touching the filesystem.
	if platform != "" && platform != "none" && platform != "general" {
		if _, ok := platforms[platform]; !ok {
			return fmt.Errorf("unknown platform %q. Run `git-community-standards list` to see valid options", platform)
		}
	}

	fmt.Println("Applying community standards...")
	if !override {
		printInfo("Override mode: OFF (existing files will be skipped)")
	} else {
		printWarn("Override mode: ON (existing files will be replaced)")
	}

	if err := applySpecs(generalFiles, override); err != nil {
		return err
	}

	if platform == "" || platform == "none" || platform == "general" {
		printSuccess("General community docs applied successfully.")
		return nil
	}

	specs := platforms[platform]
	if err := applySpecs(specs, override); err != nil {
		return err
	}
	printSuccess("Platform %q templates applied successfully.", platform)
	return nil
}

func applySpecs(specs []fileSpec, override bool) error {
	for _, spec := range specs {
		if !override {
			if _, err := os.Stat(spec.LocalPath); err == nil {
				printWarn("Skipping %s (already exists). Use `apply override` to replace it.", spec.LocalPath)
				continue
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed checking %s: %w", spec.LocalPath, err)
			}
		}

		printStep("Downloading %s...", spec.LocalPath)
		if err := downloadAndWrite(spec); err != nil {
			if spec.Optional {
				printWarn("Skipping optional file %s: %v", spec.LocalPath, err)
				continue
			}
			return fmt.Errorf("failed downloading %s: %w", spec.RemotePath, err)
		}
	}
	return nil
}

func downloadAndWrite(spec fileSpec) error {
	paths := append([]string{spec.RemotePath}, spec.Fallbacks...)
	var body []byte
	var usedRemotePath string
	var err error

	for _, remotePath := range paths {
		body, err = fetchRemote(remotePath)
		if err == nil {
			usedRemotePath = remotePath
			break
		}
	}
	if err != nil {
		return err
	}

	if usedRemotePath != spec.RemotePath {
		printInfo("Using fallback source %s for %s.", usedRemotePath, spec.LocalPath)
	}

	parentDir := filepath.Dir(spec.LocalPath)
	if parentDir != "." {
		if err := os.MkdirAll(parentDir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(spec.LocalPath, body, 0o644)
}

func fetchRemote(remotePath string) ([]byte, error) {
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

func listCategories() {
	printHeader("General docs (always applied):")
	printPlain("  general: README, LICENSE, CONTRIBUTING, CODE_OF_CONDUCT, RELEASE-NOTES, SECURITY\n\n")
	printHeader("Available platforms (applied on top of general docs):")
	for _, platform := range platformOrder {
		description := platformDescriptions[platform]
		if description == "" {
			description = "No description available."
		}
		printStep("%s: %s", platform, description)
	}
}

func printUsage() {
	printHeader("Usage:")
	printPlain("  git community-standards list\n")
	printPlain("  git community-standards apply\n")
	printPlain("  git community-standards apply <github|gitlab|bitbucket|none>\n")
	printPlain("  git community-standards apply override\n")
	printPlain("  git community-standards apply <github|gitlab|bitbucket> override\n")
	printPlain("  git community-standards --version\n")
	printPlain("  git community-standards -v\n\n")
	printInfo("If no platform (or `none`) is given, only the general community docs are applied.")
}
