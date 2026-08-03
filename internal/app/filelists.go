package app

// generalFiles are the community docs applied for every platform.
var generalFiles = []FileSpec{
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
	"bitbucket": "Bitbucket templates (issue templates, pull request template)",
}

var platforms = map[string][]FileSpec{
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

// baseURLForTesting is unexported and used by tests to assert URL construction.
func baseURLForTesting() string {
	return baseURL
}
