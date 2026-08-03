package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeFetcher returns the given contents for each remote path that maps to a
// local file, so tests run without touching the network.
func fakeFetcher(contents map[string]string, calls *[]string) Fetcher {
	return func(remotePath string) ([]byte, error) {
		*calls = append(*calls, remotePath)
		body, ok := contents[remotePath]
		if !ok {
			return nil, errors.New("not found: " + remotePath)
		}
		return []byte(body), nil
	}
}

func TestSupportedPlatforms(t *testing.T) {
	got := SupportedPlatforms()
	want := []string{"github", "gitlab", "bitbucket"}
	if len(got) != len(want) {
		t.Fatalf("len(SupportedPlatforms()) = %d, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("SupportedPlatforms()[%d] = %q, want %q", i, got[i], w)
		}
	}
}

func TestIsPlatform(t *testing.T) {
	for _, p := range []string{"github", "gitlab", "bitbucket"} {
		if !IsPlatform(p) {
			t.Errorf("IsPlatform(%q) = false, want true", p)
		}
	}
	if IsPlatform("gitea") {
		t.Error("IsPlatform(\"gitea\") = true, want false")
	}
}

func TestParseApplyArgs(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantPlat  string
		wantOver  bool
		wantErr   bool
		errSubstr string
	}{
		{"empty", []string{}, "", false, false, ""},
		{"github", []string{"github"}, "github", false, false, ""},
		{"gitlab override", []string{"gitlab", "override"}, "gitlab", true, false, ""},
		{"override only", []string{"override"}, "", true, false, ""},
		{"none", []string{"none"}, "none", false, false, ""},
		{"general", []string{"general"}, "general", false, false, ""},
		{"invalid platform", []string{"gitea"}, "", false, true, "unknown platform"},
		{"too many args", []string{"github", "gitlab", "bitbucket"}, "", false, true, "usage"},
		{"two platforms", []string{"github", "gitlab"}, "", false, true, "usage"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plat, over, err := ParseApplyArgs(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.errSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if plat != tc.wantPlat {
				t.Errorf("platform = %q, want %q", plat, tc.wantPlat)
			}
			if over != tc.wantOver {
				t.Errorf("override = %v, want %v", over, tc.wantOver)
			}
		})
	}
}

func TestBaseURLIsGitHubRaw(t *testing.T) {
	if !strings.HasPrefix(baseURLForTesting(), "https://raw.githubusercontent.com/") {
		t.Fatalf("baseURL = %q, want a raw.githubusercontent.com URL", baseURLForTesting())
	}
	if strings.Contains(baseURLForTesting(), "gist") {
		t.Fatalf("baseURL still points at gist: %q", baseURLForTesting())
	}
	if !strings.HasSuffix(baseURLForTesting(), "/docs/templates") {
		t.Fatalf("baseURL = %q, want it to end with /docs/templates", baseURLForTesting())
	}
}

func TestApplyGeneralOnlyWritesGeneralFiles(t *testing.T) {
	dir := t.TempDir()
	contents := map[string]string{
		"general/README.md":          "# Hello",
		"general/CONTRIBUTING.md":    "## Contribute",
		"general/CODE_OF_CONDUCT.md": "## CoC",
		"general/RELEASE-NOTES.md":   "# Release Notes",
		"general/SECURITY.md":        "# Security",
	}
	var calls []string
	a := &App{
		Version:  "test",
		Fetcher:  fakeFetcher(contents, &calls),
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
		BaseDir:  dir,
		useColor: false,
	}

	if err := a.applyCommand(nil); err != nil {
		t.Fatalf("applyCommand(nil) error: %v", err)
	}

	// General-only should NOT download any platform files.
	for _, c := range calls {
		if strings.HasPrefix(c, "github/") || strings.HasPrefix(c, "gitlab/") || strings.HasPrefix(c, "bitbucket/") {
			t.Errorf("downloaded platform file %q during general-only apply", c)
		}
	}

	// Confirm expected general files were written.
	wantFiles := []string{"README.md", "SECURITY.md", "RELEASE-NOTES.md"}
	for _, f := range wantFiles {
		path := filepath.Join(dir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", f, err)
		}
	}
}

