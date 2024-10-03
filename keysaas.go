package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	operatorv1 "github.com/smugug/keysaas/pkg/apis/keysaascontroller/v1"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/smugug/keysaas/pkg/utils/constants"
	appsv1 "k8s.io/api/apps/v1"
	v2 "k8s.io/api/autoscaling/v2"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiutil "k8s.io/apimachinery/pkg/util/intstr"
)

func (c *KeysaasController) deployKeysaas(ctx context.Context, foo *operatorv1.Keysaas) (string, string, string, error) {
	c.logger.Info("KeysaasController.go : Inside deployKeysaas")
	var keysaasPodName, serviceURIToReturn string
	var err error = nil
	deleteFuncs := []func(foo *operatorv1.Keysaas) error{}
	err = c.createPersistentVolume(foo)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating PV")
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deletePersistentVolume)
	// if foo.Spec.PvcVolumeName == "" {
	err = c.createPersistentVolumeClaim(foo)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating PVC")
		for i := len(deleteFuncs) - 1; i >= 0; i-- {
			deleteFuncs[i](foo)
		}
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deletePersistentVolumeClaim)
	// }
	servicePort, err := c.createService(foo)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating Service")
		for i := len(deleteFuncs) - 1; i >= 0; i-- {
			deleteFuncs[i](foo)
		}
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deleteService)

	err = c.createServiceMonitor(foo)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating Service Monitor")
		for i := len(deleteFuncs) - 1; i >= 0; i-- {
			deleteFuncs[i](foo)
		}
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deleteServiceMonitor)

	err = c.createIngress(foo)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating Ingress")
		for i := len(deleteFuncs) - 1; i >= 0; i-- {
			deleteFuncs[i](foo)
		}
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deleteIngress)

	secretName := ""
	adminPassword := ""
	secretName, adminPassword = c.getSecret(foo)
	if adminPassword == "" {
		adminPassword = c.generatePassword()
		err = c.createSecret(foo, adminPassword)
		if err != nil {
			c.logger.Error(err, "KeysaasController.go : Error creating Secret")
			for i := len(deleteFuncs) - 1; i >= 0; i-- {
				deleteFuncs[i](foo)
			}
			return "", "", "", err
		}
		deleteFuncs = append(deleteFuncs, c.deleteSecret)
	}

	keysaasPodName, err = c.createDeployment(foo, adminPassword)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating Deployment")
		for i := len(deleteFuncs) - 1; i >= 0; i-- {
			deleteFuncs[i](foo)
		}
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deleteDeployment)

	err = c.createHorizontalPodAutoscaler(foo)
	if err != nil {
		c.logger.Error(err, "KeysaasController.go : Error creating Ingress")
		for i := len(deleteFuncs) - 1; i >= 0; i-- {
			deleteFuncs[i](foo)
		}
		return "", "", "", err
	}
	deleteFuncs = append(deleteFuncs, c.deleteHorizontalPodAutoscaler)

	// Wait couple of seconds more just to give the Pod some more time.
	time.Sleep(time.Second * 2)

	// plugins := foo.Spec.Plugins

	// supportedPlugins, unsupportedPlugins = c.util.GetSupportedPlugins(plugins)
	// if len(supportedPlugins) > 0 {
	// 	namespace := getNamespace(foo)
	// 	erredPlugins = c.util.EnsurePluginsInstalled(foo, supportedPlugins, keysaasPodName, namespace, constants.PLUGIN_MAP)
	// }
	// if len(erredPlugins) > 0 {
	// 	err = errors.New("Error Installing Supported Plugin")
	// }

	if foo.Spec.DomainName == "" {
		serviceURIToReturn = foo.Name + "." + constants.BASE_URL
	} else {
		serviceURIToReturn = foo.Spec.DomainName + ":" + servicePort
	}

	c.logger.Info("KeysaasController.go : KeysaasController.go: Returning from deployKeysaas")

	return serviceURIToReturn, keysaasPodName, secretName, err
}

func (c *KeysaasController) generatePassword() string {
	mina := 97
	maxa := 122
	minA := 65
	maxA := 90
	min0 := 48
	max0 := 57
	length := 8

	password := make([]string, length)

	i := 0
	for i < length {
		charSet := rand.Intn(3)
		if charSet == 0 {
			passwordInt := rand.Intn(maxa-mina) + mina
			password[i] = string(passwordInt)
		}
		if charSet == 1 {
			passwordInt := rand.Intn(maxA-minA) + minA
			password[i] = string(passwordInt)
		}
		if charSet == 2 {
			passwordInt := rand.Intn(max0-min0) + min0
			password[i] = string(passwordInt)
		}
		i++
	}
	passwordString := strings.Join(password, "")
	c.logger.Info("KeysaasController.go : KeysaasController.go  : Generated Password", "password", passwordString)

	return passwordString
}

