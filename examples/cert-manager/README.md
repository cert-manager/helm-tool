# Cert Manager

## Parameters

### Global

|property|description|type|default|
|--|--|--|--|
|`global.imagePullSecrets`|<p>Reference to one or more secrets to be used when pulling images<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/</p>|`array`|<pre>[]</pre>|
|`global.commonLabels`|<p>Labels to apply to all resources<br> Please note that this does not add labels to the resources created dynamically by the controllers.<br> For these resources, you have to add the labels in the template in the cert-manager custom resource:</p><p> eg. podTemplate/ ingressTemplate in ACMEChallengeSolverHTTP01Ingress<br>    Ref: https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEChallengeSolverHTTP01Ingress<br> eg. secretTemplate in CertificateSpec<br>    Ref: https://cert-manager.io/docs/reference/api-docs/#cert-manager.io/v1.CertificateSpec</p>|`object`|<pre>{}</pre>|
|`global.priorityClassName`|<p>Optional priority class to be used for the cert-manager pods</p>|`string`|<pre>""</pre>|
|`global.rbac.create`|<p>Create RBAC rules</p>|`bool`|<pre>true</pre>|
|`global.rbac.aggregateClusterRoles`|<p>Aggregate ClusterRoles to Kubernetes default user-facing roles. Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles</p>|`bool`|<pre>true</pre>|
|`global.logLevel`|<p>Set the verbosity of cert-manager. Range of 0 - 6 with 6 being the most verbose.</p>|`number`|<pre>2</pre>|
|`global.leaderElection.namespace`|<p>Override the namespace used for the leader election lease</p>|`string`|<pre>kube-system</pre>|
|`global.leaderElection.leaseDuration`|<p>The duration that non-leader candidates will wait after observing a<br> leadership renewal until attempting to acquire leadership of a led but<br> unrenewed leader slot. This is effectively the maximum duration that a<br> leader can be stopped before it is replaced by another candidate.</p>|`string`|<pre>undefined</pre>|
|`global.leaderElection.renewDeadline`|<p>The interval between attempts by the acting master to renew a leadership<br> slot before it stops leading. This must be less than or equal to the<br> lease duration.</p>|`string`|<pre>undefined</pre>|
|`global.leaderElection.retryPeriod`|<p>The duration the clients should wait between attempting acquisition and<br> renewal of a leadership.</p>|`string`|<pre>undefined</pre>|
|`installCRDs`|<p>Install the CRDs</p>|`bool`|<pre>false</pre>|

### Controller

