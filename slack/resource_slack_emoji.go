package slack

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSlackEmoji() *schema.Resource {
	return &schema.Resource{
		Create: resourceSlackEmojiCreate,
		Read:   resourceSlackEmojiRead,
		Delete: resourceSlackEmojiDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func downloadFileReadCloser(path string) (io.ReadCloser, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func resourceSlackEmojiCreate(d *schema.ResourceData, meta interface{}) error {
	token := meta.(slackProvider).token

	name := d.Get("name").(string)
	sourceURL := d.Get("source_url").(string)

	sourceFile, err := downloadFileReadCloser(sourceURL)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	err = uploadEmoji(sourceFile, name, token)
	if err != nil {
		return err
	}

	return resourceSlackEmojiRead(d, meta)
}

func resourceSlackEmojiDelete(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*slackProvider).token

	url := "https://slack.com/api/emoji.remove"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("received emoji.remove response with code %d", res.StatusCode)
	}

	return nil
}

func resourceSlackEmojiRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	client := meta.(*slackProvider).apiClient

	emoji, err := client.GetEmoji()
	if err != nil {
		return err
	}

	resolvedEmoji := resolveEmojiAliases(emoji)
	if url, ok := resolvedEmoji[name]; ok {
		d.SetId(name)
		d.Set("url", url)
		return nil
	}

	return nil
}

func uploadEmoji(emoji io.Reader, name string, token string) error {
	url := "https://slack.com/api/emoji.add"

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	imageField, err := mw.CreateFormFile("image", "emoji.jpg")
	if err != nil {
		return err
	}
	io.Copy(imageField, emoji)

	nameField, err := mw.CreateFormField("name")
	if err != nil {
		return err
	}
	io.WriteString(nameField, name)

	modeField, err := mw.CreateFormField("mode")
	if err != nil {
		return err
	}
	io.WriteString(modeField, "data")

	contentType := mw.FormDataContentType()
	mw.Close()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", contentType)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	ioutil.ReadAll(res.Body)
	res.Body.Close()

	return nil
}
