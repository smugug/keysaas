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

	"github.com/smugug/keysaas/pkg/utils/constants"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiutil "k8s.io/apimachinery/pkg/util/intstr"
)

const TEST_PORT = 8000

func (c *KeysaasController) deployKeysaas(ctx context.Context, foo *operatorv1.Keysaas) (string, string, string, error) {
	c.logger.Info("KeysaasController.go : Inside deployKeysaas")
	var keysaasPodName, serviceURIToReturn string

	// c.createPersistentVolume(foo)
	if foo.Spec.PvcVolumeName == "" {
		c.createPersistentVolumeClaim(foo)
	}
	servicePort := c.createService(foo)

	c.createIngress(foo)

	err, keysaasPodName, secretName := c.createDeployment(foo)

	if err != nil {
		return serviceURIToReturn, keysaasPodName, secretName, err
	}

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
		serviceURIToReturn = foo.Name + ":" + servicePort
	} else {
		serviceURIToReturn = foo.Spec.DomainName + ":" + servicePort
	}

	c.logger.Info("KeysaasController.go : KeysaasController.go: Returning from deployKeysaas")

	return serviceURIToReturn, keysaasPodName, secretName, err
}

func (c *KeysaasController) generatePassword(keysaasPort int) string {
	seed := keysaasPort
	rand.Seed(int64(seed))
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

func (c *KeysaasController) createIngress(foo *operatorv1.Keysaas) {

	keysaasName := foo.Name

	keysaasTLSCertSecretName := ""
	tls := foo.Spec.Tls

	c.logger.Info("KeysaasController.go : TLS", "tls", tls)
	if len(tls) > 0 {
		keysaasTLSCertSecretName = keysaasName + "-domain-cert"
	}

	keysaasPath := "/"

	keysaasDomainName := getDomainName(foo)
	if keysaasDomainName == "" {
		keysaasPath = keysaasPath + keysaasName
	}

	keysaasServiceName := keysaasName
	keysaasPort := KEYSAAS_PORT

	specObj := getIngressSpec(keysaasPort, keysaasDomainName, keysaasPath,
		keysaasTLSCertSecretName, keysaasServiceName, tls)

	ingress := getIngress(foo, specObj, keysaasName, tls)

	namespace := getNamespace(foo)
	ingressesClient := c.kubeclientset.NetworkingV1().Ingresses(namespace)

	c.logger.Info("KeysaasController.go : Creating Ingress...")
	result, err := ingressesClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : Ingress already exists", "ingress name", ingress.GetObjectMeta().GetName())
	} else if err != nil {
		panic(err)
	}
	c.logger.Info("KeysaasController.go : Created Ingress.", "ingress name", result.GetObjectMeta().GetName())
}

func getIngress(foo *operatorv1.Keysaas, specObj networkingv1.IngressSpec, keysaasName, tls string) *networkingv1.Ingress {

	var ingress *networkingv1.Ingress

	if len(tls) > 0 {
		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: keysaasName,
				Annotations: map[string]string{
					"spec.ingressClassName": "haproxy",
					// "nginx.ingress.kubernetes.io/rewrite-target": "/",
					"certmanager.k8s.io/issuer":              keysaasName,
					"certmanager.k8s.io/acme-challenge-type": "http01",
				},
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: constants.API_VERSION,
						Kind:       constants.KEYSAAS_KIND,
						Name:       foo.Name,
						UID:        foo.UID,
					},
				},
			},
			Spec: specObj,
		}
	} else {
		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: keysaasName,
				Annotations: map[string]string{
					// "spec.ingressClassName":                      "nginx",
					// "nginx.ingress.kubernetes.io/ssl-redirect":   "false",
					// "nginx.ingress.kubernetes.io/rewrite-target": "/",
					"spec.ingressClassName": "haproxy",
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
			},
			Spec: specObj,
			/*
					Spec: extensionsv1beta1.IngressSpec{
						TLS: []extensionsv1beta1.IngressTLS{
							{
								Hosts: []string{keysaasDomainName},
								SecretName: keysaasTLSCertSecretName,
							},
						},
						Rules: []extensionsv1beta1.IngressRule{
							{
								Host: keysaasDomainName,
								IngressRuleValue: extensionsv1beta1.IngressRuleValue{
									HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
										Paths: []extensionsv1beta1.HTTPIngressPath{
											{
												Path: keysaasPath,
												Backend: extensionsv1beta1.IngressBackend{
												ServiceName: keysaasServiceName,
												ServicePort: apiutil.FromInt(keysaasPort),
											},
										},
									},
								},
							},
						},
					},
				},
			*/
		}
	}
	return ingress
}

