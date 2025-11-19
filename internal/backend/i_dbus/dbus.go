package i_dbus

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	dbus "github.com/godbus/dbus/v5"
	"github.com/moneronodo/sshui/internal/base"
	dbus_model "github.com/moneronodo/sshui/internal/model/dbus"
)

func DbusSignal(s *dbus.Signal) dbus_model.DbusSignal {
	if !strings.HasPrefix(s.Name, "com.moneronodo.embeddedInterface") {
		return nil
	}
	switch strings.Split(s.Name, "com.moneronodo.embeddedInterface.")[1] {
	case "factoryResetStarted":
		return dbus_model.FactoryResetStarted{}
	case "factoryResetCompleted":
		return dbus_model.FactoryResetCompleted{}
	case "factoryResetRequested":
		return dbus_model.FactoryResetRequested{}
	case "powerButtonPressDetected":
		return dbus_model.PowerButtonPressDetected{}
	case "powerButtonReleaseDetected":
		return dbus_model.PowerButtonReleaseDetected{}
	case "moneroLWSListAccountsCompleted":
		return dbus_model.MoneroLWSListAccountsCompleted{}
	case "moneroLWSListRequestsCompleted":
		return dbus_model.MoneroLWSListRequestsCompleted{}
	case "moneroLWSAccountAdded":
		return dbus_model.MoneroLWSAccountAdded{}
	case "connectionStatusChanged":
		return dbus_model.ConnectionStatusChanged{}
	case "startRecoveryNotification":
		return dbus_model.StartRecoveryNotification{
			Message: s.Body[0].(string),
		}
	case "serviceManagerNotification":
		return dbus_model.ServiceManagerNotification{
			Message: s.Body[0].(string),
		}
	case "hardwareStatusReadyNotification":
		return dbus_model.HardwareStatusReadyNotification{
			Message: s.Body[0].(string),
		}
	case "serviceStatusReadyNotification":
		return dbus_model.ServiceStatusReadyNotification{
			Message: s.Body[0].(string),
		}
	case "passwordChangeStatus":
		return dbus_model.PasswordChangeStatus{
			Status: s.Body[0].(int),
		}
	default:
		return nil
	}
}

func Signals(prog *tea.Program) {
	conn, err := dbus.SystemBus()
	if err != nil {
		spew.Fprintln(base.Dump, "Dbus: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath("/com/monero/nodo"),
		dbus.WithMatchInterface("com.moneronodo.embeddedInterface"),
	); err != nil {
		spew.Fprintln(base.Dump, "Dbus: ", err)
		os.Exit(1)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		sig := DbusSignal(v)
		msg := dbus_model.DbusSignalMsg{
			Signal: sig,
		}
		prog.Send(msg)
	}
}

func Call(notification string, args ...any) {
	spew.Fprintf(base.Dump, "Call %s\n", notification)
	conn, err := dbus.SystemBus()
	if err != nil {
		spew.Fdump(base.Dump, err)
	}
	defer conn.Close()

	obj := conn.Object("com.monero.nodo", "/com/monero/nodo")
	call := obj.Call("com.moneronodo.embeddedInterface."+notification, 0, args...)
	if call.Err != nil {
		spew.Fdump(base.Dump, call.Err)
	} else {
		spew.Fdump(base.Dump, call)
	}
}