func getNamespace(foo *operatorv1.Keysaas) string {
	namespace := apiv1.NamespaceDefault
	if foo.Namespace != "" {
		namespace = foo.Namespace
	}
	return namespace
}

func (c *KeysaasController) createIngress(foo *operatorv1.Keysaas) error {

	keysaasName := foo.Name

	keysaasTLSCertSecretName := ""
	tls := foo.Spec.Tls

	c.logger.Info("KeysaasController.go : TLS", "tls", tls)
	if len(tls) > 0 {
		keysaasTLSCertSecretName = keysaasName + "-domain-cert"
	}

	keysaasPath := constants.KEYCLOAK_PATH

	keysaasDomainName := getDomainName(foo)
	if keysaasDomainName == "" {
		// keysaasPath = keysaasPath + keysaasName ????
		// if no domain then use sub-domain
		keysaasDomainName = keysaasName + "." + constants.BASE_URL
	}

	keysaasServiceName := keysaasName
	keysaasPort := constants.KEYCLOAK_DEFAULT_HTTP_PORT

	specObj := getIngressSpec(keysaasPort, keysaasDomainName, keysaasPath, keysaasTLSCertSecretName, keysaasServiceName, tls)

	ingress := getIngress(foo, specObj, keysaasName, tls)

	namespace := getNamespace(foo)
	ingressesClient := c.kubeclientset.NetworkingV1().Ingresses(namespace)

	c.logger.Info("KeysaasController.go : Creating Ingress...")
	result, err := ingressesClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : Ingress already exists", "ingress name", ingress.GetObjectMeta().GetName())
	} else if err != nil {
		c.logger.Error(err, "KeysaasController.go : ")
		return err
	}
	c.logger.Info("KeysaasController.go : Created Ingress.", "ingress name", result.GetObjectMeta().GetName())
	return nil
}

func (c *KeysaasController) deleteIngress(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.kubeclientset.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func getIngress(foo *operatorv1.Keysaas, specObj networkingv1.IngressSpec, keysaasName, tls string) *networkingv1.Ingress {

	var ingress *networkingv1.Ingress

	if len(tls) > 0 {
		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: keysaasName,
				Annotations: map[string]string{
					"cert-manager.io/cluster-issuer": "letsencrypt-staging",
					/// passthrough so that haproxy won't validate stuff beforehand since we already have tls on keycloak
					// "haproxy.org/ssl-passthrough": "false",
					// "spec.ingressClassName": "haproxy",4
					// "certmanager.k8s.io/issuer":              keysaasName,
					// "certmanager.k8s.io/acme-challenge-type": "http01",
				},
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: constants.API_VERSION,
						Kind:       constants.KEYSAAS_KIND,
						Name:       foo.Name,
						UID:        foo.UID,
					},
				},
				Labels: map[string]string{
					"app": keysaasName,
				},
			},
			Spec: specObj,
		}
	} else {
		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        keysaasName,
				Annotations: map[string]string{
					// "spec.ingressClassName": "haproxy",
					// "nginx.ingress.kubernetes.io/ssl-redirect":   "false",
					// "nginx.ingress.kubernetes.io/rewrite-target": "/",
					// "nginx.ingress.kubernetes.io/rewrite-target": "/",
				},
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: constants.API_VERSION,
						Kind:       constants.KEYSAAS_KIND,
						Name:       foo.Name,
						UID:        foo.UID,
					},
				},
				Labels: map[string]string{
					"app": keysaasName,
				},
			},
			Spec: specObj,
		}
	}
	return ingress
}

func getIngressSpec(keysaasPort int, keysaasDomainName, keysaasPath, keysaasTLSCertSecretName,
	keysaasServiceName, tls string) networkingv1.IngressSpec {

	var specObj networkingv1.IngressSpec
	pathType := networkingv1.PathTypePrefix
	ingressClass := "traefik"
	if len(tls) > 0 {
		specObj = networkingv1.IngressSpec{
			IngressClassName: &ingressClass,
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{keysaasDomainName},
					SecretName: keysaasDomainName, //autogenerated by cert-manager
				},
			},
			Rules: []networkingv1.IngressRule{
				{
					Host: keysaasDomainName,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: keysaasPath,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: keysaasServiceName,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(constants.KEYCLOAK_DEFAULT_HTTP_PORT), //int32(constants.KEYCLOAK_DEFAULT_HTTPS_PORT), //port of service
											},
										},
										//apiutil.FromInt(keysaasPort),
									},
									PathType: &pathType,
								},
							},
						},
					},
				},
			},
		}
	} else {
		specObj = networkingv1.IngressSpec{
			IngressClassName: &ingressClass,
			Rules: []networkingv1.IngressRule{
				{
					Host: keysaasDomainName,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: keysaasPath,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: keysaasServiceName,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(constants.KEYCLOAK_DEFAULT_HTTP_PORT),
											},
										},
										//apiutil.FromInt(keysaasPort),
									},
									PathType: &pathType,
								},
							},
						},
					},
				},
			},
		}
	}

	return specObj
}

