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
	_ "github.com/ysku/my-k8s-custom-controller/pkg/logging"
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
	var factory informers.SharedInformerFactory
	if conf.Namespace != "" {
		factory = informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(conf.Namespace))
	} else {
		factory = informers.NewSharedInformerFactory(clientset, 0)
	}

	stop := make(chan struct{})
	defer close(stop)

	controllers := map[bool]*controller.LoggingController{
		conf.EnablePod: controller.NewPodLoggingController(factory),
		conf.EnableDeployment: controller.NewDeploymentLoggingController(factory),
		conf.EnableJob: controller.NewJobLoggingController(factory),
		conf.EnablePersistentVolume: controller.NewPersistentVolumeLoggingController(factory),
		conf.EnableService: controller.NewServiceLoggingController(factory),
		conf.EnableHpa: controller.NewHPALoggingController(factory),
		conf.EnablePdb: controller.NewPDBLoggingController(factory),
	}

	for enable, c := range controllers {
		if enable {
			log.Printf("[main] starting controller %s\n", c.Name)
			go c.Run(1, stop)
		}
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
	cmd.Flags().StringVar(&conf.Namespace, "namespace", "", "target namespace, all namespaces by default")
	cmd.Flags().BoolVar(&conf.EnablePod, "enable-pod", true, "enable watching pod")
	cmd.Flags().BoolVar(&conf.EnableDeployment, "enable-deployment", false, "enable watching deployment")
	cmd.Flags().BoolVar(&conf.EnableJob, "enable-job", false, "enable watching job")
	cmd.Flags().BoolVar(&conf.EnablePersistentVolume, "enable-persistent-volume", false, "enable watching persistent volume")
	cmd.Flags().BoolVar(&conf.EnableService, "enable-service", false, "enable watching service")
	cmd.Flags().BoolVar(&conf.EnableHpa, "enable-hpa", false, "enable watching horizontal pod autoscaler")
	cmd.Flags().BoolVar(&conf.EnablePdb, "enable-pdb", false, "enable watching pod disruption budget")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
