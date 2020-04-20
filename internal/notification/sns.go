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

func (a *AwsSns) PublishChange(oldObj, newObj interface{}, cluster string) {
	oldString := mapToString(oldObj.(*v1.ConfigMap).Data)
	newString := mapToString(newObj.(*v1.ConfigMap).Data)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "ClusterName: %s\n\nOld ConfigMap Data:\n%s\n\nNew ConfigMap Data:\n%s", cluster, oldString, newString)
	message := b.String()
	_, err := a.SnsClient.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: a.SnsTopic,
	})
	if err != nil {
		klog.Error(err.Error())
	}
}

func (a *AwsSns) PublishDelete(obj interface{}, cluster string) {
	oldString := mapToString(obj.(*v1.ConfigMap).Data)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "Cluster: %s\n\nDeleted ConfigMap Data:\n%s", cluster, oldString)
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
