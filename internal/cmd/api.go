package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/briheet/gxAssign/internal/cmdutil"
	"github.com/spf13/cobra"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the Restful API",
		RunE: func(cmd *cobra.Command, args []string) error {
			port = 4000

			if os.Getenv("PORT") != "" {
				port, _ := strconv.Atoi(os.Getenv("PORT"))
			}

			logger := cmdutil.NewLogger("api")

			return nil
		},
	}

	return cmd
}
