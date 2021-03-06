package open_resource_discovery_test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/open_resource_discovery"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
)

const (
	ordDocURI     = "/open-resource-discovery/v1/documents/example1"
	baseURL       = "http://localhost:8080"
	packageORDID  = "ns:package:PACKAGE_ID:v1"
	productORDID  = "ns:PRODUCT_ID"
	product2ORDID = "ns:PRODUCT_ID2"
	bundleORDID   = "ns:consumptionBundle:BUNDLE_ID:v1"
	vendorORDID   = "sap"
	api1ORDID     = "ns:apiResource:API_ID:v2"
	api2ORDID     = "ns:apiResource:API_ID2:v1"
	event1ORDID   = "ns:eventResource:EVENT_ID:v1"
	event2ORDID   = "ns2:eventResource:EVENT_ID:v1"

	appID     = "testApp"
	whID      = "testWh"
	tenantID  = "testTenant"
	packageID = "testPkg"
	bundleID  = "testBndl"
	api1ID    = "testApi1"
	api2ID    = "testApi2"
	event1ID  = "testEvent1"
	event2ID  = "testEvent2"

	cursor      = "cursor"
	policyLevel = "sap"
)

var (
	packageLinksFormat = removeWhitespace(`[
        {
          "type": "terms-of-service",
          "url": "https://example.com/en/legal/terms-of-use.html"
        },
        {
          "type": "client-registration",
          "url": "%s/ui/public/showRegisterForm"
        }
      ]`)

	linksFormat = removeWhitespace(`[
        {
		  "description": "lorem ipsum dolor nem",
          "title": "Link Title",
          "url": "https://example.com/2018/04/11/testing/"
        },
		{
		  "description": "lorem ipsum dolor nem",
          "title": "Link Title",
          "url": "%s/testing/relative"
        }
      ]`)

	packageLabels = removeWhitespace(`{
        "label-key-1": [
          "label-val"
        ],
		"pkg-label": [
          "label-val"
        ]
      }`)

	labels = removeWhitespace(`{
        "label-key-1": [
          "label-value-1",
          "label-value-2"
        ]
      }`)

	mergedLabels = removeWhitespace(`{
        "label-key-1": [
          "label-val",
		  "label-value-1",
          "label-value-2"
        ],
		"pkg-label": [
          "label-val"
        ]
      }`)

	credentialExchangeStrategiesFormat = removeWhitespace(`[
        {
		  "callbackUrl": "%s/credentials/relative",
          "customType": "ns:credential-exchange:v1",
		  "type": "custom"
        },
		{
		  "callbackUrl": "http://example.com/credentials",
          "customType": "ns:credential-exchange2:v3",
		  "type": "custom"
        }
      ]`)

	apiResourceLinksFormat = removeWhitespace(`[
        {
          "type": "console",
          "url": "https://example.com/shell/discover"
        },
		{
          "type": "console",
          "url": "%s/shell/discover/relative"
        }
      ]`)

	changeLogEntries = removeWhitespace(`[
        {
		  "date": "2020-04-29",
		  "description": "lorem ipsum dolor sit amet",
		  "releaseStatus": "active",
		  "url": "https://example.com/changelog/v1",
          "version": "1.0.0"
        }
      ]`)
)

func fixWellKnownConfig() *open_resource_discovery.WellKnownConfig {
	return &open_resource_discovery.WellKnownConfig{
		Schema: "../spec/v1/generated/Configuration.schema.json",
		OpenResourceDiscoveryV1: open_resource_discovery.OpenResourceDiscoveryV1{
			Documents: []open_resource_discovery.DocumentDetails{
				{
					URL:                 ordDocURI,
					SystemInstanceAware: true,
					AccessStrategies: []open_resource_discovery.AccessStrategy{
						{
							Type: open_resource_discovery.OpenAccessStrategy,
						},
					},
				},
			},
		},
	}
}

func fixORDDocument() *open_resource_discovery.Document {
	return fixORDDocumentWithBaseURL("")
}

