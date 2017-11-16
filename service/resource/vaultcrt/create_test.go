package vaultcrt

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultcrt/vaultcrttest"
	"github.com/giantswarm/vaultrole/vaultroletest"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func Test_Resource_VaultCrt_newCreateChange(t *testing.T) {
	testCases := []struct {
		Obj            interface{}
		CurrentState   interface{}
		DesiredState   interface{}
		ExpectedSecret *apiv1.Secret
	}{
		// Test 0 ensures a non-nil current state results in the create state to be
		// empty.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
					TTL:       "24h",
				},
			},
			CurrentState:   &apiv1.Secret{},
			DesiredState:   &apiv1.Secret{},
			ExpectedSecret: nil,
		},

		// Test 1 is the same 1 but with different content for the current state.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
					TTL:       "24h",
				},
			},
			CurrentState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Labels: map[string]string{
						"clusterID":        "foobar",
						"clusterComponent": "api",
					},
				},
				StringData: map[string]string{
					"ca":  "",
					"crt": "",
					"key": "",
				},
			},
			DesiredState:   &apiv1.Secret{},
			ExpectedSecret: nil,
		},

		// Test 2 ensures an empty current state results in a create state that
		// equals the desired state. NOTE that the secret data is extended with
		// actual certificate content, which in this case is some fake content from
		// the fake VaultCrt service.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
					TTL:       "24h",
				},
			},
			CurrentState: nil,
			DesiredState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Labels: map[string]string{
						"clusterID":        "foobar",
						"clusterComponent": "api",
					},
				},
				StringData: map[string]string{
					"ca":  "",
					"crt": "",
					"key": "",
				},
			},
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Labels: map[string]string{
						"clusterID":        "foobar",
						"clusterComponent": "api",
					},
				},
				StringData: map[string]string{
					"ca":  "test CA",
					"crt": "test crt",
					"key": "test key",
				},
			},
		},
	}

	var err error
	var newResource *Resource
	{
		c := DefaultConfig()

		c.CurrentTimeFactory = func() time.Time { return time.Time{} }
		c.K8sClient = fake.NewSimpleClientset()
		c.Logger = microloggertest.New()
		c.VaultCrt = vaultcrttest.New()
		c.VaultRole = vaultroletest.New()

		c.ExpirationThreshold = 24 * time.Hour
		c.Namespace = "default"

		newResource, err = New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, tc := range testCases {
		result, err := newResource.newCreateChange(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		secret := result.(*apiv1.Secret)
		if !reflect.DeepEqual(tc.ExpectedSecret, secret) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedSecret, secret)
		}
	}
}

func Test_Resource_VaultCrt_shouldCertBeRenewed(t *testing.T) {
	testCases := []struct {
		CurrentTime    time.Time
		Secret         *apiv1.Secret
		TTL            time.Duration
		Threshold      time.Duration
		ErrorMatcher   func(err error) bool
		ExpectedResult bool
	}{
		// Test 0 ensures that a zero value input results in an error.
		{
			CurrentTime:    time.Time{},
			Secret:         &apiv1.Secret{},
			TTL:            0,
			Threshold:      0,
			ErrorMatcher:   IsMissingAnnotation,
			ExpectedResult: false,
		},

		// Test 1 ensures using an update timestamp which is after the current time
		// does not cause certificates to be renewed.
		{
			CurrentTime: time.Unix(9, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(10, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(11, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(14, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(15, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(16, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(17, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(20, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(21, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(24, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(25, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(26, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
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
			CurrentTime: time.Unix(345322, 0).In(time.UTC),
			Secret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Annotations: map[string]string{
						UpdateTimestampAnnotation: time.Unix(10, 0).In(time.UTC).Format(UpdateTimestampLayout),
					},
				},
			},
			TTL:            10 * time.Second,
			Threshold:      5 * time.Second,
			ErrorMatcher:   nil,
			ExpectedResult: true,
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
			c.VaultRole = vaultroletest.New()

			c.ExpirationThreshold = 24 * time.Hour
			c.Namespace = "default"

			newResource, err = New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		result, err := newResource.shouldCertBeRenewed(tc.Secret, tc.TTL, tc.Threshold)
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
