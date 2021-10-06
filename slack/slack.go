package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/htenjo/gh_statistics/config"
	"github.com/htenjo/gh_statistics/github"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

var (
	DangerStyle  = "danger"
	PrimaryStyle = "primary"
)

func NewPlainTextBlock(text string) PlanTextBlock {
	return PlanTextBlock{
		Type: "plain_text",
		Text: text,
	}
}

func NewHeader(text string) HeaderBlock {
	return HeaderBlock{
		Type: "header",
		Text: NewPlainTextBlock(text),
	}
}

func NewActions() ButtonSection {
	return ButtonSection{
		Type:     "actions",
		Elements: []ButtonBlock{},
	}
}

func SendSlackMessage(messageTitle string, prInfo *[]github.RepoPR) {
	message := WebhookMessage{}
	header := NewHeader(messageTitle)
	actions := NewActions()

	redPrs := getButtonsByFlag(prInfo, github.Red)
	yellowPrs := getButtonsByFlag(prInfo, github.Yellow)
	GreenPrs := getButtonsByFlag(prInfo, github.Green)

	actions.Elements = append(actions.Elements, getButtonElements(&redPrs)...)
	actions.Elements = append(actions.Elements, getButtonElements(&yellowPrs)...)
	actions.Elements = append(actions.Elements, getButtonElements(&GreenPrs)...)
	maxNotifications := math.Min(float64(20), float64(len(actions.Elements)))
	actions.Elements = actions.Elements[0:int(maxNotifications)]

	message.Blocks = append(message.Blocks, header, actions)
	byteResponse, _ := json.MarshalIndent(message, "", "  ")
	log.Printf("%v", string(byteResponse))
	sendNotification(byteResponse)
}

func getButtonsByFlag(prInfo *[]github.RepoPR, flag github.PrReviewFlag) []github.PullRequestDetail {
	var redPrs []github.PullRequestDetail

	for _, pr := range *prInfo {
		for _, info := range pr.Prs {
			if info.ReviewFlag == flag {
				redPrs = append(redPrs, info)
			}
		}
	}

	return redPrs
}

func getButtonElements(prDetails *[]github.PullRequestDetail) []ButtonBlock {
	var buttonBlocks []ButtonBlock

	for _, pr := range *prDetails {
		buttonBlock := ButtonBlock{
			Type: "button",
			Text: PlanTextBlock{
				Type: "plain_text",
				Text: pr.Title,
			},
			Url: pr.HtmlUrl,
		}

		if pr.ReviewFlag == github.Red {
			buttonBlock.Style = &DangerStyle
		} else if pr.ReviewFlag == github.Green {
			buttonBlock.Style = &PrimaryStyle
		}

		buttonBlocks = append(buttonBlocks, buttonBlock)
	}

	return buttonBlocks
}

func sendNotification(message []byte) {
	resp, err := http.Post(config.SlackWebhookUrl(), "application/json", bytes.NewReader(message))

	if err != nil {
		log.Fatal(err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(bodyText))
}