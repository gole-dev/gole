package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gole-dev/gole/cmd/gole/internal/base"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the gole tools",
	Long:  "Upgrade the gole tools. Example: gole upgrade",
	Run:   Run,
}

// Run upgrade the gole tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/gole-dev/gole/cmd/gole",
		"github.com/gole-dev/gole/cmd/protoc-gen-go-gin",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/envoyproxy/protoc-gen-validate",
		"github.com/google/gnostic",
		"github.com/google/gnostic/cmd/protoc-gen-openapi",
	)
	if err != nil {
		fmt.Println(err)
	}
}
