# Cert Manager

## Parameters

### Global


<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>global.imagePullSecrets</td>
<td>

<p>

Reference to one or more secrets to be used when pulling images

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
</pre>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>global.commonLabels</td>
<td>

<p>

Labels to apply to all resources<br/>
Please note that this does not add labels to the resources created dynamically by the controllers. For these resources, you have to add the labels in the template in the cert-manager custom resource:

</p>
<p>

eg. podTemplate/ ingressTemplate in ACMEChallengeSolverHTTP01Ingress

</p>

<pre lang="yaml">
ref: https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEChallengeSolverHTTP01Ingress
</pre>

<p>

eg. secretTemplate in CertificateSpec

</p>

<pre lang="yaml">
ref: https://cert-manager.io/docs/reference/api-docs/#cert-manager.io/v1.CertificateSpec
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>global.priorityClassName</td>
<td>

<p>

Optional priority class to be used for the cert-manager pods

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>global.rbac.create</td>
<td>

<p>

Create RBAC rules

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>global.rbac.aggregateClusterRoles</td>
<td>

<p>

Aggregate ClusterRoles to Kubernetes default user-facing roles. ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>global.logLevel</td>
<td>

<p>

Set the verbosity of cert-manager. Range of 0 - 6 with 6 being the most verbose.

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">2</pre>

</td>

</tr>
<tr>

<td>global.leaderElection.namespace</td>
<td>

<p>

Override the namespace used for the leader election lease

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">kube-system</pre>

</td>

</tr>
<tr>

<td>global.leaderElection.leaseDuration</td>
<td>

<p>

The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>global.leaderElection.renewDeadline</td>
<td>

<p>

The interval between attempts by the acting master to renew a leadership slot before it stops leading. This must be less than or equal to the lease duration.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>global.leaderElection.retryPeriod</td>
<td>

<p>

The duration the clients should wait between attempting acquisition and renewal of a leadership.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>installCRDs</td>
<td>

<p>

Install the CRDs

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
</table>

### Controller


<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>replicaCount</td>
<td>

<p>

Number of replicas to run of the cert-manager controller

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">1</pre>

</td>

</tr>
<tr>

<td>strategy</td>
<td>

<p>

Update strategy to use, for example:

</p>

<pre lang="yaml">
type: RollingUpdate
rollingUpdate:
  maxSurge: 0
  maxUnavailable: 1
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>podDisruptionBudget.minAvailable</td>
<td>


<p>

minAvailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>podDisruptionBudget.maxUnavailable</td>
<td>


<p>

maxUnavailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>featureGates</td>
<td>

<p>

Comma separated list of feature gates that should be enabled on the controller pod.

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>maxConcurrentChallenges</td>
<td>

<p>

The maximum number of challenges that can be scheduled as 'processing' at once

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">60</pre>

</td>

</tr>
<tr>

<td>image.registry</td>
<td>

<p>

Registry to pull the image from

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>image.repository</td>
<td>

<p>

Image name, this can be the full image including registry or the short name excluding the registry. The registy can also be set in the `registry` property

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">quay.io/jetstack/cert-manager-controller</pre>

</td>

</tr>
<tr>

<td>image.tag</td>
<td>

<p>

Override the image tag to deploy by setting this variable. If no value is set, the chart's appVersion will be used.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>image.digest</td>
<td>

<p>

Setting a digest will override any tag

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>image.pullPolicy</td>
<td>

<p>

Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">IfNotPresent</pre>

</td>

</tr>
<tr>

<td>clusterResourceNamespace</td>
<td>

<p>

Override the namespace used to store DNS provider credentials etc. for ClusterIssuer resources. By default, the same namespace as cert-manager is deployed within is used. This namespace will not be automatically created by the Helm chart.

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>namespace</td>
<td>

<p>

This namespace allows you to define where the services will be installed into if not set then they will use the namespace of the release. This is helpful when installing cert manager as a chart dependency (sub chart)

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>serviceAccount.create</td>
<td>

