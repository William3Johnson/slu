package version

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	parent_cmd "github.com/sikalabs/slu/cmd/static_api"

	"github.com/spf13/cobra"
)

var FlagPretty bool
var FlagSetGitClean bool
var FlagSetGitDirty bool

type VersionJSON struct {
	GitRef                 string `json:"git_ref"`
	GitCommit              string `json:"git_commit"`
	GitTreeState           string `json:"git_tree_state"`
	BuildTimestampUnix     int    `json:"build_timestamp_unix"`
	BuildTimestampRFC3339  string `json:"build_timestamp_rfc3339"`
	BuildTimestampUnixDate string `json:"build_timestamp_unixdate"`
}

var Cmd = &cobra.Command{
	Use:     "version",
	Short:   "Generate version file",
	Aliases: []string{"v"},
	Args:    cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		r, err := git.PlainOpen(".")
		if err != nil {
			panic(err)
		}
		w, err := r.Worktree()
		if err != nil {
			panic(err)
		}
		s, err := w.Status()
		if err != nil {
			panic(err)
		}
		head, err := r.Head()
		if err != nil {
			panic(err)
		}
		var gitTreeState string
		if s.IsClean() {
			gitTreeState = "clean"
		} else {
			gitTreeState = "dirty"
		}
		if FlagSetGitClean && FlagSetGitDirty {
			log.Fatalln("can't use --set-git-clean and --set-git-dirty together")
		}
		if FlagSetGitClean {
			gitTreeState = "clean"
		}
		if FlagSetGitDirty {
			gitTreeState = "dirty"
		}
		t := time.Now()
		var data []byte
		v := VersionJSON{
			GitRef:                 head.Name().Short(),
			GitCommit:              head.Hash().String(),
			GitTreeState:           gitTreeState,
			BuildTimestampUnix:     int(t.Unix()),
			BuildTimestampRFC3339:  t.Format(time.RFC3339),
			BuildTimestampUnixDate: t.Format(time.UnixDate),
		}
		if FlagPretty {
			data, err = json.MarshalIndent(v, "", "  ")
		} else {
			data, err = json.Marshal(v)
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	},
}

func init() {
	parent_cmd.Cmd.AddCommand(Cmd)
	Cmd.Flags().BoolVarP(
		&FlagPretty,
		"pretty",
		"p",
		false,
		"Pretty output",
	)
	Cmd.Flags().BoolVar(
		&FlagSetGitClean,
		"set-git-clean",
		false,
		"Manually set Git tree state to clean",
	)
	Cmd.Flags().BoolVar(
		&FlagSetGitDirty,
		"set-git-dirty",
		false,
		"Manually set Git tree state to dirty",
	)
}
