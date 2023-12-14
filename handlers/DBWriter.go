package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"s3-queue-reader/interfaces"
	"s3-queue-reader/models"
	"s3-queue-reader/repositories"
	"strconv"
)

func NewDBWriter(dbHandler interfaces.IDbHandler, queueService interfaces.IQueueService, resultedQueueName string) *DBWriter {
	dbArtistsRepo := repositories.DBArtistsRepository{dbHandler}

	return &DBWriter{
		DBHandler:           dbHandler,
		DBArtistsRepository: dbArtistsRepo,
		QueueService:        queueService,
		ResultedQueueName:   resultedQueueName,
	}
}

type DBWriter struct {
	DBHandler           interfaces.IDbHandler
	DBArtistsRepository repositories.DBArtistsRepository
	QueueService        interfaces.IQueueService
	ResultedQueueName   string
}

func (h *DBWriter) Run() {
	fmt.Println("dbwriter run")

	//read message
	messages := h.QueueService.ConsumeQueue(h.ResultedQueueName)
	for data := range messages {
		//process message
		log.Printf("Received a message: %d bytes", len(data.Body))

		var resultMessage = models.ResultMessage{}
		err := json.Unmarshal(data.Body, &resultMessage)
		if err != nil {
			fmt.Println("Couldn't unmarshall data")
			continue
		}

		//save data
		h.UpdateArtistsData(resultMessage.Data)
		fmt.Println("Message received", resultMessage)
	}
}

func (h *DBWriter) UpdateArtistsData(rows map[string]int64) {
	for id, amount := range rows {
		ArtistId, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Error ArtistId conversion ", err.Error())
		}

		h.DBArtistsRepository.CreateOrUpdateArtist(int64(ArtistId), amount)
	}
}
