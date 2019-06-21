provider "slack" {
    alias = "source"
    token = "${var.source_slack_token}"
}

provider "slack" {
    alias = "destination"
    token = "${var.destination_slack_token}"
}

data "slack_emojis" "all" {
    provider = slack.source
}

resource "slack_emoji" "destination" {
    provider = slack.destination
    count = length(data.slack_emojis.all.names)

    name = element(tolist(data.slack_emojis.all.names), count.index)
    source_url = lookup(data.slack_emojis.all.urls, element(tolist(data.slack_emojis.all.names), count.index))
}

output "emoji" {
    value = "${data.slack_emojis.all.urls}"
}
