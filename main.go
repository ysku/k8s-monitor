package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
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

	if conf.EnableDeployment {
		deploymentController := controller.NewDeploymentLoggingController(factory)
		go deploymentController.Run(1, stop)
	}
	if conf.EnableJob {
		jobController := controller.NewJobLoggingController(factory)
		go jobController.Run(1, stop)
	}
	if conf.EnablePersistentVolume {
		persistentVolumeController := controller.NewPersistentVolumeLoggingController(factory)
		go persistentVolumeController.Run(1, stop)
	}
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

func validate(cmd *cobra.Command, args []string) {
	if err := conf.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "[main] error: %v", err)
		os.Exit(1)
	}
}

func main() {
	cmd := &cobra.Command{
		Use:          AppName,
		Short:        "platform-agent version: " + version.Version,
		PreRun:       validate,
		Run:          run,
		SilenceUsage: true,
		Version:      version.Version,
	}
	conf = config.NewConfig()
	cmd.Flags().BoolVar(&conf.EnableDeployment, "enable-deployment", false, "enable watching deployment")
	cmd.Flags().BoolVar(&conf.EnableJob, "enable-job", false, "enable watching job")
	cmd.Flags().BoolVar(&conf.EnablePersistentVolume, "enable-persistent-volume", false, "enable watching persistent volume")
	cmd.Flags().BoolVar(&conf.EnablePod, "enable-pod", false, "enable watching pod")
	cmd.Flags().BoolVar(&conf.EnableService, "enable-service", false, "enable watching service")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
