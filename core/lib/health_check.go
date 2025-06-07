package lib

import "fmt"

const HEALTH_UP = "UP"
const HEALTH_DOWN = "DOWN"

type HealthChecker interface {
	Check() error
}

type MockUpHelthChecker struct{}

func (m MockUpHelthChecker) Check() error {
	return nil
}

type MockDownHelthChecker struct{}

func (m MockDownHelthChecker) Check() error {
	return fmt.Errorf("mock down health check error")
}
