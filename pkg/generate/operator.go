package generate

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/leg100/stok/version"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type Operator struct {
	Name      string
	Namespace string
	Image     string
}

func (o *Operator) Generate(out io.Writer) error {
	resources := []interface{}{
		o.deployment(),
		o.serviceAccount(),
		o.clusterRole(),
		o.clusterRoleBinding(),
	}

	var sb strings.Builder
	for _, res := range resources {
		sb.WriteString("---\n")
		y, err := toYaml(res)
		if err != nil {
			return err
		}
		sb.WriteString(string(y))
	}
	fmt.Fprint(out, sb.String())

	return nil
}

// Operator's ClusterRole.
//
// Some of these permissions are necessary for the operator's
// metrics service:
// * services: c/d/g/l/p/u/w
// * deployments: g
// * replicasets: g
//
// The workspace controller manages:
// * roles
// * rolebindings
// * pvcs
func (o *Operator) clusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: o.Name,
			Labels: map[string]string{
				// Name of the application
				"app":                    "stok",
				"app.kubernetes.io/name": "stok",

				// Name of higher-level application this app is part of
				"app.kubernetes.io/part-of": "stok",

				// The tool being used to manage the operation of an application
				"app.kubernetes.io/managed-by": "stok-cli",

				// Unique name of instance within application
				"app.kubernetes.io/instance": "stok",

				// Current version of application
				"version":                   version.Version,
				"app.kubernetes.io/version": version.Version,

				// Component within architecture
				"component":                   "operator",
				"app.kubernetes.io/component": "operator",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"pods",
					"persistentvolumeclaims",
					"configmaps",
					"secrets",
					"services",
					"serviceaccounts",
				},
				Verbs: []string{
					"create",
					"delete",
					"get",
					"list",
					"patch",
					"update",
					"watch",
				},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{
					"deployments",
					"replicasets",
				},
				Verbs: []string{
					"get",
				},
			},
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{
					"roles",
					"rolebindings",
				},
				Verbs: []string{
					"create",
					"delete",
					"get",
					"list",
					"patch",
					"update",
					"watch",
				},
			},
			{
				APIGroups: []string{"stok.goalspike.com"},
				Resources: []string{"*"},
				Verbs: []string{
					"create",
					"delete",
					"get",
					"list",
					"patch",
					"update",
					"watch",
				},
			},
		},
	}
}

func (o *Operator) clusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: o.Name,
			Labels: map[string]string{
				// Name of the application
				"app":                    "stok",
				"app.kubernetes.io/name": "stok",

				// Name of higher-level application this app is part of
				"app.kubernetes.io/part-of": "stok",

				// The tool being used to manage the operation of an application
				"app.kubernetes.io/managed-by": "stok-cli",

				// Unique name of instance within application
				"app.kubernetes.io/instance": "stok",

				// Current version of application
				"version":                   version.Version,
				"app.kubernetes.io/version": version.Version,

				// Component within architecture
				"component":                   "operator",
				"app.kubernetes.io/component": "operator",
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      o.Name,
				Namespace: o.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     o.Name,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func (o *Operator) serviceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.Namespace,
			Labels: map[string]string{
				// Name of the application
				"app":                    "stok",
				"app.kubernetes.io/name": "stok",

				// Name of higher-level application this app is part of
				"app.kubernetes.io/part-of": "stok",

				// The tool being used to manage the operation of an application
				"app.kubernetes.io/managed-by": "stok-cli",

				// Unique name of instance within application
				"app.kubernetes.io/instance": "stok",

				// Current version of application
				"version":                   version.Version,
				"app.kubernetes.io/version": version.Version,

				// Component within architecture
				"component":                   "operator",
				"app.kubernetes.io/component": "operator",
			},
		},
	}
}

func (o *Operator) deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.Namespace,
			Labels: map[string]string{
				// Name of the application
				"app":                    "stok",
				"app.kubernetes.io/name": "stok",

				// Name of higher-level application this app is part of
				"app.kubernetes.io/part-of": "stok",

				// The tool being used to manage the operation of an application
				"app.kubernetes.io/managed-by": "stok-cli",

				// Unique name of instance within application
				"app.kubernetes.io/instance": "stok",

				// Current version of application
				"version":                   version.Version,
				"app.kubernetes.io/version": version.Version,

				// Component within architecture
				"component":                   "operator",
				"app.kubernetes.io/component": "operator",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component": "operator",
					"app.kubernetes.io/name":      "stok",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": "operator",
						"app.kubernetes.io/name":      "stok",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: o.Name,
					Containers: []corev1.Container{
						{
							Name:            "stok-operator",
							Image:           o.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"stok"},
							Args:            []string{"operator"},
							Env: []corev1.EnvVar{
								{
									Name:  "WATCH_NAMESPACE",
									Value: "",
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name:  "OPERATOR_NAME",
									Value: "stok",
								},
								{
									Name:  "RUNNER_IMAGE",
									Value: o.Image,
								},
							},
							TerminationMessagePolicy: "FallbackToLogsOnError",
						},
					},
				},
			},
		},
	}
}

// Convert struct to YAML, leveraging JSON struct tags by first converting to JSON
func toYaml(obj interface{}) ([]byte, error) {
	j, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return nil, err
	}

	return y, nil
}

func int32Ptr(i int32) *int32 { return &i }
