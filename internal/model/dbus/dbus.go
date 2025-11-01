package dbus

type DbusSignalMsg struct {
	Signal DbusSignal
}

type DbusSignal any

type StartRecoveryNotification struct {
	Message string
}
type ServiceManagerNotification struct {
	Message string
}
type HardwareStatusReadyNotification struct {
	Message string
}
type ServiceStatusReadyNotification struct {
	Message string
}
type PasswordChangeStatus struct {
	Status int
}

type FactoryResetStarted struct{}
type FactoryResetCompleted struct{}
type FactoryResetRequested struct{}
type PowerButtonPressDetected struct{}
type PowerButtonReleaseDetected struct{}
type MoneroLWSListAccountsCompleted struct{}
type MoneroLWSListRequestsCompleted struct{}
type MoneroLWSAccountAdded struct{}
type ConnectionStatusChanged struct{}
