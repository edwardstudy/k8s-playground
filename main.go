package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var podGVR = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "pods",
}

var podGVK = schema.GroupVersionKind{
	Group:   "",
	Version: "v1",
	Kind:    "Pod",
}

type unstructuredClient struct {
	dynamicClient dynamic.Interface
}

func NewClient(dynamicClient dynamic.Interface) *unstructuredClient {
	return &unstructuredClient{
		dynamicClient: dynamicClient,
	}
}

func (u *unstructuredClient) createWithObject(pod corev1.Pod) error {
	return nil
}

func (u *unstructuredClient) createWithYaml(pod corev1.Pod) error {
	bytes, err := yaml.Marshal(pod)
	if err != nil {
		return err
	}

	fmt.Printf("yaml string: %s\n", string(bytes))

	obj := &unstructured.Unstructured{}
	decoder := serializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := decoder.Decode(bytes, &podGVK, obj); err != nil {
		return err
	}

	fmt.Printf("obj: %v\n", obj)

	_, err = u.dynamicClient.Resource(podGVR).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client := NewClient(dynamicClient)
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fake",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "container1",
					Image: "busybox:1.25",
				},
			},
		},
	}

	fmt.Printf("create with yaml\n")
	err = client.createWithYaml(pod)
	if err != nil {
		panic(err)
	}
}
