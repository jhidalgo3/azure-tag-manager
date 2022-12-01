package azure

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nordcloud/azure-tag-manager/internal/azure/session"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ResourceGroupScanner represents resource group scanner that scans all resources in a resource group
type ResourceGroupScanner struct {
	Session         *session.AzureSession
	ResourcesClient *armresources.Client
	GroupsClient    *armresources.ResourceGroupsClient
}

// Scanner represents generic scanner of Azure resource groups
type Scanner interface {
	GetResources() ([]Resource, error)
	GetResourcesByResourceGroup(string) ([]Resource, error)
	GetGroups() ([]string, error)
	GetResourceGroupTags(string) (map[string]*string, error)
}

// String converts string v to the string pointer
func String(v string) *string {
	return &v
}

// GetResourceGroupTags returns a map of key value tags of a reource group rg
func (r ResourceGroupScanner) GetResourceGroupTags(rg string) (map[string]*string, error) {
	result, err := r.GroupsClient.Get(context.Background(), rg, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "GetResourceGroupTags(rg=%s): Get() failed", rg)
	}
	return result.Tags, nil
}

// NewResourceGroupScanner creates ResourceGroupScanner with Azure Serssion s
func NewResourceGroupScanner(s *session.AzureSession) *ResourceGroupScanner {

	resClient, _ := armresources.NewClient(s.SubscriptionID, s.Credential, nil)
	//resClient := resources.NewClient(s.SubscriptionID)
	//resClient.Authorizer = s.Authorizer

	grClient, _ := armresources.NewResourceGroupsClient(s.SubscriptionID, s.Credential, nil)
	//grClient := resources.NewGroupsClient(s.SubscriptionID)
	//grClient.Authorizer = s.Authorizer

	scanner := &ResourceGroupScanner{
		Session:         s,
		ResourcesClient: resClient,
		GroupsClient:    grClient,
	}

	return scanner
}

// ScanResourceGroup returns a list of resources and their tags from a resource group rg
func (r ResourceGroupScanner) ScanResourceGroup(rg string) []Resource {
	tab := make([]Resource, 0)

	pager := r.ResourcesClient.NewListByResourceGroupPager(rg, nil)
	for pager.More() {
		resp, _ := pager.NextPage(context.Background())
		if resp.ResourceListResult.Value != nil {
			for _, resource := range resp.ResourceListResult.Value {
				//resourceGroups = append(resourceGroups, resp.ResourceGroupListResult.Value...)
				tab = append(tab, Resource{
					Platform:      "azure",
					ID:            *resource.ID,
					Name:          resource.Name,
					Region:        *resource.Location,
					Tags:          resource.Tags,
					ResourceGroup: String(rg),
				})

				log.Info(*resource.Name, " ", *resource.ID, " tags: ", resource.Tags)
			}
		}
	}

	//return resourceGroups, pager.Err()

	/*for list, err := r.ResourcesClient.ListByResourceGroupComplete(context.Background(), rg, "", "", nil); list.NotDone(); err = list.NextWithContext(context.Background()) {
		if err != nil {
			log.Fatal(err)
		}
		resource := list.Value()
		tab = append(tab, Resource{
			Platform:      "azure",
			ID:            *resource.ID,
			Name:          resource.Name,
			Region:        *resource.Location,
			Tags:          resource.Tags,
			ResourceGroup: String(rg),
		})
	}*/
	return tab
}

// GetResources retruns list of resources in resource group
func (r ResourceGroupScanner) GetResources() ([]Resource, error) {
	var wg sync.WaitGroup

	groups, err := r.GetGroups()
	if err != nil {
		return nil, errors.Wrap(err, "GetResources(): GetGroups() failed")
	}

	tab := make([]Resource, 0)
	out := make(chan []Resource)
	for _, rg := range groups {
		wg.Add(1)
		go func(rg string) {
			defer wg.Done()
			out <- r.ScanResourceGroup(rg)
		}(rg)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	for s := range out {
		tab = append(tab, s...)
	}

	return tab, nil
}

// GetGroups returns list of resource groups in a subscription
func (r ResourceGroupScanner) GetGroups() ([]string, error) {
	tab := make([]string, 0)

	pager := r.GroupsClient.NewListPager(nil)

	//var resourceGroups []*armresources.ResourceGroup

	for pager.More() {
		resp, _ := pager.NextPage(context.Background())
		if resp.ResourceGroupListResult.Value != nil {
			for _, resource := range resp.ResourceGroupListResult.Value {
				//resourceGroups = append(resourceGroups, resp.ResourceGroupListResult.Value...)
				rgName := *resource.Name
				tab = append(tab, rgName)

				log.Info("rgName: ", rgName)
			}
		}
	}

	/*for list, err := r.GroupsClient.ListComplete(context.Background(), "", nil); list.NotDone(); err = list.NextWithContext(context.Background()) {
		if err != nil {
			return nil, errors.Wrap(err, "GetGroups(): GroupsClient.ListComplete failed")
		}
		rgName := *list.Value().Name
		tab = append(tab, rgName)
	}*/

	return tab, nil
}

// GetResourcesByResourceGroup returns resources in a resource group rg
func (r ResourceGroupScanner) GetResourcesByResourceGroup(rg string) ([]Resource, error) {
	tab := make([]Resource, 0)

	pager := r.ResourcesClient.NewListByResourceGroupPager(rg, nil)
	for pager.More() {
		resp, _ := pager.NextPage(context.Background())
		if resp.ResourceListResult.Value != nil {
			for _, resource := range resp.ResourceListResult.Value {
				//resourceGroups = append(resourceGroups, resp.ResourceGroupListResult.Value...)
				tab = append(tab, Resource{
					Platform:      "azure",
					ID:            *resource.ID,
					Name:          resource.Name,
					Region:        *resource.Location,
					Tags:          resource.Tags,
					ResourceGroup: String(rg),
				})

				log.Info(*resource.Name, " ", *resource.ID, " tags: ", resource.Tags)
			}
		}
	}

	/*for list, err := r.ResourcesClient.ListByResourceGroupComplete(context.Background(), rg, "", "", nil); list.NotDone(); err = list.NextWithContext(context.Background()) {
		if err != nil {
			return nil, errors.Wrapf(err, "GetResourcesByResourceGroup(rg=%q): ListByResourceGroupComplete() failed", rg)
		}
		resource := list.Value()
		tab = append(tab, Resource{
			Platform: "azure",
			ID:       *resource.ID,
			Name:     resource.Name,
			Region:   *resource.Location,
			Tags:     resource.Tags,
		})
	}*/
	return tab, nil
}