func fixSanitizedORDDocument() *open_resource_discovery.Document {
	sanitizedDoc := fixORDDocumentWithBaseURL(baseURL)

	sanitizedDoc.APIResources[0].Tags = json.RawMessage(`["testTag","apiTestTag"]`)
	sanitizedDoc.APIResources[0].Countries = json.RawMessage(`["BG","EN","US"]`)
	sanitizedDoc.APIResources[0].LineOfBusiness = json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`)
	sanitizedDoc.APIResources[0].Industry = json.RawMessage(`["automotive","finance","test"]`)
	sanitizedDoc.APIResources[0].Labels = json.RawMessage(mergedLabels)

	sanitizedDoc.APIResources[1].Tags = json.RawMessage(`["testTag","ZGWSAMPLE"]`)
	sanitizedDoc.APIResources[1].Countries = json.RawMessage(`["BG","EN","BR"]`)
	sanitizedDoc.APIResources[1].LineOfBusiness = json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`)
	sanitizedDoc.APIResources[1].Industry = json.RawMessage(`["automotive","finance","test"]`)
	sanitizedDoc.APIResources[1].Labels = json.RawMessage(mergedLabels)

	sanitizedDoc.EventResources[0].Tags = json.RawMessage(`["testTag","eventTestTag"]`)
	sanitizedDoc.EventResources[0].Countries = json.RawMessage(`["BG","EN","US"]`)
	sanitizedDoc.EventResources[0].LineOfBusiness = json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`)
	sanitizedDoc.EventResources[0].Industry = json.RawMessage(`["automotive","finance","test"]`)
	sanitizedDoc.EventResources[0].Labels = json.RawMessage(mergedLabels)

	sanitizedDoc.EventResources[1].Tags = json.RawMessage(`["testTag","eventTestTag2"]`)
	sanitizedDoc.EventResources[1].Countries = json.RawMessage(`["BG","EN","BR"]`)
	sanitizedDoc.EventResources[1].LineOfBusiness = json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`)
	sanitizedDoc.EventResources[1].Industry = json.RawMessage(`["automotive","finance","test"]`)
	sanitizedDoc.EventResources[1].Labels = json.RawMessage(mergedLabels)

	return sanitizedDoc
}

