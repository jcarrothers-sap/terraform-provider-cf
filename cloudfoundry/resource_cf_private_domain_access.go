package cloudfoundry

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-cf/cloudfoundry/cfapi"
)

func resourcePrivateDomainAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateDomainAccessCreate,
		Read:   resourcePrivateDomainAccessRead,
		Delete: resourcePrivateDomainAccessDelete,
		Importer: &schema.ResourceImporter{
			State: ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func computeID(org, domain string) string {
	return fmt.Sprintf("%s/%s", org, domain)
}

func parseID(ID string) (org string, domain string, err error) {
	parts := strings.Split(ID, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("unable to parse ID '%s', expected format is '<org-guid>/<domain-guid>'", ID)
	} else {
		org = parts[0]
		domain = parts[1]
	}
	return
}

// PrivateDomainAccessImport -
// Checks that user-given ID matches <guid>/<guid> format
func PrivateDomainAccessImport(d *schema.ResourceData, meta interface{}) (res []*schema.ResourceData, err error) {
	// session := meta.(*cfapi.Session)
	// if session == nil {
	// 	err = fmt.Errorf("client is nil")
	// 	return
	// }
	// dm := session.DomainManager()
	id := d.Id()

	// var org, domain string
	// if org, domain, err = parseID(id); err != nil {
	if _, _, err = parseID(id); err != nil {
		return
	}

	// var found bool
	// found, err = dm.HasPrivateDomainAccess(org, domain)
	// if err != nil {
	// 	return
	// }

	// if !found {
	// 	err = fmt.Errorf("organization '%s' has no access to private domain '%s'", org, domain)
	// 	return
	// }
	return schema.ImportStatePassthrough(d, meta)
}

func resourcePrivateDomainAccessCreate(d *schema.ResourceData, meta interface{}) (err error) {
	session := meta.(*cfapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	domain := d.Get("domain").(string)
	org := d.Get("org").(string)

	dm := session.DomainManager()
	if err = dm.CreatePrivateDomainAccess(org, domain); err != nil {
		return
	}

	d.SetId(computeID(org, domain))
	return nil
}

func resourcePrivateDomainAccessRead(d *schema.ResourceData, meta interface{}) (err error) {
	session := meta.(*cfapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	id := d.Id()
	// id in read hook comes from create or import callback which ensure id's validity
	var org, domain string
	org, domain, _ = parseID(id)

	dm := session.DomainManager()
	var found bool
	if found, err = dm.HasPrivateDomainAccess(org, domain); err != nil || !found {
		d.SetId("")
		return err
	}

	d.Set("org", org)
	d.Set("domain", domain)
	return
}

func resourcePrivateDomainAccessDelete(d *schema.ResourceData, meta interface{}) (err error) {
	session := meta.(*cfapi.Session)
	if session == nil {
		return fmt.Errorf("client is nil")
	}

	dm := session.DomainManager()
	id := d.Id()

	// id in read hook comes from create or import callback which ensure id's validity
	var org, domain string
	org, domain, _ = parseID(id)

	err = dm.DeletePrivateDomainAccess(org, domain)
	return
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