<p>

Specifies whether a service account should be created

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>serviceAccount.name</td>
<td>

<p>

The name of the service account to use.<br/>
If not set and create is true, a name is generated using the fullname template

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>serviceAccount.annotations</td>
<td>

<p>

Optional additional annotations to add to the controller's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>serviceAccount.labels</td>
<td>

<p>

Automount API credentials for a Service Account.<br/>
Optional additional labels to add to the controller's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>serviceAccount.automountServiceAccountToken</td>
<td>

<p>

Service account token wil be automatically mounted in Pods

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>automountServiceAccountToken</td>
<td>

<p>

Automounting API credentials for a particular pod

</p>


</td>
<td>bool</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>enableCertificateOwnerRef</td>
<td>

<p>

When this flag is enabled, secrets will be automatically removed when the certificate resource is deleted

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>config</td>
<td>

<p>

Used to configure options for the controller pod.<br/>
This allows setting options that'd usually be provided via flags. An APIVersion and Kind must be specified in your values.yaml file.<br/>
Flags will override options that are set here.

</p>
<p>

For example:

</p>

<pre lang="yaml">
apiVersion: controller.config.cert-manager.io/v1alpha1
kind: ControllerConfiguration
logging:
  verbosity: 2
  format: text
leaderElectionConfig:
  namespace: kube-system
kubernetesAPIQPS: 9000
kubernetesAPIBurst: 9000
numberOfConcurrentWorkers: 200
featureGates:
  AdditionalCertificateOutputFormats: true
  DisallowInsecureCSRUsageDefinition: true
  ExperimentalCertificateSigningRequestControllers: true
  ExperimentalGatewayAPISupport: true
  LiteralCertificateSubject: true
  SecretsFilteredCaching: true
  ServerSideApply: true
  StableCertificateRequestName: true
  UseCertificateRequestBasicConstraints: true
  ValidateCAA: true
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>dns01RecursiveNameservers</td>
<td>

<p>

Comma separated string with host and port of the recursive nameservers cert-manager should query

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>dns01RecursiveNameserversOnly</td>
<td>

<p>

Forces cert-manager to only use the recursive nameservers for verification. Enabling this option could cause the DNS01 self check to take longer due to caching performed by the recursive nameservers

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>extraArgs</td>
<td>

<p>

Additional command line flags to pass to cert-manager controller binary. To see all available flags run docker run quay.io/jetstack/cert-manager-controller:<version> --help

</p>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>extraEnv</td>
<td>

<p>

Additional environment variables

</p>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>resources</td>
<td>

<p>

Resources the controller will be given

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>securityContext</td>
<td>

<p>

Pod Security Context

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
</pre>


</td>
<td>object</td>
<td>

<pre lang="yaml">runAsNonRoot: true
seccompProfile:
  type: RuntimeDefault</pre>

</td>

</tr>
<tr>

<td>containerSecurityContext</td>
<td>

<p>

Container Security Context to be set on the controller component container

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
</pre>


</td>
<td>object</td>
<td>

<pre lang="yaml">allowPrivilegeEscalation: false
capabilities:
  drop:
    - ALL
readOnlyRootFilesystem: true</pre>

</td>

</tr>
<tr>

<td>volumes</td>
<td>

<p>

Volumes to mount to the controller pod

</p>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>volumeMounts</td>
<td>

<p>

Volumes specified in `volumes` to mount to the controller container

</p>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>deploymentAnnotations</td>
<td>

<p>

Optional additional annotations to add to the controller Deployment

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>podAnnotations</td>
<td>

<p>

Optional additional annotations to add to the controller Pods

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>podLabels</td>
<td>

<p>

Optional additional labels to add to the controller Pods

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>serviceAnnotations</td>
<td>

<p>

Optional annotations to add to the controller Service

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>serviceLabels</td>
<td>

<p>

Optional additional labels to add to the controller Service

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>podDnsPolicy</td>
<td>