func (c *KeysaasController) createPersistentVolume(foo *operatorv1.Keysaas) error {
	c.logger.Info("KeysaasController.go : Inside createPersistentVolume")
	deploymentName := foo.Name
	hostPathType := apiv1.HostPathDirectoryOrCreate
	persistentVolume := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: constants.API_VERSION,
					Kind:       constants.KEYSAAS_KIND,
					Name:       foo.Name,
					UID:        foo.UID,
				},
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			StorageClassName: "standard",
			Capacity: apiv1.ResourceList{
				"storage": resource.MustParse("50Mi"),
			},
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			// MINIKUBE ONLY SUPPORTS CERTAIN HOSTPATHS. MAKE SURE TO CHECK IT OUT
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: "/tmp/hostpath_pv/themes/" + deploymentName,
					Type: &hostPathType,
				},
			},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimDelete,
		},
	}
	c.logger.Info("/tmp/hostpath_pv/themes/" + deploymentName)
	persistentVolumeClient := c.kubeclientset.CoreV1().PersistentVolumes()

	c.logger.Info("KeysaasController.go : Creating persistentVolume...")
	result, err := persistentVolumeClient.Create(context.TODO(), persistentVolume, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : PersistentVolume already exists", "pv name", persistentVolume.GetObjectMeta().GetName())
	} else if err != nil {
		c.logger.Error(err, "KeysaasController.go : ")
		return err
	}
	c.logger.Info("KeysaasController.go : Created persistentVolume", "pv name", result.GetObjectMeta().GetName())
	return nil
}
func (c *KeysaasController) deletePersistentVolume(foo *operatorv1.Keysaas) error {
	return c.kubeclientset.CoreV1().PersistentVolumes().Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func (c *KeysaasController) createPersistentVolumeClaim(foo *operatorv1.Keysaas) error {
	c.logger.Info("KeysaasController.go : Inside createPersistentVolumeClaim")

	/// STANDARD HAS DYNAMIC PROVISION (AUTO CREATING PV), MANUAL DOESN'T
	standardStorageClassName := "standard"
	deploymentName := foo.Name
	pvcTheme := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: constants.API_VERSION,
					Kind:       constants.KEYSAAS_KIND,
					Name:       foo.Name,
					UID:        foo.UID,
				},
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			StorageClassName: &standardStorageClassName,
			Resources: apiv1.VolumeResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse("50Mi"),
				},
			},
		},
	}
	namespace := getNamespace(foo)
	persistentVolumeClaimClient := c.kubeclientset.CoreV1().PersistentVolumeClaims(namespace)

	c.logger.Info("KeysaasController.go : Creating persistentVolumeClaim for Themes...")
	_, err := persistentVolumeClaimClient.Create(context.TODO(), pvcTheme, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : persistentVolumeClaim for theme already exists", "pv name", pvcTheme.GetObjectMeta().GetName())
	} else if err != nil {
		return err
	}
	c.logger.Info("KeysaasController.go : Created persistentVolumeClaims")
	return nil
}

