package cronjob

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var CronJobRunnser *cron.Cron

func NewCronJobRunner() {
	if CronJobRunnser == nil {
		CronJobRunnser = cron.New()
		CronJobRunnser.Start()
	}
	fmt.Println("Cron runner started")
}

func AddNewJob(duration string, f func()) (int, error) {
	if CronJobRunnser == nil {
		CronJobRunnser = cron.New()
	}

	entryID, err := CronJobRunnser.AddFunc(duration, f)
	return int(entryID), err
}
