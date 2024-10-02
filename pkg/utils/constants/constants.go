package constants

const (
	API_VERSION                      = "keysaascontroller.keysaas/v1"
	KEYSAAS_KIND                     = "Keysaas"
	CONTAINER_NAME                   = "keysaas"
	TIMEOUT                          = 600 // 10mins
	KEYCLOAK_DEFAULT_HTTP_PORT       = 8080
	KEYCLOAK_DEFAULT_HTTPS_PORT      = 8443
	KEYCLOAK_DEFAULT_JGROUP_PORT     = 7600
	KEYCLOAK_DEFAULT_MANAGEMENT_PORT = 9000 // for /health and /metrics
	KEYCLOAK_PATH                    = "/"  //"/keycloak/(.*)"
	KEYCLOAK_IMAGE                   = "bitnami/keycloak:25.0.6"
	KEYCLOAK_THEME_PROVIDER_LOCATION = "/opt/bitnami/keycloak/providers"
	KEYCLOAK_CERT_LOCATION           = "/certificate"
	BASE_URL                         = "kubernetes.local"
	WAIT_FOR_POD_TIMEOUT             = 75 //*4seconds = 5 min
	POSTGRES_DATABASE_NAME           = "keysaas"
)