func TestApplyPlatformDownloadsPlatformFiles(t *testing.T) {
	dir := t.TempDir()
	contents := map[string]string{}
	// Seed general + github content.
	for _, rp := range []string{
		"general/README.md", "general/LICENSE.md", "general/CONTRIBUTING.md",
		"general/CODE_OF_CONDUCT.md", "general/RELEASE-NOTES.md",
		"github/FUNDING.yml", "github/bug_report.md", "github/feature_request.md",
		"github/PULL_REQUEST_TEMPLATE.md",
	} {
		contents[rp] = "x"
	}
	var calls []string
	a := &App{
		Version:  "test",
		Fetcher:  fakeFetcher(contents, &calls),
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
		BaseDir:  dir,
		useColor: false,
	}

	if err := a.applyCommand([]string{"github"}); err != nil {
		t.Fatalf("applyCommand error: %v", err)
	}

	for _, want := range []string{"github/FUNDING.yml", "github/bug_report.md", "github/feature_request.md", "github/PULL_REQUEST_TEMPLATE.md"} {
		found := false
		for _, c := range calls {
			if c == want {
				found = true
			}
		}
		if !found {
			t.Errorf("expected fetch of %q; got %v", want, calls)
		}
	}

	// Spot-check a platform file landed in the right place.
	dst := filepath.Join(dir, ".github/FUNDING.yml")
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("expected %s to exist: %v", dst, err)
	}
}

func TestApplyInvalidPlatformFailsBeforeDownload(t *testing.T) {
	dir := t.TempDir()
	var calls []string
	a := &App{
		Version:  "test",
		Fetcher:  fakeFetcher(map[string]string{}, &calls),
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
		BaseDir:  dir,
		useColor: false,
	}

	err := a.applyCommand([]string{"gitea"})
	if err == nil {
		t.Fatal("expected error for invalid platform")
	}
	if !strings.Contains(err.Error(), "unknown platform") {
		t.Fatalf("error = %q, want it to mention unknown platform", err.Error())
	}
	// No file should have been downloaded before the error.
	if len(calls) != 0 {
		t.Errorf("expected zero fetches before validation error, got %v", calls)
	}
}

func TestApplyOverrideReplacesExisting(t *testing.T) {
	dir := t.TempDir()
	contents := map[string]string{"general/README.md": "# v2", "general/CONTRIBUTING.md": "# c", "general/CODE_OF_CONDUCT.md": "# c", "general/RELEASE-NOTES.md": "# r", "general/SECURITY.md": "# s"}
	// Pre-create the README with old content and do NOT pass override -> must skip.
	readmePath := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# old"), 0o644); err != nil {
		t.Fatalf("seeding README: %v", err)
	}
	var calls []string
	a := &App{
		Version:  "test",
		Fetcher:  fakeFetcher(contents, &calls),
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
		BaseDir:  dir,
		useColor: false,
	}

	if err := a.applyCommand(nil); err != nil {
		t.Fatalf("applyCommand error: %v", err)
	}

	// README should be untouched (still "# old").
	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("reading README: %v", err)
	}
	if string(data) != "# old" {
		t.Errorf("README was overwritten without override; got %q", string(data))
	}

	// Now with override -> README should be replaced.
	if err := a.applyCommand([]string{"override"}); err != nil {
		t.Fatalf("applyCommand override error: %v", err)
	}
	data, err = os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("reading README: %v", err)
	}
	if string(data) != "# v2" {
		t.Errorf("README not replaced with override; got %q, want %q", string(data), "# v2")
	}
}

func TestListCategoriesOutput(t *testing.T) {
	var buf bytes.Buffer
	a := &App{
		Version:  "test",
		Fetcher:  fakeFetcher(map[string]string{}, &[]string{}),
		Stdout:   &buf,
		Stderr:   &bytes.Buffer{},
		useColor: false,
	}
	a.listCategories()
	out := buf.String()
	if !strings.Contains(out, "General docs (always applied)") {
		t.Errorf("list output missing general header: %q", out)
	}
	if !strings.Contains(out, "Available platforms") {
		t.Errorf("list output missing platforms header: %q", out)
	}
	for _, p := range []string{"github", "gitlab", "bitbucket"} {
		if !strings.Contains(out, p) {
			t.Errorf("list output missing platform %q: %q", p, out)
		}
	}
}

func TestRunVersion(t *testing.T) {
	var buf bytes.Buffer
	a := &App{
		Version:  "1.2.3",
		Fetcher:  fakeFetcher(map[string]string{}, &[]string{}),
		Stdout:   &buf,
		Stderr:   &bytes.Buffer{},
		useColor: false,
	}
	code := a.Run([]string{"--version"})
	if code != 0 {
		t.Fatalf("Run(--version) exit code = %d, want 0", code)
	}
	if !strings.Contains(buf.String(), "1.2.3") {
		t.Errorf("version output = %q, want it to contain version", buf.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	a := &App{
		Version:  "test",
		Fetcher:  fakeFetcher(map[string]string{}, &[]string{}),
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
		useColor: false,
	}
	code := a.Run([]string{"bogus"})
	if code != 1 {
		t.Fatalf("Run(bogus) exit code = %d, want 1", code)
	}
}
