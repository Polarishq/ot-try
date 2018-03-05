package deployment

import (
	"os"

	"fmt"
	"strings"

	"github.com/Polarishq/middleware/framework/log"
)

//Environment variable keys
const (
	EnvUser             = "USER"
	EnvName             = "ENVIRONMENT_NAME"
	EnvCleanupBlacklist = "CleanupBlacklist"
)

// findEnvironmentName finds out what environment are we operating under
func findEnvironmentName() (string, error) {
	var envName string
	prefixes := [2]string{EnvName, EnvUser}
	log.WithFields(log.Fields{
		"prefixes": prefixes,
	}).Info("Looking up environment name")
	for _, prefix := range prefixes {
		value, ok := os.LookupEnv(prefix)
		if ok && value != "" {
			log.WithFields(log.Fields{
				"env":   prefix,
				"value": value,
			}).Info("Found environment name")
			envName = value
			break
		} else {
			log.WithField("env", prefix).Warning("Environment variable not set, trying next")
		}
	}

	if envName == "" {
		return envName, fmt.Errorf("Unable to determine ENVIRONMENT_NAME")
	}

	return envName, nil
}

// GetResource returns the name of a given resouce with the appended env name, name+envname
func GetResource(resourceName string) (string, error) {
	envName, err := findEnvironmentName()
	if err != nil {
		return envName, err
	}

	// Ensure the env and resource resourceName are separated by '-'
	if resourceName[:1] != "-" {
		resourceName = "-" + resourceName
	}

	return envName + resourceName, nil
}

//ShouldCleanup determines if resource cleanup is allowed in the current environment
func ShouldCleanup() (bool, error) {
	envName, err := findEnvironmentName()
	if err != nil {
		return false, err
	}

	sBlacklist, ok := os.LookupEnv(EnvCleanupBlacklist)
	if ok {
		blacklist := strings.Split(sBlacklist, ",")
		log.WithFields(log.Fields{
			"blacklist": blacklist,
			"envName":   envName,
		}).Info("Retrieved blacklist")

		for _, env := range blacklist {
			env = strings.TrimSpace(env)
			if env == envName {
				log.WithFields(log.Fields{
					"blacklist": blacklist,
					"envName":   envName,
				}).Info("Environment is blacklisted")
				return false, nil
			}
		}
		log.WithFields(log.Fields{
			"blacklist": blacklist,
			"envName":   envName,
		}).Info("Environment is not blacklisted")
		return true, nil
	}

	//Blacklist not found, allow cleanup
	return true, nil
}
