package cmd

import (
	"context"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Execute(context context.Context) int {
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

	return 0
}
