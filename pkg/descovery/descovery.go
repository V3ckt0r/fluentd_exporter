package discovery

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	wg sync.WaitGroup
	// mutex is used to define a critical section of code.
	mutex sync.Mutex
)

// get home dir of user
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		h = h + "/.kube/config"
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func GetNamespaces(clientset *kubernetes.Clientset) []string {
	var namespaces []string
	ns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("Error getting namespaces, %v", err)
	}
	for _, n := range ns.Items {
		namespaces = append(namespaces, n.ObjectMeta.Name)
		log.Debugf("Namespaces: ", n.ObjectMeta.Name)
	}
	return namespaces
}

// Get services with the app: fluentd tag
func GetServices(namespace string, clientset *kubernetes.Clientset) ([]Service, error) {
	var services []Service
	servs, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Could not get services in namespace, %v", err)
	}
	// look through services
	for _, s := range servs.Items {
		if s.ObjectMeta.Labels["app"] == "fluentd" {
			retService := &Service{
				Name:      s.ObjectMeta.Name,
				ClusterIP: s.Spec.ClusterIP,
				Port:      s.Spec.Ports[0].Port,
			}
			services = append(services, *retService)
			log.Debugf("Found Service match: ", s.ObjectMeta.Name)
		}
	}
	return services, nil
}

func GetAllServices(clientset *kubernetes.Clientset) ([]Service, error) {
	services := make([]Service, 0)
	client := GetClient()
	namespaces := GetNamespaces(client)

	wg.Add(len(namespaces))

	for _, n := range namespaces {
		go func(n string) {
			defer wg.Done()
			service, err := GetServices(n, client)
			if err != nil {
				log.Errorf("Error getting services: %v", err)
			} else {
				for _, s := range service {
					// Only allow one goroutine through this
					// critical section at a time. Needed to avoid race condition.
					mutex.Lock()
					services = append(services, s)
					mutex.Unlock()
				}
			}
		}(n)
	}
	wg.Wait()
	if len(services) <= 0 {
		log.Errorf("No services found")
		return services, errors.New("No services found")
	}
	return services, nil
}

// service definition
type Service struct {
	Name      string
	ClusterIP string
	Port      int32
}

func GetClient() *kubernetes.Clientset {
	// creates the in-cluster config, this is default intended setup
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