<p>

DNS policy to use within the controller pod

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>podDnsConfig</td>
<td>

<p>

Optional DNS settings, useful if you have a public and private DNS zone for the same domain on Route 53. What follows is an example of ensuring cert-manager can access an ingress or DNS TXT records at all times. NOTE: This requires Kubernetes 1.10 or `CustomPodDNS` feature gate enabled for the cluster to work.

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>nodeSelector</td>
<td>

<p>

Node selector to limit the nodes the controller can schedule on

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">kubernetes.io/os: linux</pre>

</td>

</tr>
<tr>

<td>ingressShim.defaultIssuerName</td>
<td>

<p>

Optional default issuer to use for ingress resources

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>ingressShim.defaultIssuerKind</td>
<td>

<p>

Optional default issuer kind to use for ingress resources

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>ingressShim.defaultIssuerGroup</td>
<td>

<p>

Optional default issuer group to use for ingress resources

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>http_proxy</td>
<td>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>https_proxy</td>
<td>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>no_proxy</td>
<td>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>affinity</td>
<td>

<p>

A Kubernetes Affinity, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#affinity-v1-core for example:

</p>

<pre lang="yaml">
affinity:
  nodeAffinity:
   requiredDuringSchedulingIgnoredDuringExecution:
     nodeSelectorTerms:
     - matchExpressions:
       - key: foo.bar.com/role
         operator: In
         values:
         - master
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>tolerations</td>
<td>

<p>

A list of Kubernetes Tolerations, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#toleration-v1-core for example:

</p>

<pre lang="yaml">
tolerations:
- key: foo.bar.com/role
  operator: Equal
  value: master
  effect: NoSchedule
</pre>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>topologySpreadConstraints</td>
<td>

<p>

A list of Kubernetes TopologySpreadConstraints, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#topologyspreadconstraint-v1-core for example:

</p>

<pre lang="yaml">
topologySpreadConstraints:
- maxSkew: 2
  topologyKey: topology.kubernetes.io/zone
  whenUnsatisfiable: ScheduleAnyway
  labelSelector:
    matchLabels:
      app.kubernetes.io/instance: cert-manager
      app.kubernetes.io/component: controller
</pre>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>livenessProbe.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>livenessProbe.initialDelaySeconds</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">10</pre>

</td>

</tr>
<tr>

<td>livenessProbe.periodSeconds</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">10</pre>

</td>

</tr>
<tr>

<td>livenessProbe.timeoutSeconds</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">15</pre>

</td>

</tr>
<tr>

<td>livenessProbe.successThreshold</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">1</pre>

</td>

</tr>
<tr>

<td>livenessProbe.failureThreshold</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">8</pre>

</td>

</tr>
<tr>

<td>enableServiceLinks</td>
<td>

<p>

enableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
</table>

### Prometheus


<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>prometheus.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.enabled</td>
<td>

<p>

Create a ServiceMonitor resource to scrape the metrics endpoint

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.prometheusInstance</td>
<td>

</td>
<td>string</td>
<td>

<pre lang="yaml">default</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.targetPort</td>
<td>

<p>

The port to scrape metrics from

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">9402</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.path</td>
<td>

<p>

Path to scrape metrics from

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">/metrics</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.interval</td>
<td>

<p>

Interval to scrape metrics

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">60s</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.scrapeTimeout</td>
<td>

<p>

Timeout for each metrics scrape

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">30s</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.labels</td>
<td>

<p>

Labels to add to the ServiceMonitor resource

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.annotations</td>
<td>

<p>

Annotations to add to the ServiceMonitor resource

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.honorLabels</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>prometheus.servicemonitor.endpointAdditionalProperties</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.enabled</td>
<td>

<p>

Create a PodMonitor resource to scrape the metrics endpoint

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.prometheusInstance</td>
<td>

</td>
<td>string</td>
<td>

<pre lang="yaml">default</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.path</td>
<td>

