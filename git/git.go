package git

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/viper"
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/constants"
)

func init() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault(constants.Const.RepoDir, filepath.Join(dirname, ".liferay", "cli", "sources", "localdev"))
	viper.SetDefault(constants.Const.RepoRemote, "https://github.com/liferay/liferay-localdev.git")
	viper.SetDefault(constants.Const.RepoBranch, "main")
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
		s.FinalMSG = fmt.Sprintf(ansicolor.Good + " Synced 'localdev' sources.\n")
		s.Start()
		defer s.Stop()
	}

	repoBranch := viper.GetString(constants.Const.RepoBranch)
	repoDir := viper.GetString(constants.Const.RepoDir)
	remoteUrl := viper.GetString(constants.Const.RepoRemote)
	repo, err := git.PlainOpen(repoDir)

	cloned := false

	if err != nil {
		os.MkdirAll(repoDir, os.ModePerm)

		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			Depth:         1,
			RemoteName:    "origin",
			SingleBranch:  true,
			URL:           remoteUrl,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", repoBranch)),
		})

		if err != nil {
			if s != nil {
				s.FinalMSG = fmt.Sprintf(ansicolor.Bad+" 'localdev' sources error %s.\n", err)
				s.Stop()
				os.Exit(1)
			} else {
				log.Fatalf("'localdev' sources error %s.\n", err)
			}
		}

		tagOrBranchAndHash := tagOrBranchAndHash(repo)

		if s != nil {
			s.FinalMSG = fmt.Sprintf(ansicolor.Good+" Cloned 'localdev' [%s] sources to %s\n", tagOrBranchAndHash, repoDir)
		} else {
			fmt.Printf("Cloned 'localdev' [%s] sources to %s\n", tagOrBranchAndHash, repoDir)
		}

		cloned = true
	}

	repoConfig, err := repo.Config()

	if err != nil {
		log.Fatalf("'localdev' sources error %s.\n", err)
	}

	coreSection := repoConfig.Raw.Section("core")

	if !coreSection.HasOption("autocrlf") || !coreSection.HasOption("filemode") {
		coreSection.SetOption("autocrlf", "false")
		coreSection.SetOption("filemode", "false")
		repo.SetConfig(repoConfig)
	}

	remote, err := repo.Remote("origin")

	if err != nil {
		log.Fatalf("'localdev' sources error %s.\n", err)
	}

	if remote.Config().URLs[0] != remoteUrl {
		remote.Config().URLs[0] = remoteUrl
		remoteConfig := remote.Config()
		repo.DeleteRemote("origin")
		repo.CreateRemote(remoteConfig)
	}

	if !cloned {
		worktree, err := repo.Worktree()

		if err != nil {
			log.Fatal("worktree error: ", err)
		}

		if err = repo.Fetch(&git.FetchOptions{
			Depth:      1,
			RemoteName: "origin",
			RefSpecs: []config.RefSpec{
				config.RefSpec("+refs/heads/*:refs/remotes/origin/*")},
			Force: true,
		}); err != nil {
			if err != git.NoErrAlreadyUpToDate && err != transport.ErrEmptyUploadPackRequest {
				log.Fatalf("Fetch error %s\n", err)
			}
		}

		hash, err := repo.ResolveRevision(plumbing.Revision(fmt.Sprintf("refs/remotes/origin/%s", repoBranch)))
		if err != nil {
			log.Fatalf("Error resolving name %s", err)
		}

		if err = worktree.Reset(&git.ResetOptions{
			Commit: *hash,
			Mode:   git.HardReset,
		}); err != nil {
			if err == git.NoErrAlreadyUpToDate || err == transport.ErrEmptyUploadPackRequest {
				if s != nil {
					s.FinalMSG = fmt.Sprintf(ansicolor.Good+" 'localdev' sources %s.\n", git.NoErrAlreadyUpToDate)
				} else {
					fmt.Printf("'localdev' sources %s.\n", git.NoErrAlreadyUpToDate)
				}

				return
			}

			if s != nil {
				s.FinalMSG = fmt.Sprintf(ansicolor.Bad+" 'localdev' sources error %s.\n", err)
				s.Stop()
				os.Exit(1)
			} else {
				log.Fatalf("'localdev' sources error %s.\n", err)
			}
		}

		tagOrBranchAndHash := tagOrBranchAndHash(repo)

		if s != nil {
			s.FinalMSG = fmt.Sprintf(ansicolor.Good+" 'localdev' sources updated to [%s]\n", tagOrBranchAndHash)
		} else {
			fmt.Printf("'localdev' sources updated to [%s]\n", tagOrBranchAndHash)
		}
	}
}

func tagOrBranchAndHash(repo *git.Repository) string {
	head, err := repo.Head()

	if err != nil {
		log.Fatal(err)
	}

	itr, err := repo.Tags()

	if err != nil {
		log.Fatal(err)
	}

	var tag *plumbing.Reference
	err = itr.ForEach(func(t *plumbing.Reference) error {
		tag = t
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if tag == nil {
		tag = head
	}

	return tag.Name().Short() + "-g" + head.Hash().String()[0:7]
}
