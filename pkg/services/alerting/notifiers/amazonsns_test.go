package notifiers

import (
	"testing"

	"github.com/grafana/grafana/pkg/components/simplejson"
	m "github.com/grafana/grafana/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAmazonSnsNotifier(t *testing.T) {
	Convey("Amazon SNS notifier tests", t, func() {

		Convey("Parsing alert notification from settings", func() {
			Convey("empty settings should return error", func() {
				json := `{ }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "ops",
					Type:     "amazonsns",
					Settings: settingsJSON,
				}

				_, err := NewAmazonSnsNotifier(model)
				So(err, ShouldNotBeNil)
			})

			Convey("from settings", func() {
				json := `
				{
          "url": "http://google.com"
				}`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "ops",
					Type:     "amazonsns",
					Settings: settingsJSON,
				}

				not, err := NewAmazonSnsNotifier(model)
				snsNotifier := not.(*AmazonSnsNotifier)

				So(err, ShouldBeNil)
				So(snsNotifier.Name, ShouldEqual, "ops")
				So(snsNotifier.Type, ShouldEqual, "amazonsns")
				So(snsNotifier.Url, ShouldEqual, "http://google.com")
				So(snsNotifier.Region, ShouldEqual, "")
			})

			Convey("from settings with Region", func() {
				json := `
				{
          "url": "http://google.com",
          "region": "ap-northeast-xx"
				}`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "ops",
					Type:     "amazonsns",
					Settings: settingsJSON,
				}

				not, err := NewAmazonSnsNotifier(model)
				snsNotifier := not.(*AmazonSnsNotifier)

				So(err, ShouldBeNil)
				So(snsNotifier.Name, ShouldEqual, "ops")
				So(snsNotifier.Type, ShouldEqual, "amazonsns")
				So(snsNotifier.Url, ShouldEqual, "http://google.com")
				So(snsNotifier.Region, ShouldEqual, "ap-northeast-xx")
			})
		})
	})
}