func fixORDDocumentWithBaseURL(baseUrl string) *open_resource_discovery.Document {
	true := true
	return &open_resource_discovery.Document{
		Schema:                "./spec/v1/generated/Document.schema.json",
		OpenResourceDiscovery: "1.0-rc.1",
		Description:           "Test Document",
		DescribedSystemInstance: &model.Application{
			BaseURL: str.Ptr(baseURL),
			Labels:  json.RawMessage(labels),
		},
		ProviderSystemInstance: nil,
		Packages: []*model.PackageInput{
			{
				OrdID:            packageORDID,
				Vendor:           str.Ptr(vendorORDID),
				Title:            "PACKAGE 1 TITLE",
				ShortDescription: "lorem ipsum",
				Description:      "lorem ipsum dolor set",
				Version:          "1.1.2",
				PackageLinks:     json.RawMessage(fmt.Sprintf(packageLinksFormat, baseUrl)),
				Links:            json.RawMessage(fmt.Sprintf(linksFormat, baseUrl)),
				LicenseType:      str.Ptr("licence"),
				Tags:             json.RawMessage(`["testTag"]`),
				Countries:        json.RawMessage(`["BG","EN"]`),
				Labels:           json.RawMessage(packageLabels),
				PolicyLevel:      policyLevel,
				PartOfProducts:   json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
				LineOfBusiness:   json.RawMessage(`["lineOfBusiness"]`),
				Industry:         json.RawMessage(`["automotive","finance"]`),
			},
		},
		ConsumptionBundles: []*model.BundleCreateInput{
			{
				Name:                         "BUNDLE TITLE",
				Description:                  str.Ptr("lorem ipsum dolor nsq sme"),
				OrdID:                        str.Ptr(bundleORDID),
				ShortDescription:             str.Ptr("lorem ipsum"),
				Links:                        json.RawMessage(fmt.Sprintf(linksFormat, baseUrl)),
				Labels:                       json.RawMessage(labels),
				CredentialExchangeStrategies: json.RawMessage(fmt.Sprintf(credentialExchangeStrategiesFormat, baseUrl)),
			},
		},
		Products: []*model.ProductInput{
			{
				OrdID:            productORDID,
				Title:            "PRODUCT TITLE",
				ShortDescription: "lorem ipsum",
				Vendor:           vendorORDID,
				Parent:           str.Ptr(product2ORDID),
				PPMSObjectID:     str.Ptr("12391293812"),
				Labels:           json.RawMessage(labels),
			},
		},
		APIResources: []*model.APIDefinitionInput{
			{
				OrdID:               str.Ptr(api1ORDID),
				OrdBundleID:         str.Ptr(bundleORDID),
				OrdPackageID:        str.Ptr(packageORDID),
				Name:                "API TITLE",
				Description:         str.Ptr("lorem ipsum dolor sit amet"),
				TargetURL:           "https://exmaple.com/test/v1",
				ShortDescription:    str.Ptr("lorem ipsum"),
				SystemInstanceAware: &true,
				ApiProtocol:         str.Ptr("odata-v2"),
				Tags:                json.RawMessage(`["apiTestTag"]`),
				Countries:           json.RawMessage(`["BG","US"]`),
				Links:               json.RawMessage(fmt.Sprintf(linksFormat, baseUrl)),
				APIResourceLinks:    json.RawMessage(fmt.Sprintf(apiResourceLinksFormat, baseUrl)),
				ReleaseStatus:       str.Ptr("active"),
				SunsetDate:          nil,
				Successor:           nil,
				ChangeLogEntries:    json.RawMessage(changeLogEntries),
				Labels:              json.RawMessage(labels),
				Visibility:          str.Ptr("public"),
				Disabled:            &true,
				PartOfProducts:      json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
				LineOfBusiness:      json.RawMessage(`["lineOfBusiness2"]`),
				Industry:            json.RawMessage(`["automotive","test"]`),
				ResourceDefinitions: []*model.APIResourceDefinition{
					{
						Type:      "openapi-v3",
						MediaType: "application/json",
						URL:       fmt.Sprintf("%s/odata/1.0/catalog.svc/$value?type=json", baseURL),
						AccessStrategy: []model.AccessStrategy{
							{
								Type: "open",
							},
						},
					},
					{
						Type:      "openapi-v3",
						MediaType: "text/yaml",
						URL:       "https://test.com/odata/1.0/catalog",
						AccessStrategy: []model.AccessStrategy{
							{
								Type: "open",
							},
						},
					},
				},
				VersionInput: &model.VersionInput{
					Value: "2.1.2",
				},
			},
			{
				OrdID:               str.Ptr(api2ORDID),
				OrdBundleID:         str.Ptr(bundleORDID),
				OrdPackageID:        str.Ptr(packageORDID),
				Name:                "Gateway Sample Service",
				Description:         str.Ptr("lorem ipsum dolor sit amet"),
				TargetURL:           fmt.Sprintf("%s/some-api/v1", baseURL),
				ShortDescription:    str.Ptr("lorem ipsum"),
				SystemInstanceAware: &true,
				ApiProtocol:         str.Ptr("odata-v2"),
				Tags:                json.RawMessage(`["ZGWSAMPLE"]`),
				Countries:           json.RawMessage(`["BR"]`),
				Links:               json.RawMessage(fmt.Sprintf(linksFormat, baseUrl)),
				APIResourceLinks:    json.RawMessage(fmt.Sprintf(apiResourceLinksFormat, baseUrl)),
				ReleaseStatus:       str.Ptr("deprecated"),
				SunsetDate:          str.Ptr("2020-12-08T15:47:04+0000"),
				Successor:           str.Ptr(api1ORDID),
				ChangeLogEntries:    json.RawMessage(changeLogEntries),
				Labels:              json.RawMessage(labels),
				Visibility:          str.Ptr("public"),
				Disabled:            nil,
				PartOfProducts:      json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
				LineOfBusiness:      json.RawMessage(`["lineOfBusiness2"]`),
				Industry:            json.RawMessage(`["automotive","test"]`),
				ResourceDefinitions: []*model.APIResourceDefinition{
					{
						Type:      "edmx",
						MediaType: "application/xml",
						URL:       "https://TEST:443//odata/$metadata",
						AccessStrategy: []model.AccessStrategy{
							{
								Type: "open",
							},
						},
					},
				},
				VersionInput: &model.VersionInput{
					Value: "1.1.0",
				},
			},
		},
		EventResources: []*model.EventDefinitionInput{
			{
				OrdID:               str.Ptr(event1ORDID),
				OrdBundleID:         str.Ptr(bundleORDID),
				OrdPackageID:        str.Ptr(packageORDID),
				Name:                "EVENT TITLE",
				Description:         str.Ptr("lorem ipsum dolor sit amet"),
				ShortDescription:    str.Ptr("lorem ipsum"),
				SystemInstanceAware: &true,
				ChangeLogEntries:    json.RawMessage(changeLogEntries),
				Links:               json.RawMessage(fmt.Sprintf(linksFormat, baseUrl)),
				Tags:                json.RawMessage(`["eventTestTag"]`),
				Countries:           json.RawMessage(`["BG","US"]`),
				ReleaseStatus:       str.Ptr("active"),
				SunsetDate:          nil,
				Successor:           nil,
				Labels:              json.RawMessage(labels),
				Visibility:          str.Ptr("public"),
				Disabled:            &true,
				PartOfProducts:      json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
				LineOfBusiness:      json.RawMessage(`["lineOfBusiness2"]`),
				Industry:            json.RawMessage(`["automotive","test"]`),
				ResourceDefinitions: []*model.EventResourceDefinition{
					{
						Type:      "asyncapi-v2",
						MediaType: "application/json",
						URL:       "http://localhost:8080/asyncApi2.json",
						AccessStrategy: []model.AccessStrategy{
							{
								Type: "open",
							},
						},
					},
				},
				VersionInput: &model.VersionInput{
					Value: "2.1.2",
				},
			},
			{
				OrdID:               str.Ptr(event2ORDID),
				OrdBundleID:         str.Ptr(bundleORDID),
				OrdPackageID:        str.Ptr(packageORDID),
				Name:                "EVENT TITLE 2",
				Description:         str.Ptr("lorem ipsum dolor sit amet"),
				ShortDescription:    str.Ptr("lorem ipsum"),
				SystemInstanceAware: &true,
				ChangeLogEntries:    json.RawMessage(changeLogEntries),
				Links:               json.RawMessage(fmt.Sprintf(linksFormat, baseUrl)),
				Tags:                json.RawMessage(`["eventTestTag2"]`),
				Countries:           json.RawMessage(`["BR"]`),
				ReleaseStatus:       str.Ptr("deprecated"),
				SunsetDate:          str.Ptr("2020-12-08T15:47:04+0000"),
				Successor:           str.Ptr(event2ORDID),
				Labels:              json.RawMessage(labels),
				Visibility:          str.Ptr("public"),
				Disabled:            nil,
				PartOfProducts:      json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
				LineOfBusiness:      json.RawMessage(`["lineOfBusiness2"]`),
				Industry:            json.RawMessage(`["automotive","test"]`),
				ResourceDefinitions: []*model.EventResourceDefinition{
					{
						Type:      "asyncapi-v2",
						MediaType: "application/json",
						URL:       fmt.Sprintf("%s/api/eventCatalog.json", baseURL),
						AccessStrategy: []model.AccessStrategy{
							{
								Type: "open",
							},
						},
					},
				},
				VersionInput: &model.VersionInput{
					Value: "1.1.0",
				},
			},
		},
		Tombstones: []*model.TombstoneInput{
			{
				OrdID:       api2ORDID,
				RemovalDate: "2020-12-02T14:12:59Z",
			},
		},
		Vendors: []*model.VendorInput{
			{
				OrdID:  vendorORDID,
				Title:  "SAP",
				Type:   "sap",
				Labels: json.RawMessage(labels),
			},
		},
	}
}

