package cronjob

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

var CronJobRunnser *cron.Cron

func NewCronJobRunner() {
	if CronJobRunnser == nil {
		CronJobRunnser = cron.New()
		CronJobRunnser.Start()
	}
	fmt.Println("Cron runner started")
}

func AddNewJob(
	userID uint,
	period string,
	massage string,
	senderNumber string,
	receiverNumbers []string,
	repeatationCount uint,
) (int, error) {
	if CronJobRunnser == nil {
		NewCronJobRunner()
	}
	if repeatationCount <= 0 {
		return 0, errors.New("not enough repeatation count")
	}

	db, _ := database.GetDatabaseConnection()
	db.Create(&domain.CronJob{
		UserID:           userID,
		Period:           period,
		RepeatationCount: repeatationCount,
		Massage:          massage,
		SenderNumber:     senderNumber,
		ReceiverNumbers:  receiverNumbers,
	})
	entryID, err := CronJobRunnser.AddFunc(period, func() {
		rabbitmq.NewMassage(rabbitmq.SMSBody{
			Sender:    senderNumber,
			Receivers: receiverNumbers,
			Massage:   massage,
		})
		// TODO
		// db, _ := database.GetDatabaseConnection()
		// db.First("id = ?", )
	})
	return int(entryID), err
}