func (c *KeysaasController) deletePersistentVolumeClaim(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.kubeclientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func (c *KeysaasController) createDeployment(foo *operatorv1.Keysaas, adminPassword string) (string, error) {
	c.logger.Info("KeysaasController.go : Inside createDeployment")

	namespace := getNamespace(foo)
	deploymentsClient := c.kubeclientset.AppsV1().Deployments(namespace)

	deploymentName := foo.Name

	// keysaasPort := KEYSAAS_PORT

	//image := "lmecld/nginxforkeysaas:9.0"
	//image := "lmecld/nginxforkeysaas6:latest"
	//	image = "lmecld/nginxforkeysaas10:latest"

	// claimName := foo.Spec.PvcVolumeName
	// if claimName == "" {
	claimName := foo.Name

	//MySQL Service IP and Port
	postgresqlServiceName := "keysaas-postgresql" //foo.Spec.MySQLServiceName
	c.logger.Info("KeysaasController.go : Postgreql Service name", "service name", postgresqlServiceName)
	postgresConfigMapName := "keysaas-postgresql-secret"
	postgresConfigMap, err := c.kubeclientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), postgresConfigMapName, metav1.GetOptions{})
	if err != nil {
		// c.logger.Error(err, "KeysaasController.go : Error getting MYSQL details")
		return "", err
	}
	postgresqlUserName, exists := postgresConfigMap.Data["POSTGRES_USER"]
	if !exists {
		// Handle the case where the key does not exist
		return "", errors.New("postgres_user doesn't exist in configmap")
	}
	postgresqlUserPassword, exists := postgresConfigMap.Data["POSTGRES_PASSWORD"]
	if !exists {
		// Handle the case where the key does not exist
		return "", errors.New("postgres_password doesn't exist in configmap")
	}

	c.logger.Info("KeysaasController.go : postgresql Password", "postgresql password", postgresqlUserPassword)

	postgresqlServiceClient := c.kubeclientset.CoreV1().Services(namespace)
	postgresqlServiceResult, err := postgresqlServiceClient.Get(context.TODO(), postgresqlServiceName, metav1.GetOptions{})
	if err != nil {
		// c.logger.Error(err, "KeysaasController.go : Error getting MYSQL details")
		return "", err
	}
	postgresDatabaseName := constants.POSTGRES_DATABASE_NAME
	postgresqlHostIP := postgresqlServiceName
	postgresqlHostNumericIP := postgresqlServiceResult.Spec.ClusterIP
	postgresqlServicePortInt := postgresqlServiceResult.Spec.Ports[0].Port
	c.logger.Info("KeysaasController.go : postgresql Service Port int", "port in", postgresqlServicePortInt)
	postgresqlServicePort := fmt.Sprint(postgresqlServicePortInt)
	c.logger.Info("KeysaasController.go : postgresql Service Port", "service port", postgresqlServicePort)
	c.logger.Info("KeysaasController.go : postgresql Host IP", "host ip", postgresqlHostIP)

	err = c.createDatabaseSchema(namespace, postgresqlServiceName, postgresqlHostNumericIP, postgresqlServicePort, postgresDatabaseName, postgresqlUserName, postgresqlUserPassword, deploymentName)
	if err != nil {
		// c.logger.Error(err, "KeysaasController.go : Error creating database schema")
		return "", err
	}

	HOST_NAME := ""
	if foo.Spec.DomainName == "" {
		HOST_NAME = deploymentName + "." + constants.BASE_URL
	} else {
		HOST_NAME = foo.Spec.DomainName
	}

	c.logger.Info("KeysaasController.go : HOST and PORT", "HOST_NAME", HOST_NAME, "PORT")

	defaultIfEmpty(&foo.Spec.LimitsCpu, "500m")
	defaultIfEmpty(&foo.Spec.LimitsMemory, "512Mi")
	/// WARNING
	/// KEYCLOAK PRODUCTION (start) REQUIRES HTTPS/TLS AND HOST NAME BY DEFAULT
	/// KEYCLOAK TESTING (start-dev)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: constants.API_VERSION,
					Kind:       constants.KEYSAAS_KIND,
					Name:       foo.Name,
					UID:        foo.UID,
				},
			},
			Labels: map[string]string{
				"app": deploymentName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: apiv1.PodSpec{
					//COPY DEFAULT THEMES
					// InitContainers: []apiv1.Container{
					// 	{
					// 		Name:    "init-themes",
					// 		Image:   constants.KEYCLOAK_IMAGE,
					// 		Command: []string{"/bin/sh", "-c", fmt.Sprintf("cp -r %s/* /mnt/themes/", constants.KEYCLOAK_THEME_LOCATION)},
					// 		VolumeMounts: []apiv1.VolumeMount{
					// 			{
					// 				Name:      volumeName,
					// 				MountPath: "/mnt/themes",
					// 			},
					// 		},
					// 	},
					// },
					Containers: []apiv1.Container{
						{
							Name:  constants.CONTAINER_NAME,
							Image: constants.KEYCLOAK_IMAGE,
							// Lifecycle: &apiv1.Lifecycle{
							// 	PostStart: &apiv1.LifecycleHandler{
							// 		Exec: &apiv1.ExecAction{
							// 			Command: []string{"echo meow"},
							// 			// Command: []string{"/bin/sh", "-c", "/usr/local/scripts/keysaasinstall.sh; sleep 5; /usr/sbin/nginx -s reload"},
							// 			//Command: []string{"/bin/sh", "-c", "/usr/local/scripts/keysaasinstall.sh"},
							// 		},
							// 	},
							// },
							/*
								ReadinessProbe: &apiv1.Probe{
									Handler: apiv1.Handler{
										TCPSocket: &apiv1.TCPSocketAction{
											Port: apiutil.FromInt(80),
										},
									},
									InitialDelaySeconds: 5,
									TimeoutSeconds:      60,
									PeriodSeconds:       2,
								},*/
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse(foo.Spec.LimitsCpu),
									apiv1.ResourceMemory: resource.MustParse(foo.Spec.LimitsMemory),
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "KEYCLOAK_PRODUCTION",
									Value: "true",
								},
								{
									Name:  "KEYCLOAK_DATABASE_VENDOR",
									Value: "postgresql",
								},
								{
									Name:  "KEYCLOAK_DATABASE_NAME",
									Value: postgresDatabaseName,
								},
								{
									Name:  "KEYCLOAK_DATABASE_USER",
									Value: postgresqlUserName,
								},
								{
									Name:  "KEYCLOAK_DATABASE_PASSWORD",
									Value: postgresqlUserPassword,
								},
								{
									Name:  "KEYCLOAK_DATABASE_HOST",
									Value: postgresqlHostIP,
								},
								{
									Name:  "KEYCLOAK_DATABASE_PORT",
									Value: postgresqlServicePort,
								},
								{
									Name:  "KEYCLOAK_DATABASE_SCHEMA",
									Value: deploymentName,
								},
								{
									Name:  "KEYCLOAK_ADMIN",
									Value: foo.Spec.KeysaasUsername,
								},
								{
									Name:  "KEYCLOAK_ADMIN_PASSWORD",
									Value: foo.Spec.KeysaasPassword, //adminPassword,
								},
								{
									Name:  "KEYCLOAK_ENABLE_HTTPS",
									Value: "false",
								},
								// needed since it's behind an ingress controller
								{
									Name:  "KC_PROXY_HEADERS",
									Value: "xforwarded",
								},
								// {
								// 	Name:  "KEYCLOAK_HTTPS_PORT",
								// 	Value: strconv.Itoa(constants.KEYCLOAK_DEFAULT_HTTPS_PORT),
								// },
								// {
								// 	Name:  "KEYCLOAK_HTTPS_USE_PEM",
								// 	Value: "true",
								// },
								// {
								// 	Name:  "KEYCLOAK_HTTPS_CERTIFICATE_FILE",
								// 	Value: constants.KEYCLOAK_CERT_LOCATION + "/tls.crt",
								// },
								// {
								// 	Name:  "KEYCLOAK_HTTPS_CERTIFICATE_KEY_FILE",
								// 	Value: constants.KEYCLOAK_CERT_LOCATION + "/tls.key",
								// },
								{
									Name:  "KC_METRICS_ENABLED",
									Value: "true",
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: int32(constants.KEYCLOAK_DEFAULT_HTTP_PORT),
								},
								{
									ContainerPort: int32(constants.KEYCLOAK_DEFAULT_MANAGEMENT_PORT),
								},
								// {
								// 	ContainerPort: int32(constants.KEYCLOAK_DEFAULT_HTTPS_PORT),
								// },
								// {
								// 	ContainerPort: int32(constants.KEYCLOAK_DEFAULT_JGROUP_PORT),
								// },
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "theme-volume",
									MountPath: constants.KEYCLOAK_THEME_PROVIDER_LOCATION,
								},
								// {
								// 	Name:      "cert-volume",
								// 	MountPath: constants.KEYCLOAK_CERT_LOCATION,
								// 	ReadOnly:  true,
								// },
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "theme-volume",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: claimName,
								},
							},
						},
						// {
						// 	Name: "cert-volume",
						// 	VolumeSource: apiv1.VolumeSource{
						// 		Secret: &apiv1.SecretVolumeSource{
						// 			SecretName: "kubern	etes-tls",
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}

	// Create Deployment
	c.logger.Info("KeysaasController.go : Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : deployment already exists", "deployment name", deployment.GetObjectMeta().GetName())
	} else if err != nil {
		// c.logger.Error(err, "KeysaasController.go : ")
		return "", err
	}
	c.logger.Info("KeysaasController.go : Created deployment", "deployment name", result.GetObjectMeta().GetName())

	/*
		podname, _ := c.util.GetPodFullName(constants.TIMEOUT, foo.Name, foo.Namespace)
		keysaasPodName, podReady := c.util.WaitForPod(constants.TIMEOUT, podname, foo.Namespace)
	*/

	keysaasPodName, podReady := c.waitForPod(foo)

	if podReady {
		return keysaasPodName, nil
	} else {
		err := errors.New("keysaas pod timeout")
		return keysaasPodName, err
	}
}

