package fetcher

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	watchtools "k8s.io/client-go/tools/watch"

	"github.com/claion-org/claiflow/pkg/client/internal/service"
	"github.com/claion-org/claiflow/pkg/client/k8s"
	"github.com/claion-org/claiflow/pkg/client/log"
)

func (f *Fetcher) UpgradeClient(serviceId string, args map[string]interface{}) (err error) {
	log.Debugf("client upgrade: start\n")

	t := time.Now()

	// service processing status update
	for {
		up := service.CreateUpdateService(serviceId, 1, 0, service.StepStatusProcessing, service.Result{}, t, time.Time{})
		if err := f.serverAPI.UpdateServices(context.Background(), up); err != nil {
			log.Errorf("client upgrade: failed to update service status('processing'): error=%v\n", err)
			f.RetryHandshake()
			continue
		}
		break
	}

	// fetcher polling stop
	defer func() {
		if err != nil {
			log.Errorf("client upgrade: failed to upgrade: error=%v\n", err)

			for {
				up := service.CreateUpdateService(serviceId, 1, 0, service.StepStatusFail, service.Result{Err: err}, t, time.Now())
				if err := f.serverAPI.UpdateServices(context.Background(), up); err != nil {
					log.Errorf("client upgrade: failed to update services: error=%v\n", err)
					f.RetryHandshake()
					continue
				}
				break
			}
		}
	}()
	log.Debugf("client upgrade: stop polling\n")

	// check arguments
	var imageTag string
	var timeout time.Duration
	var envConfigs map[string]string
	if args != nil {
		imageTagInf, ok := args["image_tag"]
		if !ok || imageTagInf == nil {
			return fmt.Errorf("failed to find image_tag argument")
		}

		imageTag, ok = imageTagInf.(string)
		if !ok {
			return fmt.Errorf("failed type assertion for image_tag argument: interface to string")
		}

		if imageTag == "" {
			return fmt.Errorf("image_tag argument is empty")
		}

		envConfigsInf, ok := args["env_configs"]
		if ok && envConfigsInf != nil {
			cm, ok := envConfigsInf.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unsupported env_configs argument type(%T): want(map[string]interface{})", envConfigsInf)
			}

			envConfigs = make(map[string]string)
			for k, v := range cm {
				envConfigs[k] = fmt.Sprintf("%v", v)
			}
		}

		timeoutInf, ok := args["timeout"]
		if ok && timeoutInf != nil {
			switch v := timeoutInf.(type) {
			case float64: // encoding/json/decode.go:53
				timeout = time.Second * time.Duration(v)
			case string:
				timeoutInt, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("failed convert string to int for timeout argument")
				}
				timeout = time.Second * time.Duration(timeoutInt)
			default:
				return fmt.Errorf("unsupported timeout argument type(%T): supported type(string, int)", v)
			}
		}
	} else {
		return fmt.Errorf("argument is empty")
	}

	// clean up the remaining services before upgrade(timeout:30s)
	log.Debugf("client upgrade: waiting remain service proccess: waiting_timeout=30s\n")
	for cnt := 0; cnt < 10; cnt++ {
		<-time.After(time.Second * 3)
		remainServices := f.RemainServices()
		if len(remainServices) == 0 {
			break
		}

		buf := bytes.Buffer{}
		buf.WriteString("remain services:")
		for uuid, status := range remainServices {
			buf.WriteString(fmt.Sprintf("\n\tuuid: %s, status: %d", uuid, status))
		}
		log.Infof(buf.String() + "\n")
	}
	log.Debugf("client upgrade: end remain service proccess\n")

	// get namespace
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return err
	}
	ns := string(namespace)

	// get pod_name
	podName, err := os.Hostname()
	if err != nil {
		return err
	}
	if podName == "" {
		return fmt.Errorf("pod name is empty")
	}

	log.Debugf("client upgrade: found pod info: namespace=%s, pod_name=%s\n", ns, podName)

	// get k8s client
	k8sClient, err := k8s.GetClient()
	if err != nil {
		return err
	}

	// get self-pod -> owner replicaset -> owner deployment
	deploymentObj, err := findDeploymentFromPod(k8sClient, ns, podName)
	if err != nil {
		return err
	}
	prevDeploymentObj := deploymentObj.DeepCopy()

	recoverFunc, err := applyEnvConfingsToDeployment(k8sClient, deploymentObj, envConfigs)
	if err != nil {
		return err
	}

	prevImage := deploymentObj.Spec.Template.Spec.Containers[0].Image
	upgradeImage := replaceImageTag(prevImage, imageTag)

	// patch deployment's image to upgrade image
	log.Debugf("client upgrade: request to patch deployment's image with %s\n", upgradeImage)
	deploymentObj.Spec.Template.Spec.Containers[0].Image = upgradeImage
	if deploymentObj.Spec.Template.Annotations == nil {
		deploymentObj.Spec.Template.Annotations = make(map[string]string)
	}
	deploymentObj.Spec.Template.Annotations["claiflow/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = k8sClient.ResourcePatch(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{
		"namespace":  ns,
		"name":       deploymentObj.Name,
		"patch_type": "json",
		"patch_data": []map[string]interface{}{
			{
				"op":    "replace",
				"path":  "/spec/template",
				"value": &deploymentObj.Spec.Template,
			},
		},
	})
	if err != nil {
		if err := recoverFunc(); err != nil {
			log.Warnf("failed to recover env config. error=%s\n", err)
		}
		return err
	}
	log.Debugf("client upgrade: patched deployment\n")

	// watch my deployment
	watchInf, err := k8sClient.ResourceWatch(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{"namespace": ns, "name": deploymentObj.Name})
	if err != nil {
		return err
	}

	var watchCtx context.Context
	var watchCancel context.CancelFunc
	if timeout == 0 {
		watchCtx, watchCancel = context.WithCancel(context.Background())
	} else {
		watchCtx, watchCancel = context.WithTimeout(context.Background(), timeout)
	}
	defer watchCancel()

	log.Debugf("client upgrade: watch the status of the rollout until it's done\n")
	if _, watchErr := watchtools.UntilWithoutRetry(watchCtx, watchInf, func(e watch.Event) (bool, error) {
		switch t := e.Type; t {
		case watch.Added, watch.Modified:
			deploymentObj := e.Object.(*appsv1.Deployment)

			if deploymentObj.Generation <= deploymentObj.Status.ObservedGeneration {
				var cond *appsv1.DeploymentCondition
				for i := range deploymentObj.Status.Conditions {
					c := deploymentObj.Status.Conditions[i]
					if c.Type == appsv1.DeploymentProgressing {
						cond = &c
					}
				}

				if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
					return false, fmt.Errorf("deployment %q exceeded its progress deadline", deploymentObj.Name)
				} else if deploymentObj.Spec.Replicas != nil && deploymentObj.Status.UpdatedReplicas < *deploymentObj.Spec.Replicas {
					log.Debugf("client upgrade: waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated\n", deploymentObj.Name, deploymentObj.Status.UpdatedReplicas, *deploymentObj.Spec.Replicas)
					return false, nil
				} else if deploymentObj.Status.Replicas > deploymentObj.Status.UpdatedReplicas {
					log.Debugf("client upgrade: waiting for deployment %q rollout to finish: %d old replicas are pending termination\n", deploymentObj.Name, deploymentObj.Status.Replicas-deploymentObj.Status.UpdatedReplicas)
					return false, nil
				} else if deploymentObj.Status.AvailableReplicas < deploymentObj.Status.UpdatedReplicas {
					log.Debugf("client upgrade: waiting for deployment %q rollout to finish: %d of %d updated replicas are available\n", deploymentObj.Name, deploymentObj.Status.AvailableReplicas, deploymentObj.Status.UpdatedReplicas)
					return false, nil
				} else {
					log.Debugf("client upgrade: deployment %q successfully rolled out\n", deploymentObj.Name)
					return true, nil
				}
			}
			log.Debugf("client upgrade: waiting for deployment spec update to be observed\n")

			return false, nil

		case watch.Deleted:
			return true, fmt.Errorf("object has been deleted")

		default:
			return true, fmt.Errorf("internal error: unexpected_event=%#v", e)
		}
	}); watchErr != nil {
		// recover env config
		if err := recoverFunc(); err != nil {
			log.Warnf("failed to recover env config. error=%s\n", err)
		}

		// patch my deployment's prev-image
		deploymentObj.Spec.Template.Spec.Containers[0].Image = prevImage
		if _, err := k8sClient.ResourcePatch(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{
			"namespace":  ns,
			"name":       deploymentObj.Name,
			"patch_type": "json",
			"patch_data": []map[string]interface{}{
				{
					"op":    "replace",
					"path":  "/spec/template",
					"value": &prevDeploymentObj.Spec.Template,
				},
			},
		}); err != nil {
			return fmt.Errorf("failed to patch deployment prev-image : watch_error: {%v}, patch_error: {%v}", watchErr, err)
		}

		return watchErr
	}

	// upgrade success
	for {
		up := service.CreateUpdateService(serviceId, 1, 0, service.StepStatusSuccess, service.Result{Body: "client upgrade will be complete"}, t, time.Now())
		if err := f.serverAPI.UpdateServices(context.Background(), up); err != nil {
			log.Errorf("client upgrade: failed to update services: error=%v", err)
			f.RetryHandshake()
			continue
		}
		break
	}

	return nil
}