func getIngressSpec(keysaasPort int, keysaasDomainName, keysaasPath, keysaasTLSCertSecretName,
	keysaasServiceName, tls string) networkingv1.IngressSpec {

	var specObj networkingv1.IngressSpec
	pathType := networkingv1.PathTypePrefix
	fmt.Println("CUNNNNNNNNNNYYYYYYYYYYYYYYYY HAPROXYYYYYYYYYYYYy")
	ingressClass := "haproxy"
	if len(tls) > 0 {
		specObj = networkingv1.IngressSpec{
			IngressClassName: &ingressClass,
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{keysaasDomainName},
					SecretName: keysaasTLSCertSecretName,
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
												Number: 80, //port of service
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
												Number: 80,
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

func (c *KeysaasController) createPersistentVolume(foo *operatorv1.Keysaas) {
	c.logger.Info("KeysaasController.go : Inside createPersistentVolume")

	deploymentName := foo.Name
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
			StorageClassName: "manual",
			Capacity: apiv1.ResourceList{
				//					map[string]resource.Quantity{
				"storage": resource.MustParse("1Gi"),
				//					},
			},
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				//					{
				"ReadWriteOnce",
				//					},
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: "/mnt/keysaas-data",
				},
			},
		},
	}

	persistentVolumeClient := c.kubeclientset.CoreV1().PersistentVolumes()

	c.logger.Info("KeysaasController.go : Creating persistentVolume...")
	result, err := persistentVolumeClient.Create(context.TODO(), persistentVolume, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : PersistentVolume already exists", "pv name", persistentVolume.GetObjectMeta().GetName())
	} else if err != nil {
		panic(err)
	}
	c.logger.Info("KeysaasController.go : Created persistentVolume", "pv name", result.GetObjectMeta().GetName())
}

func (c *KeysaasController) createPersistentVolumeClaim(foo *operatorv1.Keysaas) {
	c.logger.Info("KeysaasController.go : Inside createPersistentVolumeClaim")

	storageClassName := "standard"
	deploymentName := foo.Name
	persistentVolumeClaim := &apiv1.PersistentVolumeClaim{
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
				//					{
				"ReadWriteOnce",
				//					},
			},
			StorageClassName: &storageClassName,
			Resources: apiv1.VolumeResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse("1Gi"),
					//							map[string]resource.Quantity{
					//							"storage": resource.MustParse("1Gi"),
					//						},
				},
			},
		},
	}

	namespace := getNamespace(foo)
	persistentVolumeClaimClient := c.kubeclientset.CoreV1().PersistentVolumeClaims(namespace)

	c.logger.Info("KeysaasController.go : Creating persistentVolumeClaim...")
	result, err := persistentVolumeClaimClient.Create(context.TODO(), persistentVolumeClaim, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : persistentVolumeClaim already exists", "pv name", persistentVolumeClaim.GetObjectMeta().GetName())
	} else if err != nil {
		panic(err)
	}
	c.logger.Info("KeysaasController.go : Created persistentVolumeClaim", "pvc name", result.GetObjectMeta().GetName())
}

