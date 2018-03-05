package backend

import (
	"github.com/Polarishq/middleware/framework/log"
	"github.com/Polarishq/ot-try/models"
)

// Health check endpoint ...
func Health() (models.Health, error) {
	healthy := false
	h := models.Health{Healthy: &healthy}

	si := new(models.ServiceInfo)
	si.Name = &Constants.ComponentName
	si.Version = &Constants.Version
	h.ServiceInfo = si
	h.Revision = &GitShaCommit
	healthy = true
	log.Infof("service=%s healthy=%t Revision=%s", *h.ServiceInfo.Name, *h.Healthy, GitShaCommit)

	return h, nil
}
