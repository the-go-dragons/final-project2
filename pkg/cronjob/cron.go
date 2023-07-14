package cronjob

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

var CronJobRunner *cron.Cron

func NewCronJobRunner() {
	if CronJobRunner == nil {
		CronJobRunner = cron.New()
		CronJobRunner.Start()
	}
	fmt.Println("Cron runner started")
}

type JobWithRepetitions struct {
	EntryID      cron.EntryID
	Job          cron.Job
	Repetitions  uint
	Counter      uint
	StopOnFinish bool
}

func (j *JobWithRepetitions) Run() {
	j.Counter++
	j.Job.Run()

	if j.StopOnFinish && j.Counter >= j.Repetitions {
		CronJobRunner.Remove(j.EntryID)
	}
}

func AddNewJob(
	user domain.User,
	period string,
	massage string,
	senderNumber string,
	receiverNumbers string,
	repetitionCount uint,
	smsS usecase.SMSService,
) (int, error) {
	if CronJobRunner == nil {
		NewCronJobRunner()
	}
	if repetitionCount <= 0 {
		return 0, errors.New("not enough repetition count")
	}

	parsedSchedule, err := cron.ParseStandard(period)
	if err != nil {
		return 0, err
	}

	job := &JobWithRepetitions{
		Job: cron.FuncJob(func() { // Wrap the actual job function
			rabbitmq.NewMassage(rabbitmq.SMSBody{
				Sender:    senderNumber,
				Receivers: receiverNumbers,
				Massage:   massage,
			})
			smsHistoryRecord := domain.SMSHistory{
				UserId:          user.ID,
				User:            user,
				SenderNumber:    senderNumber,
				ReceiverNumbers: receiverNumbers,
				Content:         massage,
			}
			smsS.CreateSMS(smsHistoryRecord)

		}),
		Repetitions:  repetitionCount,
		Counter:      0,
		StopOnFinish: true,
	}

	entryID := CronJobRunner.Schedule(parsedSchedule, job)

	job.EntryID = entryID

	return int(entryID), nil
}
