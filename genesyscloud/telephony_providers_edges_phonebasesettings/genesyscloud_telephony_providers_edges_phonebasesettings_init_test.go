package telephony_providers_edges_phonebasesettings

import (
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerDataSources map[string]*schema.Resource
var providerResources map[string]*schema.Resource

type registerTestInstance struct {
	resourceMapMutex   sync.RWMutex
	datasourceMapMutex sync.RWMutex
}

func (r *registerTestInstance) registerTestResources() {
	r.resourceMapMutex.Lock()
	defer r.resourceMapMutex.Unlock()
	providerResources[ResourceType] = ResourcePhoneBaseSettings()
}

func (r *registerTestInstance) registerTestDataSources() {
	r.datasourceMapMutex.Lock()
	defer r.datasourceMapMutex.Unlock()
	providerDataSources[ResourceType] = DataSourcePhoneBaseSettings()
}

func initTestresources() {
	providerDataSources = make(map[string]*schema.Resource)
	providerResources = make(map[string]*schema.Resource)
	regInstance := &registerTestInstance{}
	regInstance.registerTestDataSources()
	regInstance.registerTestResources()
}

func TestMain(m *testing.M) {
	// Run setup function before starting the test suite
	initTestresources()
	// Run the actual test suite
	m.Run()
}
