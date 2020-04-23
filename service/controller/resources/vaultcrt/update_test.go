package vaultcrt

import (
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/certs"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultcrt/vaultcrttest"
	apiv1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_Resource_VaultCrt_shouldCertBeRenewed_expiration(t *testing.T) {
	testCases := []struct {
		CurrentTime    time.Time
		CustomObject   v1alpha1.CertConfig
		Secret         *apiv1.Secret
		TTL            time.Duration
		Threshold      time.Duration
		ErrorMatcher   func(err error) bool
		ExpectedResult bool
	}{
		// Test 0 ensures that a zero value input results in an error.
		{
			CurrentTime:    time.Time{},
			CustomObject:   v1alpha1.CertConfig{},
			Secret:         &apiv1.Secret{},
			TTL:            0,
			Threshold:      0,
			ErrorMatcher:   IsMissingAnnotation,
			ExpectedResult: false,
		},

		// Test 1 ensures using an update timestamp which is after the current time
		// does not cause certificates to be renewed.
		{
			CurrentTime:  time.Unix(9, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 2 ensures using an update timestamp which is equal to the current
		// time does not cause certificates to be renewed.
		{
			CurrentTime:  time.Unix(10, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 3 ensures using an update timestamp which is before the current time
		// does not cause certificates to be renewed.
		{
			CurrentTime:  time.Unix(11, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 4 is the same as 3 but with a different current time.
		{
			CurrentTime:  time.Unix(14, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 5 is the same as 3 but with a different current time.
		{
			CurrentTime:  time.Unix(15, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 6 is the same as 3 but with a different current time where the
		// certificates are expected to be renewed.
		//
		// NOTE that the tests move on the timeline of the current time and the
		// expected result flips here.
		//
		//     (update timestamp + TTL - threshold < current time) == true
		//     (10               + 10  - 5         < 16          ) == true
		//
		{
			CurrentTime:  time.Unix(16, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 7 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(17, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 8 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(20, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 9 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(21, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 10 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(24, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 11 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(25, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 12 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(26, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 13 is the same as 6 but with a different current time.
		{
			CurrentTime:  time.Unix(345322, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 14 is the same as 13 but with a service account cert config which
		// should result into not being regenerated at all.
		{
			CurrentTime: time.Unix(345322, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterComponent: string(certs.ServiceAccountCert),
					},
				},
			},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 15 is the same as 14 but with a the general flag being used to
		// disable regeneration.
		{
			CurrentTime: time.Unix(345322, 0).In(time.UTC),
			CustomObject: v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						DisableRegeneration: true,
					},
				},
			},
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "hash",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},
	}

	for i, tc := range testCases {
		var err error
		var newResource *Resource
		{
			c := DefaultConfig()

			c.CurrentTimeFactory = func() time.Time { return tc.CurrentTime }
			c.K8sClient = fake.NewSimpleClientset()
			c.Logger = microloggertest.New()
			c.VaultCrt = vaultcrttest.New()

			c.ExpirationThreshold = 24 * time.Hour
			c.Namespace = "default"

			newResource, err = New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		result, err := newResource.shouldCertBeRenewed(tc.CustomObject, tc.Secret, tc.Secret, tc.TTL, tc.Threshold)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if tc.ExpectedResult != result {
				t.Fatalf("case %d expected %t got %t", i, tc.ExpectedResult, result)
			}
		}
	}
}

func Test_Resource_VaultCrt_shouldCertBeRenewed_hash(t *testing.T) {
	testCases := []struct {
		CustomObject   v1alpha1.CertConfig
		CurrentSecret  *apiv1.Secret
		DesiredSecret  *apiv1.Secret
		ErrorMatcher   func(err error) bool
		ExpectedResult bool
	}{
		// Test 0 ensures that a zero value input results in an error.
		{
			CustomObject:   v1alpha1.CertConfig{},
			CurrentSecret:  &apiv1.Secret{},
			DesiredSecret:  &apiv1.Secret{},
			ErrorMatcher:   IsMissingAnnotation,
			ExpectedResult: false,
		},

		// Test 1 ensures using different config hashes the secret should be
		// updated.
		{
			CustomObject: v1alpha1.CertConfig{},
			CurrentSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "current",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			DesiredSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "desired",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 2 ensures using equal config hashes the secret should not be
		// updated.
		{
			CustomObject: v1alpha1.CertConfig{},
			CurrentSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "same",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			DesiredSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "same",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 3 ensures having no config hash annotation for the current state and
		// having a config hash value for the desired state results in updating the
		// secret.
		{
			CustomObject: v1alpha1.CertConfig{},
			CurrentSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			DesiredSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "new",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 4 is the same as 3 but with a service account cert config which
		// should result into not being regenerated at all.
		{
			CustomObject: v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterComponent: string(certs.ServiceAccountCert),
					},
				},
			},
			CurrentSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			DesiredSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "new",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 5 is the same as 4 but with a the general flag being used to disable
		// regeneration.
		{
			CustomObject: v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						DisableRegeneration: true,
					},
				},
			},
			CurrentSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			DesiredSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						ConfigHashAnnotation:      "new",
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},
	}

	for i, tc := range testCases {
		var err error
		var newResource *Resource
		{
			c := DefaultConfig()

			c.CurrentTimeFactory = func() time.Time { return time.Time{} }
			c.K8sClient = fake.NewSimpleClientset()
			c.Logger = microloggertest.New()
			c.VaultCrt = vaultcrttest.New()

			c.ExpirationThreshold = 24 * time.Hour
			c.Namespace = "default"

			newResource, err = New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		result, err := newResource.shouldCertBeRenewed(tc.CustomObject, tc.CurrentSecret, tc.DesiredSecret, 10*time.Second, 5*time.Second)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if tc.ExpectedResult != result {
				t.Fatalf("case %d expected %t got %t", i, tc.ExpectedResult, result)
			}
		}
	}
}
