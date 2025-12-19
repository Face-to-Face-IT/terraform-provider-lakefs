package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	userName := fmt.Sprintf("testuser%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig(userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakefs_user.test", "id", userName),
					resource.TestCheckResourceAttrSet("lakefs_user.test", "creation_date"),
				),
			},
			{
				ResourceName:      "lakefs_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUserResourceConfig(userName string) string {
	return fmt.Sprintf(`
resource "lakefs_user" "test" {
  id = %[1]q
}
`, userName)
}

func TestAccGroupResource(t *testing.T) {
	groupName := fmt.Sprintf("testgroup%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupResourceConfig(groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakefs_group.test", "id", groupName),
					resource.TestCheckResourceAttrSet("lakefs_group.test", "creation_date"),
				),
			},
			{
				ResourceName:      "lakefs_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroupResourceConfig(groupName string) string {
	return fmt.Sprintf(`
resource "lakefs_group" "test" {
  id = %[1]q
}
`, groupName)
}

func TestAccPolicyResource(t *testing.T) {
	policyName := fmt.Sprintf("testpolicy%d", time.Now().UnixNano())
	statement := `[{"effect":"allow","action":["fs:ReadObject"],"resource":"arn:lakefs:fs:::repository/*"}]`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyResourceConfig(policyName, statement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakefs_policy.test", "id", policyName),
					resource.TestCheckResourceAttrSet("lakefs_policy.test", "creation_date"),
				),
			},
			{
				ResourceName:            "lakefs_policy.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"statement"}, // JSON key order may differ
			},
		},
	})
}

func testAccPolicyResourceConfig(policyName, statement string) string {
	return fmt.Sprintf(`
resource "lakefs_policy" "test" {
  id        = %[1]q
  statement = %[2]q
}
`, policyName, statement)
}

func TestAccGroupMembershipResource(t *testing.T) {
	userName := fmt.Sprintf("memuser%d", time.Now().UnixNano())
	groupName := fmt.Sprintf("memgroup%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMembershipResourceConfig(userName, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakefs_group_membership.test", "user_id", userName),
					resource.TestCheckResourceAttr("lakefs_group_membership.test", "group_id", groupName),
				),
			},
		},
	})
}

func testAccGroupMembershipResourceConfig(userName, groupName string) string {
	return fmt.Sprintf(`
resource "lakefs_user" "test" {
  id = %[1]q
}

resource "lakefs_group" "test" {
  id = %[2]q
}

resource "lakefs_group_membership" "test" {
  group_id = lakefs_group.test.id
  user_id  = lakefs_user.test.id
}
`, userName, groupName)
}

func TestAccUserPolicyAttachmentResource(t *testing.T) {
	userName := fmt.Sprintf("upuser%d", time.Now().UnixNano())
	policyName := fmt.Sprintf("uppolicy%d", time.Now().UnixNano())
	statement := `[{"effect":"allow","action":["fs:ReadObject"],"resource":"arn:lakefs:fs:::repository/*"}]`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPolicyAttachmentResourceConfig(userName, policyName, statement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakefs_user_policy_attachment.test", "user_id", userName),
					resource.TestCheckResourceAttr("lakefs_user_policy_attachment.test", "policy_id", policyName),
				),
			},
		},
	})
}

func testAccUserPolicyAttachmentResourceConfig(userName, policyName, statement string) string {
	return fmt.Sprintf(`
resource "lakefs_user" "test" {
  id = %[1]q
}

resource "lakefs_policy" "test" {
  id        = %[2]q
  statement = %[3]q
}

resource "lakefs_user_policy_attachment" "test" {
  user_id   = lakefs_user.test.id
  policy_id = lakefs_policy.test.id
}
`, userName, policyName, statement)
}

func TestAccUserCredentialsResource(t *testing.T) {
	userName := fmt.Sprintf("creduser%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserCredentialsResourceConfig(userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakefs_user_credentials.test", "user_id", userName),
					resource.TestCheckResourceAttrSet("lakefs_user_credentials.test", "access_key_id"),
					resource.TestCheckResourceAttrSet("lakefs_user_credentials.test", "secret_access_key"),
				),
			},
		},
	})
}

func testAccUserCredentialsResourceConfig(userName string) string {
	return fmt.Sprintf(`
resource "lakefs_user" "test" {
  id = %[1]q
}

resource "lakefs_user_credentials" "test" {
  user_id = lakefs_user.test.id
}
`, userName)
}
