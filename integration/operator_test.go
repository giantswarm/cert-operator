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

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	certOperatorValuesFile = "/tmp/cert-operator-install.yaml"
	defaultTimeout         = 25
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
		t.Errorf("could not install cert-resource-lab, %v", err)
	}

	secretName := fmt.Sprintf("%s-api", os.Getenv("CLUSTER_NAME"))
	err = waitFor(cs.secretFunc("default", secretName))
	if err != nil {
		t.Errorf("could not find expected secret '%s': %v", secretName, err)
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
	cs.K8sCs.ExtensionsV1beta1().ThirdPartyResources().Delete("certificate.giantswarm.io", &metav1.DeleteOptions{})
}

func (cs *clients) createGSNamespace() error {
	// check if the namespace already exists
	_, err := cs.K8sCs.CoreV1().Namespaces().Get("giantswarm", metav1.GetOptions{})
	if err == nil {
		return nil
	}

	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "giantswarm",
		},
	}
	_, err = cs.K8sCs.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	return waitFor(cs.activeNamespaceFunc("giantswarm"))
}

func (cs *clients) installVault() error {
	if err := runCmd("helm registry install quay.io/giantswarm/vaultlab-chart:stable -- --set vaultToken=${VAULT_TOKEN} -n vault"); err != nil {
		return microerror.Mask(err)
	}

	return waitFor(cs.runningPodFunc("default", "app=vault"))
}

func (cs *clients) installCertOperator() error {
	certOperatorChartValuesEnv := os.ExpandEnv(certOperatorChartValues)
	if err := ioutil.WriteFile(certOperatorValuesFile, []byte(certOperatorChartValuesEnv), os.ModePerm); err != nil {
		return microerror.Mask(err)
	}
	if err := runCmd("helm registry install quay.io/giantswarm/cert-operator-chart@1.0.0-${CIRCLE_SHA1} -- -n cert-operator --values " + certOperatorValuesFile); err != nil {
		return microerror.Mask(err)
	}

	return waitFor(cs.certConfigFunc())
}

func runCmd(cmdStr string) error {
	cmdEnv := os.ExpandEnv(cmdStr)
	fields := strings.Fields(cmdEnv)
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	return cmd.Run()
}

func waitFor(f func() error) error {
	timeout := time.After(defaultTimeout * time.Second)
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	for {
		select {
		case <-timeout:
			ticker.Stop()
			return microerror.Mask(waitTimeoutError)
		case <-ticker.C:
			if err := f(); err == nil {
				return nil
			}
		}
	}
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
