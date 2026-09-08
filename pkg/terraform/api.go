package terraform

// API interface to interact with Terraform.
type API interface {
	Init(env string) error
	Run(env string, command string, extraArgs []string) error
}
