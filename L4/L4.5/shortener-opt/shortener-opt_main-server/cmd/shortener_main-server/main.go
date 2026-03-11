package main

import (
	"context"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/app"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/zlog"
)

// Profiling is enabled via environment variables:
// CPU_PROFILE - path to write CPU profile
// MEM_PROFILE - path to write memory profile
// TRACE_FILE - path to write trace

func main() {
	env := os.Getenv("ENV")
	if env == "local" {
		zlog.InitConsole()
	} else {
		zlog.Init()
	}

	lg := zlog.Logger.With().Str("component", "main").Logger()

	startProfiling(&lg)

	lg.Info().Msgf("%s shortener-opt_main-server started", pkgConst.AppStart)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.New(env)
	if err != nil {
		lg.Error().Err(err).Msgf("%s application initialization failed", pkgConst.Error)
		os.Exit(1)
	}

	if err := application.Run(cancel); err != nil {
		lg.Error().Err(err).Msgf("%s application stopped with error", pkgConst.Error)
		os.Exit(1)
	}

	application.GracefullShutdown(ctx, cancel)

	stopProfiling(&lg)

	lg.Info().Msgf("%s application stopped gracefully", pkgConst.AppStop)
}

func startProfiling(lg *zlog.Zerolog) {
	if cpuProfile := os.Getenv("CPU_PROFILE"); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			lg.Error().Err(err).Msg("failed to create CPU profile file")
		} else {
			if err := pprof.StartCPUProfile(f); err != nil {
				lg.Error().Err(err).Msg("failed to start CPU profile")
				f.Close()
			} else {
				lg.Info().Str("file", cpuProfile).Msg("CPU profiling started")
			}
		}
	}

	if traceFile := os.Getenv("TRACE_FILE"); traceFile != "" {
		f, err := os.Create(traceFile)
		if err != nil {
			lg.Error().Err(err).Msg("failed to create trace file")
		} else {
			if err := trace.Start(f); err != nil {
				lg.Error().Err(err).Msg("failed to start trace")
				f.Close()
			} else {
				lg.Info().Str("file", traceFile).Msg("trace started")
			}
		}
	}
}

func stopProfiling(lg *zlog.Zerolog) {
	if cpuProfile := os.Getenv("CPU_PROFILE"); cpuProfile != "" {
		pprof.StopCPUProfile()
		lg.Info().Str("file", cpuProfile).Msg("CPU profiling stopped")
	}

	if traceFile := os.Getenv("TRACE_FILE"); traceFile != "" {
		trace.Stop()
		lg.Info().Str("file", traceFile).Msg("trace stopped")
	}

	if memProfile := os.Getenv("MEM_PROFILE"); memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			lg.Error().Err(err).Msg("failed to create memory profile file")
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
			lg.Info().Str("file", memProfile).Msg("memory profile written")
		}
	}
}
