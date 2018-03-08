package notifiers

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

func init() {
	alerting.RegisterNotifier(&alerting.NotifierPlugin{
		Type:        "amazonsns",
		Name:        "Amazon SNS",
		Description: "Sends notifications to Amazon SNS Topic",
		Factory:     NewAmazonSnsNotifier,
		OptionsTemplate: `
      <h3 class="page-heading">Amazon SNS settings</h3>
      <div class="gf-form">
        <span class="gf-form-label width-10">Region</span>
        <input type="text" required class="gf-form-input max-width-26" ng-model="ctrl.model.settings.region" placeholder="ap-northeast-1"></input>
      </div>
      <div class="gf-form">
        <span class="gf-form-label width-10">Topic ARN</span>
        <input type="text" required class="gf-form-input max-width-26" ng-model="ctrl.model.settings.url" placeholder="arn:aws:sns:REGION:ACCOUNTID:TOPICNAME"></input>
      </div>
    `,
	})

}

func NewAmazonSnsNotifier(model *m.AlertNotification) (alerting.Notifier, error) {
	url := model.Settings.Get("url").MustString()
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find url property in settings"}
	}

	region := model.Settings.Get("region").MustString()

	return &AmazonSnsNotifier{
		NotifierBase: NewNotifierBase(model.Id, model.IsDefault, model.Name, model.Type, model.Settings),
		Url:          url,
		Region:       region,
		log:          log.New("alerting.notifier.amazonSns"),
	}, nil
}

type AmazonSnsNotifier struct {
	NotifierBase
	Url    string
	Region string
	log    log.Logger
}

func (this *AmazonSnsNotifier) ShouldNotify(context *alerting.EvalContext) bool {
	return defaultShouldNotify(context)
}

func (this *AmazonSnsNotifier) Notify(evalContext *alerting.EvalContext) error {
	this.log.Info("Executing Amazon SNS notification", "ruleId", evalContext.Rule.Id, "notification", this.Name)

	ruleUrl, err := evalContext.GetRuleUrl()
	if err != nil {
		this.log.Error("Failed get rule link", "error", err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(this.Region)},
	)
	client := sns.New(sess)

	body := map[string]interface{}{
		"title":       evalContext.GetNotificationTitle(),
		"message":     evalContext.Rule.Message,
		"ruleId":      evalContext.Rule.Id,
		"ruleName":    evalContext.Rule.Name,
		"state":       evalContext.Rule.State,
		"evalMatches": evalContext.EvalMatches,
		"ruleUrl":     ruleUrl,
	}
	data, _ := json.Marshal(&body)

	input := &sns.PublishInput{
		Message:  aws.String(string(data)),
		Subject:  aws.String("Grafana Alert"),
		TopicArn: aws.String(this.Url),
	}

	req, _ := client.PublishRequest(input)
	err = req.Send()
	if err != nil {
		this.log.Error("Failed to send SNS Topic", "error", err, "name", this.Name)
		return err
	}

	return nil
}
