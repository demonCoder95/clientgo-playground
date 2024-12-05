package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func createLocalClient() (kubernetes.Interface, error) {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}

	return client, nil
}

func createConfigMap(client kubernetes.Interface, namespace string, name string, data map[string]string) error {
	_, err := client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to create a configmap: %v", err)
	}

	return nil
}

func main() {
	client, err := createLocalClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	var tc testCase
	tc.data.groups = [][]string{{"CollaboratorPowerUser"}}
	tc.data.resources = []string{"secrets"}
	tc.data.verbs = []string{"get", "list", "watch"}
	tc.data.namespaces = []string{"kube-system"}
	tc.run(context.TODO(), client, false)
	fmt.Println(tc.output)

}
