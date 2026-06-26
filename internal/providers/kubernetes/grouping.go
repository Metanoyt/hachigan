package kubernetes

const (
	labelAppKubernetesName = "app.kubernetes.io/name"
	labelApp               = "app"
)

func applicationName(labels map[string]string, fallback string) string {
	if labels != nil {
		if value := labels[labelAppKubernetesName]; value != "" {
			return value
		}
		if value := labels[labelApp]; value != "" {
			return value
		}
	}
	return fallback
}