func (c *KeysaasController) deleteDeployment(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.kubeclientset.AppsV1().Deployments(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func (c *KeysaasController) createHorizontalPodAutoscaler(foo *operatorv1.Keysaas) error {
	c.logger.Info("KeysaasController.go : Inside createHorizontalPodAutoscaler")
	deploymentName := foo.Name
	namespace := getNamespace(foo)
	hpaClient := c.kubeclientset.AutoscalingV2().HorizontalPodAutoscalers(namespace)
	hpa := &v2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace,
		},
		Spec: v2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       deploymentName,
			},
			MinReplicas: int32Ptr(1),
			MaxReplicas: 2,
			Metrics: []v2.MetricSpec{
				v2.MetricSpec{
					Type: v2.ResourceMetricSourceType,
					Resource: &v2.ResourceMetricSource{
						Name: apiv1.ResourceCPU,
						Target: v2.MetricTarget{
							Type:               v2.UtilizationMetricType,
							AverageUtilization: int32Ptr(90),
						},
					},
				},
				v2.MetricSpec{
					Type: v2.ResourceMetricSourceType,
					Resource: &v2.ResourceMetricSource{
						Name: apiv1.ResourceMemory,
						Target: v2.MetricTarget{
							Type:               v2.UtilizationMetricType,
							AverageUtilization: int32Ptr(90),
						},
					},
				},
			},
		},
	}
	result, err := hpaClient.Create(context.TODO(), hpa, metav1.CreateOptions{})
	if err != nil {
		// c.logger.Error(err, "KeysaasController.go : ")
		return err
	}
	c.logger.Info("KeysaasController.go : Created HorizontalPodAutoscaler", "HPA name", result.GetObjectMeta().GetName())
	return nil
}