func findDeploymentFromPod(k8sClient *k8s.Client, ns, podName string) (*appsv1.Deployment, error) {
	podJson, err := k8sClient.ResourceGet(schema.GroupVersion{Version: "v1"}, "pods", map[string]interface{}{"namespace": ns, "name": podName})
	if err != nil {
		return nil, err
	}
	podObj := new(corev1.Pod)
	if err := json.Unmarshal([]byte(podJson), podObj); err != nil {
		return nil, err
	}

	// find owner replicaset
	var replicasetName string
	for _, ownerRef := range podObj.OwnerReferences {
		if ownerRef.Kind == "ReplicaSet" {
			replicasetName = ownerRef.Name
		}
	}
	if replicasetName == "" {
		return nil, fmt.Errorf("failed to find replicaset name")
	}

	log.Debugf("client upgrade: found owner replicaset info: namespace=%s, name=%s\n", ns, replicasetName)

	replicasetJson, err := k8sClient.ResourceGet(schema.GroupVersion{Group: "apps", Version: "v1"}, "replicasets", map[string]interface{}{"namespace": ns, "name": replicasetName})
	if err != nil {
		return nil, err
	}
	replicasetObj := new(appsv1.ReplicaSet)
	if err := json.Unmarshal([]byte(replicasetJson), replicasetObj); err != nil {
		return nil, err
	}

	// find owner deployment
	var deploymentName string
	for _, ownerRef := range replicasetObj.OwnerReferences {
		if ownerRef.Kind == "Deployment" {
			deploymentName = ownerRef.Name
		}
	}
	if deploymentName == "" {
		return nil, fmt.Errorf("failed to find replicaset name")
	}

	log.Debugf("client upgrade: found owner deployment info: namespace=%s, name=%s\n", ns, deploymentName)

	deploymentJson, err := k8sClient.ResourceGet(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{"namespace": ns, "name": deploymentName})
	if err != nil {
		return nil, err
	}
	deploymentObj := new(appsv1.Deployment)
	if err := json.Unmarshal([]byte(deploymentJson), deploymentObj); err != nil {
		return nil, err
	}

	return deploymentObj, nil
}