func (c *KeysaasController) createDeployment(foo *operatorv1.Keysaas) (error, string, string) {

	c.logger.Info("KeysaasController.go : Inside createDeployment")

	namespace := getNamespace(foo)
	deploymentsClient := c.kubeclientset.AppsV1().Deployments(namespace)

	deploymentName := foo.Name

	keysaasPort := KEYSAAS_PORT

	image := "crccheck/hello-world:latest"
	//image := "lmecld/nginxforkeysaas:9.0"
	//image := "lmecld/nginxforkeysaas6:latest"
	//	image = "lmecld/nginxforkeysaas10:latest"

	volumeName := "keysaas-data"

	claimName := foo.Spec.PvcVolumeName
	if claimName == "" {
		claimName = foo.Name
	}

	secretName := ""
	adminPassword := ""
	secretName, adminPassword = c.getSecret(foo)
	if adminPassword == "" {
		adminPassword = c.generatePassword(KEYSAAS_PORT)
		secretName = c.createSecret(foo, adminPassword)
	}

	//MySQL Service IP and Port
	mysqlServiceName := "keysaas2-mysql" //foo.Spec.MySQLServiceName
	c.logger.Info("KeysaasController.go : MySQL Service name", "service name", mysqlServiceName)

	mysqlUserName := "user1" //foo.Spec.MySQLUserName
	c.logger.Info("KeysaasController.go : MySQL Username", "mysql username", mysqlUserName)

	// passwordLocation := "mysql-secret.mysql-password" //foo.Spec.MySQLUserPassword
	// secretPasswordSplitted := strings.Split(passwordLocation, ".")
	// mysqlSecretName := secretPasswordSplitted[0]
	// mysqlPasswordField := secretPasswordSplitted[1]

	// secretsClient := c.kubeclientset.CoreV1().Secrets(namespace)
	// secret, err := secretsClient.Get(context.TODO(), mysqlSecretName, metav1.GetOptions{})

	// if err != nil {
	// 	c.logger.Info("KeysaasController.go : Error, secret not found from in namespace", "secret", mysqlSecretName, "namespace", namespace, "error", err)
	// }
	// if _, ok := secret.Data[mysqlPasswordField]; !ok {
	// 	c.logger.Info("KeysaasController.go : Error, secret does not have the field", "secret", mysqlSecretName, "password field", mysqlPasswordField)
	// }
	mysqlUserPassword := "password1" //string(secret.Data[mysqlPasswordField])

	c.logger.Info("KeysaasController.go : MySQL Password", "mysql password", mysqlUserPassword)

	keysaasAdminEmail := foo.Spec.KeysaasAdminEmail
	c.logger.Info("KeysaasController.go : Keysaas Admin Email", "mysql email", keysaasAdminEmail)

	mysqlServiceClient := c.kubeclientset.CoreV1().Services(namespace)
	mysqlServiceResult, err := mysqlServiceClient.Get(context.TODO(), mysqlServiceName, metav1.GetOptions{})

	if err != nil {
		c.logger.Info("KeysaasController.go : Error getting MySQL Service details", "error", err)
		return err, "", secretName
	}

	mysqlHostIP := mysqlServiceName
	mysqlServicePortInt := mysqlServiceResult.Spec.Ports[0].Port
	c.logger.Info("KeysaasController.go : MySQL Service Port int", "port in", mysqlServicePortInt)
	mysqlServicePort := fmt.Sprint(mysqlServicePortInt)
	c.logger.Info("KeysaasController.go : MySQL Service Port", "service port", mysqlServicePort)
	c.logger.Info("KeysaasController.go : MySQL Host IP", "host ip", mysqlHostIP)

	HOST_NAME := ""
	if foo.Spec.DomainName == "" {
		HOST_NAME = deploymentName + ":" + strconv.Itoa(KEYSAAS_PORT)
	} else {
		HOST_NAME = foo.Spec.DomainName
	}

	c.logger.Info("KeysaasController.go : HOST_NAME", "HOST_NAME", HOST_NAME)

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
					Containers: []apiv1.Container{
						{
							Name:  constants.CONTAINER_NAME,
							Image: image,
							// Lifecycle: &apiv1.Lifecycle{
							// 	PostStart: &apiv1.LifecycleHandler{
							// 		Exec: &apiv1.ExecAction{
							// 			Command: []string{"echo meow"},
							// 			// Command: []string{"/bin/sh", "-c", "/usr/local/scripts/keysaasinstall.sh; sleep 5; /usr/sbin/nginx -s reload"},
							// 			//Command: []string{"/bin/sh", "-c", "/usr/local/scripts/keysaasinstall.sh"},
							// 		},
							// 	},
							// },
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: int32(TEST_PORT),
								},
							},
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
							Env: []apiv1.EnvVar{
								{
									Name:  "APPLICATION_NAME",
									Value: deploymentName,
								},
								{
									Name:  "MYSQL_DATABASE",
									Value: "keysaas",
								},
								{
									Name:  "MYSQL_USER",
									Value: mysqlUserName,
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: mysqlUserPassword,
								},
								{
									Name:  "MYSQL_HOST",
									Value: mysqlHostIP,
									/*ValueFrom: &apiv1.EnvVarSource{
									  FieldRef: &apiv1.ObjectFieldSelector{
									      FieldPath: "status.hostIP",
									  },
									},*/
								},
								{
									Name:  "MYSQL_PORT",
									Value: mysqlServicePort,
								},
								{
									Name:  "MYSQL_TABLE_PREFIX",
									Value: "mdl_",
								},
								{
									Name:  "KEYSAAS_ADMIN_PASSWORD",
									Value: adminPassword,
								},
								{
									Name:  "KEYSAAS_ADMIN_EMAIL",
									Value: keysaasAdminEmail,
								},
								{
									Name:  "KEYSAAS_PORT",
									Value: strconv.Itoa(keysaasPort),
								},
								{
									Name:  "HOST_NAME",
									Value: HOST_NAME,
									/*ValueFrom: &apiv1.EnvVarSource{
									  FieldRef: &apiv1.ObjectFieldSelector{
									      FieldPath: "status.hostIP",
									  },
									},*/
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/opt/keysaasdata",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: volumeName,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: claimName,
								},
							},
						},
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
		panic(err)
	}
	c.logger.Info("KeysaasController.go : Created deployment", "deployment name", result.GetObjectMeta().GetName())

	/*
		podname, _ := c.util.GetPodFullName(constants.TIMEOUT, foo.Name, foo.Namespace)
		keysaasPodName, podReady := c.util.WaitForPod(constants.TIMEOUT, podname, foo.Namespace)
	*/

	keysaasPodName, podReady := c.waitForPod(foo)

	if podReady {
		return nil, keysaasPodName, secretName
	} else {
		err1 := errors.New("Keysaas Pod Timeout")
		return err1, keysaasPodName, secretName
	}
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