func (c *KeysaasController) deleteHorizontalPodAutoscaler(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.kubeclientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
	return nil
}

func (c *KeysaasController) getSecret(foo *operatorv1.Keysaas) (string, string) {
	c.logger.Info("KeysaasController.go : Inside getSecret")
	secretName := foo.Name

	namespace := getNamespace(foo)
	secretsClient := c.kubeclientset.CoreV1().Secrets(namespace)

	c.logger.Info("KeysaasController.go : Getting secrets..")
	result, err := secretsClient.Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		c.logger.Info("KeysaasController.go : error getting secret", "error", err)
		//panic(err)
	}
	if result != nil {
		c.logger.Info("KeysaasController.go : Getting Secret", "result", result.GetObjectMeta().GetName())

		adminPasswordByteArray := result.Data["adminPassword"]
		adminPassword := string(adminPasswordByteArray)

		c.logger.Info("KeysaasController.go : Admin Password", "admin password", adminPassword)

		return secretName, adminPassword

	} else {
		return "", ""
	}
}

func (c *KeysaasController) createSecret(foo *operatorv1.Keysaas, adminPassword string) error {

	c.logger.Info("KeysaasController.go : Inside createSecret")
	secretName := foo.Name

	c.logger.Info("KeysaasController.go : Secret Name", "secret name", secretName)
	c.logger.Info("KeysaasController.go : Admin Password", "admin password", adminPassword)

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: constants.API_VERSION,
					Kind:       constants.KEYSAAS_KIND,
					Name:       foo.Name,
					UID:        foo.UID,
				},
			},
			Labels: map[string]string{
				"secret": secretName,
			},
		},
		Data: map[string][]byte{
			"adminPassword": []byte(adminPassword),
		},
	}

	namespace := getNamespace(foo)
	secretsClient := c.kubeclientset.CoreV1().Secrets(namespace)

	c.logger.Info("KeysaasController.go : Creating secrets..")
	result, err := secretsClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : secret already exists", "secret name", secret.GetObjectMeta().GetName())
	} else if err != nil {
		// c.logger.Error(err, "KeysaasController.go : ")
		return err
	}
	c.logger.Info("KeysaasController.go : Created Secret", "secret name", result.GetObjectMeta().GetName())
	return nil
}

func (c *KeysaasController) deleteSecret(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.kubeclientset.CoreV1().Secrets(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func (c *KeysaasController) createService(foo *operatorv1.Keysaas) (string, error) {

	c.logger.Info("KeysaasController.go : Inside createService")
	deploymentName := foo.Name
	namespace := getNamespace(foo)
	serviceClient := c.kubeclientset.CoreV1().Services(namespace)

	serviceObj, servicePort := getServiceSpec(deploymentName, foo.Spec.DomainName)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: constants.API_VERSION,
					Kind:       constants.KEYSAAS_KIND,
					Name:       deploymentName,
					UID:        foo.UID,
				},
			},
			Annotations: map[string]string{
				"prometheus.io/port":   "metrics",
				"prometheus.io/scrape": "true",
			},
			Labels: map[string]string{
				"app": deploymentName,
			},
		},
		Spec: serviceObj,
	}

	result1, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : service already exists", "service name", service.GetObjectMeta().GetName())
	} else if err != nil {
		// c.logger.Error(err, "KeysaasController.go : ")
		return "", err
	}
	c.logger.Info("KeysaasController.go : Created service", "service name", result1.GetObjectMeta().GetName())

	//nodePort1 := result1.Spec.Ports[0].NodePort
	//nodePort := fmt.Sprint(nodePort1)
	//servicePort := fmt.Sprint(keysaasPort)

	// Parse ServiceIP and Port
	serviceIP := result1.Spec.ClusterIP
	c.logger.Info("KeysaasController.go : Keysaas Service IP", "keysaas ip", serviceIP)

	//servicePortInt := result1.Spec.Ports[0].Port
	//servicePort := fmt.Sprint(servicePortInt)

	serviceURI := serviceIP + ":" + servicePort

	c.logger.Info("KeysaasController.go : Service URI", "service uri", serviceURI)

	return servicePort, nil
}