|property|description|type|default|
|--|--|--|--|
|`replicaCount`|<p>Number of replicas to run of the cert-manager controller</p>|`number`|<pre>1</pre>|
|`strategy`|<p>Update strategy to use, for example:</p><pre>type: RollingUpdate<br>rollingUpdate:<br>  maxSurge: 0<br>  maxUnavailable: 1</pre>|`object`|<pre>{}</pre>|
|`podDisruptionBudget.minAvailable`|<p>minAvailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)</p>|`number`|<pre>undefined</pre>|
|`podDisruptionBudget.maxUnavailable`|<p>maxUnavailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)</p>|`number`|<pre>undefined</pre>|
|`featureGates`|<p>Comma separated list of feature gates that should be enabled on the<br> controller pod.</p>|`string`|<pre>""</pre>|
|`maxConcurrentChallenges`|<p>The maximum number of challenges that can be scheduled as 'processing' at once</p>|`number`|<pre>60</pre>|
|`image.registry`|<p>Registry to pull the image from</p>|`string`|<pre>undefined</pre>|
|`image.repository`|<p>Image name, this can be the full image including registry or the short name<br> excluding the registry. The registy can also be set in the `registry` property</p>|`string`|<pre>quay.io/jetstack/cert-manager-controller</pre>|
|`image.tag`|<p>Override the image tag to deploy by setting this variable.<br> If no value is set, the chart's appVersion will be used.</p>|`string`|<pre>undefined</pre>|
|`image.digest`|<p>Setting a digest will override any tag</p>|`string`|<pre>undefined</pre>|
|`image.pullPolicy`|<p>Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy</p>|`string`|<pre>IfNotPresent</pre>|
|`clusterResourceNamespace`|<p>Override the namespace used to store DNS provider credentials etc. for ClusterIssuer<br> resources. By default, the same namespace as cert-manager is deployed within is<br> used. This namespace will not be automatically created by the Helm chart.</p>|`string`|<pre>""</pre>|
|`namespace`|<p>This namespace allows you to define where the services will be installed into<br> if not set then they will use the namespace of the release<br> This is helpful when installing cert manager as a chart dependency (sub chart)</p>|`string`|<pre>""</pre>|
|`serviceAccount.create`|<p>Specifies whether a service account should be created</p>|`bool`|<pre>true</pre>|
|`serviceAccount.name`|<p>The name of the service account to use.<br> If not set and create is true, a name is generated using the fullname template</p>|`string`|<pre>undefined</pre>|
|`serviceAccount.annotations`|<p>Optional additional annotations to add to the controller's ServiceAccount</p>|`object`|<pre>undefined</pre>|
|`serviceAccount.labels`|<p>Automount API credentials for a Service Account.<br> Optional additional labels to add to the controller's ServiceAccount</p>|`object`|<pre>undefined</pre>|
|`serviceAccount.automountServiceAccountToken`|<p>Service account token wil be automatically mounted in Pods</p>|`bool`|<pre>true</pre>|
|`automountServiceAccountToken`|<p>Automounting API credentials for a particular pod</p>|`bool`|<pre>undefined</pre>|
|`enableCertificateOwnerRef`|<p>When this flag is enabled, secrets will be automatically removed when the certificate resource is deleted</p>|`bool`|<pre>false</pre>|
|`config`|<p>Used to configure options for the controller pod.<br> This allows setting options that'd usually be provided via flags.<br> An APIVersion and Kind must be specified in your values.yaml file.<br> Flags will override options that are set here.</p><p> For example:</p><pre>apiVersion: controller.config.cert-manager.io/v1alpha1<br>kind: ControllerConfiguration<br>logging:<br>  verbosity: 2<br>  format: text<br>leaderElectionConfig:<br>  namespace: kube-system<br>kubernetesAPIQPS: 9000<br>kubernetesAPIBurst: 9000<br>numberOfConcurrentWorkers: 200<br>featureGates:<br>  AdditionalCertificateOutputFormats: true<br>  DisallowInsecureCSRUsageDefinition: true<br>  ExperimentalCertificateSigningRequestControllers: true<br>  ExperimentalGatewayAPISupport: true<br>  LiteralCertificateSubject: true<br>  SecretsFilteredCaching: true<br>  ServerSideApply: true<br>  StableCertificateRequestName: true<br>  UseCertificateRequestBasicConstraints: true<br>  ValidateCAA: true</pre>|`object`|<pre>{}</pre>|
|`dns01RecursiveNameservers`|<p>Comma separated string with host and port of the recursive nameservers cert-manager should query</p>|`string`|<pre>""</pre>|
|`dns01RecursiveNameserversOnly`|<p>Forces cert-manager to only use the recursive nameservers for verification.<br> Enabling this option could cause the DNS01 self check to take longer due to caching performed by the recursive nameservers</p>|`bool`|<pre>false</pre>|
|`extraArgs`|<p>Additional command line flags to pass to cert-manager controller binary.<br> To see all available flags run docker run quay.io/jetstack/cert-manager-controller:<version> --help</p>|`array`|<pre>[]</pre>|
|`extraEnv`|<p>Additional environment variables</p>|`array`|<pre>[]</pre>|
|`resources`|<p>Resources the controller will be given</p>|`object`|<pre>{}</pre>|
|`securityContext`|<p>Pod Security Context<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>runAsNonRoot: true<br>seccompProfile:<br>  type: RuntimeDefault</pre>|
|`containerSecurityContext`|<p>Container Security Context to be set on the controller component container<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>allowPrivilegeEscalation: false<br>capabilities:<br>  drop:<br>    - ALL<br>readOnlyRootFilesystem: true</pre>|
|`volumes`|<p>Volumes to mount to the controller pod</p>|`array`|<pre>[]</pre>|
|`volumeMounts`|<p>Volumes specified in `volumes` to mount to the controller container</p>|`array`|<pre>[]</pre>|
|`deploymentAnnotations`|<p>Optional additional annotations to add to the controller Deployment</p>|`object`|<pre>undefined</pre>|
|`podAnnotations`|<p>Optional additional annotations to add to the controller Pods</p>|`object`|<pre>undefined</pre>|
|`podLabels`|<p>Optional additional labels to add to the controller Pods</p>|`object`|<pre>{}</pre>|
|`serviceAnnotations`|<p>Optional annotations to add to the controller Service</p>|`object`|<pre>undefined</pre>|
|`serviceLabels`|<p>Optional additional labels to add to the controller Service</p>|`object`|<pre>undefined</pre>|
|`podDnsPolicy`|<p>DNS policy to use within the controller pod</p>|`string`|<pre>undefined</pre>|
|`podDnsConfig`|<p>Optional DNS settings, useful if you have a public and private DNS zone for<br> the same domain on Route 53. What follows is an example of ensuring<br> cert-manager can access an ingress or DNS TXT records at all times.<br> NOTE: This requires Kubernetes 1.10 or `CustomPodDNS` feature gate enabled for<br> the cluster to work.</p>|`object`|<pre>undefined</pre>|
|`nodeSelector`|<p>Node selector to limit the nodes the controller can schedule on</p>|`object`|<pre>kubernetes.io/os: linux</pre>|
|`ingressShim.defaultIssuerName`|<p>Optional default issuer to use for ingress resources</p>|`string`|<pre>undefined</pre>|
|`ingressShim.defaultIssuerKind`|<p>Optional default issuer kind to use for ingress resources</p>|`string`|<pre>undefined</pre>|
|`ingressShim.defaultIssuerGroup`|<p>Optional default issuer group to use for ingress resources</p>|`string`|<pre>undefined</pre>|
|`http_proxy`||`string`|<pre>undefined</pre>|
|`https_proxy`||`string`|<pre>undefined</pre>|
|`no_proxy`||`string`|<pre>undefined</pre>|
|`affinity`|<p>A Kubernetes Affinity, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#affinity-v1-core for example:</p><pre>affinity:<br>  nodeAffinity:<br>   requiredDuringSchedulingIgnoredDuringExecution:<br>     nodeSelectorTerms:<br>     - matchExpressions:<br>       - key: foo.bar.com/role<br>         operator: In<br>         values:<br>         - master</pre>|`object`|<pre>{}</pre>|
|`tolerations`|<p>A list of Kubernetes Tolerations, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#toleration-v1-core for example:</p><pre>tolerations:<br>- key: foo.bar.com/role<br>  operator: Equal<br>  value: master<br>  effect: NoSchedule</pre>|`array`|<pre>[]</pre>|
|`topologySpreadConstraints`|<p>A list of Kubernetes TopologySpreadConstraints, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#topologyspreadconstraint-v1-core for example:</p><pre>topologySpreadConstraints:<br>- maxSkew: 2<br>  topologyKey: topology.kubernetes.io/zone<br>  whenUnsatisfiable: ScheduleAnyway<br>  labelSelector:<br>    matchLabels:<br>      app.kubernetes.io/instance: cert-manager<br>      app.kubernetes.io/component: controller</pre>|`array`|<pre>[]</pre>|
|`livenessProbe.enabled`||`bool`|<pre>true</pre>|
|`livenessProbe.initialDelaySeconds`||`number`|<pre>10</pre>|
|`livenessProbe.periodSeconds`||`number`|<pre>10</pre>|
|`livenessProbe.timeoutSeconds`||`number`|<pre>15</pre>|
|`livenessProbe.successThreshold`||`number`|<pre>1</pre>|
|`livenessProbe.failureThreshold`||`number`|<pre>8</pre>|
|`enableServiceLinks`|<p>enableServiceLinks indicates whether information about services should be<br> injected into pod's environment variables, matching the syntax of Docker<br> links.</p>|`bool`|<pre>false</pre>|

