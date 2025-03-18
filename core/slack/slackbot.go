package slack

// import (
// 	"github.com/chunhui2001/go-starter/core/config"
// 	"github.com/slack-go/slack"
// )

// var (
// 	logger = config.Log
// )

// func init() {

// }

// func NewClient(slackToken string) *slack.Client {
// 	return slack.New(slackToken, slack.OptionDebug(false))
// }

// func SlackMessage(slackClient *slack.Client, channelID string, title string, message map[string]string) {
// 	fields := make([]slack.AttachmentField, 0, len(message))

// 	for k, v := range message {
// 		fields = append(fields, slack.AttachmentField{
// 			Title: k,
// 			Value: v,
// 		})
// 	}

// 	attachment := slack.Attachment{
// 		Pretext: title,
// 		Color: "#36a64f", // 绿色
// 		//Color: "ff000" // 红色
// 	}

// 	_, _, err := slackClient.PostMessage(
// 		channelID,
// 		slack.MsgOptionAttachments(attachment),
// 	)

// 	if err != nil {
// 		logger.Errorf(`发送Slack预警消息-失败: channelID=%s, title=%s, ErrorMessage=%s`, channelID, title, err)
// 	} else {
// 		logger.Infof(`发送Slack预警消息-成功: channelID=%s, title=%s`, channelID, title)
// 	}
// }

// // https://knock.app/blog/how-to-render-tables-in-slack-markdown
// func SlackTable(slackClient *slack.Client, channelID string, title string) {

// 	titleBlock := slack.SectionBlock{
// 		Type:      slack.MBTHeader,
// 		Text:      slack.NewTextBlockObject("plain_text", "💰 Our Savings", true, false),
// 		Fields:    nil,
// 		Accessory: nil,
// 	}

// 	headerBlock := slack.SectionBlock{
// 		Type:      slack.MBTSection,
// 		Fields:    []*slack.TextBlockObject{
// 			slack.NewTextBlockObject("mrkdwn", "*Month*", false, false),
// 			slack.NewTextBlockObject("mrkdwn", "*Savings*", false, false),
// 		},
// 		Accessory: nil,
// 	}

// 	row1 := slack.SectionBlock{
// 		Type:      slack.MBTSection,
// 		Fields:    []*slack.TextBlockObject{
// 			slack.NewTextBlockObject("mrkdwn", "January", false, false),
// 			slack.NewTextBlockObject("mrkdwn", "$250", false, false),
// 		},
// 		Accessory: nil,
// 	}

// 	row2 := slack.SectionBlock{
// 		Type:      slack.MBTSection,
// 		Fields:    []*slack.TextBlockObject{
// 			slack.NewTextBlockObject("mrkdwn", "February", false, false),
// 			slack.NewTextBlockObject("mrkdwn", "$80", false, false),
// 		},
// 		Accessory: nil,
// 	}

// 	row3 := slack.SectionBlock{
// 		Type:      slack.MBTSection,
// 		Fields:    []*slack.TextBlockObject{
// 			slack.NewTextBlockObject("mrkdwn", "March", false, false),
// 			slack.NewTextBlockObject("mrkdwn", "$420", false, false),
// 		},
// 		Accessory: nil,
// 	}

// 	// 构造行数据
// 	blocks := []slack.Block{
// 		titleBlock,
// 		headerBlock,
// 		slack.NewDividerBlock(),
// 		row1,
// 		slack.NewDividerBlock(),
// 		row2,
// 		slack.NewDividerBlock(),
// 		row3,
// 		slack.NewDividerBlock(),
// 	}

// 	// 发送消息
// 	_, _, err := slackClient.PostMessage(
// 		channelID,
// 		slack.MsgOptionBlocks(blocks...),
// 	)

// 	if err != nil {
// 		logger.Errorf(`发送Slack预警消息-失败: channelID=%s, title=%s, ErrorMessage=%s`, channelID, title, err)
// 	} else {
// 		logger.Infof(`发送Slack预警消息-成功: channelID=%s, title=%s`, channelID, title)
// 	}
// }


// func SlackMrkdwn(slackClient *slack.Client, channelID string, title string) {

// 	// 表头
// 	header := "*Month*\t\t*Savings*\t\t*Expenses*"

// 	// 表内容
// 	rows := []string{
// 		"January\t\t$250\t\t$150",
// 		"February\t\t$80\t\t$50",
// 		"March\t\t$420\t\t$200",
// 	}

// 	// 拼接表格内容
// 	text := header + "\n" + "```\n" + formatRows(rows) + "\n```"

// 	// 创建 SectionBlock
// 	section := slack.NewSectionBlock(
// 		slack.NewTextBlockObject("mrkdwn", text, false, false),
// 		nil,
// 		nil,
// 	)

// 	// 发送消息
// 	_, _, err := slackClient.PostMessage(
// 		channelID,
// 		slack.MsgOptionBlocks(section),
// 	)

// 	if err != nil {
// 		logger.Errorf(`发送Slack预警消息-失败: channelID=%s, title=%s, ErrorMessage=%s`, channelID, title, err)
// 	} else {
// 		logger.Infof(`发送Slack预警消息-成功: channelID=%s, title=%s`, channelID, title)
// 	}
// }
