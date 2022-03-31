package vaultcrt

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/v6/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultcrt/vaultcrttest"
	apiv1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	fakectrl "sigs.k8s.io/controller-runtime/pkg/client/fake" //nolint:staticcheck // v0.6.4 has a deprecation on pkg/client/fake that was removed in later versions

	"github.com/giantswarm/cert-operator/pkg/project"
)

func Test_Resource_VaultCrt_GetDesiredState(t *testing.T) {
	testCases := []struct {
		Obj            interface{}
		ExpectedSecret *apiv1.Secret
	}{
		// Test 0 ensures the desired state is always the same placeholder state.
		{
			Obj: &v1alpha1.CertConfig{
				ObjectMeta: apismetav1.ObjectMeta{
					Labels: map[string]string{
						"cert-operator.giantswarm.io/version": project.Version(),
					},
				},
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID:        "foobar",
						ClusterComponent: "api",
					},
				},
			},
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Annotations: map[string]string{
						ConfigHashAnnotation:      "001ad3d32b3f7d64e00ec0a3d5592fbb791849c2",
						UpdateTimestampAnnotation: (time.Time{}).Format(UpdateTimestampLayout),
					},
					Labels: map[string]string{
						"giantswarm.io/cluster":               "foobar",
						"giantswarm.io/certificate":           "api",
						"cert-operator.giantswarm.io/version": project.Version(),
					},
				},
				StringData: map[string]string{
					"ca":  "",
					"crt": "",
					"key": "",
				},
			},
		},

		// Test 1 is the same as 0 but with a different custom object.
		{
			Obj: &v1alpha1.CertConfig{
				ObjectMeta: apismetav1.ObjectMeta{
					Labels: map[string]string{
						"cert-operator.giantswarm.io/version": project.Version(),
					},
				},
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID:        "al9qy",
						ClusterComponent: "worker",
					},
				},
			},
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
					Annotations: map[string]string{
						ConfigHashAnnotation:      "4bf7b5296ba01161f182de54b243e1400ae6660e",
						UpdateTimestampAnnotation: (time.Time{}).Format(UpdateTimestampLayout),
					},
					Labels: map[string]string{
						"giantswarm.io/cluster":               "al9qy",
						"giantswarm.io/certificate":           "worker",
						"cert-operator.giantswarm.io/version": project.Version(),
					},
				},
				StringData: map[string]string{
					"ca":  "",
					"crt": "",
					"key": "",
				},
			},
		},
	}

	var err error
	var newResource *Resource
	{
		c := DefaultConfig()
		scheme := runtime.NewScheme()
		_ = capi.AddToScheme(scheme)

		c.CurrentTimeFactory = func() time.Time { return time.Time{} }
		c.K8sClient = fake.NewSimpleClientset()
		c.CtrlClient = fakectrl.NewClientBuilder().WithScheme(scheme).Build()
		c.Logger = microloggertest.New()
		c.VaultCrt = vaultcrttest.New()

		c.ExpirationThreshold = 24 * time.Hour
		c.Namespace = "default"

		newResource, err = New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, tc := range testCases {
		result, err := newResource.GetDesiredState(context.TODO(), tc.Obj)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		secret := result.(*apiv1.Secret)
		if !reflect.DeepEqual(tc.ExpectedSecret, secret) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedSecret, secret)
		}
	}
}
