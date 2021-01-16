package main

import "testing"

const full = `Battery 0: Unknown, 0%, rate information unavailable\n
Battery 1: Not charging, 100%`

const discharging = `Battery 0: Unknown, 0%, rate information unavailable
Battery 1: Discharging, 99%, 02:15:24 remaining`

const charging = `Battery 0: Unknown, 0%, rate information unavailable
Battery 1: Charging, 99%, 02:15:24 remaining`

type MockFullRunner struct{}

func (m MockFullRunner) Run(cmd []string) string {
	return full
}

type MockDischargingRunner struct{}

func (m MockDischargingRunner) Run(cmd []string) string {
	return discharging
}

type MockChargingRunner struct{}

func (m MockChargingRunner) Run(cmd []string) string {
	return charging
}

func TestGetBattery(t *testing.T) {
	for _, tc := range []struct {
		Runner Runner
		Expect Battery
	}{
		{MockFullRunner{}, Battery{1.0, false}},
		{MockDischargingRunner{}, Battery{0.99, false}},
		{MockChargingRunner{}, Battery{0.99, true}},
	} {
		b, err := GetBattery(tc.Runner)
		if err != nil {
			t.Errorf("Expected not erros, but got %s", err)
		}
		if b.Percent != tc.Expect.Percent {
			t.Errorf("Expected %f, but got %f", tc.Expect.Percent, b.Percent)
		}
		if b.Charging != tc.Expect.Charging {
			t.Errorf("Expected %v, but got %v", tc.Expect.Charging, b.Charging)
		}
	}

}

func TestBatteryLabel(t *testing.T) {
	for _, tc := range []struct {
		Battery Battery
		Expect  string
	}{
		{Battery{0.1, true}, ""},
		{Battery{0.1, false}, ""},
		{Battery{0.3, false}, ""},
		{Battery{0.5, false}, ""},
		{Battery{0.8, false}, ""},
		{Battery{0.9, false}, ""},
	} {
		if l := tc.Battery.Label(); l != tc.Expect {
			t.Errorf("Expected %s, but got %s", tc.Expect, l)
		}

	}
}

func TestBatteryNotificationLevel(t *testing.T) {
	for _, tc := range []struct {
		Battery Battery
		Expect  string
	}{
		{Battery{0.2, false}, "critical"},
		{Battery{0.3, false}, "normal"},
		{Battery{0.4, false}, "low"},
		{Battery{0.5, false}, ""},
	} {
		if l := tc.Battery.NotificationLevel(); l != tc.Expect {
			t.Errorf("Expected %s, but got %s", tc.Expect, l)
		}

	}
}

func TestBatteryColor(t *testing.T) {
	for _, tc := range []struct {
		Battery Battery
		Expect  string
	}{
		{Battery{0.2, false}, "#EF3340"},
		{Battery{0.3, false}, "#FF8000"},
		{Battery{0.4, false}, "#EFEFEF"},
		{Battery{0.5, false}, ""},
	} {
		if l := tc.Battery.Color(); l != tc.Expect {
			t.Errorf("Expected %s, but got %s", tc.Expect, l)
		}

	}
}
