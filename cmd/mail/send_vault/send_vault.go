package send_vault

import (
	"log"
	"strconv"

	mail_cmd "github.com/sikalabs/slu/cmd/mail"
	"github.com/sikalabs/slu/lib/vault_smtp"
	"github.com/sikalabs/slu/utils/mail_utils"
	"github.com/sikalabs/slu/utils/stdin_utils"
	"github.com/spf13/cobra"
)

var FlagTo string
var FlagSubject string
var FlagMessage string

var Cmd = &cobra.Command{
	Use:     "send-vault",
	Short:   "Send mail (with credentials from vault)",
	Aliases: []string{"send-v", "s-v", "sv"},
	Args:    cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		var err error
		host, port, user, password, err := vault_smtp.GetSMTPSecrets("default")
		if err != nil {
			log.Fatal(err)
		}
		message := FlagMessage
		if message == "-" {
			message = stdin_utils.ReadAll()
		}
		err = mail_utils.SendSimpleMail(
			host,
			strconv.Itoa(port),
			user,
			password,
			user,
			FlagTo,
			FlagSubject,
			message,
		)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	mail_cmd.Cmd.AddCommand(Cmd)
	Cmd.Flags().StringVarP(
		&FlagTo,
		"to",
		"t",
		"",
		"to (john@acme.com)",
	)
	Cmd.MarkFlagRequired("to")
	Cmd.Flags().StringVarP(
		&FlagSubject,
		"subject",
		"s",
		"",
		"Email subject",
	)
	Cmd.MarkFlagRequired("subject")
	Cmd.Flags().StringVarP(
		&FlagMessage,
		"message",
		"m",
		"",
		"email message (\"-\" for stdin)",
	)
	Cmd.MarkFlagRequired("message")
}
