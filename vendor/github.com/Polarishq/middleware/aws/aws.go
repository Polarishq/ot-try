package aws

import (
	"github.com/Polarishq/middleware/framework/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSess get a new aws session
func NewSess(region string) (*session.Session, error) {
	logFields := log.Fields{"action": "NewSess"}

	log.FuncStart(logFields)

	// todo;
	awsSess := session.New(&aws.Config{
		Region: &region,
		Logger: aws.LoggerFunc(func(args ...interface{}) {
			log.Debug(args...)
		}),
		// TODO add config knob to turn it on/off because it's really verbose
		LogLevel: aws.LogLevel(aws.LogOff),
	})

	*awsSess.Config.LogLevel = aws.LogDebugWithHTTPBody

	/* check that credentials were located */
	_, err := awsSess.Config.Credentials.Get()
	if err != nil {
		log.WithAddFields(logFields,
			log.Fields{"step": "check_credentials", "status": "fail"}).Error(err.Error())
		return nil, err
	}

	log.FuncSucc(logFields)

	return awsSess, nil
}
