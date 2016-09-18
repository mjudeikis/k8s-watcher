package main

import (
	"fmt"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
	"net/http"

	"time"

	log "github.com/Sirupsen/logrus"
)

func main() {
	// Skip 0-th argument containing the binary's name.
	log.SetLevel(log.DebugLevel)

	var option WatcherOptions
	option.init()
	watcher := option.Validate()

	log.Debugf("Credentials and details [%s] [%s:%s][%s]", watcher.RestConfig.Host, watcher.RestConfig.Username, watcher.RestConfig.Password, watcher.RestConfig.BearerToken)
	config := &restclient.Config{
		Host: "https://10.2.2.2:8443",
		//WTF, why token is not picked when doing this???
		BearerToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6IndhdGNoZXItdG9rZW4tcjJuODUiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoid2F0Y2hlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImU1YTNkN2M1LTdkMzQtMTFlNi05MjgzLTUyNTQwMGM1ODNhZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OndhdGNoZXIifQ.ah3gCWagFUYFmQWSUVqXC0jZYPc2Y_jiPNMeJyqHbb252ReP_HcBS8dKUX3tedIJVGw3xOjssoHCX6swnjPP10MOW8ROHaQQioXvbOzUvWKxw0_cTd6_2Q4u_EjUtLrcPN0_Xgjsi5D3uukCtGMe5a3J6NEROVAUxgGIKInf888-cBvLimfTQFIw-pegUnd0AYqTozBIEMvX2ak4FzUPp_zODYY7-iXEMAMCgM929KZYqgtdZGYZ5NlhMvZJFmSls-mPq1vgtaLB3Q11WnarfAqrpuXMJf2UoY6ONL1yOho0ZQgobceDkSdDDCVVUC2Umr3PoIT13d20ZtJxxz01lA",
		//BearerToken: fmt.Sprintf("%s", watcher.RestConfig.BearerToken),
		Insecure: watcher.RestConfig.Insecure,
	}

	//kube clien
	kubeClient, err := client.New(config)

	if err != nil {
		log.Fatalln("Client not created sucessfully:", err)
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
	fog.Info("Pod deleted: " + pod.ObjectMeta.Name)
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
