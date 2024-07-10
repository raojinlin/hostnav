package plugins

import (
	"context"

	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesPlugin struct {
	option *hostnav.KubernetesOption
	client *kubernetes.Clientset
}

func NewKubernetesPlugin() *KubernetesPlugin {
	return &KubernetesPlugin{option: &hostnav.KubernetesOption{}}
}

func (p *KubernetesPlugin) Init(option interface{}) error {
	err := hostnav.MapToStruct(option, p.option)
	if err != nil {
		return err
	}
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", p.option.KubeConfig)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return err
	}

	p.client = client
	return nil
}

func (p *KubernetesPlugin) List(options *ListOptions) ([]terminal.Host, error) {
	var result []terminal.Host
	for _, namespace := range p.option.Namespaces {
		podList, err := p.client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, pod := range podList.Items {
			isReady := false
			for _, podCondition := range pod.Status.Conditions {
				if podCondition.Type == v1.ContainersReady && podCondition.Status == v1.ConditionTrue {
					isReady = true
				}
			}

			// filter out pods that are not running
			if pod.Status.Phase == v1.PodRunning && isReady {
				for _, container := range pod.Spec.Containers {
					result = append(result, terminal.Host{
						Type: terminal.TerminalTypeContainer,
						ContainerInfo: terminal.Pod{
							Name:       pod.Name,
							Namespace:  pod.Namespace,
							KubeConfig: p.option.KubeConfig,
							Container: terminal.Container{
								Name:    container.Name,
								Command: "/bin/sh",
								Id:      container.Name,
							},
						},
					})
				}
			}

		}
	}

	return result, nil
}

func (p *KubernetesPlugin) Name() string {
	return "kubernetes"
}

func (p *KubernetesPlugin) Cache() bool {
	return false
}
