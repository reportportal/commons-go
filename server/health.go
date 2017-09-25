package server

//HealthCheck specifies interface for health checks
type HealthCheck interface {
	Check() error
}