func replaceImageTag(image, tag string) string {
	imageName := image

	index := strings.LastIndex(image, ":")
	if index != -1 {
		imageName = image[:index]
	}

	return fmt.Sprintf("%s:%s", imageName, tag)
}

func applyEnvConfingsToDeployment(k8sClient *k8s.Client, deploymentObj *appsv1.Deployment, envConfigs map[string]string) (func() error, error) {
	if len(envConfigs) <= 0 {
		return nil, nil
	}

	if len(deploymentObj.Spec.Template.Spec.Containers) <= 0 {
		return nil, fmt.Errorf("deployment containers is empty")
	}

	// precedence : last, env than envFrom
	// 1. envFrom.configMapRef or envFrom.secretRef
	// 2. env.name + env.value
	// 3. env.name + (env.valueFrom.configMapKeyRef or env.valueFrom.secretKeyRef)
	generateResourceCacheKey := func(resourceType, namespace, name string) string { // generate key: (configmap or secret)|namespace|name
		return fmt.Sprintf("%s|%s|%s", resourceType, namespace, name)
	}

	// finding environment variables used in deployment
	defaultEnvPath, envMap, err := findEnvironmentVariablesFromContainer(k8sClient, deploymentObj.Namespace, &deploymentObj.Spec.Template.Spec.Containers[0], generateResourceCacheKey)
	if err != nil {
		return nil, err
	}

	log.Debugf("default env path: %v", defaultEnvPath)

	// extract json patch actions
	generateResourceCacheKeyWithNamespace := func(resourceType, name string) string {
		return generateResourceCacheKey(resourceType, deploymentObj.Namespace, name)
	}
	patchActions := extractPatchActionsFromEnvironmentVariables(envMap, envConfigs, *defaultEnvPath, generateResourceCacheKeyWithNamespace)

	recoverPatchActions := make(map[string]*patchAction)
	// apply environment variable changes
	for cacheKey, patchActionElem := range patchActions {
		if cacheKey == "" { // deployment.spec.template.spec.container[0].env
			for i, jpa := range patchActionElem.jsonPatchAction {
				act := fmt.Sprintf("%v", jpa["op"])
				envName := fmt.Sprintf("%v", jpa["path"])
				envValue := fmt.Sprintf("%v", jpa["value"])

				switch act {
				case "add":
					deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{Name: envName, Value: envValue})
				case "replace":
					for i, envvar := range deploymentObj.Spec.Template.Spec.Containers[0].Env {
						if envvar.Name == envName {
							envvar.Value = envValue
							deploymentObj.Spec.Template.Spec.Containers[0].Env[i] = envvar
							break
						}
					}
				case "remove":
					for i, envvar := range deploymentObj.Spec.Template.Spec.Containers[0].Env {
						if envvar.Name == envName {
							deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env[:i], deploymentObj.Spec.Template.Spec.Containers[0].Env[i:]...)
							break
						}
					}
				default:
					return nil, fmt.Errorf("unknown json patch operation. op=%s", act)
				}
				log.Debugf("changed env in deployment. key=%s, action=%s, updated_value=%s\n", envName, act, envValue)

				// attach recover action
				if pas, ok := recoverPatchActions[cacheKey]; ok {
					pas.recoverJsonPatchAction = append(pas.recoverJsonPatchAction, jpa)
				} else {
					recoverPatchActions[cacheKey] = &patchAction{typ: patchActionElem.typ, resourceType: patchActionElem.resourceType, recoverJsonPatchAction: []map[string]interface{}{patchActionElem.recoverJsonPatchAction[i]}}
				}
			}
		} else { // configmap or secret patch
			// patchAction
			cacheKeySplits := strings.Split(cacheKey, "|")
			if len(cacheKeySplits) < 3 {
				return nil, fmt.Errorf("cache_key is not valid. cache_key=%s", cacheKey)
			}

			if _, err := k8sClient.ResourcePatch(schema.GroupVersion{Version: "v1"}, cacheKeySplits[0], map[string]interface{}{
				"namespace":  cacheKeySplits[1],
				"name":       cacheKeySplits[2],
				"patch_type": "json",
				"patch_data": patchActionElem.jsonPatchAction,
			}); err != nil {
				return nil, err
			}

			// env.valueFrom -> deployment
			for i, jpa := range patchActionElem.jsonPatchAction {
				act := fmt.Sprintf("%v", jpa["op"])
				envKey := strings.TrimPrefix(fmt.Sprintf("%v", jpa["path"]), "/data/")
				envValue := fmt.Sprintf("%v", jpa["value"])

				if patchActionElem.typ == "env" {
					switch act {
					case "add":
						evs := &corev1.EnvVarSource{}
						if cacheKeySplits[0] == "configmaps" {
							evs.ConfigMapKeyRef = &corev1.ConfigMapKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: cacheKeySplits[2]}, Key: envKey}
						} else {
							evs.SecretKeyRef = &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: cacheKeySplits[2]}, Key: envKey}
						}
						deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{Name: envKey, ValueFrom: evs})
					case "replace":
						log.Debugf("no replace action")
					case "remove":
						for i, envvar := range deploymentObj.Spec.Template.Spec.Containers[0].Env {
							if envvar.Name == envKey {
								deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env[:i], deploymentObj.Spec.Template.Spec.Containers[0].Env[i:]...)
								break
							}
						}
					default:
						return nil, fmt.Errorf("unknown json patch operation. op=%s", act)
					}
				}

				// attach recover action
				if pas, ok := recoverPatchActions[cacheKey]; ok {
					pas.recoverJsonPatchAction = append(pas.recoverJsonPatchAction, jpa)
				} else {
					recoverPatchActions[cacheKey] = &patchAction{typ: patchActionElem.typ, resourceType: patchActionElem.resourceType, recoverJsonPatchAction: []map[string]interface{}{patchActionElem.recoverJsonPatchAction[i]}}
				}
				log.Debugf("changed env in deployment. key=%s, action=%s, updated_value=%s\n", envKey, act, envValue)
			}
			log.Debugf("patched resources for env. resource_type=%s, namespace=%s, name=%s\n", cacheKeySplits[0], cacheKeySplits[1], cacheKeySplits[2])
		}
	}

	return func() error {
		for cacheKey, patchActionElem := range patchActions {
			if cacheKey == "" { // deployment.spec.template.spec.container[0].env
				for _, jpa := range patchActionElem.recoverJsonPatchAction {
					act := fmt.Sprintf("%v", jpa["op"])
					envName := fmt.Sprintf("%v", jpa["path"])
					envValue := fmt.Sprintf("%v", jpa["value"])

					switch act {
					case "add":
						deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{Name: envName, Value: envValue})
					case "replace":
						for i, envvar := range deploymentObj.Spec.Template.Spec.Containers[0].Env {
							if envvar.Name == envName {
								envvar.Value = envValue
								deploymentObj.Spec.Template.Spec.Containers[0].Env[i] = envvar
								break
							}
						}
					case "remove":
						for i, envvar := range deploymentObj.Spec.Template.Spec.Containers[0].Env {
							if envvar.Name == envName {
								deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env[:i], deploymentObj.Spec.Template.Spec.Containers[0].Env[i:]...)
								break
							}
						}
					default:
						return fmt.Errorf("unknown json patch operation. op=%s", act)
					}
					log.Debugf("changed env in deployment. key=%s, action=%s, updated_value=%s\n", envName, act, envValue)
				}
			} else { // configmap or secret patch
				// patchAction
				cacheKeySplits := strings.Split(cacheKey, "|")
				if len(cacheKeySplits) < 3 {
					return fmt.Errorf("cache_key is not valid. cache_key=%s", cacheKey)
				}

				if _, err := k8sClient.ResourcePatch(schema.GroupVersion{Version: "v1"}, cacheKeySplits[0], map[string]interface{}{
					"namespace":  cacheKeySplits[1],
					"name":       cacheKeySplits[2],
					"patch_type": "json",
					"patch_data": patchActionElem.recoverJsonPatchAction,
				}); err != nil {
					return err
				}

				// env.valueFrom -> deployment
				for _, jpa := range patchActionElem.recoverJsonPatchAction {
					act := fmt.Sprintf("%v", jpa["op"])
					envKey := strings.TrimPrefix(fmt.Sprintf("%v", jpa["path"]), "/data/")
					envValue := fmt.Sprintf("%v", jpa["value"])

					if patchActionElem.typ == "env" {
						switch act {
						case "add":
							evs := &corev1.EnvVarSource{}
							if cacheKeySplits[0] == "configmaps" {
								evs.ConfigMapKeyRef = &corev1.ConfigMapKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: cacheKeySplits[2]}, Key: envKey}
							} else {
								evs.SecretKeyRef = &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: cacheKeySplits[2]}, Key: envKey}
							}
							deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{Name: envKey, ValueFrom: evs})
						case "replace":
							log.Debugf("no replace action")
						case "remove":
							for i, envvar := range deploymentObj.Spec.Template.Spec.Containers[0].Env {
								if envvar.Name == envKey {
									deploymentObj.Spec.Template.Spec.Containers[0].Env = append(deploymentObj.Spec.Template.Spec.Containers[0].Env[:i], deploymentObj.Spec.Template.Spec.Containers[0].Env[i:]...)
									break
								}
							}
						default:
							return fmt.Errorf("unknown json patch operation. op=%s", act)
						}
					}
					log.Debugf("changed env in deployment. key=%s, action=%s, updated_value=%s\n", envKey, act, envValue)
				}
				log.Debugf("patched resources for env. resource_type=%s, namespace=%s, name=%s\n", cacheKeySplits[0], cacheKeySplits[1], cacheKeySplits[2])
			}
		}
		return nil
	}, nil
}

