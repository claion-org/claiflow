package k8s

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var supportedMethodList = []string{
	// v1
	"kubernetes.configmaps.delete.v1",
	"kubernetes.configmaps.get.v1",
	"kubernetes.configmaps.list.v1",
	"kubernetes.configmaps.patch.v1",
	"kubernetes.endpoints.delete.v1",
	"kubernetes.endpoints.get.v1",
	"kubernetes.endpoints.list.v1",
	"kubernetes.endpoints.patch.v1",
	"kubernetes.events.delete.v1",
	"kubernetes.events.get.v1",
	"kubernetes.events.list.v1",
	"kubernetes.events.patch.v1",
	"kubernetes.limitranges.delete.v1",
	"kubernetes.limitranges.get.v1",
	"kubernetes.limitranges.list.v1",
	"kubernetes.limitranges.patch.v1",
	"kubernetes.namespaces.delete.v1",
	"kubernetes.namespaces.get.v1",
	"kubernetes.namespaces.list.v1",
	"kubernetes.namespaces.patch.v1",
	"kubernetes.nodes.delete.v1",
	"kubernetes.nodes.get.v1",
	"kubernetes.nodes.list.v1",
	"kubernetes.nodes.patch.v1",
	"kubernetes.persistentvolumeclaims.delete.v1",
	"kubernetes.persistentvolumeclaims.get.v1",
	"kubernetes.persistentvolumeclaims.list.v1",
	"kubernetes.persistentvolumeclaims.patch.v1",
	"kubernetes.persistentvolumes.delete.v1",
	"kubernetes.persistentvolumes.get.v1",
	"kubernetes.persistentvolumes.list.v1",
	"kubernetes.persistentvolumes.patch.v1",
	"kubernetes.pods.delete.v1",
	"kubernetes.pods.exec.v1",
	"kubernetes.pods.get.v1",
	"kubernetes.pods.list.v1",
	"kubernetes.pods.patch.v1",
	"kubernetes.podtemplates.delete.v1",
	"kubernetes.podtemplates.get.v1",
	"kubernetes.podtemplates.list.v1",
	"kubernetes.podtemplates.patch.v1",
	"kubernetes.replicationcontrollers.delete.v1",
	"kubernetes.replicationcontrollers.get.v1",
	"kubernetes.replicationcontrollers.list.v1",
	"kubernetes.replicationcontrollers.patch.v1",
	"kubernetes.resourcequotas.delete.v1",
	"kubernetes.resourcequotas.get.v1",
	"kubernetes.resourcequotas.list.v1",
	"kubernetes.resourcequotas.patch.v1",
	"kubernetes.secrets.delete.v1",
	"kubernetes.secrets.get.v1",
	"kubernetes.secrets.list.v1",
	"kubernetes.secrets.patch.v1",
	"kubernetes.serviceaccounts.delete.v1",
	"kubernetes.serviceaccounts.get.v1",
	"kubernetes.serviceaccounts.list.v1",
	"kubernetes.serviceaccounts.patch.v1",
	"kubernetes.services.delete.v1",
	"kubernetes.services.get.v1",
	"kubernetes.services.list.v1",
	"kubernetes.services.patch.v1",

	// apps/v1
	"kubernetes.controllerrevisions.delete.apps/v1",
	"kubernetes.controllerrevisions.get.apps/v1",
	"kubernetes.controllerrevisions.list.apps/v1",
	"kubernetes.controllerrevisions.patch.apps/v1",
	"kubernetes.daemonsets.delete.apps/v1",
	"kubernetes.daemonsets.get.apps/v1",
	"kubernetes.daemonsets.list.apps/v1",
	"kubernetes.daemonsets.patch.apps/v1",
	"kubernetes.deployments.delete.apps/v1",
	"kubernetes.deployments.get.apps/v1",
	"kubernetes.deployments.list.apps/v1",
	"kubernetes.deployments.patch.apps/v1",
	"kubernetes.replicasets.delete.apps/v1",
	"kubernetes.replicasets.get.apps/v1",
	"kubernetes.replicasets.list.apps/v1",
	"kubernetes.replicasets.patch.apps/v1",
	"kubernetes.statefulsets.delete.apps/v1",
	"kubernetes.statefulsets.get.apps/v1",
	"kubernetes.statefulsets.list.apps/v1",
	"kubernetes.statefulsets.patch.apps/v1",

	// networking.k8s.io/v1
	"kubernetes.ingressclasses.delete.networking.k8s.io/v1",
	"kubernetes.ingressclasses.get.networking.k8s.io/v1",
	"kubernetes.ingressclasses.list.networking.k8s.io/v1",
	"kubernetes.ingressclasses.patch.networking.k8s.io/v1",
	"kubernetes.ingresses.delete.networking.k8s.io/v1",
	"kubernetes.ingresses.get.networking.k8s.io/v1",
	"kubernetes.ingresses.list.networking.k8s.io/v1",
	"kubernetes.ingresses.patch.networking.k8s.io/v1",
	"kubernetes.networkpolicies.delete.networking.k8s.io/v1",
	"kubernetes.networkpolicies.get.networking.k8s.io/v1",
	"kubernetes.networkpolicies.list.networking.k8s.io/v1",
	"kubernetes.networkpolicies.patch.networking.k8s.io/v1",

	// storage.k8s.io/v1
	"kubernetes.csidrivers.delete.storage.k8s.io/v1",
	"kubernetes.csidrivers.get.storage.k8s.io/v1",
	"kubernetes.csidrivers.list.storage.k8s.io/v1",
	"kubernetes.csidrivers.patch.storage.k8s.io/v1",
	"kubernetes.csinodes.delete.storage.k8s.io/v1",
	"kubernetes.csinodes.get.storage.k8s.io/v1",
	"kubernetes.csinodes.list.storage.k8s.io/v1",
	"kubernetes.csinodes.patch.storage.k8s.io/v1",
	"kubernetes.csistoragecapacities.delete.storage.k8s.io/v1",
	"kubernetes.csistoragecapacities.get.storage.k8s.io/v1",
	"kubernetes.csistoragecapacities.list.storage.k8s.io/v1",
	"kubernetes.csistoragecapacities.patch.storage.k8s.io/v1",
	"kubernetes.storageclasses.delete.storage.k8s.io/v1",
	"kubernetes.storageclasses.get.storage.k8s.io/v1",
	"kubernetes.storageclasses.list.storage.k8s.io/v1",
	"kubernetes.storageclasses.patch.storage.k8s.io/v1",
	"kubernetes.volumeattachments.delete.storage.k8s.io/v1",
	"kubernetes.volumeattachments.get.storage.k8s.io/v1",
	"kubernetes.volumeattachments.list.storage.k8s.io/v1",
	"kubernetes.volumeattachments.patch.storage.k8s.io/v1",

	// monitoring.coreos.com/v1
	"kubernetes.alertmanagers.delete.monitoring.coreos.com/v1",
	"kubernetes.alertmanagers.get.monitoring.coreos.com/v1",
	"kubernetes.alertmanagers.list.monitoring.coreos.com/v1",
	"kubernetes.alertmanagers.patch.monitoring.coreos.com/v1",
	"kubernetes.podmonitors.delete.monitoring.coreos.com/v1",
	"kubernetes.podmonitors.get.monitoring.coreos.com/v1",
	"kubernetes.podmonitors.list.monitoring.coreos.com/v1",
	"kubernetes.podmonitors.patch.monitoring.coreos.com/v1",
	"kubernetes.probes.delete.monitoring.coreos.com/v1",
	"kubernetes.probes.get.monitoring.coreos.com/v1",
	"kubernetes.probes.list.monitoring.coreos.com/v1",
	"kubernetes.probes.patch.monitoring.coreos.com/v1",
	"kubernetes.prometheuses.delete.monitoring.coreos.com/v1",
	"kubernetes.prometheuses.get.monitoring.coreos.com/v1",
	"kubernetes.prometheuses.list.monitoring.coreos.com/v1",
	"kubernetes.prometheuses.patch.monitoring.coreos.com/v1",
	"kubernetes.prometheusrules.delete.monitoring.coreos.com/v1",
	"kubernetes.prometheusrules.get.monitoring.coreos.com/v1",
	"kubernetes.prometheusrules.list.monitoring.coreos.com/v1",
	"kubernetes.prometheusrules.patch.monitoring.coreos.com/v1",
	"kubernetes.servicemonitors.delete.monitoring.coreos.com/v1",
	"kubernetes.servicemonitors.get.monitoring.coreos.com/v1",
	"kubernetes.servicemonitors.list.monitoring.coreos.com/v1",
	"kubernetes.servicemonitors.patch.monitoring.coreos.com/v1",
	"kubernetes.thanosrulers.delete.monitoring.coreos.com/v1",
	"kubernetes.thanosrulers.get.monitoring.coreos.com/v1",
	"kubernetes.thanosrulers.list.monitoring.coreos.com/v1",
	"kubernetes.thanosrulers.patch.monitoring.coreos.com/v1",

	// monitoring.coreos.com/v1alpha1
	"kubernetes.alertmanagerconfigs.delete.monitoring.coreos.com/v1alpha1",
	"kubernetes.alertmanagerconfigs.get.monitoring.coreos.com/v1alpha1",
	"kubernetes.alertmanagerconfigs.list.monitoring.coreos.com/v1alpha1",
	"kubernetes.alertmanagerconfigs.patch.monitoring.coreos.com/v1alpha1",

	// batch/v1
	"kubernetes.cronjobs.delete.batch/v1",
	"kubernetes.cronjobs.get.batch/v1",
	"kubernetes.cronjobs.list.batch/v1",
	"kubernetes.cronjobs.patch.batch/v1",
	"kubernetes.jobs.delete.batch/v1",
	"kubernetes.jobs.get.batch/v1",
	"kubernetes.jobs.list.batch/v1",
	"kubernetes.jobs.patch.batch/v1",

	// admissionregistration.k8s.io/v1
	"kubernetes.mutatingwebhookconfigurations.delete.admissionregistration.k8s.io/v1",
	"kubernetes.mutatingwebhookconfigurations.get.admissionregistration.k8s.io/v1",
	"kubernetes.mutatingwebhookconfigurations.list.admissionregistration.k8s.io/v1",
	"kubernetes.mutatingwebhookconfigurations.patch.admissionregistration.k8s.io/v1",
	"kubernetes.validatingwebhookconfigurations.delete.admissionregistration.k8s.io/v1",
	"kubernetes.validatingwebhookconfigurations.get.admissionregistration.k8s.io/v1",
	"kubernetes.validatingwebhookconfigurations.list.admissionregistration.k8s.io/v1",
	"kubernetes.validatingwebhookconfigurations.patch.admissionregistration.k8s.io/v1",

	// apiextensions.k8s.io/v1
	"kubernetes.customresourcedefinitions.delete.apiextensions.k8s.io/v1",
	"kubernetes.customresourcedefinitions.get.apiextensions.k8s.io/v1",
	"kubernetes.customresourcedefinitions.list.apiextensions.k8s.io/v1",
	"kubernetes.customresourcedefinitions.patch.apiextensions.k8s.io/v1",

	// apiregistration.k8s.io/v1
	"kubernetes.apiservices.delete.apiregistration.k8s.io/v1",
	"kubernetes.apiservices.get.apiregistration.k8s.io/v1",
	"kubernetes.apiservices.list.apiregistration.k8s.io/v1",
	"kubernetes.apiservices.patch.apiregistration.k8s.io/v1",

	// autoscaling/v2
	"kubernetes.horizontalpodautoscalers.delete.autoscaling/v2",
	"kubernetes.horizontalpodautoscalers.get.autoscaling/v2",
	"kubernetes.horizontalpodautoscalers.list.autoscaling/v2",
	"kubernetes.horizontalpodautoscalers.patch.autoscaling/v2",

	// certificates.k8s.io/v1
	"kubernetes.certificatesigningrequests.delete.certificates.k8s.io/v1",
	"kubernetes.certificatesigningrequests.get.certificates.k8s.io/v1",
	"kubernetes.certificatesigningrequests.list.certificates.k8s.io/v1",
	"kubernetes.certificatesigningrequests.patch.certificates.k8s.io/v1",

	// coordination.k8s.io/v1
	"kubernetes.leases.delete.coordination.k8s.io/v1",
	"kubernetes.leases.get.coordination.k8s.io/v1",
	"kubernetes.leases.list.coordination.k8s.io/v1",
	"kubernetes.leases.patch.coordination.k8s.io/v1",

	// discovery.k8s.io/v1
	"kubernetes.endpointslices.delete.discovery.k8s.io/v1",
	"kubernetes.endpointslices.get.discovery.k8s.io/v1",
	"kubernetes.endpointslices.list.discovery.k8s.io/v1",
	"kubernetes.endpointslices.patch.discovery.k8s.io/v1",

	// node.k8s.io/v1
	"kubernetes.runtimeclasses.delete.node.k8s.io/v1",
	"kubernetes.runtimeclasses.get.node.k8s.io/v1",
	"kubernetes.runtimeclasses.list.node.k8s.io/v1",
	"kubernetes.runtimeclasses.patch.node.k8s.io/v1",

	// policy/v1
	"kubernetes.poddisruptionbudgets.delete.policy/v1",
	"kubernetes.poddisruptionbudgets.get.policy/v1",
	"kubernetes.poddisruptionbudgets.list.policy/v1",
	"kubernetes.poddisruptionbudgets.patch.policy/v1",

	// rbac.authorization.k8s.io/v1
	"kubernetes.clusterrolebindings.delete.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterrolebindings.get.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterrolebindings.list.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterrolebindings.patch.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterroles.delete.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterroles.get.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterroles.list.rbac.authorization.k8s.io/v1",
	"kubernetes.clusterroles.patch.rbac.authorization.k8s.io/v1",
	"kubernetes.rolebindings.delete.rbac.authorization.k8s.io/v1",
	"kubernetes.rolebindings.get.rbac.authorization.k8s.io/v1",
	"kubernetes.rolebindings.list.rbac.authorization.k8s.io/v1",
	"kubernetes.rolebindings.patch.rbac.authorization.k8s.io/v1",
	"kubernetes.roles.delete.rbac.authorization.k8s.io/v1",
	"kubernetes.roles.get.rbac.authorization.k8s.io/v1",
	"kubernetes.roles.list.rbac.authorization.k8s.io/v1",
	"kubernetes.roles.patch.rbac.authorization.k8s.io/v1",

	// scheduling.k8s.io/v1
	"kubernetes.priorityclasses.delete.scheduling.k8s.io/v1",
	"kubernetes.priorityclasses.get.scheduling.k8s.io/v1",
	"kubernetes.priorityclasses.list.scheduling.k8s.io/v1",
	"kubernetes.priorityclasses.patch.scheduling.k8s.io/v1",
}

func InitCommandFuncs() (map[string]func(map[string]interface{}) ([]byte, error), error) {
	funcList := make(map[string]func(map[string]interface{}) ([]byte, error))

	for _, method := range supportedMethodList {
		mlist := strings.SplitN(method, ".", 4)

		if len(mlist) != 4 {
			return nil, fmt.Errorf("there is not enough method(%s) for k8s. want(4) but got(%d)", method, len(mlist))
		}

		gv, err := schema.ParseGroupVersion(mlist[3])
		if err != nil {
			return nil, err
		}

		resource := mlist[1]
		verb := mlist[2]

		fn := func(resource, verb string) func(map[string]interface{}) ([]byte, error) {
			return func(args map[string]interface{}) ([]byte, error) {
				c, err := GetClient()
				if err != nil {
					return nil, err
				}

				res, err := c.ResourceRequest(gv, resource, verb, args)
				return []byte(res), err
			}
		}

		funcList[method] = fn(resource, verb)
	}

	return funcList, nil
}
