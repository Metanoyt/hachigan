package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hachigan/hachigan/internal/config"
	"github.com/hachigan/hachigan/internal/domain"
)

type Client struct {
	clientset   kubernetes.Interface
	contextName string
	namespaces  map[string]struct{}
}

func New(cfg config.Config) (*Client, error) {
	kubeconfig := cfg.Kubeconfig
	if kubeconfig == "" {
		if home, err := os.UserHomeDir(); err == nil {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	overrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("build kubernetes config: %w", err)
	}
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	nsFilter := make(map[string]struct{}, len(cfg.Namespaces))
	for _, namespace := range cfg.Namespaces {
		if namespace != "" {
			nsFilter[namespace] = struct{}{}
		}
	}

	return &Client{clientset: clientset, contextName: rawConfig.CurrentContext, namespaces: nsFilter}, nil
}

func (c *Client) ClusterSummary(ctx context.Context) (domain.ClusterSummary, error) {
	namespaces, err := c.Namespaces(ctx)
	if err != nil {
		return domain.ClusterSummary{}, err
	}
	workloads, err := c.Workloads(ctx)
	if err != nil {
		return domain.ClusterSummary{}, err
	}
	services, err := c.Services(ctx)
	if err != nil {
		return domain.ClusterSummary{}, err
	}
	return domain.ClusterSummary{
		ContextName:    c.contextName,
		NamespaceCount: len(namespaces),
		WorkloadCount:  len(workloads),
		ServiceCount:   len(services),
	}, nil
}

func (c *Client) Namespaces(ctx context.Context) ([]string, error) {
	list, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces: %w", err)
	}
	names := make([]string, 0, len(list.Items))
	for _, ns := range list.Items {
		if c.includeNamespace(ns.Name) {
			names = append(names, ns.Name)
		}
	}
	sort.Strings(names)
	return names, nil
}

func (c *Client) Workloads(ctx context.Context) ([]domain.Workload, error) {
	var out []domain.Workload
	namespaces, err := c.Namespaces(ctx)
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaces {
		deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("list deployments in %s: %w", namespace, err)
		}
		for _, item := range deployments.Items {
			desired := int(item.Status.Replicas)
			ready := int(item.Status.ReadyReplicas)
			out = append(out, domain.Workload{Name: item.Name, Namespace: namespace, Kind: "Deployment", Desired: desired, Ready: ready, Health: domain.DeriveHealth(desired, ready)})
		}

		statefulSets, err := c.clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("list statefulsets in %s: %w", namespace, err)
		}
		for _, item := range statefulSets.Items {
			desired := int(item.Status.Replicas)
			ready := int(item.Status.ReadyReplicas)
			out = append(out, domain.Workload{Name: item.Name, Namespace: namespace, Kind: "StatefulSet", Desired: desired, Ready: ready, Health: domain.DeriveHealth(desired, ready)})
		}

		daemonSets, err := c.clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("list daemonsets in %s: %w", namespace, err)
		}
		for _, item := range daemonSets.Items {
			desired := int(item.Status.DesiredNumberScheduled)
			ready := int(item.Status.NumberReady)
			out = append(out, domain.Workload{Name: item.Name, Namespace: namespace, Kind: "DaemonSet", Desired: desired, Ready: ready, Health: domain.DeriveHealth(desired, ready)})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Namespace+"/"+out[i].Name < out[j].Namespace+"/"+out[j].Name
	})
	return out, nil
}

func (c *Client) Services(ctx context.Context) ([]domain.Service, error) {
	var out []domain.Service
	namespaces, err := c.Namespaces(ctx)
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaces {
		services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("list services in %s: %w", namespace, err)
		}
		for _, item := range services.Items {
			out = append(out, domain.Service{
				Name:      item.Name,
				Namespace: namespace,
				Type:      string(item.Spec.Type),
				ClusterIP: item.Spec.ClusterIP,
				Ports:     servicePorts(item.Spec.Ports),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Namespace+"/"+out[i].Name < out[j].Namespace+"/"+out[j].Name
	})
	return out, nil
}

func (c *Client) Applications(ctx context.Context) ([]domain.Application, error) {
	grouped := map[string]*domain.Application{}
	namespaces, err := c.Namespaces(ctx)
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaces {
		deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, item := range deployments.Items {
			app := ensureApp(grouped, applicationName(item.Labels, item.Name), namespace)
			app.Workloads = append(app.Workloads, domain.WorkloadRef{Name: item.Name, Namespace: namespace, Kind: "Deployment"})
			app.Health = domain.AggregateHealth(app.Health, domain.DeriveHealth(int(item.Status.Replicas), int(item.Status.ReadyReplicas)))
		}
		statefulSets, err := c.clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, item := range statefulSets.Items {
			app := ensureApp(grouped, applicationName(item.Labels, item.Name), namespace)
			app.Workloads = append(app.Workloads, domain.WorkloadRef{Name: item.Name, Namespace: namespace, Kind: "StatefulSet"})
			app.Health = domain.AggregateHealth(app.Health, domain.DeriveHealth(int(item.Status.Replicas), int(item.Status.ReadyReplicas)))
		}
		daemonSets, err := c.clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, item := range daemonSets.Items {
			app := ensureApp(grouped, applicationName(item.Labels, item.Name), namespace)
			app.Workloads = append(app.Workloads, domain.WorkloadRef{Name: item.Name, Namespace: namespace, Kind: "DaemonSet"})
			app.Health = domain.AggregateHealth(app.Health, domain.DeriveHealth(int(item.Status.DesiredNumberScheduled), int(item.Status.NumberReady)))
		}
		services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, item := range services.Items {
			app := ensureApp(grouped, applicationName(item.Labels, item.Name), namespace)
			app.Services = append(app.Services, domain.ServiceRef{Name: item.Name, Namespace: namespace})
		}
	}
	apps := make([]domain.Application, 0, len(grouped))
	for _, app := range grouped {
		if len(app.Workloads) == 0 {
			app.Health = domain.HealthUnknown
		}
		apps = append(apps, *app)
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Namespace+"/"+apps[i].Name < apps[j].Namespace+"/"+apps[j].Name
	})
	return apps, nil
}

func (c *Client) includeNamespace(namespace string) bool {
	if len(c.namespaces) == 0 {
		return true
	}
	_, ok := c.namespaces[namespace]
	return ok
}

func ensureApp(grouped map[string]*domain.Application, name, namespace string) *domain.Application {
	key := namespace + "/" + name
	if app, ok := grouped[key]; ok {
		return app
	}
	grouped[key] = &domain.Application{Name: name, Namespace: namespace, Health: domain.HealthUnknown}
	return grouped[key]
}

func servicePorts(ports []corev1.ServicePort) []string {
	out := make([]string, 0, len(ports))
	for _, port := range ports {
		value := strconv.Itoa(int(port.Port))
		if port.Protocol != "" {
			value += "/" + string(port.Protocol)
		}
		out = append(out, value)
	}
	return out
}