func fixApplicationPage() *model.ApplicationPage {
	return &model.ApplicationPage{
		Data: []*model.Application{
			{
				Name:   "testApp",
				Tenant: tenantID,
				BaseEntity: &model.BaseEntity{
					ID:    appID,
					Ready: true,
				},
			},
		},
		PageInfo: &pagination.Page{
			StartCursor: cursor,
			EndCursor:   cursor,
			HasNextPage: false,
		},
		TotalCount: 1,
	}
}

func fixWebhooks() []*model.Webhook {
	return []*model.Webhook{
		{
			ID:            whID,
			TenantID:      str.Ptr(tenantID),
			ApplicationID: str.Ptr(appID),
			Type:          model.WebhookTypeOpenResourceDiscovery,
			URL:           str.Ptr(baseURL),
		},
	}
}

func fixVendors() []*model.Vendor {
	return []*model.Vendor{
		{
			OrdID:         vendorORDID,
			TenantID:      tenantID,
			ApplicationID: appID,
			Title:         "SAP",
			Type:          "sap",
			Labels:        json.RawMessage(labels),
		},
	}
}

func fixProducts() []*model.Product {
	return []*model.Product{
		{
			OrdID:            productORDID,
			TenantID:         tenantID,
			ApplicationID:    appID,
			Title:            "PRODUCT TITLE",
			ShortDescription: "lorem ipsum",
			Vendor:           vendorORDID,
			Parent:           str.Ptr(product2ORDID),
			PPMSObjectID:     str.Ptr("12391293812"),
			Labels:           json.RawMessage(labels),
		},
	}
}