type envPath struct {
	typ          string // env, envFrom
	resourceType string // "", configmaps, secrets
	resourceName string // resource(configmap|secret) name
}

type envValue struct {
	typ                   string // env or envFrom
	value                 string // only env type
	resourceCacheKey      string // (configmap or secret)|namespace|name
	resourceDataKey       string // data.key in configmap or secret
	resourceDataKeyPrefix string // data.key prefix
}

func findEnvironmentVariablesFromContainer(k8sClient *k8s.Client, namespace string, container *corev1.Container, generateResourceCacheKey func(resourceType, namespace, name string) string) (*envPath, map[string]*envValue, error) {
	if container == nil {
		return nil, nil, fmt.Errorf("deployment or container is empty")
	}

	resourceCacheMap := make(map[string]interface{}) // key: (configmap or secret)|namespace|name
	defaultEnvPath := envPath{typ: "env"}            // precedence order by desc: secret, configmap, basic

	// precedence : last, env than envFrom
	// 1. envFrom.configMapRef or envFrom.secretRef
	// 2. env.name + env.value
	// 3. env.name + (env.valueFrom.configMapKeyRef or env.valueFrom.secretKeyRef)
	envMap := make(map[string]*envValue) // currently used client environment variables. key: env name, value: env path

	for _, ef := range container.EnvFrom {
		if ef.ConfigMapRef != nil && ef.ConfigMapRef.Name != "" {
			var cm *corev1.ConfigMap
			cmInf, ok := resourceCacheMap[generateResourceCacheKey("configmaps", namespace, ef.ConfigMapRef.Name)]
			if !ok {
				cmJson, err := k8sClient.ResourceGet(schema.GroupVersion{Version: "v1"}, "configmaps", map[string]interface{}{"namespace": namespace, "name": ef.ConfigMapRef.Name})
				if err != nil {
					log.Warnf("failed to get configmap. namespace=%s, name=%s, error=%v\n", namespace, ef.ConfigMapRef.Name, err)
					continue
				}

				configMapObj := &corev1.ConfigMap{}
				if err := json.Unmarshal([]byte(cmJson), configMapObj); err != nil {
					log.Warnf("failed to json unmarshal. namespace=%s, name=%s, error=%v\n", namespace, ef.ConfigMapRef.Name, err)
					continue
				}
				resourceCacheMap[generateResourceCacheKey("configmaps", namespace, ef.ConfigMapRef.Name)] = configMapObj
				cm = configMapObj

				if defaultEnvPath.typ == "env" && defaultEnvPath.resourceType == "" {
					defaultEnvPath = envPath{typ: "envFrom", resourceType: "configmaps", resourceName: ef.ConfigMapRef.Name}
				}
			} else {
				cm, ok = cmInf.(*corev1.ConfigMap)
				if !ok {
					log.Warnf("failed type assertion to *corev1.ConfigMap. got=%T", cmInf)
					continue
				}
			}

			for k, v := range cm.Data {
				envMap[ef.Prefix+k] = &envValue{value: v, typ: "envFrom", resourceCacheKey: generateResourceCacheKey("configmaps", namespace, ef.ConfigMapRef.Name)}
			}
		}

		if ef.SecretRef != nil && ef.SecretRef.Name != "" {
			var secret *corev1.Secret
			secretInf, ok := resourceCacheMap[generateResourceCacheKey("secrets", namespace, ef.SecretRef.Name)]
			if !ok {
				scJson, err := k8sClient.ResourceGet(schema.GroupVersion{Version: "v1"}, "secrets", map[string]interface{}{"namespace": namespace, "name": ef.SecretRef.Name})
				if err != nil {
					log.Warnf("failed to get secret. namespace=%s, name=%s, error=%v\n", namespace, ef.SecretRef.Name, err)
					continue
				}

				secretObj := &corev1.Secret{}
				if err := json.Unmarshal([]byte(scJson), secretObj); err != nil {
					log.Warnf("failed to json unmarshal. namespace=%s, name=%s, error=%v\n", namespace, ef.SecretRef.Name, err)
					continue
				}
				resourceCacheMap[generateResourceCacheKey("secrets", namespace, ef.SecretRef.Name)] = secretObj
				secret = secretObj

				if (defaultEnvPath.typ == "env" && defaultEnvPath.resourceType == "") || (defaultEnvPath.typ == "envFrom" && defaultEnvPath.resourceType == "configmaps") {
					defaultEnvPath = envPath{typ: "envFrom", resourceType: "secrets", resourceName: ef.SecretRef.Name}
				}
			} else {
				secret, ok = secretInf.(*corev1.Secret)
				if !ok {
					log.Warnf("failed type assertion to *corev1.Secret. got=%T", secretInf)
					continue
				}
			}

			for k, v := range secret.Data {
				envMap[ef.Prefix+k] = &envValue{value: string(v), typ: "envFrom", resourceCacheKey: generateResourceCacheKey("secrets", namespace, ef.SecretRef.Name), resourceDataKeyPrefix: ef.Prefix}
			}
		}
	}
	for _, e := range container.Env {
		if e.Value != "" {
			envMap[e.Name] = &envValue{value: e.Value, typ: "env"}
		} else if e.ValueFrom != nil {
			if e.ValueFrom.ConfigMapKeyRef != nil && e.ValueFrom.ConfigMapKeyRef.Name != "" && e.ValueFrom.ConfigMapKeyRef.Key != "" {
				var cm *corev1.ConfigMap
				cmInf, ok := resourceCacheMap[generateResourceCacheKey("configmaps", namespace, e.ValueFrom.ConfigMapKeyRef.Name)]
				if !ok {
					cmJson, err := k8sClient.ResourceGet(schema.GroupVersion{Version: "v1"}, "configmaps", map[string]interface{}{"namespace": namespace, "name": e.ValueFrom.ConfigMapKeyRef.Name})
					if err != nil {
						log.Warnf("failed to get configmap. namespace=%s, name=%s, error=%v\n", namespace, e.ValueFrom.ConfigMapKeyRef.Name, err)
						continue
					}

					configMapObj := &corev1.ConfigMap{}
					if err := json.Unmarshal([]byte(cmJson), configMapObj); err != nil {
						log.Warnf("failed to json unmarshal. namespace=%s, name=%s, error=%v\n", namespace, e.ValueFrom.ConfigMapKeyRef.Name, err)
						continue
					}
					resourceCacheMap[generateResourceCacheKey("configmaps", namespace, e.ValueFrom.ConfigMapKeyRef.Name)] = configMapObj
					cm = configMapObj

					if defaultEnvPath.typ == "env" && defaultEnvPath.resourceType == "" {
						defaultEnvPath = envPath{typ: "env", resourceType: "configmaps", resourceName: e.ValueFrom.ConfigMapKeyRef.Name}
					}
				} else {
					cm, ok = cmInf.(*corev1.ConfigMap)
					if !ok {
						log.Warnf("failed type assertion to *corev1.ConfigMap. got=%T", cmInf)
						continue
					}
				}

				envMap[e.Name] = &envValue{value: cm.Data[e.ValueFrom.ConfigMapKeyRef.Key], typ: "env", resourceCacheKey: generateResourceCacheKey("configmaps", namespace, e.ValueFrom.ConfigMapKeyRef.Name), resourceDataKey: e.ValueFrom.ConfigMapKeyRef.Key}
			}

			if e.ValueFrom.SecretKeyRef != nil && e.ValueFrom.SecretKeyRef.Name != "" && e.ValueFrom.SecretKeyRef.Key != "" {
				var secret *corev1.Secret
				secretInf, ok := resourceCacheMap[generateResourceCacheKey("secrets", namespace, e.ValueFrom.SecretKeyRef.Name)]
				if !ok {
					scJson, err := k8sClient.ResourceGet(schema.GroupVersion{Version: "v1"}, "secrets", map[string]interface{}{"namespace": namespace, "name": e.ValueFrom.SecretKeyRef.Name})
					if err != nil {
						log.Warnf("failed to get secret. namespace=%s, name=%s, error=%v\n", namespace, e.ValueFrom.SecretKeyRef.Name, err)
						continue
					}

					secretObj := &corev1.Secret{}
					if err := json.Unmarshal([]byte(scJson), secretObj); err != nil {
						log.Warnf("failed to json unmarshal. namespace=%s, name=%s, error=%v\n", namespace, e.ValueFrom.SecretKeyRef.Name, err)
						continue
					}
					resourceCacheMap[generateResourceCacheKey("secrets", namespace, e.ValueFrom.SecretKeyRef.Name)] = secretObj
					secret = secretObj

					if defaultEnvPath.typ == "env" && defaultEnvPath.resourceType == "" || (defaultEnvPath.typ == "env" && defaultEnvPath.resourceType == "configmaps") {
						defaultEnvPath = envPath{typ: "env", resourceType: "secrets", resourceName: e.ValueFrom.SecretKeyRef.Name}
					}
				} else {
					secret, ok = secretInf.(*corev1.Secret)
					if !ok {
						log.Warnf("failed type assertion to *corev1.Secret. got=%T", secretInf)
						continue
					}
				}

				envMap[e.Name] = &envValue{value: string(secret.Data[e.ValueFrom.SecretKeyRef.Key]), typ: "env", resourceCacheKey: generateResourceCacheKey("secrets", namespace, e.ValueFrom.SecretKeyRef.Name), resourceDataKey: e.ValueFrom.SecretKeyRef.Key}
			}
		}
	}

	return &defaultEnvPath, envMap, nil
}

