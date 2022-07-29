package setup_runner

import (
	"log"
	"os"

	gitlab_ci_cmd "github.com/sikalabs/slu/cmd/gitlab_ci"
	"github.com/sikalabs/slu/utils/gitlab_ci_utils/setup_runner_utils"

	"github.com/sikalabs/slu/lib/vault_gitlab_ci"
	"github.com/spf13/cobra"
)

var FlagGitlabUrl string
var FlagRegistrationToken string
var FlagGitlabName string
var FlagConcurency int

var Cmd = &cobra.Command{
	Use:     "setup-runner",
	Short:   "Setup Gitlab Runner in Docker",
	Aliases: []string{"sr"},
	Run: func(c *cobra.Command, args []string) {
		var err error

		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}

		gitlabUrl := FlagGitlabUrl
		registrationToken := FlagRegistrationToken

		if FlagGitlabName != "" {
			gitlabUrl, registrationToken, err = vault_gitlab_ci.GetGitlabCiSecrets(FlagGitlabName)
			if err != nil {
				log.Fatal(err)
			}
		}

		if gitlabUrl == "" {
			log.Fatal("flags gitlab-url and registration-token OR flag gitlab (for Vault) are required")
		}

		err = setup_runner_utils.SetupGitlabRunnerDocker(gitlabUrl, registrationToken, hostname, FlagConcurency, false)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	gitlab_ci_cmd.Cmd.AddCommand(Cmd)
	Cmd.PersistentFlags().StringVarP(
		&FlagGitlabUrl,
		"gitlab-url",
		"u",
		"",
		"Gitlab URL",
	)
	Cmd.PersistentFlags().StringVarP(
		&FlagRegistrationToken,
		"registration-token",
		"t",
		"",
		"Gitlab Registration Token",
	)
	Cmd.MarkFlagsRequiredTogether("gitlab-url", "registration-token")
	Cmd.PersistentFlags().StringVarP(
		&FlagGitlabName,
		"gitlab",
		"g",
		"",
		"Gitlab name in Vault (slu/gitlab-ci/$name)",
	)
	Cmd.MarkFlagsMutuallyExclusive("gitlab-url", "gitlab")
	Cmd.PersistentFlags().IntVarP(
		&FlagConcurency,
		"concurency",
		"c",
		1,
		"Set maximun concurency",
	)
}
