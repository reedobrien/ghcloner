package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	git "gopkg.in/src-d/go-git.v4"
	// git "gopkg.in/src-d/go-git.v4"
)

func main() {
	ghtoken := flag.String("token", "", "Github API Token")
	ghorg := flag.String("org", "", "Github org")

	flag.Parse()

	if *ghtoken == "" {
		log.Fatal("A github token is required.")
	}
	if *ghorg == "" {
		log.Fatal("A github organzation is required.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *ghtoken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 40},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, *ghorg, opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	// _, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
	// 	URL:      "https://github.com/src-d/go-git",
	// 	Progress: os.Stdout,
	// })

	for _, p := range allRepos {
		fmt.Println(*p.Name, *p.SSHURL)
		r, err := git.PlainClone("repos"+"/"+*ghorg+"/"+*p.Name, false, &git.CloneOptions{
			URL:               *p.SSHURL,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Progress:          os.Stdout,
		})
		if err != nil {
			log.Printf("Error cloning %s: %s", *ghorg+"/"+*p.Name, err)
		}
		ref, err := r.Head() // Ignoring error, I know.
		if err != nil {
			fmt.Println("Failed to clone", *p.Name)
			continue
		}
		fmt.Println("Cloned", *p.Name, ref.Hash())

	}
}
