package cmd

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	err := godotenv.Load()
	if err != nil {
		return 1
	}

	cpuProfile := false

	rootcmd := &cobra.Command{
		Use:   "check",
		Short: "check",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cpuProfile {
				return nil
			}

			file, perr := os.Create("cpu.pprof")
			if perr != nil {
				return perr
			}

			_ = pprof.StartCPUProfile(file)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !cpuProfile {
				return nil
			}

			pprof.StopCPUProfile()

			file, merr := os.Create("mem.pprof")
			if merr != nil {
				return merr
			}
			defer file.Close()

			runtime.GC()

			err := pprof.WriteHeapProfile(file)
			return err
		},
	}

	rootcmd.PersistentFlags().BoolVarP(&cpuProfile, "profile", "p", false, "record cpu profile")
	rootcmd.AddCommand(APICmd(ctx))

	// Debug profiling and runtime metrics
	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	err = rootcmd.Execute()
	if err != nil {
		return 1
	}

	return 0
}
