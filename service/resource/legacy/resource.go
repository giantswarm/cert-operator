package legacy

import (
	"context"
	"fmt"

	"github.com/giantswarm/cert-operator/service/crt"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

const (
	Name = "legacy"
)

type Config struct {
	CrtService *crt.Service
	Logger     micrologger.Logger
}

func DefaultConfig() Config {
	return Config{
		CrtService: nil,
		Logger:     nil,
	}
}

type Resource struct {
	crtService *crt.Service
	logger     micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.CrtService == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.CrtService must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	newResource := &Resource{
		crtService: config.CrtService,
		logger: config.Logger.With(
			"resource", Name,
		),
	}

	return newResource, nil
}

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	cert := *obj.(*certificatetpr.CustomObject)
	r.logger.Log("debug", fmt.Sprintf("creating certificate '%s'", cert.Spec.CommonName))

	err := r.crtService.IssueAndWait(cert.Spec)
	if err != nil {
		return nil, microerror.Maskf(err, "could not issue cert")
	}

	r.logger.Log("info", fmt.Sprintf("certificate '%s' issued", cert.Spec.CommonName))

	return nil, nil
}

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	cert := *obj.(*certificatetpr.CustomObject)
	r.logger.Log("debug", fmt.Sprintf("deleting certificate '%s'", cert.Spec.CommonName))

	err := r.crtService.DeleteCertificateAndWait(cert.Spec)
	if err != nil {
		return nil, microerror.Maskf(err, "could not delete certificate")
	}

	r.logger.Log("info", fmt.Sprintf("certificate '%s' deleted", cert.Spec.CommonName))

	return nil, nil
}

func (r *Resource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	return nil, nil, nil, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	return nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	return nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
