package main

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	//"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
	"net/http"
	"os"

	"time"

	log "github.com/Sirupsen/logrus"
)

func main() {
	// Skip 0-th argument containing the binary's name.
	log.SetLevel(log.DebugLevel)

	var option WatcherOptions
	config := option.init(os.Args)

	log.Debugf("Credentials and details [%s] [%s:%s][%s]", config.Host, config.Username, config.Password, config.BearerToken)
	kubeClient, err := client.New(config)

	if err != nil {
		log.Fatalln("Client not created sukubeClient, err := *client.New(option.Watcher.Config)cessfully:", err)
	}

	namespaces, err := kubeClient.Namespaces().List(api.ListOptions{})
	log.Info(err)
	log.Info("namespaces")
	for _, v := range namespaces.Items {
		log.Info(v.GetGenerateName)
	}

	//Create a cache to store Pods
	var podsStore cache.Store

	//Watch for Pods
	podsStore = watchPods(kubeClient, podsStore)
	//Keep alive
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func podCreated(obj interface{}) {
	pod := obj.(*api.Pod)
	log.Info("Pod created: " + pod.ObjectMeta.Name)
}
func podDeleted(obj interface{}) {
	pod := obj.(*api.Pod)
	log.Info("Pod deleted: " + pod.ObjectMeta.Name)
}

func watchPods(client *client.Client, store cache.Store) cache.Store {
	//Define what we want to look for (Pods)
	watchlist := cache.NewListWatchFromClient(client, "pods", api.NamespaceAll, fields.Everything())
	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.Pod{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    podCreated,
			DeleteFunc: podDeleted,
		},
	)
	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)
	return eStore
}
