package fastly

import (
	"fmt"
	"testing"

	gofastly "github.com/fastly/go-fastly/v6/fastly"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFastlyServiceDynamicSnippetContent_create(t *testing.T) {
	var service gofastly.ServiceDetail
	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	expectedNumberOfSnippets := "1"
	expectedSnippetName := "dyn_hit_test"
	expectedSnippetType := "hit"
	expectedSnippetPriority := "100"
	expectedSnippetContent := "if ( req.url ) {\n set req.http.my-snippet-test-header = \"true\";\n}"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceVCLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDynamicSnippetContentConfigWithDynamicSnippet(serviceName, expectedSnippetName, expectedSnippetContent, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					testAccCheckFastlyServiceDynamicSnippetContentRemoteState(&service, serviceName, expectedSnippetName, expectedSnippetContent),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", expectedNumberOfSnippets),
					resource.TestCheckTypeSetElemNestedAttrs("fastly_service_vcl.foo", "dynamicsnippet.*", map[string]string{
						"name":     expectedSnippetName,
						"type":     expectedSnippetType,
						"priority": expectedSnippetPriority,
					}),
					resource.TestCheckResourceAttr("fastly_service_dynamic_snippet_content.content", "content", expectedSnippetContent),
				),
			},
			{
				ResourceName:            "fastly_service_dynamic_snippet_content.content",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manage_snippets"},
			},
		},
	})
}

func TestAccFastlyServiceDynamicSnippetContent_update(t *testing.T) {
	var service gofastly.ServiceDetail
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	dynamicSnippetName := fmt.Sprintf("dynamic snippet %s", acctest.RandString(10))

	expectedRemoteItems := "if ( req.url ) {\n set req.http.my-snippet-test-header = \"true\";\n}"

	expectedRemoteItemsAfterUpdate := "if ( req.url ) {\n set req.http.my-updated-test-header = \"true\";\n}"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceVCLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDynamicSnippetContentConfigWithDynamicSnippet(name, dynamicSnippetName, expectedRemoteItems, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					testAccCheckFastlyServiceDynamicSnippetContentRemoteState(&service, name, dynamicSnippetName, expectedRemoteItems),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", "1"),
					resource.TestCheckResourceAttr("fastly_service_dynamic_snippet_content.content", "content", expectedRemoteItems),
				),
			},
			{
				Config: testAccServiceDynamicSnippetContentConfigWithDynamicSnippet(name, dynamicSnippetName, expectedRemoteItemsAfterUpdate, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					testAccCheckFastlyServiceDynamicSnippetContentRemoteState(&service, name, dynamicSnippetName, expectedRemoteItemsAfterUpdate),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", "1"),
					resource.TestCheckResourceAttr("fastly_service_dynamic_snippet_content.content", "content", expectedRemoteItemsAfterUpdate),
				),
			},
		},
	})
}

func TestAccFastlyServiceDynamicSnippetContent_external_snippet_is_removed(t *testing.T) {
	var service gofastly.ServiceDetail
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	externalDynamicSnippetName := fmt.Sprintf("existing dynamic snippet %s", acctest.RandString(10))
	externalContent := "if ( req.url ) {\n set req.http.my-snippet-test-header = \"true\";\n}"

	managedDynamicSnippetName := fmt.Sprintf("dynamic snippet %s", acctest.RandString(10))
	managedContent := "if ( req.url ) {\n set req.http.my-updated-test-header = \"true\";\n}"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceVCLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDynamicSnippetContentConfigWithDynamicSnippet(name, managedDynamicSnippetName, managedContent, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", "1"),
				),
			},
			{
				PreConfig: func() {
					createDynamicSnippetThroughApi(t, &service, externalDynamicSnippetName, gofastly.SnippetTypeHit, externalContent)
				},
				Config: testAccServiceDynamicSnippetContentConfigWithDynamicSnippet(name, managedDynamicSnippetName, managedContent, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					testAccCheckFastlyServiceDynamicSnippetContentRemoteState(&service, name, managedDynamicSnippetName, managedContent),
					testAccCheckFastlyServiceDynamicSnippetContentRemoteStateDoesntExist(&service, name, externalDynamicSnippetName),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", "1"),
					resource.TestCheckResourceAttr("fastly_service_dynamic_snippet_content.content", "content", managedContent),
				),
			},
		},
	})
}

