package runner

type Controller struct {
}

type ControllerParams struct {
}

func NewController(params ControllerParams) *Controller {
	return &Controller{}
}

func (c *Controller) Run() error {
	return nil
}