<p>

Path to scrape metrics from

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">/metrics</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.interval</td>
<td>

<p>

Interval to scrape metrics

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">60s</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.scrapeTimeout</td>
<td>

<p>

Timeout for each metrics scrape

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">30s</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.labels</td>
<td>

<p>

Labels to add to the PodMonitor resource

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.annotations</td>
<td>

<p>

Annotations to add to the PodMonitor resource

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.honorLabels</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>prometheus.podmonitor.endpointAdditionalProperties</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
</table>

### Webhook


<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>webhook.replicaCount</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">1</pre>

</td>

</tr>
<tr>

<td>webhook.timeoutSeconds</td>
<td>

<p>

Seconds the API server should wait for the webhook to respond before treating the call as a failure.<br/>
Value must be between 1 and 30 seconds. See:<br/>
https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/validating-webhook-configuration-v1/

</p>
<p>

We set the default to the maximum value of 30 seconds. Here's why: Users sometimes report that the connection between the K8S API server and the cert-manager webhook server times out. If *this* timeout is reached, the error message will be "context deadline exceeded", which doesn't help the user diagnose what phase of the HTTPS connection timed out. For example, it could be during DNS resolution, TCP connection, TLS negotiation, HTTP negotiation, or slow HTTP response from the webhook server. So by setting this timeout to its maximum value the underlying timeout error message has more chance of being returned to the end user.

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">30</pre>

</td>

</tr>
<tr>

<td>webhook.config</td>
<td>

<p>

Used to configure options for the webhook pod.<br/>
This allows setting options that'd usually be provided via flags. An APIVersion and Kind must be specified in your values.yaml file. Flags will override options that are set here. Example config:

</p>

<pre lang="yaml">
apiVersion: webhook.config.cert-manager.io/v1alpha1
kind: WebhookConfiguration
# The port that the webhook should listen on for requests.
# In GKE private clusters, by default kubernetes apiservers are allowed to
# talk to the cluster nodes only on 443 and 10250. so configuring
# securePort: 10250, will work out of the box without needing to add firewall
# rules or requiring NET_BIND_SERVICE capabilities to bind port numbers <1000.
# This should be uncommented and set as a default by the chart once we graduate
# the apiVersion of WebhookConfiguration past v1alpha1.
securePort: 10250
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.strategy</td>
<td>

<p>

Deployment strategy, for example:

</p>

<pre lang="yaml">
type: RollingUpdate
rollingUpdate:
  maxSurge: 0
  maxUnavailable: 1
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.securityContext</td>
<td>

<p>

Pod Security Context to be set on the webhook component Pod. Rref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">runAsNonRoot: true
seccompProfile:
  type: RuntimeDefault</pre>

</td>

</tr>
<tr>

<td>webhook.podDisruptionBudget.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>webhook.podDisruptionBudget.minAvailable</td>
<td>


<p>

minAvailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.podDisruptionBudget.maxUnavailable</td>
<td>


<p>

maxUnavailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.containerSecurityContext</td>
<td>

<p>

Container Security Context to be set on the webhook component container

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
</pre>


</td>
<td>object</td>
<td>

<pre lang="yaml">allowPrivilegeEscalation: false
capabilities:
  drop:
    - ALL
readOnlyRootFilesystem: true</pre>

</td>

</tr>
<tr>

<td>webhook.deploymentAnnotations</td>
<td>

<p>

Optional additional annotations to add to the webhook Deployment

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.podAnnotations</td>
<td>

<p>

Optional additional annotations to add to the webhook Pods

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.serviceAnnotations</td>
<td>

<p>

Optional additional annotations to add to the webhook Service

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.mutatingWebhookConfigurationAnnotations</td>
<td>

<p>

Optional additional annotations to add to the webhook MutatingWebhookConfiguration

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.validatingWebhookConfigurationAnnotations</td>
<td>

<p>

