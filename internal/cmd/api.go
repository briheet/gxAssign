package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/briheet/gxAssign/internal/api"
	"github.com/briheet/gxAssign/internal/cmdutil"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int
	var dbName string

	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the Restful API",
		RunE: func(cmd *cobra.Command, args []string) error {
			port = 4000
			dbName = "gxAssign"

			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			logger := cmdutil.NewLogger("api")
			defer func() { _ = logger.Sync() }()

			if os.Getenv("DB_NAME") != "" {
				dbName = os.Getenv("DB_NAME")
			}

			api := api.NewAPI(ctx, logger, dbName)
			srv := api.Server(port)

			go func() {
				_ = srv.ListenAndServe()
			}()

			logger.Info("started api", zap.Int("port", port))
			logger.Info("started db", zap.String("db", dbName))

			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			return nil
		},
	}

	return cmd
}