func (c *KeysaasController) createSecret(foo *operatorv1.Keysaas, adminPassword string) string {

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
		panic(err)
	}
	c.logger.Info("KeysaasController.go : Created Secret", "secret name", result.GetObjectMeta().GetName())
	return secretName
}

func (c *KeysaasController) createService(foo *operatorv1.Keysaas) string {

	c.logger.Info("KeysaasController.go : Inside createService")
	deploymentName := foo.Name
	keysaasPort := KEYSAAS_PORT

	namespace := getNamespace(foo)
	serviceClient := c.kubeclientset.CoreV1().Services(namespace)

	serviceObj, servicePort := getServiceSpec(keysaasPort, deploymentName, foo.Spec.DomainName)
	service := &apiv1.Service{
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
		Spec: serviceObj,

		/*
			Spec: apiv1.ServiceSpec{
				Ports: []apiv1.ServicePort{
					{
						Name:       "my-port",
						Port:       int32(keysaasPort),
						TargetPort: apiutil.FromInt(keysaasPort),
						//NodePort:   int32(KEYSAAS_PORT),
						Protocol:   apiv1.ProtocolTCP,
					},
				},
				Selector: map[string]string{
					"app": deploymentName,
				},
				//Type: apiv1.ServiceTypeNodePort,
				Type: apiv1.ServiceTypeClusterIP,
				//Type: apiv1.ServiceTypeLoadBalancer,
			},*/
	}

	result1, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if kerrors.IsAlreadyExists(err) {
		c.logger.Info("KeysaasController.go : service already exists", "service name", service.GetObjectMeta().GetName())
	} else if err != nil {
		panic(err)
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

	return servicePort
}

func getDomainName(foo *operatorv1.Keysaas) string {
	return foo.Spec.DomainName

	/*
		if len(foo.Spec.DomainName) > 0 {
			return foo.Spec.DomainName
		} else {

			return foo.Name
		}
	*/
}

func getServiceSpec(keysaasPort int, deploymentName, domainName string) (apiv1.ServiceSpec, string) {

	var serviceObj apiv1.ServiceSpec

	var servicePort string

	if domainName == "" {
		serviceObj = apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "my-port",
					Port:       80,                         //internally exposed port
					TargetPort: apiutil.FromInt(TEST_PORT), //port on pod
					NodePort:   int32(keysaasPort),         //externally exposed port
					Protocol:   apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": deploymentName,
			},
			Type: apiv1.ServiceTypeNodePort,
			//Type: apiv1.ServiceTypeClusterIP,
			//Type: apiv1.ServiceTypeLoadBalancer,
		}
		servicePort = strconv.Itoa(keysaasPort)
	} else {
		serviceObj = apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name: "my-port",
					//Port:       int32(keysaasPort),
					Port:       80,
					TargetPort: apiutil.FromInt(TEST_PORT),
					Protocol:   apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": deploymentName,
			},
			//Type: apiv1.ServiceTypeNodePort,
			Type: apiv1.ServiceTypeClusterIP,
			//Type: apiv1.ServiceTypeLoadBalancer,
		}
		servicePort = strconv.Itoa(80)
	}
	return serviceObj, servicePort
}

