package main

import (
	"context"
	"net/http"

	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var metric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "openshift_cluster_cpu_usage",
	Help: "CPU usage in the OpenShift cluster",
})

func collectMetrics() {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("openshift-monitoring").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod name: %s, Status: %s\n", pod.ObjectMeta.Name, pod.Status.Phase)
		// TODO: Collect metrics from Prometheus endpoint and update gauge metric
	}
}

func init() {
	prometheus.MustRegister(metric)
}

func main() {
	go func() {
		for {
			collectMetrics()
			time.Sleep(time.Minute)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
