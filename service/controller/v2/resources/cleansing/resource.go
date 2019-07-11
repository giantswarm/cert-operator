package cleansing

type Resource struct {
}

func (Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	panic("implement me")
}

func (Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	panic("implement me")
}

func (Resource) Name() string {
	panic("implement me")
}
