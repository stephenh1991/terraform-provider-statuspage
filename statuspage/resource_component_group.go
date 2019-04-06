package statuspage

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	sp "github.com/yannh/statuspage-go-sdk"
)

func resourceComponentGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)

	// convert []interface{} to []string
	var components []string
	for _, c := range d.Get("components").([]interface{}) {
		components = append(components, c.(string))
	}

	componentGroup, err := sp.CreateComponentGroup(
		client,
		d.Get("page_id").(string),
		&sp.ComponentGroup{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Components:  components,
		},
	)

	if err != nil {
		log.Printf("[WARN] Statuspage Failed creating component group: %s\n", err)
		return err
	}

	log.Printf("[INFO] Statuspage Created component group: %s\n", componentGroup.ID)
	d.SetId(componentGroup.ID)

	return resourceComponentGroupRead(d, m)
}

func resourceComponentGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	componentGroup, err := sp.GetComponentGroup(client, d.Get("page_id").(string), d.Id())
	if err != nil {
		log.Printf("[ERROR] Statuspage could not find component group with ID: %s\n", d.Id())
		return err
	}

	if componentGroup == nil {
		log.Printf("[INFO] Statuspage could not find component group with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Statuspage read component group: %s\n", componentGroup.ID)

	d.Set("name", componentGroup.Name)
	d.Set("description", componentGroup.Description)
	d.Set("components", componentGroup.Components)

	return nil
}

func resourceComponentGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	componentGroupID := d.Id()

	var components []string
	for _, c := range d.Get("components").([]interface{}) {
		components = append(components, c.(string))
	}

	_, err := sp.UpdateComponentGroup(
		client,
		d.Get("page_id").(string),
		componentGroupID,
		&sp.ComponentGroup{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Components:  components,
		},
	)
	if err != nil {
		log.Printf("[WARN] Statuspage Failed creating component group: %s\n", err)
		return err
	}

	d.SetId(componentGroupID)

	return resourceComponentGroupRead(d, m)
}

func resourceComponentGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)

	return sp.DeleteComponentGroup(client, d.Get("page_id").(string), d.Id())
}

func resourceComponentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentGroupCreate,
		Read:   resourceComponentGroupRead,
		Update: resourceComponentGroupUpdate,
		Delete: resourceComponentGroupDelete,

		Schema: map[string]*schema.Schema{
			"page_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"components": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}