// func (c *KeysaasController) handlePluginDeployment(keysaas *operatorv1.Keysaas) (string, []string, []string, []string) {

// 	installedPlugins := keysaas.Status.InstalledPlugins
// 	specPlugins := keysaas.Spec.Plugins
// 	unsupportedPlugins := keysaas.Status.UnsupportedPlugins

// 	c.logger.Info("KeysaasController.go : Spec Plugins", "spec plugins", specPlugins)
// 	c.logger.Info("KeysaasController.go : Installed Plugins", "installed plugins", installedPlugins)
// 	var addList []string
// 	var removeList []string

// 	// addList = specList - installedList - unsupportedPlugins
// 	addList = c.getDiff(specPlugins, installedPlugins)
// 	c.logger.Info("KeysaasController.go : Plugins to install", "plugins to install", addList)

// 	if unsupportedPlugins != nil {
// 		addList = c.getDiff(addList, unsupportedPlugins)
// 	}

// 	// removeList = installedList - specList
// 	removeList = c.getDiff(installedPlugins, specPlugins)
// 	c.logger.Info("KeysaasController.go : Plugins to remove", "plugins to remove", removeList)

// 	var podName string
// 	var supportedPlugins, unsupportedPlugins1 []string
// 	supportedPlugins, unsupportedPlugins1 = c.util.GetSupportedPlugins(addList)

// 	var erredPlugins []string
// 	if len(supportedPlugins) > 0 {
// 		podName = keysaas.Status.PodName
// 		namespace := getNamespace(keysaas)
// 		//podname, _ := c.util.GetPodFullName(constants.TIMEOUT, podName, namespace)
// 		erredPlugins = c.util.EnsurePluginsInstalled(keysaas, supportedPlugins, podName, namespace, constants.PLUGIN_MAP)
// 	}
// 	if len(removeList) > 0 {
// 		c.logger.Info("KeysaasController.go : ============= Plugin removal not implemented yet ===============")
// 	}

// 	/*
// 	   if len(supportedPlugins) > 0 || len(removeList) > 0 {
// 	   	return podName, supportedPlugins, unsupportedPlugins
// 	   } else {
// 	      return podName, supportedPlugins, unsupportedPlugins
// 	   }*/

// 	return podName, supportedPlugins, unsupportedPlugins1, erredPlugins
// }

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
	TIMEOUT_COUNT := 150 // 150*4(sleep time)=600=10 minutes; this should be made configurable
	for {
		pods := c.getPods(namespace, deploymentName)
		for _, d := range pods.Items {
			//my-hello-5fb5bb554-8l22r sp
			parts := strings.Split(d.Name, "-")
			parts = parts[:len(parts)-2]
			podDepName := strings.Join(parts, "")
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
	pods, err := c.kubeclientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		//LabelSelector: deploymentName,
		//LabelSelector: metav1.LabelSelector{
		//	MatchLabels: map[string]string{
		//	"app": deploymentName,
		//},
		//},
	})
	c.logger.Info("Number of pods in a cluster", "number", len(pods.Items))
	if err != nil {
		c.logger.Info("Error getting pod list", "error", err)
	}
	return pods
}

func int32Ptr(i int32) *int32 { return &i }