### Prometheus

|property|description|type|default|
|--|--|--|--|
|`prometheus.enabled`||`bool`|<pre>true</pre>|
|`prometheus.servicemonitor.enabled`|<p>Create a ServiceMonitor resource to scrape the metrics endpoint</p>|`bool`|<pre>false</pre>|
|`prometheus.servicemonitor.prometheusInstance`||`string`|<pre>default</pre>|
|`prometheus.servicemonitor.targetPort`|<p>The port to scrape metrics from</p>|`number`|<pre>9402</pre>|
|`prometheus.servicemonitor.path`|<p>Path to scrape metrics from</p>|`string`|<pre>/metrics</pre>|
|`prometheus.servicemonitor.interval`|<p>Interval to scrape metrics</p>|`string`|<pre>60s</pre>|
|`prometheus.servicemonitor.scrapeTimeout`|<p>Timeout for each metrics scrape</p>|`string`|<pre>30s</pre>|
|`prometheus.servicemonitor.labels`|<p>Labels to add to the ServiceMonitor resource</p>|`object`|<pre>{}</pre>|
|`prometheus.servicemonitor.annotations`|<p>Annotations to add to the ServiceMonitor resource</p>|`object`|<pre>{}</pre>|
|`prometheus.servicemonitor.honorLabels`||`bool`|<pre>false</pre>|
|`prometheus.servicemonitor.endpointAdditionalProperties`||`object`|<pre>{}</pre>|
|`prometheus.podmonitor.enabled`|<p>Create a PodMonitor resource to scrape the metrics endpoint</p>|`bool`|<pre>false</pre>|
|`prometheus.podmonitor.prometheusInstance`||`string`|<pre>default</pre>|
|`prometheus.podmonitor.path`|<p>Path to scrape metrics from</p>|`string`|<pre>/metrics</pre>|
|`prometheus.podmonitor.interval`|<p>Interval to scrape metrics</p>|`string`|<pre>60s</pre>|
|`prometheus.podmonitor.scrapeTimeout`|<p>Timeout for each metrics scrape</p>|`string`|<pre>30s</pre>|
|`prometheus.podmonitor.labels`|<p>Labels to add to the PodMonitor resource</p>|`object`|<pre>{}</pre>|
|`prometheus.podmonitor.annotations`|<p>Annotations to add to the PodMonitor resource</p>|`object`|<pre>{}</pre>|
|`prometheus.podmonitor.honorLabels`||`bool`|<pre>false</pre>|
|`prometheus.podmonitor.endpointAdditionalProperties`||`object`|<pre>{}</pre>|

