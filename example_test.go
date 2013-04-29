package clog_test

import (
	"github.com/matrixik/clog"

	"math"
	"os"
)

var Log *clog.Clog = clog.NewClog()

func Example() {
	Log.AddOutput(os.Stdout, clog.LevelWarn)
	dailyFile := clog.NewDailyFile("/opt/logs/myprocess_%s.log")
	Log.AddOutput(dailyFile, clog.LevelTrace)
	Log.Info("Pi is %v", math.Pi)
}
