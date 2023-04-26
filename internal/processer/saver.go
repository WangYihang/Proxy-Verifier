package processer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/WangYihang/Proxy-Verifier/internal"
	"github.com/WangYihang/Proxy-Verifier/internal/model"
	logger "github.com/sirupsen/logrus"
)

func Saver(resultQueue chan *model.Result, numWorkers int) {
	var numFinishedWorkers int = 0

	outputFile, err := os.OpenFile(internal.Options.OutputFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Error(err)
	}

	for {
		if numFinishedWorkers == numWorkers {
			logger.Infof("all %d workers finished\n", numFinishedWorkers)
			break
		}
		result := <-resultQueue

		if result == nil {
			fmt.Printf("%d / %d (%.2f%%) workers finished\r", numFinishedWorkers, numWorkers, float64(numFinishedWorkers)/float64(numWorkers)*100)
			numFinishedWorkers += 1
			continue
		}

		internal.State.TaskDone(result.Error == "")
		if result.Error == "" {
			fmt.Printf("%-96s\n", result.String())
		} else {
			fmt.Printf("%-96s\r", result.String())
		}

		// Dump result to file as json format
		resultString, err := json.Marshal(result)
		if err != nil {
			logger.Error(err)
		}
		_, err = outputFile.WriteString(fmt.Sprintf("%s\n", resultString))
		if err != nil {
			logger.Error(err)
		}

		// Flush to disk
		err = outputFile.Sync()
		if err != nil {
			logger.Error(err)
		}
	}
	outputFile.Close()
}
