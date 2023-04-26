package internal

import (
	"net/url"
	"os"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/google/uuid"
	flags "github.com/jessevdk/go-flags"
	logger "github.com/sirupsen/logrus"
)

var State *model.State
var Options *model.MainOptions

func Init() {
	Options = &model.MainOptions{}
	// Parse options
	_, err := flags.ParseArgs(Options, os.Args)
	if err != nil {
		os.Exit(1)
	}

	if Options.MeasurementId == "" {
		Options.MeasurementId = uuid.New().String()
	}

	// Validate options
	Options.TargetUrlObject, err = url.ParseRequestURI(Options.TargetUrl)
	if err != nil {
		os.Exit(1)
	}

	// Create state
	State = model.NewState(Options.MonitorInterval)
	// Setup logger
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})
}