Optional additional annotations to add to the webhook ValidatingWebhookConfiguration

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.extraArgs</td>
<td>

<p>

Additional command line flags to pass to cert-manager webhook binary. To see all available flags run docker run quay.io/jetstack/cert-manager-webhook:<version> --help

</p>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>webhook.featureGates</td>
<td>

<p>

Comma separated list of feature gates that should be enabled on the webhook pod.

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>webhook.resources</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.livenessProbe</td>
<td>

<p>

Liveness probe values

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes
</pre>

<p>



</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">failureThreshold: 3
initialDelaySeconds: 60
periodSeconds: 10
successThreshold: 1
timeoutSeconds: 1</pre>

</td>

</tr>
<tr>

<td>webhook.readinessProbe</td>
<td>

<p>

Readiness probe values

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes
</pre>

<p>



</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">failureThreshold: 3
initialDelaySeconds: 5
periodSeconds: 5
successThreshold: 1
timeoutSeconds: 1</pre>

</td>

</tr>
<tr>

<td>webhook.nodeSelector</td>
<td>


</td>
<td>object</td>
<td>

<pre lang="yaml">kubernetes.io/os: linux</pre>

</td>

</tr>
<tr>

<td>webhook.affinity</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.tolerations</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>webhook.topologySpreadConstraints</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>webhook.podLabels</td>
<td>

<p>

Optional additional labels to add to the Webhook Pods

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.serviceLabels</td>
<td>

<p>

Optional additional labels to add to the Webhook Service

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.image.registry</td>
<td>

<p>

Registry to pull the image from

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.image.repository</td>
<td>

<p>

Image name, this can be the full image including registry or the short name excluding the registry. The registy can also be set in the `registry` property

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">quay.io/jetstack/cert-manager-webhook</pre>

</td>

</tr>
<tr>

<td>webhook.image.tag</td>
<td>

<p>

Override the image tag to deploy by setting this variable. If no value is set, the chart's appVersion will be used.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.image.digest</td>
<td>

<p>

Setting a digest will override any tag

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.image.pullPolicy</td>
<td>

<p>

Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">IfNotPresent</pre>

</td>

</tr>
<tr>

<td>webhook.serviceAccount.create</td>
<td>

<p>

Specifies whether a service account should be created

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>webhook.serviceAccount.name</td>
<td>

<p>

The name of the service account to use.<br/>
If not set and create is true, a name is generated using the fullname template

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.serviceAccount.annotations</td>
<td>

<p>

Optional additional annotations to add to the controller's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.serviceAccount.labels</td>
<td>

<p>

Optional additional labels to add to the webhook's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.serviceAccount.automountServiceAccountToken</td>
<td>

<p>

Automount API credentials for a Service Account.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>webhook.securePort</td>
<td>

<p>

The port that the webhook should listen on for requests. In GKE private clusters, by default kubernetes apiservers are allowed to talk to the cluster nodes only on 443 and 10250. so configuring securePort: 10250, will work out of the box without needing to add firewall rules or requiring NET_BIND_SERVICE capabilities to bind port numbers <1000

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">10250</pre>

</td>

</tr>
<tr>

<td>webhook.hostNetwork</td>
<td>

<p>

Specifies if the webhook should be started in hostNetwork mode.

</p>
<p>

Required for use in some managed kubernetes clusters (such as AWS EKS) with custom. CNI (such as calico), because control-plane managed by AWS cannot communicate with pods' IP CIDR and admission webhooks are not working

</p>
<p>

Since the default port for the webhook conflicts with kubelet on the host network, `webhook.securePort` should be changed to an available port if running in hostNetwork mode.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>webhook.serviceType</td>
<td>

<p>

Specifies how the service should be handled. Useful if you want to expose the webhook to outside of the cluster. In some cases, the control plane cannot reach internal services.

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">ClusterIP</pre>

</td>

</tr>
<tr>

<td>webhook.loadBalancerIP</td>
<td>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>webhook.url</td>
<td>

<p>

