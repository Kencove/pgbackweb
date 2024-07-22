package main

import (
	"github.com/eduardolat/pgbackweb/internal/cron"
	"github.com/eduardolat/pgbackweb/internal/logger"
	"github.com/eduardolat/pgbackweb/internal/service"
	"github.com/google/uuid"
)

func initSchedule(cr *cron.Cron, servs *service.Service) {
	/*
		Initial executions
	*/

	servs.ExecutionsService.SoftDeleteExpiredExecutions()
	servs.AuthService.DeleteOldSessions()

	/*
		Schedules
	*/

	err := cr.UpsertJob(uuid.New(), "UTC", "* * * * *", func() {
		servs.ExecutionsService.SoftDeleteExpiredExecutions()
	})
	if err != nil {
		logger.FatalError(
			"error scheduling soft deletion of expired executions",
			logger.KV{"error": err},
		)
	}

	err = cr.UpsertJob(uuid.New(), "UTC", "* * * * *", func() {
		servs.AuthService.DeleteOldSessions()
	})
	if err != nil {
		logger.FatalError(
			"error scheduling deletion of old sessions", logger.KV{"error": err},
		)
	}

	err = servs.BackupsService.ScheduleAll()
	if err != nil {
		logger.FatalError("error scheduling all backups", logger.KV{"error": err})
	}
}