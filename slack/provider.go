package slack

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nlopes/slack"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Slack API token.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"slack_emojis": dataSourceSlackEmojis(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"slack_emoji": resourceSlackEmoji(),
		},
		ConfigureFunc: providerConfigure,
	}
}

type slackProvider struct {
	apiClient *slack.Client
	token     string
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)

	return &slackProvider{slack.New(token), token}, nil
}