Overrides the mutating webhook and validating webhook so they reach the webhook service using the `url` field instead of a service.

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>webhook.networkPolicy.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>webhook.networkPolicy.ingress</td>
<td>


</td>
<td>array</td>
<td>

<pre lang="yaml">- from:
    - ipBlock:
        cidr: 0.0.0.0/0</pre>

</td>

</tr>
<tr>

<td>webhook.networkPolicy.egress</td>
<td>


</td>
<td>array</td>
<td>

<pre lang="yaml">- ports:
    - port: 80
      protocol: TCP
    - port: 443
      protocol: TCP
    - port: 53
      protocol: TCP
    - port: 53
      protocol: UDP
    - port: 6443
      protocol: TCP
  to:
    - ipBlock:
        cidr: 0.0.0.0/0</pre>

</td>

</tr>
<tr>

<td>webhook.volumes</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>webhook.volumeMounts</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>webhook.enableServiceLinks</td>
<td>

<p>

enableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
</table>

### CA Injector


<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>cainjector.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>cainjector.replicaCount</td>
<td>

</td>
<td>number</td>
<td>

<pre lang="yaml">1</pre>

</td>

</tr>
<tr>

<td>cainjector.config</td>
<td>

<p>

Used to configure options for the cainjector pod.<br/>
This allows setting options that'd usually be provided via flags. An APIVersion and Kind must be specified in your values.yaml file. Flags will override options that are set here. For example:

</p>

<pre lang="yaml">
apiVersion: cainjector.config.cert-manager.io/v1alpha1
kind: CAInjectorConfiguration
logging:
 verbosity: 2
 format: text
leaderElectionConfig:
 namespace: kube-system
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>cainjector.strategy</td>
<td>

<p>

Deployment strategy, for example:

</p>

<pre lang="yaml">
type: RollingUpdate
rollingUpdate:
  maxSurge: 0
  maxUnavailable: 1
</pre>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>cainjector.securityContext.runAsNonRoot</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>cainjector.securityContext.seccompProfile.type</td>
<td>

</td>
<td>string</td>
<td>

<pre lang="yaml">RuntimeDefault</pre>

</td>

</tr>
<tr>

<td>cainjector.podDisruptionBudget.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
<tr>

<td>cainjector.podDisruptionBudget.minAvailable</td>
<td>


<p>

minAvailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.podDisruptionBudget.maxUnavailable</td>
<td>


<p>

maxUnavailable can either be set to an integer (e.g. 1) or a percentage value (e.g. 25%)

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.containerSecurityContext</td>
<td>

<p>

Container Security Context to be set on the cainjector component container

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
</pre>


</td>
<td>object</td>
<td>

<pre lang="yaml">allowPrivilegeEscalation: false
capabilities:
  drop:
    - ALL
readOnlyRootFilesystem: true</pre>

</td>

</tr>
<tr>

<td>cainjector.deploymentAnnotations</td>
<td>

<p>

Optional additional annotations to add to the cainjector Deployment

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.podAnnotations</td>
<td>

<p>

Optional additional annotations to add to the cainjector Pods

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.extraArgs</td>
<td>

<p>

Additional command line flags to pass to cert-manager cainjector binary. To see all available flags run docker run quay.io/jetstack/cert-manager-cainjector:<version> --help

</p>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>cainjector.featureGates</td>
<td>

<p>

Comma separated list of feature gates that should be enabled on the cainjector pod.

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">""</pre>

</td>

</tr>
<tr>

<td>cainjector.resources</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>cainjector.nodeSelector</td>
<td>


</td>
<td>object</td>
<td>

<pre lang="yaml">kubernetes.io/os: linux</pre>

</td>

</tr>
<tr>

<td>cainjector.affinity</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>cainjector.tolerations</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>cainjector.topologySpreadConstraints</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>cainjector.podLabels</td>
<td>

<p>

Optional additional labels to add to the CA Injector Pods

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>cainjector.image.registry</td>
<td>

