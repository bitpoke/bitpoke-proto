/*
Copyright 2019 Pressinfra SRL.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	// Zap exposes the set zap logger for downstream use
	restoreLogger      = func() {}
	restoreLoggerMutex = &sync.Mutex{}

	// SetLogger is an alias for compatibility with sigs.k8s.io/controller-runtime/pkg/runtime/log
	SetLogger = logf.SetLogger

	// Log is an alias for compatibility with sigs.k8s.io/controller-runtime/pkg/runtime/log
	Log = logf.Log

	// KBLog is an alias for compatibility with sigs.k8s.io/controller-runtime/pkg/runtime/log
	KBLog = logf.KBLog
)

// ZapLogger is rewritten to allow exposing the underlying zap logger
// It should be compatible wit the one in sigs.k8s.io/controller-runtime/pkg/runtime/log
func ZapLogger(development bool, opts ...zap.Option) logr.Logger {
	return ZapLoggerTo(os.Stderr, development)
}

// ZapLoggerTo is rewritten to allow exposing the underlying zap logger
// It should be compatible wit the one in sigs.k8s.io/controller-runtime/pkg/runtime/log
func ZapLoggerTo(destWriter io.Writer, development bool, opts ...zap.Option) logr.Logger {
	// this basically mimics New<type>Config, but with a custom sink
	sink := zapcore.AddSync(destWriter)

	var enc zapcore.Encoder
	var lvl zap.AtomicLevel
	if development {
		encCfg := zap.NewDevelopmentEncoderConfig()
		enc = zapcore.NewConsoleEncoder(encCfg)
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
		opts = append(opts, zap.Development(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		encCfg := zap.NewProductionEncoderConfig()
		enc = zapcore.NewJSONEncoder(encCfg)
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
		opts = append(opts, zap.AddStacktrace(zap.WarnLevel),
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSampler(core, time.Second, 100, 100)
			}))
	}
	opts = append(opts, zap.AddCallerSkip(1), zap.ErrorOutput(sink))
	log := zap.New(zapcore.NewCore(&logf.KubeAwareEncoder{Encoder: enc, Verbose: development}, sink, lvl))
	log = log.WithOptions(opts...)
	restoreLoggerMutex.Lock()
	defer restoreLoggerMutex.Unlock()
	restoreLogger()
	restoreLogger = zap.ReplaceGlobals(log)
	return zapr.NewLogger(log)
}
