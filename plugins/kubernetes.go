package plugins

import (
	"context"

	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesPlugin struct {
	config *jmfzf.KubernetesConfig
	client *kubernetes.Clientset
}

func NewKubernetesPlugin(options interface{}) (jmfzf.Plugin, error) {
	var config jmfzf.KubernetesConfig
	if options != nil {
		err := jmfzf.MapToStruct(options, &config)
		if err != nil {
			return nil, err
		}
	}
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return &KubernetesPlugin{config: &config, client: client}, nil
}

func (plugin *KubernetesPlugin) List(options *jmfzf.ListOptions) ([]terminal.Host, error) {
	var result []terminal.Host
	for _, namespace := range plugin.config.Namespaces {
		podsInterface := plugin.client.CoreV1().Pods(namespace)
		podList, err := podsInterface.List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, pod := range podList.Items {
			// filter out pods that are not running
			if pod.Status.Phase == v1.PodRunning {
				for _, container := range pod.Spec.Containers {
					result = append(result, terminal.Host{
						Type: terminal.TerminalTypeContainer,
						ContainerInfo: terminal.Pod{
							Name:       pod.Name,
							Namespace:  pod.Namespace,
							KubeConfig: plugin.config.KubeConfig,
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

func (plugin *KubernetesPlugin) Name() string {
	return "kubernetes"
}