### Webhook

|property|description|type|default|
|--|--|--|--|
|`webhook.replicaCount`||`number`|<pre>1</pre>|
|`webhook.timeoutSeconds`|<p>Seconds the API server should wait for the webhook to respond before treating the call as a failure. Value must be between 1 and 30 seconds. See:<br> https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/validating-webhook-configuration-v1/</p><p> We set the default to the maximum value of 30 seconds. Here's why:<br> Users sometimes report that the connection between the K8S API server and<br> the cert-manager webhook server times out.<br> If *this* timeout is reached, the error message will be "context deadline exceeded",<br> which doesn't help the user diagnose what phase of the HTTPS connection timed out.<br> For example, it could be during DNS resolution, TCP connection, TLS<br> negotiation, HTTP negotiation, or slow HTTP response from the webhook<br> server.<br> So by setting this timeout to its maximum value the underlying timeout error<br> message has more chance of being returned to the end user.</p>|`number`|<pre>30</pre>|
|`webhook.config`|<p>Used to configure options for the webhook pod.<br> This allows setting options that'd usually be provided via flags.<br> An APIVersion and Kind must be specified in your values.yaml file.<br> Flags will override options that are set here. Example config:</p><pre>apiVersion: webhook.config.cert-manager.io/v1alpha1<br>kind: WebhookConfiguration<br># The port that the webhook should listen on for requests.<br># In GKE private clusters, by default kubernetes apiservers are allowed to<br># talk to the cluster nodes only on 443 and 10250. so configuring<br># securePort: 10250, will work out of the box without needing to add firewall<br># rules or requiring NET_BIND_SERVICE capabilities to bind port numbers <1000.<br># This should be uncommented and set as a default by the chart once we graduate<br># the apiVersion of WebhookConfiguration past v1alpha1.<br>securePort: 10250</pre>|`object`|<pre>{}</pre>|
|`webhook.strategy`|<p>Deployment strategy, for example:</p><pre>type: RollingUpdate<br>rollingUpdate:<br>  maxSurge: 0<br>  maxUnavailable: 1</pre>|`object`|<pre>{}</pre>|
|`webhook.securityContext`|<p>Pod Security Context to be set on the webhook component Pod<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>runAsNonRoot: true<br>seccompProfile:<br>  type: RuntimeDefault</pre>|
|`webhook.podDisruptionBudget.enabled`||`bool`|<pre>false</pre>|
|`webhook.podDisruptionBudget.minAvailable`|<p>minAvailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)</p>|`number`|<pre>undefined</pre>|
|`webhook.podDisruptionBudget.maxUnavailable`|<p>maxUnavailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)</p>|`number`|<pre>undefined</pre>|
|`webhook.containerSecurityContext`|<p>Container Security Context to be set on the webhook component container<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>allowPrivilegeEscalation: false<br>capabilities:<br>  drop:<br>    - ALL<br>readOnlyRootFilesystem: true</pre>|
|`webhook.deploymentAnnotations`|<p>Optional additional annotations to add to the webhook Deployment</p>|`object`|<pre>undefined</pre>|
|`webhook.podAnnotations`|<p>Optional additional annotations to add to the webhook Pods</p>|`object`|<pre>undefined</pre>|
|`webhook.serviceAnnotations`|<p>Optional additional annotations to add to the webhook Service</p>|`object`|<pre>undefined</pre>|
|`webhook.mutatingWebhookConfigurationAnnotations`|<p>Optional additional annotations to add to the webhook MutatingWebhookConfiguration</p>|`object`|<pre>undefined</pre>|
|`webhook.validatingWebhookConfigurationAnnotations`|<p>Optional additional annotations to add to the webhook ValidatingWebhookConfiguration</p>|`object`|<pre>undefined</pre>|
|`webhook.extraArgs`|<p>Additional command line flags to pass to cert-manager webhook binary.<br> To see all available flags run docker run quay.io/jetstack/cert-manager-webhook:<version> --help</p>|`array`|<pre>[]</pre>|
|`webhook.featureGates`|<p>Comma separated list of feature gates that should be enabled on the<br> webhook pod.</p>|`string`|<pre>""</pre>|
|`webhook.resources`||`object`|<pre>{}</pre>|
|`webhook.livenessProbe`|<p>Liveness probe values<br> Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes</p>|`object`|<pre>failureThreshold: 3<br>initialDelaySeconds: 60<br>periodSeconds: 10<br>successThreshold: 1<br>timeoutSeconds: 1</pre>|
|`webhook.readinessProbe`|<p>Readiness probe values<br> Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes</p>|`object`|<pre>failureThreshold: 3<br>initialDelaySeconds: 5<br>periodSeconds: 5<br>successThreshold: 1<br>timeoutSeconds: 1</pre>|
|`webhook.nodeSelector`||`object`|<pre>kubernetes.io/os: linux</pre>|
|`webhook.affinity`||`object`|<pre>{}</pre>|
|`webhook.tolerations`||`array`|<pre>[]</pre>|
|`webhook.topologySpreadConstraints`||`array`|<pre>[]</pre>|
|`webhook.podLabels`|<p>Optional additional labels to add to the Webhook Pods</p>|`object`|<pre>{}</pre>|
|`webhook.serviceLabels`|<p>Optional additional labels to add to the Webhook Service</p>|`object`|<pre>{}</pre>|
|`webhook.image.registry`|<p>Registry to pull the image from</p>|`string`|<pre>undefined</pre>|
|`webhook.image.repository`|<p>Image name, this can be the full image including registry or the short name<br> excluding the registry. The registy can also be set in the `registry` property</p>|`string`|<pre>quay.io/jetstack/cert-manager-webhook</pre>|
|`webhook.image.tag`|<p>Override the image tag to deploy by setting this variable.<br> If no value is set, the chart's appVersion will be used.</p>|`string`|<pre>undefined</pre>|
|`webhook.image.digest`|<p>Setting a digest will override any tag</p>|`string`|<pre>undefined</pre>|
|`webhook.image.pullPolicy`|<p>Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy</p>|`string`|<pre>IfNotPresent</pre>|
|`webhook.serviceAccount.create`|<p>Specifies whether a service account should be created</p>|`bool`|<pre>true</pre>|
|`webhook.serviceAccount.name`|<p>The name of the service account to use.<br> If not set and create is true, a name is generated using the fullname template</p>|`string`|<pre>undefined</pre>|
|`webhook.serviceAccount.annotations`|<p>Optional additional annotations to add to the controller's ServiceAccount</p>|`object`|<pre>undefined</pre>|
|`webhook.serviceAccount.labels`|<p>Optional additional labels to add to the webhook's ServiceAccount</p>|`object`|<pre>undefined</pre>|
|`webhook.serviceAccount.automountServiceAccountToken`|<p>Automount API credentials for a Service Account.</p>|`bool`|<pre>true</pre>|
|`webhook.securePort`|<p>The port that the webhook should listen on for requests.<br> In GKE private clusters, by default kubernetes apiservers are allowed to<br> talk to the cluster nodes only on 443 and 10250. so configuring</p><pre>securePort: 10250, will work out of the box without needing to add firewall</pre><p>rules or requiring NET_BIND_SERVICE capabilities to bind port numbers <1000</p>|`number`|<pre>10250</pre>|
|`webhook.hostNetwork`|<p>Specifies if the webhook should be started in hostNetwork mode.</p><p> Required for use in some managed kubernetes clusters (such as AWS EKS) with custom<br> CNI (such as calico), because control-plane managed by AWS cannot communicate<br> with pods' IP CIDR and admission webhooks are not working</p><p> Since the default port for the webhook conflicts with kubelet on the host<br> network, `webhook.securePort` should be changed to an available port if<br> running in hostNetwork mode.</p>|`bool`|<pre>false</pre>|
|`webhook.serviceType`|<p>Specifies how the service should be handled. Useful if you want to expose the<br> webhook to outside of the cluster. In some cases, the control plane cannot<br> reach internal services.</p>|`string`|<pre>ClusterIP</pre>|
|`webhook.loadBalancerIP`||`string`|<pre>undefined</pre>|
|`webhook.url`|<p>Overrides the mutating webhook and validating webhook so they reach the webhook<br> service using the `url` field instead of a service.</p>|`object`|<pre>{}</pre>|
|`webhook.networkPolicy.enabled`||`bool`|<pre>false</pre>|
|`webhook.networkPolicy.ingress`||`array`|<pre>- from:<br>    - ipBlock:<br>        cidr: 0.0.0.0/0</pre>|
|`webhook.networkPolicy.egress`||`array`|<pre>- ports:<br>    - port: 80<br>      protocol: TCP<br>    - port: 443<br>      protocol: TCP<br>    - port: 53<br>      protocol: TCP<br>    - port: 53<br>      protocol: UDP<br>    - port: 6443<br>      protocol: TCP<br>  to:<br>    - ipBlock:<br>        cidr: 0.0.0.0/0</pre>|
|`webhook.volumes`||`array`|<pre>[]</pre>|
|`webhook.volumeMounts`||`array`|<pre>[]</pre>|
|`webhook.enableServiceLinks`|<p>enableServiceLinks indicates whether information about services should be<br> injected into pod's environment variables, matching the syntax of Docker<br> links.</p>|`bool`|<pre>false</pre>|

