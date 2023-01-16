package graphhelper

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	// "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/ref"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type GraphHelper struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}
func (g *GraphHelper) InitializeGraphForAppAuth() error {
	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return err
	}

	g.clientSecretCredential = credential

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(g.clientSecretCredential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return err
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client using request adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	g.appClient = client

	return nil
}

func (g *GraphHelper) GetAppToken() (*string, error) {
	token, err := g.clientSecretCredential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{
			"https://graph.microsoft.com/.default",
		},
	})
	if err != nil {
		return nil, err
	}

	return &token.Token, nil
}

func (g *GraphHelper) GetUsers() (models.UserCollectionResponseable, error) {
	var topValue int32 = 25
	query := users.UsersRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "id", "mail"},
		// Get at most 25 results
		Top: &topValue,
		// Sort by display name
		Orderby: []string{"displayName"},
	}

	return g.appClient.Users().
		Get(context.Background(),
			&users.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

func (g *GraphHelper) ListMembers() (models.DirectoryObjectCollectionResponseable, error) {

	return g.appClient.GroupsById("d523466e-07dd-4fde-8856-876de5e42d86").Members().Get(context.Background(), nil)
}

func (g *GraphHelper) UpdateMembers() error {

	// reference := ref.NewRef()
	// reference.GetAdditionalData()["@odata.id"] = "https://graph.microsoft.com/v1.0/directoryObjects/95bfa28-a695-447e-aa13-7e68bc57f97d"
	requestBody := models.NewReferenceCreate()
	oidData := "https://graph.microsoft.com/v1.0/directoryObjects/95bfa28-a695-447e-aa13-7e68bc57f97d"

	requestBody.SetOdataId(&oidData)
	// requestBody := models.NewReferenceCreate()
	// requestBody.SetAdditionalData(map[string]interface{}{
	// 	"@odata.id": "https://graph.microsoft.com/v1.0/directoryObjects/95bfa28-a695-447e-aa13-7e68bc57f97d",
	// },
	// // )
	err := g.appClient.GroupsById("d523466e-07dd-4fde-8856-876de5e42d86").Members().Ref().Post(context.Background(), requestBody, nil)
	if err != nil {
		return err
	}
	return nil
}

func (g *GraphHelper) DeleteMembers() error {

	// reference := ref.NewRef()
	// reference.GetAdditionalData()["@odata.id"] = "https://graph.microsoft.com/v1.0/directoryObjects/95bfa28-a695-447e-aa13-7e68bc57f97d"
	requestBody := models.NewReferenceCreate()
	oidData := "https://graph.microsoft.com/v1.0/directoryObjects/95bfa28-a695-447e-aa13-7e68bc57f97d"

	requestBody.SetOdataId(&oidData)
	// requestBody := models.NewReferenceCreate()
	// requestBody.SetAdditionalData(map[string]interface{}{
	// 	"@odata.id": "https://graph.microsoft.com/v1.0/directoryObjects/95bfa28-a695-447e-aa13-7e68bc57f97d",
	// },
	// // )
	err := g.appClient.GroupsById("d523466e-07dd-4fde-8856-876de5e42d86").Members().Ref().Post(context.Background(), requestBody, nil)
	if err != nil {
		return err
	}
	return nil
}
