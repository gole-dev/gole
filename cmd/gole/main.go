package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/gole-dev/gole/cmd/gole/internal/cache"
	"github.com/gole-dev/gole/cmd/gole/internal/handler"
	"github.com/gole-dev/gole/cmd/gole/internal/model"
	"github.com/gole-dev/gole/cmd/gole/internal/project"
	"github.com/gole-dev/gole/cmd/gole/internal/proto"
	"github.com/gole-dev/gole/cmd/gole/internal/repo"
	"github.com/gole-dev/gole/cmd/gole/internal/run"
	"github.com/gole-dev/gole/cmd/gole/internal/service"
	"github.com/gole-dev/gole/cmd/gole/internal/task"
	"github.com/gole-dev/gole/cmd/gole/internal/upgrade"
)

var (
	// Version is the version of the compiled software.
	Version = "v0.1.0"

	rootCmd = &cobra.Command{
		Use:     "gole",
		Short:   "gole: The Go framework with what you need",
		Long:    `gole: The Go framework with what you need`,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(handler.CmdHandler)
	rootCmd.AddCommand(cache.CmdCache)
	rootCmd.AddCommand(repo.CmdRepo)
	rootCmd.AddCommand(service.CmdService)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(task.CmdTask)
	rootCmd.AddCommand(model.CmdNew)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
