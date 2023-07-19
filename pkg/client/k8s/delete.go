package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceDelete(gv schema.GroupVersion, resource string, params map[string]interface{}) error {
	var err error

	var namespace string
	var name string

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return err
	}

	if found, err := FindCastFromMap(params, "name", &name); found && err != nil {
		return err
	} else if !found {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultK8sTimeout)
	defer cancel()

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "configmaps":
			err = c.client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "endpoints":
			err = c.client.CoreV1().Endpoints(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "events":
			err = c.client.CoreV1().Events(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "namespaces":
			err = c.client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "nodes":
			err = c.client.CoreV1().Nodes().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "persistentvolumes":
			err = c.client.CoreV1().PersistentVolumes().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "persistentvolumeclaims":
			err = c.client.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "pods":
			err = c.client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "secrets":
			err = c.client.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "services":
			err = c.client.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "limitranges":
			err = c.client.CoreV1().LimitRanges(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "podtemplates":
			err = c.client.CoreV1().PodTemplates(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "replicationcontrollers":
			err = c.client.CoreV1().ReplicationControllers(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "resourcequotas":
			err = c.client.CoreV1().ResourceQuotas(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "serviceaccounts":
			err = c.client.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apps/v1":
		switch resource {
		case "daemonsets":
			err = c.client.AppsV1().DaemonSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "deployments":
			err = c.client.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "replicasets":
			err = c.client.AppsV1().ReplicaSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "statefulsets":
			err = c.client.AppsV1().StatefulSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "controllerrevisions":
			err = c.client.AppsV1().ControllerRevisions(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "networking.k8s.io/v1":
		switch resource {
		case "ingresses":
			err = c.client.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "ingressclasses":
			err = c.client.NetworkingV1().IngressClasses().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "networkpolicies":
			err = c.client.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "storage.k8s.io/v1":
		switch resource {
		case "storageclasses":
			err = c.client.StorageV1().StorageClasses().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "csidrivers":
			err = c.client.StorageV1().CSIDrivers().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "csinodes":
			err = c.client.StorageV1().CSINodes().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "csistoragecapacities":
			err = c.client.StorageV1().CSIStorageCapacities(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "volumeattachments":
			err = c.client.StorageV1().VolumeAttachments().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "monitoring.coreos.com/v1":
		switch resource {
		case "prometheuses":
			err = c.mclient.MonitoringV1().Prometheuses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "prometheusrules":
			err = c.mclient.MonitoringV1().PrometheusRules(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "servicemonitors":
			err = c.mclient.MonitoringV1().ServiceMonitors(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "alertmanagers":
			err = c.mclient.MonitoringV1().Alertmanagers(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "podmonitors":
			err = c.mclient.MonitoringV1().PodMonitors(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "probes":
			err = c.mclient.MonitoringV1().Probes(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "thanosrulers":
			err = c.mclient.MonitoringV1().ThanosRulers(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "monitoring.coreos.com/v1alpha1":
		switch resource {
		case "alertmanagerconfigs":
			err = c.mclient.MonitoringV1alpha1().AlertmanagerConfigs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "batch/v1":
		switch resource {
		case "cronjobs":
			err = c.client.BatchV1().CronJobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "jobs":
			err = c.client.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "admissionregistration.k8s.io/v1":
		switch resource {
		case "mutatingwebhookconfigurations":
			err = c.client.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "validatingwebhookconfigurations":
			err = c.client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apiextensions.k8s.io/v1":
		switch resource {
		case "customresourcedefinitions":
			err = c.apiextv1client.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apiregistration.k8s.io/v1":
		switch resource {
		case "apiservices":
			err = c.aggrev1client.ApiregistrationV1().APIServices().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "autoscaling/v2":
		switch resource {
		case "horizontalpodautoscalers":
			err = c.client.AutoscalingV2().HorizontalPodAutoscalers(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "certificates.k8s.io/v1":
		switch resource {
		case "certificatesigningrequests":
			err = c.client.CertificatesV1().CertificateSigningRequests().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "coordination.k8s.io/v1":
		switch resource {
		case "leases":
			err = c.client.CoordinationV1beta1().Leases(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "discovery.k8s.io/v1":
		switch resource {
		case "endpointslices":
			err = c.client.DiscoveryV1().EndpointSlices(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "node.k8s.io/v1":
		switch resource {
		case "runtimeclasses":
			err = c.client.NodeV1().RuntimeClasses().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "policy/v1":
		switch resource {
		case "poddisruptionbudgets":
			err = c.client.PolicyV1().PodDisruptionBudgets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "rbac.authorization.k8s.io/v1":
		switch resource {
		case "clusterrolebindings":
			err = c.client.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "clusterroles":
			err = c.client.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "rolebindings":
			err = c.client.RbacV1().RoleBindings(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		case "roles":
			err = c.client.RbacV1().Roles(namespace).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "scheduling.k8s.io/v1":
		switch resource {
		case "priorityclasses":
			err = c.client.SchedulingV1().PriorityClasses().Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	default:
		err = fmt.Errorf("unsupported group version(%s)", gv.Identifier())
	}

	if err != nil {
		return err
	}

	return nil
}