### CA Injector

|property|description|type|default|
|--|--|--|--|
|`cainjector.enabled`||`bool`|<pre>true</pre>|
|`cainjector.replicaCount`||`number`|<pre>1</pre>|
|`cainjector.config`|<p>Used to configure options for the cainjector pod.<br> This allows setting options that'd usually be provided via flags.<br> An APIVersion and Kind must be specified in your values.yaml file.<br> Flags will override options that are set here. For example:</p><pre>apiVersion: cainjector.config.cert-manager.io/v1alpha1<br>kind: CAInjectorConfiguration<br>logging:<br> verbosity: 2<br> format: text<br>leaderElectionConfig:<br> namespace: kube-system</pre>|`object`|<pre>{}</pre>|
|`cainjector.strategy`|<p>Deployment strategy, for example:</p><pre>type: RollingUpdate<br>rollingUpdate:<br>  maxSurge: 0<br>  maxUnavailable: 1</pre>|`object`|<pre>{}</pre>|
|`cainjector.securityContext.runAsNonRoot`||`bool`|<pre>true</pre>|
|`cainjector.securityContext.seccompProfile.type`||`string`|<pre>RuntimeDefault</pre>|
|`cainjector.podDisruptionBudget.enabled`||`bool`|<pre>false</pre>|
|`cainjector.podDisruptionBudget.minAvailable`|<p>minAvailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)</p>|`number`|<pre>undefined</pre>|
|`cainjector.podDisruptionBudget.maxUnavailable`|<p>maxUnavailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)</p>|`number`|<pre>undefined</pre>|
|`cainjector.containerSecurityContext`|<p>Container Security Context to be set on the cainjector component container<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>allowPrivilegeEscalation: false<br>capabilities:<br>  drop:<br>    - ALL<br>readOnlyRootFilesystem: true</pre>|
|`cainjector.deploymentAnnotations`|<p>Optional additional annotations to add to the cainjector Deployment</p>|`object`|<pre>undefined</pre>|
|`cainjector.podAnnotations`|<p>Optional additional annotations to add to the cainjector Pods</p>|`object`|<pre>undefined</pre>|
|`cainjector.extraArgs`|<p>Additional command line flags to pass to cert-manager cainjector binary.<br> To see all available flags run docker run quay.io/jetstack/cert-manager-cainjector:<version> --help</p>|`array`|<pre>[]</pre>|
|`cainjector.featureGates`|<p>Comma separated list of feature gates that should be enabled on the<br> cainjector pod.</p>|`string`|<pre>""</pre>|
|`cainjector.resources`||`object`|<pre>{}</pre>|
|`cainjector.nodeSelector`||`object`|<pre>kubernetes.io/os: linux</pre>|
|`cainjector.affinity`||`object`|<pre>{}</pre>|
|`cainjector.tolerations`||`array`|<pre>[]</pre>|
|`cainjector.topologySpreadConstraints`||`array`|<pre>[]</pre>|
|`cainjector.podLabels`|<p>Optional additional labels to add to the CA Injector Pods</p>|`object`|<pre>{}</pre>|
|`cainjector.image.registry`|<p>Registry to pull the image from</p>|`string`|<pre>undefined</pre>|
|`cainjector.image.repository`|<p>Image name, this can be the full image including registry or the short name<br> excluding the registry. The registy can also be set in the `registry` property</p>|`string`|<pre>quay.io/jetstack/cert-manager-cainjector</pre>|
|`cainjector.image.tag`|<p>Override the image tag to deploy by setting this variable.<br> If no value is set, the chart's appVersion will be used.</p>|`string`|<pre>undefined</pre>|
|`cainjector.image.digest`|<p>Setting a digest will override any tag</p>|`string`|<pre>undefined</pre>|
|`cainjector.image.pullPolicy`|<p>Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy</p>|`string`|<pre>IfNotPresent</pre>|
|`cainjector.serviceAccount.create`|<p>Specifies whether a service account should be created</p>|`bool`|<pre>true</pre>|
|`cainjector.serviceAccount.name`|<p>The name of the service account to use.<br> If not set and create is true, a name is generated using the fullname template</p>|`string`|<pre>undefined</pre>|
|`cainjector.serviceAccount.annotations`|<p>Optional additional annotations to add to the controller's ServiceAccount</p>|`object`|<pre>undefined</pre>|
|`cainjector.serviceAccount.automountServiceAccountToken`|<p>Automount API credentials for a Service Account.<br> Optional additional labels to add to the cainjector's ServiceAccount</p><pre>labels: {}</pre>|`bool`|<pre>true</pre>|
|`cainjector.volumes`|<p>Automounting API credentials for a particular pod</p><pre>automountServiceAccountToken: true</pre>|`array`|<pre>[]</pre>|
|`cainjector.volumeMounts`||`array`|<pre>[]</pre>|
|`cainjector.enableServiceLinks`|<p>enableServiceLinks indicates whether information about services should be<br> injected into pod's environment variables, matching the syntax of Docker<br> links.</p>|`bool`|<pre>false</pre>|

