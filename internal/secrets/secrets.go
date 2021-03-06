package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type RepoSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RepoSecretFromSecretsManager(ctx context.Context, secretARN string) (*RepoSecret, error) {
	secretValue, err := getSecretsManagerValue(ctx, secretARN)
	if err != nil {
		return nil, fmt.Errorf("error getting secret from secrets manager: %s", err.Error())
	}
	retVal := &RepoSecret{}
	err = json.Unmarshal([]byte(secretValue), retVal)
	return retVal, err
}

func getSecretsManagerValue(ctx context.Context, secretARN string) (string, error) {
	secretRegion, err := getRegionFromSecretARN(secretARN)
	if err != nil {
		return "", fmt.Errorf("unable to create initMeta to initialize the secret manager client: %v", err)
	}
	s := session.Must(session.NewSession())

	mgr := secretsmanager.New(s, aws.NewConfig().WithRegion(secretRegion))
	if err != nil {
		return "", fmt.Errorf("unable to create secrets manager client: %v", err)
	}

	secretVal, err := mgr.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &secretARN,
	})
	if err != nil {
		return "", err
	}

	return *secretVal.SecretString, nil
}

func getRegionFromSecretARN(secretARN string) (string, error) {
	parsedARN, err := arn.Parse(secretARN)
	if err != nil {
		return "", fmt.Errorf("cannot parse secret ARN: %v", err)
	}
	return parsedARN.Region, nil
}