func fixPackages() []*model.Package {
	return []*model.Package{
		{
			ID:               packageID,
			TenantID:         tenantID,
			ApplicationID:    appID,
			OrdID:            packageORDID,
			Vendor:           str.Ptr(vendorORDID),
			Title:            "PACKAGE 1 TITLE",
			ShortDescription: "lorem ipsum",
			Description:      "lorem ipsum dolor set",
			Version:          "1.1.2",
			PackageLinks:     json.RawMessage(fmt.Sprintf(packageLinksFormat, baseURL)),
			Links:            json.RawMessage(fmt.Sprintf(linksFormat, baseURL)),
			LicenseType:      str.Ptr("licence"),
			Tags:             json.RawMessage(`["testTag"]`),
			Countries:        json.RawMessage(`["BG","EN"]`),
			Labels:           json.RawMessage(packageLabels),
			PolicyLevel:      policyLevel,
			PartOfProducts:   json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
			LineOfBusiness:   json.RawMessage(`["lineOfBusiness"]`),
			Industry:         json.RawMessage(`["automotive","finance"]`),
		},
	}
}

func fixBundles() []*model.Bundle {
	return []*model.Bundle{
		{
			TenantID:                     tenantID,
			ApplicationID:                appID,
			Name:                         "BUNDLE TITLE",
			Description:                  str.Ptr("lorem ipsum dolor nsq sme"),
			OrdID:                        str.Ptr(bundleORDID),
			ShortDescription:             str.Ptr("lorem ipsum"),
			Links:                        json.RawMessage(fmt.Sprintf(linksFormat, baseURL)),
			Labels:                       json.RawMessage(labels),
			CredentialExchangeStrategies: json.RawMessage(fmt.Sprintf(credentialExchangeStrategiesFormat, baseURL)),
			BaseEntity: &model.BaseEntity{
				ID:    bundleID,
				Ready: true,
			},
		},
	}
}