### ACME Solver

|property|description|type|default|
|--|--|--|--|
|`acmesolver.image.registry`|<p>Image registry to pull from</p>|`string`|<pre>undefined</pre>|
|`acmesolver.image.repository`|<p>Image name, this can be the full image including registry or the short name<br> excluding the registry. The registy can also be set in the `registry` property</p>|`string`|<pre>quay.io/jetstack/cert-manager-acmesolver</pre>|
|`acmesolver.image.tag`|<p>Override the image tag to deploy by setting this variable.<br> If no value is set, the chart's appVersion will be used.</p>|`string`|<pre>undefined</pre>|
|`acmesolver.image.digest`|<p>Setting a digest will override any tag</p>|`string`|<pre>undefined</pre>|

### Startup check API
This startupapicheck is a Helm post-install hook that waits for the webhook
 endpoints to become available.
 The check is implemented using a Kubernetes Job- if you are injecting mesh
 sidecar proxies into cert-manager pods, you probably want to ensure that they
 are not injected into this Job's pod. Otherwise the installation may time out
 due to the Job never being completed because the sidecar proxy does not exit.
 See https://github.com/cert-manager/cert-manager/pull/4414 for context.

|property|description|type|default|
|--|--|--|--|
|`startupapicheck.enabled`||`bool`|<pre>true</pre>|
|`startupapicheck.securityContext`|<p>Pod Security Context to be set on the startupapicheck component Pod<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>runAsNonRoot: true<br>seccompProfile:<br>  type: RuntimeDefault</pre>|
|`startupapicheck.containerSecurityContext`|<p>Container Security Context to be set on the controller component container<br> Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</p>|`object`|<pre>allowPrivilegeEscalation: false<br>capabilities:<br>  drop:<br>    - ALL<br>readOnlyRootFilesystem: true</pre>|
|`startupapicheck.timeout`|<p>Timeout for 'kubectl check api' command</p>|`string`|<pre>1m</pre>|
|`startupapicheck.backoffLimit`|<p>Job backoffLimit</p>|`number`|<pre>4</pre>|
|`startupapicheck.jobAnnotations`|<p>Optional additional annotations to add to the startupapicheck Job</p>|`object`|<pre>helm.sh/hook: post-install<br>helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded<br>helm.sh/hook-weight: "1"</pre>|
|`startupapicheck.podAnnotations`|<p>Optional additional annotations to add to the startupapicheck Pods</p>|`object`|<pre>undefined</pre>|
|`startupapicheck.extraArgs`|<p>Additional command line flags to pass to startupapicheck binary.<br> To see all available flags run docker run quay.io/jetstack/cert-manager-ctl:<version> --help</p><p> We enable verbose logging by default so that if startupapicheck fails, users<br> can know what exactly caused the failure. Verbose logs include details of<br> the webhook URL, IP address and TCP connect errors for example.</p>|`array`|<pre>- -v</pre>|
|`startupapicheck.resources`||`object`|<pre>{}</pre>|
|`startupapicheck.nodeSelector`||`object`|<pre>kubernetes.io/os: linux</pre>|
|`startupapicheck.affinity`||`object`|<pre>{}</pre>|
|`startupapicheck.tolerations`||`array`|<pre>[]</pre>|
|`startupapicheck.podLabels`|<p>Optional additional labels to add to the startupapicheck Pods</p>|`object`|<pre>{}</pre>|
|`startupapicheck.image.registry`|<p>Image registry to pull from</p>|`string`|<pre>undefined</pre>|
|`startupapicheck.image.repository`|<p>Image name, this can be the full image including registry or the short name<br> excluding the registry. The registy can also be set in the `registry` property</p>|`string`|<pre>quay.io/jetstack/cert-manager-startupapicheck</pre>|
|`startupapicheck.image.tag`|<p>Override the image tag to deploy by setting this variable.<br> If no value is set, the chart's appVersion will be used.</p>|`string`|<pre>undefined</pre>|
|`startupapicheck.image.digest`|<p>Setting a digest will override any tag</p>|`string`|<pre>undefined</pre>|
|`startupapicheck.image.pullPolicy`|<p>Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy</p>|`string`|<pre>IfNotPresent</pre>|
|`startupapicheck.rbac.annotations`|<p>annotations for the startup API Check job RBAC and PSP resources</p>|`object`|<pre>helm.sh/hook: post-install<br>helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded<br>helm.sh/hook-weight: "-5"</pre>|
|`startupapicheck.serviceAccount.create`|<p>Specifies whether a service account should be created</p>|`bool`|<pre>true</pre>|
|`startupapicheck.serviceAccount.name`|<p>The name of the service account to use.<br> If not set and create is true, a name is generated using the fullname template</p>|`string`|<pre>undefined</pre>|
|`startupapicheck.serviceAccount.annotations`|<p>Optional additional annotations to add to the Job's ServiceAccount</p>|`object`|<pre>helm.sh/hook: post-install<br>helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded<br>helm.sh/hook-weight: "-5"</pre>|
|`startupapicheck.serviceAccount.automountServiceAccountToken`|<p>Automount API credentials for a Service Account.</p>|`bool`|<pre>true</pre>|
|`startupapicheck.serviceAccount.labels`|<p>Optional additional labels to add to the startupapicheck's ServiceAccount</p>|`object`|<pre>undefined</pre>|
|`startupapicheck.volumes`||`array`|<pre>[]</pre>|
|`startupapicheck.volumeMounts`||`array`|<pre>[]</pre>|
|`startupapicheck.enableServiceLinks`|<p>enableServiceLinks indicates whether information about services should be<br> injected into pod's environment variables, matching the syntax of Docker<br> links.</p>|`bool`|<pre>false</pre>|