func (c *KeysaasController) deleteService(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.kubeclientset.CoreV1().Services(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func (c *KeysaasController) createServiceMonitor(foo *operatorv1.Keysaas) error {
	c.logger.Info("KeysaasController.go : Inside createServiceMonitor")
	deploymentName := foo.Name
	namespace := getNamespace(foo)
	serviceMonitorClient := c.monitoringclientset.MonitoringV1().ServiceMonitors(namespace)

	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: constants.API_VERSION,
					Kind:       constants.KEYSAAS_KIND,
					Name:       deploymentName,
					UID:        foo.UID,
				},
			},
			Labels: map[string]string{
				"app":        deploymentName,
				"prometheus": "system-monitoring-prometheus",
			},
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				monitoringv1.Endpoint{
					Port:     "management",
					Path:     "/metrics",
					Interval: monitoringv1.Duration("30s"),
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{
					"customer2",
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
		},
	}
	result1, err := serviceMonitorClient.Create(context.TODO(), serviceMonitor, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : servicemonitor already exists", "servicemonitor name", serviceMonitor.GetObjectMeta().GetName())
	} else if err != nil {
		// c.logger.Error(err, "KeysaasController.go : ")
		return err
	}
	c.logger.Info("KeysaasController.go : Created servicemonitor", "servicemonitor name", result1.GetObjectMeta().GetName())
	return nil
}

func (c *KeysaasController) deleteServiceMonitor(foo *operatorv1.Keysaas) error {
	namespace := getNamespace(foo)
	return c.monitoringclientset.MonitoringV1().ServiceMonitors(namespace).Delete(context.TODO(), foo.Name, *metav1.NewDeleteOptions(0))
}

func getDomainName(foo *operatorv1.Keysaas) string {
	return foo.Spec.DomainName
}

func getServiceSpec(deploymentName, domainName string) (apiv1.ServiceSpec, string) {

	var serviceObj apiv1.ServiceSpec

	var servicePort string

	if domainName == "" {
		serviceObj = apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       constants.KEYCLOAK_DEFAULT_HTTP_PORT,                  //internally exposed port
					TargetPort: apiutil.FromInt(constants.KEYCLOAK_DEFAULT_HTTP_PORT), //port on pod
					// NodePort:   int32(keysaasPort),                                    //externally exposed port
					Protocol: apiv1.ProtocolTCP,
				},
				{
					Name:       "management",
					Port:       constants.KEYCLOAK_DEFAULT_MANAGEMENT_PORT,                  //internally exposed port
					TargetPort: apiutil.FromInt(constants.KEYCLOAK_DEFAULT_MANAGEMENT_PORT), //port on pod
					// NodePort:   int32(keysaasPort),                                      //externally exposed port
					// Protocol: apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": deploymentName,
			},
			//Type: apiv1.ServiceTypeNodePort,
			Type: apiv1.ServiceTypeClusterIP,
			//Type: apiv1.ServiceTypeLoadBalancer,
		}
		servicePort = strconv.Itoa(constants.KEYCLOAK_DEFAULT_HTTPS_PORT)
	} else {
		serviceObj = apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       constants.KEYCLOAK_DEFAULT_HTTP_PORT,                  //internally exposed port
					TargetPort: apiutil.FromInt(constants.KEYCLOAK_DEFAULT_HTTP_PORT), //port on pod
					Protocol:   apiv1.ProtocolTCP,
				},
				{
					Name:       "https",
					Port:       constants.KEYCLOAK_DEFAULT_HTTPS_PORT,                  //internally exposed port
					TargetPort: apiutil.FromInt(constants.KEYCLOAK_DEFAULT_HTTPS_PORT), //port on pod
					Protocol:   apiv1.ProtocolTCP,
				},
				{
					Name:       "jgroup",
					Port:       constants.KEYCLOAK_DEFAULT_JGROUP_PORT,                  //internally exposed port
					TargetPort: apiutil.FromInt(constants.KEYCLOAK_DEFAULT_JGROUP_PORT), //port on pod
					// Protocol: apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": deploymentName,
			},
			//Type: apiv1.ServiceTypeNodePort,
			Type: apiv1.ServiceTypeClusterIP,
			//Type: apiv1.ServiceTypeLoadBalancer,
		}
		servicePort = strconv.Itoa(constants.KEYCLOAK_DEFAULT_HTTP_PORT)
	}
	return serviceObj, servicePort
}

