package main

import "testing"

const full = `Battery 0: Unknown, 0%, rate information unavailable\n
Battery 1: Not charging, 100%`

const discharging = `Battery 0: Unknown, 0%, rate information unavailable
Battery 1: Discharging, 99%, 02:15:24 remaining`

type MockFullRunner struct {
}

func (m MockFullRunner) Run(cmd []string) string {
	return full
}

type MockDischargingRunner struct {
}

func (m MockDischargingRunner) Run(cmd []string) string {
	return discharging
}

func TestGetBattery(t *testing.T) {
	for _, tc := range []struct {
		Runner Runner
		Expect Battery
	}{
		{MockFullRunner{}, Battery{1.0}},
		{MockDischargingRunner{}, Battery{0.99}},
	} {
		if b := GetBattery(tc.Runner, 1); b != tc.Expect {
			t.Errorf("Expected 100%%, but got %s", b)
		}
	}

}

func TestBatteryLabel(t *testing.T) {
	for _, tc := range []struct {
		Battery Battery
		Expect  string
	}{
		{Battery{0.1}, ""},
		{Battery{0.3}, ""},
		{Battery{0.5}, ""},
		{Battery{0.8}, ""},
		{Battery{0.9}, ""},
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
		{Battery{0.2}, "critical"},
		{Battery{0.3}, "normal"},
		{Battery{0.4}, "low"},
		{Battery{0.5}, ""},
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
		{Battery{0.2}, "#EF3340"},
		{Battery{0.3}, "#FF8000"},
		{Battery{0.4}, "#EFEFEF"},
		{Battery{0.5}, ""},
	} {
		if l := tc.Battery.Color(); l != tc.Expect {
			t.Errorf("Expected %s, but got %s", tc.Expect, l)
		}

	}
}
