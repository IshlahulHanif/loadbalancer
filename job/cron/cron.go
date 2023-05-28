package cronjob

import "github.com/IshlahulHanif/poneglyph"

func (m Module) StartAllCronJobScheduler() (err error) {
	_, err = m.cronModule.AddFunc(m.config.CronConfig.HealthCheckAll, m.handlerHealthCheckJob)
	if err != nil {
		return poneglyph.Trace(err)
	}

	m.cronModule.Start()

	return nil
}

func (m Module) StopCronJobScheduler() {
	m.cronModule.Stop()
}
