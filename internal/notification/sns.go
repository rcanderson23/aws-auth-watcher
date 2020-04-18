package notification

import (
	"bytes"
	"fmt"
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
	oldString := mapToString(oldObj.(*v1.ConfigMap).Data)
	newString := mapToString(newObj.(*v1.ConfigMap).Data)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "Old ConfigMap Data:\n%s\n\nNew ConfigMap Data:\n%s", oldString, newString)
	message := b.String()
	_, err := a.SnsClient.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: a.SnsTopic,
	})
	if err != nil {
		klog.Error(err.Error())
	}
}

func (a *AwsSns) PublishDelete(obj interface{}) {
	oldString := mapToString(obj.(*v1.ConfigMap).Data)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "Deleted ConfigMap Data:\n%s", oldString)
	message := b.String()
	_, err := a.SnsClient.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: a.SnsTopic,
	})
	if err != nil {
		klog.Error(err.Error())
	}

}

func mapToString(data map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range data {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