func TestAccFastlyServiceDynamicSnippetContent_normal_snippet_is_not_removed(t *testing.T) {
	var service gofastly.ServiceDetail
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	normalSnippetName := fmt.Sprintf("normal dynamic snippet %s", acctest.RandString(10))
	normalContent := "if ( req.url ) {\n set req.http.my-snippet-test-header = \"true\";\n}"

	dynamicSnippetName := fmt.Sprintf("existing dynamic snippet %s", acctest.RandString(10))
	dynamicContent := "if ( req.url ) {\n set req.http.my-new-content-test-header = \"true\";\n}"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceVCLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDynamicSnippetContentConfigWithSnippet(name, normalSnippetName, normalContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "snippet.#", "1"),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", "0"),
				),
			},
			{
				Config: testAccServiceDynamicSnippetContentConfigWithSnippetAndDynamicSnippet(name, normalSnippetName, normalContent, dynamicSnippetName, dynamicContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceVCLExists("fastly_service_vcl.foo", &service),
					testAccCheckFastlyServiceDynamicSnippetContentRemoteState(&service, name, dynamicSnippetName, dynamicContent),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "snippet.#", "1"),
					resource.TestCheckResourceAttr("fastly_service_vcl.foo", "dynamicsnippet.#", "1"),
					resource.TestCheckResourceAttr("fastly_service_dynamic_snippet_content.content", "content", dynamicContent),
				),
			},
		},
	})
}

func testAccCheckFastlyServiceDynamicSnippetContentRemoteState(service *gofastly.ServiceDetail, name, dynamicSnippetName, expectedContent string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if service.Name != name {
			return fmt.Errorf("Bad name, expected (%s), got (%s)", name, service.Name)
		}

		conn := testAccProvider.Meta().(*FastlyClient).conn
		snippet, err := conn.GetSnippet(&gofastly.GetSnippetInput{
			ServiceID:      service.ID,
			ServiceVersion: service.ActiveVersion.Number,
			Name:           dynamicSnippetName,
		})

		if err != nil {
			return fmt.Errorf("[ERR] Error looking up snippet records for (%s), version (%v): %s", service.Name, service.ActiveVersion.Number, err)
		}

		dynamicSnippet, err := conn.GetDynamicSnippet(&gofastly.GetDynamicSnippetInput{
			ServiceID: service.ID,
			ID:        snippet.ID,
		})

		if err != nil {
			return fmt.Errorf("[ERR] Error looking up Dynamic snippet content for (%s), snippet (%s): %s", service.Name, snippet.ID, err)
		}

		if dynamicSnippet.Content != expectedContent {
			return fmt.Errorf("[ERR] Error matching:\nexpected: %s\ngot: %s", expectedContent, dynamicSnippet.Content)
		}

		return nil
	}
}

func testAccCheckFastlyServiceDynamicSnippetContentRemoteStateDoesntExist(service *gofastly.ServiceDetail, name, dynamicSnippetName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if service.Name != name {
			return fmt.Errorf("Bad name, expected (%s), got (%s)", name, service.Name)
		}

		conn := testAccProvider.Meta().(*FastlyClient).conn
		snippets, err := conn.ListSnippets(&gofastly.ListSnippetsInput{
			ServiceID:      service.ID,
			ServiceVersion: service.ActiveVersion.Number,
		})

		if err != nil {
			return fmt.Errorf("[ERR] Error looking up snippet records for (%s), version (%v): %s", service.Name, service.ActiveVersion.Number, err)
		}

		for _, snippet := range snippets {
			if snippet.Name == dynamicSnippetName {
				return fmt.Errorf("[ERR] Dynamic snippet (%s) exists in service (%s)", dynamicSnippetName, service.Name)
			}
		}

		return nil
	}
}