func fixAPIs() []*model.APIDefinition {
	true := true
	return []*model.APIDefinition{
		{
			ApplicationID:    appID,
			BundleID:         str.Ptr(bundleORDID),
			PackageID:        str.Ptr(packageORDID),
			Tenant:           tenantID,
			Name:             "API TITLE",
			Description:      str.Ptr("lorem ipsum dolor sit amet"),
			TargetURL:        "/test/v1",
			OrdID:            str.Ptr(api1ORDID),
			ShortDescription: str.Ptr("lorem ipsum"),
			ApiProtocol:      str.Ptr("odata-v2"),
			Tags:             json.RawMessage(`["testTag","apiTestTag"]`),
			Countries:        json.RawMessage(`["BG","EN","US"]`),
			Links:            json.RawMessage(fmt.Sprintf(linksFormat, baseURL)),
			APIResourceLinks: json.RawMessage(fmt.Sprintf(apiResourceLinksFormat, baseURL)),
			ReleaseStatus:    str.Ptr("active"),
			ChangeLogEntries: json.RawMessage(changeLogEntries),
			Labels:           json.RawMessage(mergedLabels),
			Visibility:       str.Ptr("public"),
			Disabled:         &true,
			PartOfProducts:   json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
			LineOfBusiness:   json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`),
			Industry:         json.RawMessage(`["automotive","finance","test"]`),
			Version: &model.Version{
				Value: "2.1.3",
			},
			BaseEntity: &model.BaseEntity{
				ID:    api1ID,
				Ready: true,
			},
		},
		{
			ApplicationID:    appID,
			BundleID:         str.Ptr(bundleORDID),
			PackageID:        str.Ptr(packageORDID),
			Tenant:           tenantID,
			Name:             "Gateway Sample Service",
			Description:      str.Ptr("lorem ipsum dolor sit amet"),
			TargetURL:        "/some-api/v1",
			OrdID:            str.Ptr(api2ORDID),
			ShortDescription: str.Ptr("lorem ipsum"),
			ApiProtocol:      str.Ptr("odata-v2"),
			Tags:             json.RawMessage(`["testTag","ZGWSAMPLE"]`),
			Countries:        json.RawMessage(`["BG","EN","BR"]`),
			Links:            json.RawMessage(fmt.Sprintf(linksFormat, baseURL)),
			APIResourceLinks: json.RawMessage(fmt.Sprintf(apiResourceLinksFormat, baseURL)),
			ReleaseStatus:    str.Ptr("deprecated"),
			SunsetDate:       str.Ptr("2020-12-08T15:47:04+0000"),
			Successor:        str.Ptr(api1ORDID),
			ChangeLogEntries: json.RawMessage(changeLogEntries),
			Labels:           json.RawMessage(mergedLabels),
			Visibility:       str.Ptr("public"),
			PartOfProducts:   json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
			LineOfBusiness:   json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`),
			Industry:         json.RawMessage(`["automotive","finance","test"]`),
			Version: &model.Version{
				Value: "1.1.1",
			},
			BaseEntity: &model.BaseEntity{
				ID:    api2ID,
				Ready: true,
			},
		},
	}
}

