package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"log"
	"math"
	"os"
	"s3-queue-reader/infrastructures"
	"s3-queue-reader/interfaces"
	"s3-queue-reader/models"
)

func NewQueueProcessor(queueService interfaces.IQueueService, currentQueueName, resultedQueueName string) *QueueProcessor {

	tempFileName := infrastructures.GetTempDirectory() + "/" + "processor_" + currentQueueName + ".parquet"

	return &QueueProcessor{
		CurrentQueueName:  currentQueueName,
		ResultedQueueName: resultedQueueName,
		QueueService:      queueService,
		tempFileName:      tempFileName,
	}
}

type QueueProcessor struct {
	CurrentQueueName  string
	ResultedQueueName string
	QueueService      interfaces.IQueueService
	tempFileName      string
}

func (h QueueProcessor) Run() {
	fmt.Println("processor run")

	//read message
	messages := h.QueueService.ConsumeQueue(h.CurrentQueueName)
	for data := range messages {
		//process message
		log.Printf("Received a message: %d bytes", len(data.Body))
		resultMap := h.processData(data.Body)

		jsonData, err := json.Marshal(resultMap)
		if err != nil {
			fmt.Println("Couldn't marshall data")
			continue
		}

		//send to the resulted queue
		err = h.QueueService.PushDataToQueue(jsonData, h.ResultedQueueName)
		if err != nil {
			fmt.Println("Couldn't send json data to the result queue")
			continue
		}
	}
}

func (h QueueProcessor) processData(data []byte) models.ResultMessage {
	resultMessage := models.ResultMessage{Data: make(map[string]int64)}

	err := os.WriteFile(h.tempFileName, data, 0644)
	if err != nil {
		log.Println("Can't write temp file ", h.tempFileName)
	}

	//fill data
	fr, err := local.NewLocalFileReader(h.tempFileName)
	if err != nil {
		log.Println("Can't open file ", h.tempFileName)
		return resultMessage
	}

	pr, err := reader.NewParquetReader(fr, new(models.YParketData), 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return resultMessage
	}
	num := int(pr.GetNumRows())

	dataBatchSize := 10
	maxIteration := int(math.Ceil(float64(num) / float64(dataBatchSize)))
	for i := 0; i < maxIteration; i++ {
		rows := make([]models.YParketData, dataBatchSize)
		if err = pr.Read(&rows); err != nil {
			log.Println("Read error", err)

			continue
		}

		for _, row := range rows {
			resultMessage.Data[row.Artist] = resultMessage.Data[row.Artist] + row.Amount
		}
	}

	log.Println("Import finished")
	pr.ReadStop()
	fr.Close()

	os.Remove(h.tempFileName)

	return resultMessage
}
