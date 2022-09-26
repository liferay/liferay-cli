package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/viper"
)

var gitRepo *git.Repository

func init() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault(Const.repoDir, path.Join(dirname, ".lcectl"))
	viper.SetDefault(Const.repoRemote, "https://github.com/gamerson/lxc-localdev")
	viper.SetDefault(Const.repoBranch, "master")
}

func InitGit() {
	repoDir := viper.GetString(Const.repoDir)
	repo, err := git.PlainOpen(repoDir)

	if err != nil {
		os.MkdirAll(repoDir, os.ModePerm)

		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			Depth:        1,
			SingleBranch: true,
			RemoteName: fmt.Sprintf(
				"refs/heads/%s",
				viper.GetString(Const.repoBranch)),
			URL: viper.GetString(Const.repoRemote),
		})

		if err != nil {
			log.Fatal("Clone error: ", err)
		}

		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{viper.GetString(Const.repoRemote)},
		})
	}

	worktree, err := repo.Worktree()

	if err != nil {
		log.Fatal("worktree error: ", err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})

	if err != nil && err.Error() != "already up-to-date" {
		log.Fatal("pull error: ", err)
	}

	gitRepo = repo
}