<p>

Registry to pull the image from

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.image.repository</td>
<td>

<p>

Image name, this can be the full image including registry or the short name excluding the registry. The registy can also be set in the `registry` property

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">quay.io/jetstack/cert-manager-cainjector</pre>

</td>

</tr>
<tr>

<td>cainjector.image.tag</td>
<td>

<p>

Override the image tag to deploy by setting this variable. If no value is set, the chart's appVersion will be used.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.image.digest</td>
<td>

<p>

Setting a digest will override any tag

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.image.pullPolicy</td>
<td>

<p>

Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">IfNotPresent</pre>

</td>

</tr>
<tr>

<td>cainjector.serviceAccount.create</td>
<td>

<p>

Specifies whether a service account should be created

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>cainjector.serviceAccount.name</td>
<td>

<p>

The name of the service account to use.<br/>
If not set and create is true, a name is generated using the fullname template

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.serviceAccount.annotations</td>
<td>

<p>

Optional additional annotations to add to the controller's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.serviceAccount.labels</td>
<td>

<p>

Automount API credentials for a Service Account.<br/>
Optional additional labels to add to the cainjector's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.serviceAccount.automountServiceAccountToken</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>cainjector.automountServiceAccountToken</td>
<td>

<p>

Automounting API credentials for a particular pod

</p>


</td>
<td>bool</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>cainjector.volumes</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>cainjector.volumeMounts</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>cainjector.enableServiceLinks</td>
<td>

<p>

enableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
</table>

### ACME Solver


<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>acmesolver.image.registry</td>
<td>

<p>

Image registry to pull from

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>acmesolver.image.repository</td>
<td>

<p>

Image name, this can be the full image including registry or the short name excluding the registry. The registy can also be set in the `registry` property

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">quay.io/jetstack/cert-manager-acmesolver</pre>

</td>

</tr>
<tr>

<td>acmesolver.image.tag</td>
<td>

<p>

Override the image tag to deploy by setting this variable. If no value is set, the chart's appVersion will be used.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>acmesolver.image.digest</td>
<td>

<p>

Setting a digest will override any tag

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
</table>

### Startup check API


<p>

This startupapicheck is a Helm post-install hook that waits for the webhook endpoints to become available. The check is implemented using a Kubernetes Job- if you are injecting mesh sidecar proxies into cert-manager pods, you probably want to ensure that they are not injected into this Job's pod. Otherwise the installation may time out due to the Job never being completed because the sidecar proxy does not exit. See https://github.com/cert-manager/cert-manager/pull/4414 for context.

</p>

<table>
<tr>
<th>Property</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
<tr>

<td>startupapicheck.enabled</td>
<td>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>startupapicheck.securityContext</td>
<td>

<p>

Pod Security Context to be set on the startupapicheck component Pod

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
</pre>


</td>
<td>object</td>
<td>

<pre lang="yaml">runAsNonRoot: true
seccompProfile:
  type: RuntimeDefault</pre>

</td>

</tr>
<tr>

<td>startupapicheck.containerSecurityContext</td>
<td>

<p>

Container Security Context to be set on the controller component container

</p>

<pre lang="yaml">
ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
</pre>


</td>
<td>object</td>
<td>

<pre lang="yaml">allowPrivilegeEscalation: false
capabilities:
  drop:
    - ALL
readOnlyRootFilesystem: true</pre>

</td>

</tr>
<tr>

<td>startupapicheck.timeout</td>
<td>

<p>

Timeout for 'kubectl check api' command

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">1m</pre>

</td>

</tr>
<tr>

<td>startupapicheck.backoffLimit</td>
<td>

<p>

Job backoffLimit

</p>

</td>
<td>number</td>
<td>

<pre lang="yaml">4</pre>

</td>

</tr>
<tr>

<td>startupapicheck.jobAnnotations</td>
<td>

<p>

