package flags

var (
	GitHubRepo          = FlagSet.String("github-repo", "", "GitHub repo")
	GitHubOwner         = FlagSet.String("github-owner", "", "GitHub owner")
	GitHubToken         = FlagSet.String("github-token", "", "GitHub token")
	GitHubFile          = FlagSet.String("github-file", "", "GitHub file")
	GitHubRef           = FlagSet.String("github-ref", "", "GitHub ref")
	GitHubOpenPR        = FlagSet.Bool("github-open-pr", false, "open PR on changes. Default: false")
	GitHubBaseBranch    = FlagSet.String("github-base-branch", "", "base branch for PR")
	GitHubBranch        = FlagSet.String("github-branch", "", "branch for PR.")
	GitHubCommitName    = FlagSet.String("github-commit-name", "", "commit name")
	GitHubCommitEmail   = FlagSet.String("github-commit-email", "", "commit email")
	GitHubCommitMessage = FlagSet.String("github-commit-message", "", "commit message")
	GitHubPRTitle       = FlagSet.String("github-pr-title", "", "PR title")
	GitHubPRBody        = FlagSet.String("github-pr-body", "", "PR body")
)
