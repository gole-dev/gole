package project

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/gole-dev/gole/cmd/gole/internal/base"
)

// Project is a project template.
type Project struct {
	Name string
}

// New create a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("%s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		e = os.RemoveAll(to)
		if e != nil {
			return e
		}
	}

	fmt.Printf("Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	repo := base.NewRepo(layout)
	if err := repo.CopyTo(ctx, to, p.Name, []string{".git", ".github"}); err != nil {
		return err
	}
	//e := os.Rename(
	//	path.Join(to, "cmd", "server"),
	//	path.Join(to, "cmd", p.Name),
	//)
	//if e != nil {
	//	return e
	//}
	base.Tree(to, dir)

	fmt.Printf("\nProject creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("Use the following command to start the project:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go build"))
	fmt.Println(color.WhiteString("$ ./%s\n", p.Name))
	fmt.Println("Thanks for using gole")
	fmt.Println("Tutorial: https://gole.dev")
	return nil
}
