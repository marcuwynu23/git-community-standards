package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

const baseURL = "https://raw.githubusercontent.com/marcuwynu23/git-community-standards/refs/heads/main/docs/templates"

var version = "dev"

type fileSpec struct {
	RemotePath string
	LocalPath  string
	Fallbacks  []string
	Optional   bool
}

var categoryDescriptions = map[string]string{
	"root":            "Root-level community docs (README, CONTRIBUTING, CODE_OF_CONDUCT, RELEASE-NOTES, etc.)",
	"github":          "GitHub repository metadata files under .github (for example FUNDING.yml)",
	"issue-templates": "GitHub issue templates under .github/ISSUE_TEMPLATE",
	"pr-template":     "GitHub pull request template under .github/PULL_REQUEST_TEMPLATE.md",
}

var categories = map[string][]fileSpec{
	"root": {
		{RemotePath: "README.md", LocalPath: "README.md"},
		{
			RemotePath: "LICENSE",
			LocalPath:  "LICENSE",
			Fallbacks:  []string{"LICENSE.md", "LICENCE", "LICENCE.md"},
			Optional:   true,
		},
		{RemotePath: "CONTRIBUTING.md", LocalPath: "CONTRIBUTING.md"},
		{RemotePath: "CODE_OF_CONDUCT.md", LocalPath: "CODE_OF_CONDUCT.md"},
		{RemotePath: "RELEASE-NOTES.md", LocalPath: "RELEASE-NOTES.md"},
		{
			RemotePath: "SECURITY.md",
			LocalPath:  "SECURITY.md",
			Fallbacks:  []string{"SECURITY_POLICY.md"},
			Optional:   true,
		},
	},
	"github": {
		{RemotePath: "FUNDING.yml", LocalPath: ".github/FUNDING.yml"},
	},
	"issue-templates": {
		{
			RemotePath: ".github/ISSUE_TEMPLATE/bug_report.md",
			LocalPath:  ".github/ISSUE_TEMPLATE/bug_report.md",
			Fallbacks:  []string{"ISSUE_TEMPLATE/bug_report.md", "bug_report.md"},
		},
		{
			RemotePath: ".github/ISSUE_TEMPLATE/feature_request.md",
			LocalPath:  ".github/ISSUE_TEMPLATE/feature_request.md",
			Fallbacks:  []string{"ISSUE_TEMPLATE/feature_request.md", "feature_request.md"},
		},
	},
	"pr-template": {
		{
			RemotePath: ".github/PULL_REQUEST_TEMPLATE.md",
			LocalPath:  ".github/PULL_REQUEST_TEMPLATE.md",
			Fallbacks:  []string{"PULL_REQUEST_TEMPLATE.md"},
		},
	},
}

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
		fmt.Println(version)
		return
	}

	switch args[0] {
	case "list":
		listCategories()
	case "apply":
		if err := applyCommand(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func applyCommand(args []string) error {
	if len(args) > 2 {
		return errors.New("usage: git community-standards apply [category] [override]")
	}

	override := false
	category := ""

	for _, arg := range args {
		if arg == "override" {
			override = true
			continue
		}
		if category != "" {
			return errors.New("usage: git community-standards apply [category] [override]")
		}
		category = arg
	}

	fmt.Println("Applying community standards...")
	if !override {
		fmt.Println("Override mode: OFF (existing files will be skipped)")
	} else {
		fmt.Println("Override mode: ON (existing files will be replaced)")
	}

	if category == "" {
		order := []string{"root", "github", "issue-templates", "pr-template"}
		for _, category := range order {
			if err := applyCategory(category, override); err != nil {
				return err
			}
		}
		fmt.Println("Community standards applied successfully.")
		return nil
	}

	if _, ok := categories[category]; !ok {
		return fmt.Errorf("unknown category %q. Run `git-community-standards list` to see valid categories", category)
	}

	if err := applyCategory(category, override); err != nil {
		return err
	}
	fmt.Printf("Category %q applied successfully.\n", category)
	return nil
}

func applyCategory(category string, override bool) error {
	specs := categories[category]
	for _, spec := range specs {
		if !override {
			if _, err := os.Stat(spec.LocalPath); err == nil {
				fmt.Printf("Skipping %s (already exists). Use `apply override` to replace it.\n", spec.LocalPath)
				continue
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed checking %s: %w", spec.LocalPath, err)
			}
		}

		fmt.Printf("Downloading %s...\n", spec.LocalPath)
		if err := downloadAndWrite(spec); err != nil {
			if spec.Optional {
				fmt.Printf("Skipping optional file %s: %v\n", spec.LocalPath, err)
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
		fmt.Printf("Using fallback source %s for %s.\n", usedRemotePath, spec.LocalPath)
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
	fmt.Println("Available categories:")
	keys := make([]string, 0, len(categories))
	for category := range categories {
		keys = append(keys, category)
	}
	sort.Strings(keys)
	for _, category := range keys {
		description := categoryDescriptions[category]
		if description == "" {
			description = "No description available."
		}
		fmt.Printf("- %s: %s\n", category, description)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  git community-standards list")
	fmt.Println("  git community-standards apply")
	fmt.Println("  git community-standards apply <category>")
	fmt.Println("  git community-standards apply override")
	fmt.Println("  git community-standards apply <category> override")
	fmt.Println("  git community-standards --version")
	fmt.Println("  git community-standards -v")
}
