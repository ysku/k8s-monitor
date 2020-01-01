package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/ysku/my-k8s-custom-controller/pkg/config"
	"github.com/ysku/my-k8s-custom-controller/pkg/controller"
	"github.com/ysku/my-k8s-custom-controller/pkg/version"
)

const AppName = "k8s-custom-controller"

var conf *config.Config

func NewClientSet() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}
	return clientset
}

func run(cmd *cobra.Command, args []string) {
	clientset := NewClientSet()
	factory := informers.NewSharedInformerFactory(clientset, 0)

	stop := make(chan struct{})
	defer close(stop)

	if conf.EnablePod {
		podController := controller.NewPodLoggingController(factory)
		go podController.Run(1, stop)
	}
	if conf.EnableService {
		serviceController := controller.NewServiceLoggingController(factory)
		go serviceController.Run(1, stop)
	}

	log.Print("[main] start running logging")
	select {}
}

func main() {
	log.Print("[main] Shared Informer app started")

	cmd := &cobra.Command{
		Use:          AppName,
		Short:        "platform-agent version: " + version.Version,
		Run:          run,
		SilenceUsage: true,
		Version:      version.Version,
	}
	conf := config.NewConfig()
	cmd.Flags().BoolVar(&conf.EnablePod, "enable-pod", true, "enable watching pod")
	cmd.Flags().BoolVar(&conf.EnableService, "enable-service", true, "enable watching service")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
