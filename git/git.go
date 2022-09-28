package git

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
)

func init() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault(constants.Const.RepoDir, filepath.ToSlash(path.Join(dirname, ".lcectl", "sources", "localdev")))
	viper.SetDefault(constants.Const.RepoRemote, "https://github.com/gamerson/lxc-localdev")
	viper.SetDefault(constants.Const.RepoBranch, "master")
	viper.SetDefault(constants.Const.RepoSync, true)
}

func SyncGit(verbose bool) {
	if repoSync := viper.GetBool(constants.Const.RepoSync); !repoSync {
		return
	}

	var s *spinner.Spinner

	if !verbose {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Synchronizing 'localdev' sources..."
		s.FinalMSG = fmt.Sprintf("\u2705 Synced 'localdev' sources.\n")
		s.Start()
		defer s.Stop()
	}

	repoDir := viper.GetString(constants.Const.RepoDir)
	repo, err := git.PlainOpen(repoDir)

	if err != nil {
		os.MkdirAll(repoDir, os.ModePerm)

		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			Depth:        1,
			SingleBranch: true,
			RemoteName: fmt.Sprintf(
				"refs/heads/%s",
				viper.GetString(constants.Const.RepoBranch)),
			URL: viper.GetString(constants.Const.RepoRemote),
		})

		if err != nil {
			log.Fatal("Clone error: ", err)
		}

		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{viper.GetString(constants.Const.RepoRemote)},
		})

		if err != nil {
			log.Fatal("Remote error: ", err)
		}
	}

	worktree, err := repo.Worktree()

	if err != nil {
		log.Fatal("worktree error: ", err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})

	if err != nil && err.Error() != "already up-to-date" {
		//log.Fatal("pull error: ", err)
	}
}
