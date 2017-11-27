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

	"github.com/cenk/backoff"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
)

const (
	certOperatorValuesFile = "/tmp/cert-operator-install.yaml"
	defaultDeadline        = 15
	// values required by cert-operator-chart, the envirnment variables will
	// be expanded before writing the contents to a file.
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

var cs kubernetes.Interface

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var v int
	var err error
	cs, err = getK8sClient()
	if err != nil {
		v = 1
		log.Printf("unexpected error: %v\n", err)
	}

	if err := setUp(cs); err != nil {
		v = 1
		log.Printf("unexpected error: %v\n", err)
	}

	if v == 0 {
		v = m.Run()
	}

	tearDown(cs)

	os.Exit(v)
}

func TestSecretsAreCreated(t *testing.T) {
	err := runCmd("helm registry install quay.io/giantswarm/cert-resource-lab-chart:stable -- -n cert-resource-lab --set commonDomain=${COMMON_DOMAIN} --set clusterName=${CLUSTER_NAME}")
	if err != nil {
		t.Errorf("could not install cert-resource-lab, %v", err)
	}

	secretName := fmt.Sprintf("%s-api", os.Getenv("CLUSTER_NAME"))
	err = waitFor(secretFunc(cs, "default", secretName))
	if err != nil {
		t.Errorf("could not find expected secret: %v", err)
	}
}

func setUp(cs kubernetes.Interface) error {
	if err := createGSNamespace(cs); err != nil {
		return err
	}

	if err := installVault(cs); err != nil {
		return err
	}

	if err := installCertOperator(cs); err != nil {
		return err
	}
	return nil
}

func tearDown(cs kubernetes.Interface) {
	runCmd("helm delete vault --purge")
	runCmd("helm delete cert-operator --purge")
	runCmd("helm delete cert-resource-lab --purge")
	cs.CoreV1().Namespaces().Delete("giantswarm", &metav1.DeleteOptions{})
	cs.ExtensionsV1beta1().ThirdPartyResources().Delete(certificatetpr.Name, &metav1.DeleteOptions{})
}

func createGSNamespace(cs kubernetes.Interface) error {
	// check if the namespace already exists
	_, err := cs.CoreV1().Namespaces().Get("giantswarm", metav1.GetOptions{})
	if err == nil {
		return nil
	}

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "giantswarm",
		},
	}
	_, err = cs.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	return waitFor(activeNamespaceFunc(cs, "giantswarm"))
}

func installVault(cs kubernetes.Interface) error {
	if err := runCmd("helm registry install quay.io/giantswarm/vaultlab-chart:stable -- --set vaultToken=${VAULT_TOKEN} -n vault"); err != nil {
		return microerror.Mask(err)
	}

	return waitFor(runningPodFunc(cs, "default", "app=vault"))
}

func installCertOperator(cs kubernetes.Interface) error {
	certOperatorChartValuesEnv := os.ExpandEnv(certOperatorChartValues)
	if err := ioutil.WriteFile(certOperatorValuesFile, []byte(certOperatorChartValuesEnv), os.ModePerm); err != nil {
		return err
	}
	if err := runCmd("helm registry install quay.io/giantswarm/cert-operator-chart@1.0.0-${CIRCLE_SHA1} -- -n cert-operator --values " + certOperatorValuesFile); err != nil {
		return err
	}

	return waitFor(tprFunc(cs, certificatetpr.Name))
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
	timeout := time.After(defaultDeadline * time.Second)
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	for {
		select {
		case <-timeout:
			ticker.Stop()
			return fmt.Errorf("wait timeout")
		case <-ticker.C:
			if err := f(); err == nil {
				return nil
			}
		}
	}
}

func runningPodFunc(cs kubernetes.Interface, namespace, labelSelector string) func() error {
	return func() error {
		pods, err := cs.CoreV1().Pods(namespace).List(metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return microerror.Mask(err)
		}
		if len(pods.Items) != 1 {
			return fmt.Errorf("unexpected number of pods")
		}
		pod := pods.Items[0]
		phase := pod.Status.Phase
		if phase != v1.PodRunning {
			return fmt.Errorf("unexpected pod status phase " + string(phase))
		}
		return nil
	}
}

func activeNamespaceFunc(cs kubernetes.Interface, name string) func() error {
	return func() error {
		ns, err := cs.CoreV1().Namespaces().Get(name, metav1.GetOptions{})

		if err != nil {
			return microerror.Mask(err)
		}

		phase := ns.Status.Phase
		if phase != v1.NamespaceActive {
			return fmt.Errorf("unexpected ns status phase " + string(phase))
		}

		return nil
	}
}

func secretFunc(cs kubernetes.Interface, namespace, secretName string) func() error {
	return func() error {
		_, err := cs.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
		return microerror.Mask(err)
	}
}

func tprFunc(cs kubernetes.Interface, tprName string) func() error {
	return func() error {
		// FIXME: use proper clientset call when apiextensions are in place,
		// `cs.ExtensionsV1beta1().ThirdPartyResources().Get(tprName, metav1.GetOptions{})` finding
		// the tpr is not enough for being able to create a tpo.
		return runCmd("kubectl get certificate")
	}
}
