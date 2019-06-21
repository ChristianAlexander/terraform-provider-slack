package slack

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSlackEmojis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSlackEmojisRead,

		Schema: map[string]*schema.Schema{
			"urls": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"names": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceSlackEmojisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*slackProvider).apiClient

	emoji, err := client.GetEmoji()
	if err != nil {
		return err
	}

	resolvedEmoji := resolveEmojiAliases(emoji)

	d.SetId("emoji")
	d.Set("urls", resolvedEmoji)
	d.Set("names", getMapKeys(resolvedEmoji))

	return nil
}

func getMapKeys(m map[string]string) []string {
	result := []string{}

	for k := range m {
		result = append(result, k)
	}

	return result
}

func resolveEmojiAliases(emoji map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range emoji {
		if split := strings.Split(v, ":"); len(split) == 2 && split[0] == "alias" {
			if resolvedValue := emoji[split[1]]; resolvedValue != "" {
				result[k] = resolvedValue
			}
		} else {
			result[k] = v
		}
	}
	return result
}