Optional additional annotations to add to the startupapicheck Job

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">helm.sh/hook: post-install
helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded
helm.sh/hook-weight: "1"</pre>

</td>

</tr>
<tr>

<td>startupapicheck.podAnnotations</td>
<td>

<p>

Optional additional annotations to add to the startupapicheck Pods

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>startupapicheck.extraArgs</td>
<td>

<p>

Additional command line flags to pass to startupapicheck binary. To see all available flags run docker run quay.io/jetstack/cert-manager-ctl:<version> --help

</p>
<p>

We enable verbose logging by default so that if startupapicheck fails, users can know what exactly caused the failure. Verbose logs include details of the webhook URL, IP address and TCP connect errors for example.

</p>


</td>
<td>array</td>
<td>

<pre lang="yaml">- -v</pre>

</td>

</tr>
<tr>

<td>startupapicheck.resources</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>startupapicheck.nodeSelector</td>
<td>


</td>
<td>object</td>
<td>

<pre lang="yaml">kubernetes.io/os: linux</pre>

</td>

</tr>
<tr>

<td>startupapicheck.affinity</td>
<td>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>startupapicheck.tolerations</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>startupapicheck.podLabels</td>
<td>

<p>

Optional additional labels to add to the startupapicheck Pods

</p>

</td>
<td>object</td>
<td>

<pre lang="yaml">{}</pre>

</td>

</tr>
<tr>

<td>startupapicheck.image.registry</td>
<td>

<p>

Image registry to pull from

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>startupapicheck.image.repository</td>
<td>

<p>

Image name, this can be the full image including registry or the short name excluding the registry. The registy can also be set in the `registry` property

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">quay.io/jetstack/cert-manager-startupapicheck</pre>

</td>

</tr>
<tr>

<td>startupapicheck.image.tag</td>
<td>

<p>

Override the image tag to deploy by setting this variable. If no value is set, the chart's appVersion will be used.

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>startupapicheck.image.digest</td>
<td>

<p>

Setting a digest will override any tag

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>startupapicheck.image.pullPolicy</td>
<td>

<p>

Image pull policy, see https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy

</p>

</td>
<td>string</td>
<td>

<pre lang="yaml">IfNotPresent</pre>

</td>

</tr>
<tr>

<td>startupapicheck.rbac.annotations</td>
<td>

<p>

annotations for the startup API Check job RBAC and PSP resources

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">helm.sh/hook: post-install
helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded
helm.sh/hook-weight: "-5"</pre>

</td>

</tr>
<tr>

<td>startupapicheck.serviceAccount.create</td>
<td>

<p>

Specifies whether a service account should be created

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>startupapicheck.serviceAccount.name</td>
<td>

<p>

The name of the service account to use.<br/>
If not set and create is true, a name is generated using the fullname template

</p>


</td>
<td>string</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>startupapicheck.serviceAccount.annotations</td>
<td>

<p>

Optional additional annotations to add to the Job's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">helm.sh/hook: post-install
helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded
helm.sh/hook-weight: "-5"</pre>

</td>

</tr>
<tr>

<td>startupapicheck.serviceAccount.automountServiceAccountToken</td>
<td>

<p>

Automount API credentials for a Service Account.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">true</pre>

</td>

</tr>
<tr>

<td>startupapicheck.serviceAccount.labels</td>
<td>

<p>

Optional additional labels to add to the startupapicheck's ServiceAccount

</p>


</td>
<td>object</td>
<td>

<pre lang="yaml">undefined</pre>

</td>

</tr>
<tr>

<td>startupapicheck.volumes</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>startupapicheck.volumeMounts</td>
<td>

</td>
<td>array</td>
<td>

<pre lang="yaml">[]</pre>

</td>

</tr>
<tr>

<td>startupapicheck.enableServiceLinks</td>
<td>

<p>

enableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links.

</p>

</td>
<td>bool</td>
<td>

<pre lang="yaml">false</pre>

</td>

</tr>
</table>