type patchAction struct {
	typ                    string // env, envFrom
	resourceType           string // "", configmaps, secrets
	jsonPatchAction        []map[string]interface{}
	recoverJsonPatchAction []map[string]interface{}
}

func extractPatchActionsFromEnvironmentVariables(envMap map[string]*envValue, envConfigs map[string]string, defaultEnvPath envPath, generateResourceCacheKeyWithNamespace func(resourceType, name string) string) map[string]*patchAction {
	patchActions := make(map[string]*patchAction) // key: resourceCacheKey, value: patch actions

	for k, v := range envConfigs {
		vv, ok := envMap[k]
		if ok {
			log.Debugf("found env. will update. key=%s, value=%#v, update_value=%s", k, vv, v)
			if vv.resourceCacheKey != "" {
				if pa, ok := patchActions[vv.resourceCacheKey]; ok {
					keySplit := strings.Split(vv.resourceCacheKey, "|")
					if len(keySplit) <= 0 {
						log.Warnf("resourceCacheKey is not valid. key=%s\n", vv.resourceCacheKey)
						continue
					}

					var jvalue, rvalue string
					if keySplit[0] == "configmaps" {
						jvalue = v
						rvalue = vv.value
					} else {
						jvalue = base64.StdEncoding.EncodeToString([]byte(v))
						rvalue = base64.StdEncoding.EncodeToString([]byte(vv.value))
					}

					pathKey := k
					if vv.resourceDataKey != "" {
						pathKey = vv.resourceDataKey
					}

					pa.jsonPatchAction = append(pa.jsonPatchAction, map[string]interface{}{
						"op":    "replace",
						"path":  "/data/" + pathKey,
						"value": jvalue,
					})
					pa.recoverJsonPatchAction = append(pa.recoverJsonPatchAction, map[string]interface{}{
						"op":    "replace",
						"path":  "/data/" + pathKey,
						"value": rvalue,
					})
				} else {
					keySplit := strings.Split(vv.resourceCacheKey, "|")
					if len(keySplit) <= 0 {
						log.Warnf("resourceCacheKey is not valid. key=%s\n", vv.resourceCacheKey)
						continue
					}

					var jvalue, rvalue string
					if keySplit[0] == "configmaps" {
						jvalue = v
						rvalue = vv.value
					} else {
						jvalue = base64.StdEncoding.EncodeToString([]byte(v))
						rvalue = base64.StdEncoding.EncodeToString([]byte(vv.value))
					}

					pathKey := k
					if vv.resourceDataKey != "" {
						pathKey = vv.resourceDataKey
					}

					if vv.resourceDataKeyPrefix != "" {
						pathKey = strings.TrimPrefix(pathKey, vv.resourceDataKeyPrefix)
					}

					patchActions[vv.resourceCacheKey] = &patchAction{
						typ:          vv.typ,
						resourceType: keySplit[0],
						jsonPatchAction: []map[string]interface{}{{
							"op":    "replace",
							"path":  "/data/" + pathKey,
							"value": jvalue,
						}},
						recoverJsonPatchAction: []map[string]interface{}{{
							"op":    "replace",
							"path":  "/data/" + pathKey,
							"value": rvalue,
						}},
					}
				}
			} else {
				if pa, ok := patchActions[vv.resourceCacheKey]; ok {
					pa.jsonPatchAction = append(pa.jsonPatchAction, map[string]interface{}{
						"op":    "replace",
						"path":  k, // deployment env key
						"value": v,
					})
					pa.recoverJsonPatchAction = append(pa.recoverJsonPatchAction, map[string]interface{}{
						"op":    "replace",
						"path":  k, // deployment env key
						"value": vv.value,
					})
				} else {
					patchActions[vv.resourceCacheKey] = &patchAction{
						typ:          vv.typ,
						resourceType: "",
						jsonPatchAction: []map[string]interface{}{{
							"op":    "replace",
							"path":  k, // deployment env key
							"value": v,
						}},
						recoverJsonPatchAction: []map[string]interface{}{{
							"op":    "replace",
							"path":  k, // deployment env key
							"value": vv.value,
						}},
					}
				}
			}
		} else {
			log.Debugf("not found env. will create. key=%s, update_value=%s", k, v)
			if defaultEnvPath.resourceType != "" {
				cacheKey := generateResourceCacheKeyWithNamespace(defaultEnvPath.resourceType, defaultEnvPath.resourceName)
				if pa, ok := patchActions[cacheKey]; ok {
					var jvalue string
					if defaultEnvPath.resourceType == "configmaps" {
						jvalue = v
					} else {
						jvalue = base64.StdEncoding.EncodeToString([]byte(v))
					}

					pa.jsonPatchAction = append(pa.jsonPatchAction, map[string]interface{}{
						"op":    "add",
						"path":  "/data/" + k,
						"value": jvalue,
					})
					pa.recoverJsonPatchAction = append(pa.recoverJsonPatchAction, map[string]interface{}{
						"op":   "remove",
						"path": "/data/" + k,
					})
				} else {
					keySplit := strings.Split(cacheKey, "|")
					if len(keySplit) <= 0 {
						log.Warnf("resourceCacheKey is not valid. key=%s\n", cacheKey)
						continue
					}

					var jvalue string
					if keySplit[0] == "configmaps" {
						jvalue = v
					} else {
						jvalue = base64.StdEncoding.EncodeToString([]byte(v))
					}

					patchActions[cacheKey] = &patchAction{
						typ:          defaultEnvPath.typ,
						resourceType: keySplit[0],
						jsonPatchAction: []map[string]interface{}{{
							"op":    "add",
							"path":  "/data/" + k,
							"value": jvalue,
						}},
						recoverJsonPatchAction: []map[string]interface{}{{
							"op":   "remove",
							"path": "/data/" + k,
						}},
					}
				}
			} else {
				if pa, ok := patchActions[""]; ok {
					pa.jsonPatchAction = append(pa.jsonPatchAction, map[string]interface{}{
						"op":    "add",
						"path":  k, // deployment env key
						"value": v,
					})
					pa.recoverJsonPatchAction = append(pa.recoverJsonPatchAction, map[string]interface{}{
						"op":   "remove",
						"path": k, // deployment env key
					})
				} else {
					patchActions[""] = &patchAction{
						typ:          defaultEnvPath.typ,
						resourceType: "",
						jsonPatchAction: []map[string]interface{}{{
							"op":    "add",
							"path":  k, // deployment env key
							"value": v,
						}},
						recoverJsonPatchAction: []map[string]interface{}{{
							"op":   "remove",
							"path": k, // deployment env key
						}},
					}
				}
			}
		}
	}

	return patchActions
}
