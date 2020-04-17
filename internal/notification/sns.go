package notification

import (
	"github.com/aws/aws-sdk-go/service/sns"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type Notifier interface {
	PublishChange()
	PublishDelete()
}

type AwsSns struct {
	SnsClient *sns.SNS
	SnsTopic  *string
}

func (a *AwsSns) PublishChange(oldObj, newObj interface{}) {
	message := oldObj.(*v1.ConfigMap).Data["mapRoles"]
	result, err := a.SnsClient.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: a.SnsTopic,
	})
	if err != nil {
		klog.Error(err.Error())
	}
	klog.Info(result)

}

func (a *AwsSns) PublishDelete(obj interface{}) {

}
