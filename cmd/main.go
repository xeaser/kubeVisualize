package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx := context.Background()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	kubeconfig := filepath.Join(
		homeDir, ".kube", "config",
	)
	namespace := "default"
	k8sClient, err := getClient(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	pods, err := getPods(ctx, namespace, k8sClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	log.Printf("Total Pods found: %d in namespace: %s", len(pods.Items), namespace)
}

func getClient(configLocation string) (typev1.CoreV1Interface, error) {
	kubeconfig := filepath.Clean(configLocation)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1(), nil
}

func getPods(ctx context.Context, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error) {
	pods, err := k8sClient.Pods(namespace).List(ctx, metav1.ListOptions{})

	for _, pod := range pods.Items {
		fmt.Fprintf(os.Stdout, "pod name: %v\n", pod.Name)
	}
	return pods, err
}
