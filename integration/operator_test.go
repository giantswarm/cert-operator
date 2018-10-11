// +build k8srequired

package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	certOperatorValuesFile = "/tmp/cert-operator-install.yaml"
	// certOperatorChartValues values required by cert-operator-chart, the environment
	// variables will be expanded before writing the contents to a file.
	certOperatorChartValues = `commonDomain: ${COMMON_DOMAIN}
clusterName: ${CLUSTER_NAME}
Installation:
  V1:
    Auth:
      Vault:
        Address: http://vault.default.svc.cluster.local:8200
        CA:
          TTL: 1440h
    GiantSwarm:
      CertOperator:
        CRD:
          LabelSelector: 'giantswarm.io/cluster={{ .ClusterName }}'
    Guest:
      Kubernetes:
        API:
          EndpointBase: ${COMMON_DOMAIN}
    Secret:
      CertOperator:
        SecretYaml: |
          service:
            vault:
              config:
                token: ${VAULT_TOKEN}
      Registry:
        PullSecret:
          DockerConfigJSON: "{\"auths\":{\"quay.io\":{\"auth\":\"$REGISTRY_PULL_SECRET\"}}}"
`
)

var cs *clients

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var v int
	var err error
	cs, err = newClients()
	if err != nil {
		v = 1
		log.Printf("unexpected error: %v\n", err)
	}

	if err := cs.setUp(); err != nil {
		v = 1
		log.Printf("unexpected error: %v\n", err)
	}

	if v == 0 {
		v = m.Run()
	}

	cs.tearDown()

	os.Exit(v)
}

func TestSecretsAreCreated(t *testing.T) {
	err := runCmd("helm registry install quay.io/giantswarm/cert-resource-lab-chart:stable -- -n cert-resource-lab --set commonDomain=${COMMON_DOMAIN} --set clusterName=${CLUSTER_NAME}")
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}

	{
		o := cs.secretFunc("giantswarm", fmt.Sprintf("%s-api", os.Getenv("CLUSTER_NAME")))
		b := backoff.NewExponential(30*time.Second, 5*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}
}

func (cs *clients) setUp() error {
	if err := cs.createGSNamespace(); err != nil {
		return microerror.Mask(err)
	}

	if err := cs.installVault(); err != nil {
		return microerror.Mask(err)
	}

	if err := cs.installCertOperator(); err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (cs *clients) tearDown() {
	runCmd("helm delete vault --purge")
	runCmd("helm delete cert-operator --purge")
	runCmd("helm delete cert-resource-lab --purge")
	cs.K8sCs.CoreV1().Namespaces().Delete("giantswarm", &metav1.DeleteOptions{})
}

func (cs *clients) createGSNamespace() error {
	{
		n := &apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "giantswarm",
			},
		}
		_, err := cs.K8sCs.CoreV1().Namespaces().Create(n)
		if errors.IsAlreadyExists(err) {
			return nil
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		o := cs.activeNamespaceFunc("giantswarm")
		b := backoff.NewExponential(30*time.Second, 5*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (cs *clients) installVault() error {
	err := runCmd("helm registry install quay.io/giantswarm/vaultlab-chart:stable -- --set vaultToken=${VAULT_TOKEN} -n vault")
	if err != nil {
		return microerror.Mask(err)
	}

	{
		o := cs.runningPodFunc("default", "app=vault")
		b := backoff.NewExponential(30*time.Second, 5*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (cs *clients) installCertOperator() error {
	{
		certOperatorChartValuesEnv := os.ExpandEnv(certOperatorChartValues)
		err := ioutil.WriteFile(certOperatorValuesFile, []byte(certOperatorChartValuesEnv), os.ModePerm)
		if err != nil {
			return microerror.Mask(err)
		}
		err = runCmd("helm registry install quay.io/giantswarm/cert-operator-chart@1.0.0-${CIRCLE_SHA1} -- -n cert-operator --values " + certOperatorValuesFile)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		o := cs.certConfigFunc()
		b := backoff.NewExponential(30*time.Second, 5*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func runCmd(cmdStr string) error {
	cmdEnv := os.ExpandEnv(cmdStr)
	fields := strings.Fields(cmdEnv)
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	return cmd.Run()
}

func (cs *clients) runningPodFunc(namespace, labelSelector string) func() error {
	return func() error {
		pods, err := cs.K8sCs.CoreV1().Pods(namespace).List(metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return microerror.Mask(err)
		}
		if len(pods.Items) > 1 {
			return microerror.Mask(tooManyResultsError)
		}
		pod := pods.Items[0]
		phase := pod.Status.Phase
		if phase != apiv1.PodRunning {
			return microerror.Maskf(unexpectedStatusPhase, "current status: %s", string(phase))
		}
		return nil
	}
}

func (cs *clients) activeNamespaceFunc(name string) func() error {
	return func() error {
		ns, err := cs.K8sCs.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		phase := ns.Status.Phase
		if phase != apiv1.NamespaceActive {
			return microerror.Maskf(unexpectedStatusPhase, "current status: %s", string(phase))
		}

		return nil
	}
}

func (cs *clients) secretFunc(namespace, secretName string) func() error {
	return func() error {
		_, err := cs.K8sCs.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (cs *clients) certConfigFunc() func() error {
	return func() error {
		_, err := cs.GsCs.CoreV1alpha1().
			CertConfigs("default").
			List(metav1.ListOptions{})

		return microerror.Mask(err)
	}
}
