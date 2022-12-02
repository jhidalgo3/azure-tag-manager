package session

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/pkg/errors"
)

// AzureSession stores subscription id and Authorized object
type AzureSession struct {
	SubscriptionID string

	Credential *azidentity.DefaultAzureCredential
}

// func readJSON(path string) (*map[string]interface{}, error) {
// 	data, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		log.Fatalf("failed to read file: %v", err)
// 	}

// 	contents := make(map[string]interface{})
// 	json.Unmarshal(data, &contents)
// 	return &contents, nil
// }

func NewFromAzureCredential(subscriptionId string) (*AzureSession, error) {
	// Create a credentials object.
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.Wrap(err, "Authentication failure: %+v")
	}

	sess := AzureSession{
		SubscriptionID: subscriptionId,

		Credential: cred,
	}

	return &sess, err
}

// NewFromFile creates new session from file kept in AZURE_AUTH_LOCATION.
/*func NewFromFile() (*AzureSession, error) {
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get initial session")
	}

	a, err := auth.GetSettingsFromFile()
	sess := AzureSession{
		SubscriptionID: a.GetSubscriptionID(),
		Authorizer:     authorizer,
	}

	return &sess, err
}*/