func createDynamicSnippetThroughApi(t *testing.T, service *gofastly.ServiceDetail, dynamicSnippetName string, snippetType gofastly.SnippetType, content string) {

	conn := testAccProvider.Meta().(*FastlyClient).conn

	newVersion, err := conn.CloneVersion(&gofastly.CloneVersionInput{
		ServiceID:      service.ID,
		ServiceVersion: service.ActiveVersion.Number,
	})

	if err != nil {
		t.Fatalf("[ERR] Error cloning service version (%s), version (%v): %s", service.Name, service.ActiveVersion.Number, err)
	}

	dynamicSnippet, err := conn.CreateSnippet(&gofastly.CreateSnippetInput{
		ServiceID:      service.ID,
		ServiceVersion: newVersion.Number,
		Name:           dynamicSnippetName,
		Type:           snippetType,
		Dynamic:        1,
	})

	if err != nil {
		t.Fatalf("[ERR] Error creating Dynamic snippet records for (%s), version (%v): %s", service.Name, service.ActiveVersion.Number, err)
	}

	_, err = conn.ActivateVersion(&gofastly.ActivateVersionInput{
		ServiceID:      service.ID,
		ServiceVersion: newVersion.Number,
	})

	if err != nil {
		t.Fatalf("[ERR] Error activating service version (%s), version (%v): %s", service.Name, newVersion.Number, err)
	}

	_, err = conn.UpdateDynamicSnippet(&gofastly.UpdateDynamicSnippetInput{
		ServiceID: service.ID,
		ID:        dynamicSnippet.ID,
		Content:   gofastly.String(content),
	})

	if err != nil {
		t.Fatalf("[ERR] Error update content for Dynamic snippet records for (%s), snippet id (%v): %s", service.Name, dynamicSnippet.ID, err)
	}

	latest, err := conn.GetServiceDetails(&gofastly.GetServiceInput{
		ID: service.ID,
	})

	if err != nil {
		t.Fatalf("[ERR] Error retrieving service details for (%s): %s", service.ID, err)
	}

	*service = *latest

}

func testAccServiceDynamicSnippetContentConfigWithSnippet(serviceName, snippetName, content string) string {
	backendName := fmt.Sprintf("%s.aws.amazon.com", acctest.RandString(3))
	domainName := fmt.Sprintf("fastly-test.tf-%s.com", acctest.RandString(10))

	return fmt.Sprintf(`
resource "fastly_service_vcl" "foo" {
  name = "%s"

  domain {
    name    = "%s"
    comment = "tf-testing-domain"
  }

  backend {
    address = "%s"
    name    = "tf -test backend"
  }

  snippet {
	name = "%s"
	type = "hit"
	content = %q
  }

  force_destroy = true
}`, serviceName, domainName, backendName, snippetName, content)
}

func testAccServiceDynamicSnippetContentConfigWithSnippetAndDynamicSnippet(serviceName, snippetName, snippetContent, dynamicSnippetName, dynamicSnippetContent string) string {
	backendName := fmt.Sprintf("%s.aws.amazon.com", acctest.RandString(3))
	domainName := fmt.Sprintf("fastly-test.tf-%s.com", acctest.RandString(10))

	return fmt.Sprintf(`
variable "mydynamicsnippet" {
	type = object({ name=string, content=string })
	default = {
		name = "%s"
		content = %q
	}
}

resource "fastly_service_vcl" "foo" {
  name = "%s"

  domain {
    name    = "%s"
    comment = "tf-testing-domain"
  }

  backend {
    address = "%s"
    name    = "tf -test backend"
  }

  snippet {
	name = "%s"
	type = "hit"
	content = %q
  }

  dynamicsnippet {
	name = var.mydynamicsnippet.name
	type = "hit"
  }

  force_destroy = true
}

resource "fastly_service_dynamic_snippet_content" "content" {
    service_id = fastly_service_vcl.foo.id
    snippet_id = {for s in fastly_service_vcl.foo.dynamicsnippet : s.name => s.snippet_id}[var.mydynamicsnippet.name]
    content = var.mydynamicsnippet.content
}`, dynamicSnippetName, dynamicSnippetContent, serviceName, domainName, backendName, snippetName, snippetContent)
}

func testAccServiceDynamicSnippetContentConfigWithDynamicSnippet(serviceName, dynamicSnippetName, content string, manageSnippets bool) string {
	backendName := fmt.Sprintf("%s.aws.amazon.com", acctest.RandString(3))
	domainName := fmt.Sprintf("fastly-test.tf-%s.com", acctest.RandString(10))

	return fmt.Sprintf(`
variable "mydynamicsnippet" {
	type = object({ name=string, content=string })
	default = {
		name = "%s"
		content = %q
	}
}

resource "fastly_service_vcl" "foo" {
  name = "%s"

  domain {
    name    = "%s"
    comment = "tf-testing-domain"
	}

  backend {
    address = "%s"
    name    = "tf -test backend"
  }

  dynamicsnippet {
	name = var.mydynamicsnippet.name
	type = "hit"
  }

  force_destroy = true
}

resource "fastly_service_dynamic_snippet_content" "content" {
  service_id      = fastly_service_vcl.foo.id
  snippet_id      = { for s in fastly_service_vcl.foo.dynamicsnippet : s.name => s.snippet_id }[var.mydynamicsnippet.name]
  manage_snippets = %t
  content         = var.mydynamicsnippet.content
}`, dynamicSnippetName, content, serviceName, domainName, backendName, manageSnippets)
}
