package main

// with go modules disabled

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/analogj/checkr/pkg/actions"
	"github.com/analogj/checkr/pkg/config"
	"github.com/analogj/checkr/pkg/utils"
	"github.com/analogj/checkr/pkg/version"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var goos string
var goarch string

func main() {

	config, err := config.Create()
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	app := &cli.App{
		Name:     "checkr",
		Usage:    "Github Check Suite CLI",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			ghcsRepo := "github.com/AnalogJ/checkr"

			var versionInfo string
			if len(goos) > 0 && len(goarch) > 0 {
				versionInfo = fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)
			} else {
				versionInfo = fmt.Sprintf("dev-%s", version.VERSION)
			}

			subtitle := ghcsRepo + utils.LeftPad2Len(versionInfo, " ", 65-len(ghcsRepo))

			color.New(color.FgGreen).Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
				`
			  oooooooo8 oooo                             oooo                    
			o888     88  888ooooo   ooooooooo8  ooooooo   888  ooooo oo oooooo   
			888          888   888 888oooooo8 888     888 888o888     888    888 
			888o     oo  888   888 888        888         8888 88o    888        
			 888oooo88  o888o o888o  88oooo888  88ooo888 o888o o888o o888o
			%s

			`), subtitle))

			return nil
		},

		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a Github Check Run",
				//UsageText:   "doo - does the dooing",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					if c.IsSet("pr") {
						config.Set("pr", c.Int("pr"))
					} else if c.IsSet("sha") {
						config.Set("sha", c.String("sha"))
					}

					if c.IsSet("full-name") {
						parts := strings.Split(c.String("full-name"), "/")
						config.Set("org", parts[0])
						config.Set("repo", parts[1])
					} else if c.IsSet("org") && c.IsSet("repo") {
						config.Set("org", c.String("org"))
						config.Set("repo", c.String("repo"))
					}

					payloadPath, err := utils.ExpandPath(c.String("payload-path"))
					if err != nil {
						return err
					}

					if !utils.FileExists(payloadPath) {
						return errors.New(fmt.Sprintf("payload path invalid. Please ensure that file exists: %s", payloadPath))
					}

					err = config.ValidateConfig()
					if err != nil {
						return err
					}

					runAction := actions.RunAction{Config: config}
					return runAction.Create(payloadPath)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "org, o",
						Usage: "Github repository owner/organization name",
					},
					&cli.StringFlag{
						Name:  "repo, r",
						Usage: "Github repository name",
					},
					&cli.StringFlag{
						Name:  "full-name, fn",
						Usage: "Github full repository name, eg. ORG_NAME/REPO_NAME",
					},

					&cli.StringFlag{
						Name:  "pr",
						Usage: "Github pull request number (required if sha is not provided)",
					},
					&cli.StringFlag{
						Name:  "sha",
						Usage: "Github pull request head SHA (required if pr is not provided)",
					},

					&cli.StringFlag{
						Name:  "url",
						Usage: "Provide an optional link that will be set in the Check run as the `detail_url`",
					},

					&cli.StringFlag{
						Name:     "payload-path",
						Required: true,
						Usage:    "Provide the Github Check Run compatible JSON payload. See: https://developer.github.com/v3/checks/runs/#create-a-check-run",
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}
}

//func main2() {
//
//
//
//	var createMessage = "new file added"
//	var createBranch = "AnalogJ-patch-3"
//
//	//author ghcs-test[bot] <54657905+ghcs-test[bot]@users.noreply.github.com> 1567062088 +0000
//	var authorName = "ghcs-test"                                //name can be anything
//	var authorEmail = "ghcs-test[bot]@users.noreply.github.com" //[bot]required here, but prefix number not requirerd.
//	var author = github.CommitAuthor{
//		Name:  &authorName,
//		Email: &authorEmail,
//	}
//
//	/// TOOD: write a file to a github branch using the commits api
//	created, resp, err := appClient.Repositories.CreateFile(ctx, "AnalogJ", "golang_analogj_test", "netnew11.txt", &github.RepositoryContentFileOptions{
//		Message:   &createMessage,
//		Content:   []byte("This is my content in a byte array,"),
//		Branch:    &createBranch,
//		Author:    &author,
//		Committer: &author,
//	})
//	if err != nil {
//		fmt.Printf("error: %s", err)
//	}
//	fmt.Print(created)
//
//}
