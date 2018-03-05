package framework

import (
	"fmt"
	"os"
	"strings"
	"time"

	sess "github.com/Polarishq/middleware/aws"
	"github.com/Polarishq/middleware/framework/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

const (
	//GlobalComponent is the global component name
	GlobalComponent = "novaglobals"
	//DefaultEnvironment is the default environment name
	DefaultEnvironment = "default"
	//DeployatronSeparator is the default separator for deployatron parameter names
	DeployatronSeparator = "."
	//AWSPageSize is the length of the batch size for getParameter calls
	//10 is the absolute maximum for SSM enforced by AWS
	AWSPageSize = 10
)

//ConfigStore handles pulling down the env variables at server start
type ConfigStore struct {
	Component   string
	Environment string
	SSMSess     ssmiface.SSMAPI
}

//NewConfigStore creates a new Config Store Handler instance
func NewConfigStore(component string) (*ConfigStore, error) {
	env, ok := os.LookupEnv("ENVIRONMENT_NAME")
	if !ok {
		return nil, fmt.Errorf("ENVIRONMENT_NAME env variable is not set")
	}
	awsSess, err := sess.NewSess("us-west-1")
	if err != nil {
		return nil, err
	}
	svc := ssm.New(awsSess)

	log.WithFields(log.Fields{
		"component":   component,
		"environment": env,
	}).Info("Creating new ConfigStore")

	return &ConfigStore{
		Environment: env,
		Component:   component,
		SSMSess:     svc,
	}, nil
}

func makeParamName(component string, environment string) string {
	return component + DeployatronSeparator + environment + DeployatronSeparator
}

func (cs *ConfigStore) setParamNames() []string {
	globalDefaults := makeParamName(GlobalComponent, DefaultEnvironment)
	globalOverrides := makeParamName(GlobalComponent, cs.Environment)
	componentDefaults := makeParamName(cs.Component, DefaultEnvironment)
	componentOverrides := makeParamName(cs.Component, cs.Environment)

	//Ordering matters here
	return []string{globalDefaults, globalOverrides, componentDefaults, componentOverrides}
}

//SetupEnvironmentFromSSM sets environment variables from parameters in SSM
func (cs *ConfigStore) SetupEnvironmentFromSSM() error {
	paramGroups := cs.setParamNames()
	log.WithField("paramGroups", paramGroups).Infof("Setup environment from SSM")
	for _, v := range paramGroups {
		paramNames, err := describeParameters(cs.SSMSess, v)
		if err != nil {
			return err
		}
		if len(paramNames) > 0 {
			params, err := getParameters(cs.SSMSess, paramNames)
			if err != nil {
				return err
			}
			setEnvParameters(params)
		}
	}

	return nil
}

// Function to determine if the error is due to Rate exceeded
func israteerror(e error) bool {
	if e == nil {
		return false
	}
	errstring := e.Error()
	if strings.Contains(errstring, "Rate exceeded") {
		// let's slow down a bit before retrying
		time.Sleep(100 * time.Millisecond)
		return true
	}
	return false
}

//Encapsulation to the method in order to avoid crashing from SSM
func describeParametersenc(svc ssmiface.SSMAPI, filter *ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error) {
	for {
		desribeRes, err := svc.DescribeParameters(filter)
		if israteerror(err) == false {
			return desribeRes, err
		}
	}
}

//Encapsulation to the method in order to avoid crashing from SSM
func getParametersenc(svc ssmiface.SSMAPI, getFilter *ssm.GetParametersInput) (*ssm.GetParametersOutput, error) {
	for {
		getRes, err := svc.GetParameters(getFilter)
		if israteerror(err) == false {
			return getRes, err
		}
	}
}

func describeParameters(svc ssmiface.SSMAPI, filter string) ([]*string, error) {
	filterKey := "Name"
	nextToken := " "
	paramNames := make([]*string, 0)
	log.WithField("filter", filter).Debugf("Describing parameters")
	for len(nextToken) > 0 {
		filter := ssm.DescribeParametersInput{
			Filters: []*ssm.ParametersFilter{
				{
					Key: &filterKey,
					Values: []*string{
						&filter,
					},
				},
			},
			NextToken: &nextToken,
		}
		desribeRes, err := describeParametersenc(svc, &filter)
		for _, v := range desribeRes.Parameters {
			if v.Type != nil && (*v.Type == "String" || *v.Type == "SecureString") {
				paramNames = append(paramNames, v.Name)
			}
		}
		if err != nil {
			return nil, err
		}
		if desribeRes.NextToken != nil {
			nextToken = *desribeRes.NextToken
		} else {
			nextToken = ""
		}
	}

	var names string
	for _, p := range paramNames {
		names += *p + ","
	}
	log.WithField("paramNames", names).Debugf("Retrieved parameter names")

	return paramNames, nil
}

func getParameters(svc ssmiface.SSMAPI, paramNames []*string) ([]*ssm.Parameter, error) {
	// Have to page the get parameter calls in sets of ten
	params := make([]*ssm.Parameter, 0)

	for i := 0; i < len(paramNames); i += AWSPageSize {
		var getFilter ssm.GetParametersInput
		var nameSubset []*string
		if i+AWSPageSize-1 < len(paramNames) {
			nameSubset = paramNames[i : i+AWSPageSize-1]
		} else {
			nameSubset = paramNames[i:]
		}

		getFilter = ssm.GetParametersInput{
			Names:          nameSubset,
			WithDecryption: aws.Bool(true),
		}

		var names string
		for _, p := range nameSubset {
			names += *p + ","
		}
		log.WithFields(log.Fields{
			"paramNames": names,
			"getFilter":  getFilter,
		}).Debug("Get parameters from SSM")
		getRes, err := getParametersenc(svc, &getFilter)

		if err != nil {
			return nil, err
		}
		params = append(params, getRes.Parameters...)

	}

	return params, nil
}

func setEnvParameters(params []*ssm.Parameter) {
	for _, p := range params {
		if p.Name != nil && p.Value != nil {
			parts := strings.Split(*p.Name, DeployatronSeparator)
			k := parts[2]
			v := *p.Value
			logv := *p.Value
			if *p.Type == ssm.ParameterTypeSecureString {
				logv = "*** REDACTED ***"
			}
			log.WithFields(log.Fields{
				"key":   k,
				"value": logv,
			}).Info("Setting environment variable")
			os.Setenv(k, v)
		}
	}
}
