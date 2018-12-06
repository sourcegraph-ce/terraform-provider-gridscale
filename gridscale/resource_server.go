package gridscale

import (
	"bitbucket.org/gridscale/gsclient-go"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"strings"
)

func resourceGridscaleServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleServerCreate,
		Read:   resourceGridscaleServerRead,
		Delete: resourceGridscaleServerDelete,
		Update: resourceGridscaleServerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"memory": {
				Type:         schema.TypeInt,
				Description:  "The amount of server memory in GB.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"cores": {
				Type:         schema.TypeInt,
				Description:  "The number of server cores.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"hardware_profile": {
				Type:        schema.TypeString,
				Description: "The number of server cores.",
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, profile := range HardwareProfiles {
						if v.(string) == profile {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid hardware profile. Valid hardware profiles are: %v", v.(string), strings.Join(StorageTypes, ",")))
					}
					return
				},
			},
			"storage": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				MaxItems: 8,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"controller": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lun": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"license_product_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_used_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				MaxItems: 7,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_template_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"partner_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ordering": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						//"vlan": {
						//	Type:     schema.TypeInt,
						//	Computed: true,
						//},
						//"vxlan": {
						//	Type:     schema.TypeInt,
						//	Computed: true,
						//},
						//"mcast": {
						//	Type:     schema.TypeString,
						//	Computed: true,
						//},
					},
				},
			},
			"ipv4": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"power": {
				Type:        schema.TypeBool,
				Description: "The number of server cores.",
				Optional:    true,
				Computed:    true,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"auto_recovery": {
				Type:        schema.TypeInt,
				Description: "If the server should be auto-started in case of a failure (default=true).",
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "Defines which Availability-Zone the Server is placed.",
				Optional:    true,
			},
			"console_token": {
				Type:        schema.TypeString,
				Description: "If the server should be auto-started in case of a failure (default=true).",
				Computed:    true,
			},
			"legacy": {
				Type:        schema.TypeBool,
				Description: "Legacy-Hardware emulation instead of virtio hardware. If enabled, hotplugging cores, memory, storage, network, etc. will not work, but the server will most likely run every x86 compatible operating system. This mode comes with a performance penalty, as emulated hardware does not benefit from the virtio driver infrastructure.",
				Computed:    true,
			},
			"usage_in_minutes_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"usage_in_minutes_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGridscaleServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	server, err := client.GetServer(d.Id())
	if err != nil {
		if requestError, ok := err.(*gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", server.Properties.Name)
	d.Set("memory", server.Properties.Memory)
	d.Set("cores", server.Properties.Cores)
	d.Set("hardware_profile", server.Properties.HardwareProfile)
	d.Set("location_uuid", server.Properties.LocationUuid)
	d.Set("power", server.Properties.Power)
	d.Set("current_price", server.Properties.CurrentPrice)
	d.Set("availability_zone", server.Properties.AvailablityZone)
	d.Set("auto_recovery", server.Properties.AutoRecovery)
	d.Set("console_token", server.Properties.ConsoleToken)
	d.Set("legacy", server.Properties.Legacy)
	d.Set("usage_in_minutes_memory", server.Properties.UsageInMinutesMemory)
	d.Set("usage_in_minutes_cores", server.Properties.UsageInMinutesCores)
	d.Set("labels", server.Properties.Labels)

	//Get storages
	storages := make([]interface{}, 0)
	for _, value := range server.Properties.Relations.Storages {
		storage := map[string]interface{}{
			"object_uuid":        value.ObjectUuid,
			"bootdevice":         value.BootDevice,
			"create_time":        value.CreateTime,
			"controller":         value.Controller,
			"target":             value.Target,
			"lun":                value.Lun,
			"license_product_no": value.LicenseProductNo,
			"bus":                value.Bus,
			"object_name":        value.ObjectName,
			"storage_type":       value.StorageType,
			"last_used_template": value.LastUsedTemplate,
			"capacity":           value.Capacity,
		}
		storages = append(storages, storage)
	}
	d.Set("storage", storages)

	//Get storages
	networks := make([]interface{}, 0)
	for _, value := range server.Properties.Relations.Networks {
		if !value.PublicNet {
			network := map[string]interface{}{
				"object_uuid":            value.ObjectUuid,
				"bootdevice":             value.BootDevice,
				"create_time":            value.CreateTime,
				"mac":                    value.Mac,
				"firewall":               value.Firewall,
				"firewall_template_uuid": value.FirewallTemplateUuid,
				"object_name":            value.ObjectName,
				"network_type":           value.NetworkType,
				"ordering":               value.Ordering,
				//"vlan":                   value.Vlan,
				//"vxlan":                  value.Vxlan,
				//"mcast":                  value.Mcast,
			}
			networks = append(networks, network)
		}
	}
	d.Set("network", networks)

	//Get IP addresses
	var ipv4, ipv6 string
	for _, ip := range server.Properties.Relations.PublicIps {
		if ip.Family == 4 {
			ipv4 = ip.ObjectUuid
		}
		if ip.Family == 6 {
			ipv6 = ip.ObjectUuid
		}
	}
	d.Set("ipv4", ipv4)
	d.Set("ipv6", ipv6)

	return nil
}

func resourceGridscaleServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.ServerCreateRequest{
		Name:            d.Get("name").(string),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
		LocationUuid:    d.Get("location_uuid").(string),
		HardwareProfile: d.Get("hardware_profile").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          d.Get("labels").(*schema.Set).List(),
	}

	requestBody.Relations.IsoImages = []gsclient.ServerIsoImage{}
	requestBody.Relations.Storages = []gsclient.ServerCreateRequestStorage{}
	requestBody.Relations.Networks = []gsclient.ServerCreateRequestNetwork{}
	requestBody.Relations.PublicIps = []gsclient.ServerCreateRequestIp{}

	if attr, ok := d.GetOk("storage"); ok {
		for _, value := range attr.(*schema.Set).List() {
			storage := value.(map[string]interface{})
			createStorageRequest := gsclient.ServerCreateRequestStorage{
				StorageUuid: storage["object_uuid"].(string),
				BootDevice:  storage["bootdevice"].(bool),
			}

			requestBody.Relations.Storages = append(requestBody.Relations.Storages, createStorageRequest)
		}
	}

	if attr, ok := d.GetOk("ipv4"); ok {
		if client.GetIpVersion(attr.(string)) != 4 {
			return fmt.Errorf("The IP address with UUID %v is not version 4", attr.(string))
		}
		ip := gsclient.ServerCreateRequestIp{
			IpaddrUuid: attr.(string),
		}
		requestBody.Relations.PublicIps = append(requestBody.Relations.PublicIps, ip)
	}
	if attr, ok := d.GetOk("ipv6"); ok {
		if client.GetIpVersion(attr.(string)) != 6 {
			return fmt.Errorf("The IP address with UUID %v is not version 6", attr.(string))
		}
		ip := gsclient.ServerCreateRequestIp{
			IpaddrUuid: attr.(string),
		}
		requestBody.Relations.PublicIps = append(requestBody.Relations.PublicIps, ip)
	}

	//Add public network if we have an IP
	if len(requestBody.Relations.PublicIps) > 0 {
		publicNetwork, err := client.GetNetworkPublic()
		if err != nil {
			return err
		}
		network := gsclient.ServerCreateRequestNetwork{
			NetworkUuid: publicNetwork.Properties.ObjectUuid,
		}
		requestBody.Relations.Networks = append(requestBody.Relations.Networks, network)
	}

	if attr, ok := d.GetOk("network"); ok {
		for _, value := range attr.(*schema.Set).List() {
			network := value.(map[string]interface{})
			createNetworkRequest := gsclient.ServerCreateRequestNetwork{
				NetworkUuid: network["object_uuid"].(string),
				BootDevice:  network["bootdevice"].(bool),
			}
			requestBody.Relations.Networks = append(requestBody.Relations.Networks, createNetworkRequest)
		}
	}

	response, err := client.CreateServer(requestBody)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", requestBody.Name, err)
	}

	d.SetId(response.ServerUuid)

	log.Printf("[DEBUG] The id for %s has been set to: %v", requestBody.Name, response.ServerUuid)

	power := d.Get("power").(bool)
	if power {
		client.StartServer(d.Id())
	}

	return resourceGridscaleServerRead(d, meta)
}

func resourceGridscaleServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Id()
	err := client.StopServer(id)
	if err != nil {
		return err
	}
	err = client.DeleteServer(id)

	return err
}

func resourceGridscaleServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	shutdownRequired := false

	var err error

	if d.HasChange("cores") {
		old, new := d.GetChange("cores")
		if new.(int) < old.(int) || d.Get("legacy").(bool) { //Legacy systems don't support updating the memory while running
			shutdownRequired = true
		}
	}

	if d.HasChange("memory") {
		old, new := d.GetChange("memory")
		if new.(int) < old.(int) || d.Get("legacy").(bool) { //Legacy systems don't support updating the memory while running
			shutdownRequired = true
		}
	}

	requestBody := gsclient.ServerUpdateRequest{
		Name:            d.Get("name").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          d.Get("labels").(*schema.Set).List(),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
	}

	if shutdownRequired {
		err = client.ShutdownServer(d.Id())
		if err != nil {
			return err
		}
	}

	err = client.UpdateServer(d.Id(), requestBody)
	if err != nil {
		return err
	}

	// Make sure the server in is the expected power state.
	// The StartServer and ShutdownServer functions do a check to see if the server isn't already running, so we don't need to do that here.
	if d.Get("power").(bool) {
		err = client.StartServer(d.Id())
	} else {
		err = client.ShutdownServer(d.Id())
	}
	if err != nil {
		return err
	}

	return resourceGridscaleServerRead(d, meta)

}
