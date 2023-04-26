package processer

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/WangYihang/Proxy-Verifier/internal"
	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/protocol"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
	logger "github.com/sirupsen/logrus"
)

func Loader(inputFilepath string, taskQueue chan *model.Task, numWorkers int) {
	var index int

	defer func() {
		fmt.Printf("All %d tasks loaded\n", index)
		// Signal tasks are all loaded
		for i := 0; i < numWorkers; i++ {
			taskQueue <- nil
		}
	}()

	// Detect file format
	var parsers = map[string]func(record []string) (string, string, uint16){
		"csv":  util.ParseRecordCsv,
		"xmap": util.ParseRecordXmap,
	}
	fileFormat, err := util.DetectFileFormat(inputFilepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get parser
	ParseRecordFunc := parsers[fileFormat]

	// Open file
	reader, err := os.OpenFile(inputFilepath, os.O_RDONLY, 0)
	if err != nil {
		logger.Error(err)
	}
	r := csv.NewReader(reader)

	// Read records
	for index = 0; ; index++ {
		record, err := r.Read()
		if err != nil {
			break
		}

		// Parse record
		proxyProtocol, proxyHost, proxyPort := ParseRecordFunc(record)
		if err != nil {
			logger.Error(err)
			continue
		}

		targetProtocol := internal.Options.TargetUrlObject.Scheme
		proxyProtocol2SubProxyProtocol2Handler := protocol.TargetProtocl2ProxyProtocol2SubProxyProtocol2Handler[targetProtocol]
		for _, handler := range proxyProtocol2SubProxyProtocol2Handler[proxyProtocol] {
			task := model.NewTask(index, internal.Options.MeasurementId, proxyProtocol, proxyHost, proxyPort, internal.Options.TargetUrl, internal.Options.Timeout, handler)
			// techniqueName := fmt.Sprintf("%s-via-%s-%s", targetProtocol, proxyProtocol, subProxyProtocol)
			// task := model.NewTask(index, internal.Options.MeasurementId, proxyProtocol, proxyHost, proxyPort, internal.Options.TargetUrl, internal.Options.Timeout, techniqueName, handler)
			ScheduleTask(&task, taskQueue)
		}
	}
}