func (c *KeysaasController) getDiff(leftHandSide, rightHandSide []string) []string {
	var diffList []string
	for _, inspec := range leftHandSide {
		var found = false
		for _, installed := range rightHandSide {
			if inspec == installed {
				found = true
				break
			}
		}
		if !found {
			diffList = append(diffList, inspec)
		}
	}
	return diffList
}

func (c *KeysaasController) isInitialDeployment(foo *operatorv1.Keysaas) bool {
	if foo.Status.Url == "" {
		return true
	} else {
		return false
	}
}
func (c *KeysaasController) waitForPod(foo *operatorv1.Keysaas) (string, bool) {
	var podName string
	deploymentName := foo.Name
	namespace := getNamespace(foo)
	// Check if Postgres Pod is ready or not
	podReady := false
	podTimeoutCount := 0
	TIMEOUT_COUNT := constants.WAIT_FOR_POD_TIMEOUT // 75*4(sleep time)=300=5 minutes; this should be made configurable
	for {
		pods := c.getPods(namespace, deploymentName)
		for _, d := range pods.Items {
			//my-hello-5fb5bb554-8l22r sp
			parts := strings.Split(d.Name, "-")
			parts = parts[:len(parts)-2]
			podDepName := strings.Join(parts, "-")
			if podDepName == deploymentName {
				podName = d.Name
				c.logger.Info("Keysaas Pod Name", "pod name", podName)
				///to test
				podConditions := d.Status.Conditions
				for _, podCond := range podConditions {
					if podCond.Type == apiv1.PodReady {
						if podCond.Status == apiv1.ConditionTrue {
							c.logger.Info("Keysaas Pod is running.")
							podReady = true
							break
						}
					}
				}
			}
			if podReady {
				break
			}
		}
		if podReady {
			break
		} else {
			c.logger.Info("Waiting for Keysaas Pod to get ready.")
			time.Sleep(time.Second * 4)
			podTimeoutCount = podTimeoutCount + 1
			if podTimeoutCount >= TIMEOUT_COUNT {
				podReady = false
				break
			}
		}
	}
	if podReady {
		c.logger.Info("Pod is ready.")
	} else {
		c.logger.Info("Pod timeout")
	}
	return podName, podReady
}

func (c *KeysaasController) getPods(namespace, deploymentName string) *apiv1.PodList {
	// TODO(devkulkarni): This is returning all Pods. We should change this
	// to only return Pods whose Label matches the Deployment Name.
	deploymentName = "app=" + deploymentName
	pods, err := c.kubeclientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: deploymentName,
	})
	c.logger.Info("Number of pods in a cluster", "number", len(pods.Items))
	if err != nil {
		c.logger.Info("Error getting pod list", "error", err)
	}
	return pods
}

func int32Ptr(i int32) *int32 { return &i }

func (c *KeysaasController) createDatabaseSchema(namespace, deploymentName, dbHost, dbPort, dbName, dbUser, dbPass, schemaName string) error {
	// // TODO: THIS IS FOR TESTING ONLY, MUST USE THE SECOND ONE WHEN DEPLOYING INSIDE THE CLUSTER
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", dbUser, dbPass, "192.168.49.2", "30500")
	// //connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	// c.logger.Info("KeysaasController.go : Database to connect " + connStr)
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	// // Create the schema
	// command := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schemaName)
	// result, err := db.Exec(command)
	// c.logger.Info("UWUWUUUUUUUUUUUUUUUUUUUUUUUUU", "result", result, "error", err)
	// if err != nil {
	// 	return err
	// }
	deploymentName = "keysaas-postgresql"
	pods := c.getPods(namespace, deploymentName)
	for _, d := range pods.Items {
		//my-hello-5fb5bb554-8l22r sp
		parts := strings.Split(d.Name, "-")
		parts = parts[:len(parts)-2]
		podDepName := strings.Join(parts, "-")
		if podDepName == deploymentName {
			// command := fmt.Sprintf(`psql -U %s -c "CREATE SCHEMA IF NOT EXISTS %s;" -d %s`, dbUser, schemaName, dbName)
			command2 := []string{"psql", "-U", dbUser, "-c", "DROP SCHEMA IF EXISTS " + schemaName + " CASCADE; CREATE SCHEMA " + schemaName + ";", "-d", dbName}
			success := c.util.ExecuteExecCall(d.Name, "postgresql", namespace, command2)
			c.logger.Info("KeysaasController.go : Try to create new schema at pod", "pod", d.Name, "success", success)
			return nil
		}
	}
	return errors.New("couldn't create a new schema")
}

func defaultIfEmpty(value *string, defaultValue string) {
	if *value == "" {
		*value = defaultValue
	}
}
