package godspeed

type DeploymentStrategy interface {
	Setup()
	Deploy() error
	Rollback() error
	Teardown()
}