func fixEvents() []*model.EventDefinition {
	true := true
	return []*model.EventDefinition{
		{
			Tenant:           tenantID,
			ApplicationID:    appID,
			BundleID:         str.Ptr(bundleORDID),
			PackageID:        str.Ptr(packageORDID),
			Name:             "EVENT TITLE",
			Description:      str.Ptr("lorem ipsum dolor sit amet"),
			OrdID:            str.Ptr(event1ORDID),
			ShortDescription: str.Ptr("lorem ipsum"),
			ChangeLogEntries: json.RawMessage(changeLogEntries),
			Links:            json.RawMessage(fmt.Sprintf(linksFormat, baseURL)),
			Tags:             json.RawMessage(`["testTag","eventTestTag"]`),
			Countries:        json.RawMessage(`["BG","EN","US"]`),
			ReleaseStatus:    str.Ptr("active"),
			Labels:           json.RawMessage(mergedLabels),
			Visibility:       str.Ptr("public"),
			Disabled:         &true,
			PartOfProducts:   json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
			LineOfBusiness:   json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`),
			Industry:         json.RawMessage(`["automotive","finance","test"]`),
			Version: &model.Version{
				Value: "2.1.3",
			},
			BaseEntity: &model.BaseEntity{
				ID:    event1ID,
				Ready: true,
			},
		},
		{
			Tenant:           tenantID,
			ApplicationID:    appID,
			BundleID:         str.Ptr(bundleORDID),
			PackageID:        str.Ptr(packageORDID),
			Name:             "EVENT TITLE 2",
			Description:      str.Ptr("lorem ipsum dolor sit amet"),
			OrdID:            str.Ptr(event2ORDID),
			ShortDescription: str.Ptr("lorem ipsum"),
			ChangeLogEntries: json.RawMessage(changeLogEntries),
			Links:            json.RawMessage(fmt.Sprintf(linksFormat, baseURL)),
			Tags:             json.RawMessage(`["testTag","eventTestTag2"]`),
			Countries:        json.RawMessage(`["BG","EN","BR"]`),
			ReleaseStatus:    str.Ptr("deprecated"),
			SunsetDate:       str.Ptr("2020-12-08T15:47:04+0000"),
			Successor:        str.Ptr(event2ORDID),
			Labels:           json.RawMessage(mergedLabels),
			Visibility:       str.Ptr("public"),
			PartOfProducts:   json.RawMessage(fmt.Sprintf(`["%s"]`, productORDID)),
			LineOfBusiness:   json.RawMessage(`["lineOfBusiness","lineOfBusiness2"]`),
			Industry:         json.RawMessage(`["automotive","finance","test"]`),
			Version: &model.Version{
				Value: "1.1.1",
			},
			BaseEntity: &model.BaseEntity{
				ID:    event2ID,
				Ready: true,
			},
		},
	}
}

func fixApi1SpecInputs() []*model.SpecInput {
	apiType := model.APISpecTypeOpenAPIV3
	return []*model.SpecInput{
		{
			Format:     "application/json",
			APIType:    &apiType,
			CustomType: str.Ptr(""),
			FetchRequest: &model.FetchRequestInput{
				URL: baseURL + "/odata/1.0/catalog.svc/$value?type=json",
			},
		},
		{
			Format:     "text/yaml",
			APIType:    &apiType,
			CustomType: str.Ptr(""),
			FetchRequest: &model.FetchRequestInput{
				URL: "https://test.com/odata/1.0/catalog",
			},
		},
	}
}

func fixApi2SpecInputs() []*model.SpecInput {
	apiType := model.APISpecTypeEDMX
	return []*model.SpecInput{
		{
			Format:     "application/xml",
			APIType:    &apiType,
			CustomType: str.Ptr(""),
			FetchRequest: &model.FetchRequestInput{
				URL: "https://TEST:443//odata/$metadata",
			},
		},
	}
}

func fixEvent1SpecInputs() []*model.SpecInput {
	eventType := model.EventSpecTypeAsyncAPIV2
	return []*model.SpecInput{
		{
			Format:     "application/json",
			EventType:  &eventType,
			CustomType: str.Ptr(""),
			FetchRequest: &model.FetchRequestInput{
				URL: "http://localhost:8080/asyncApi2.json",
			},
		},
	}
}

func fixEvent2SpecInputs() []*model.SpecInput {
	eventType := model.EventSpecTypeAsyncAPIV2
	return []*model.SpecInput{
		{
			Format:     "application/json",
			EventType:  &eventType,
			CustomType: str.Ptr(""),
			FetchRequest: &model.FetchRequestInput{
				URL: baseURL + "/api/eventCatalog.json",
			},
		},
	}
}

func fixTombstones() []*model.Tombstone {
	return []*model.Tombstone{
		{
			OrdID:         api2ORDID,
			TenantID:      tenantID,
			ApplicationID: appID,
			RemovalDate:   "2020-12-02T14:12:59Z",
		},
	}
}

func bundleUpdateInputFromCreateInput(in model.BundleCreateInput) model.BundleUpdateInput {
	return model.BundleUpdateInput{
		Name:                           in.Name,
		Description:                    in.Description,
		InstanceAuthRequestInputSchema: in.InstanceAuthRequestInputSchema,
		DefaultInstanceAuth:            in.DefaultInstanceAuth,
		OrdID:                          in.OrdID,
		ShortDescription:               in.ShortDescription,
		Links:                          in.Links,
		Labels:                         in.Labels,
		CredentialExchangeStrategies:   in.CredentialExchangeStrategies,
	}
}

func removeWhitespace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", ""), "\t", "")
}
