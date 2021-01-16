package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Battery handles its label and notification level based on its remaining energy
type Battery struct {
	Percent  float64
	Charging bool
}

// Represents the battery in a human-readable strung such as 42%
func (b Battery) String() string {
	return fmt.Sprintf("%d%%", int64(math.Round(b.Percent*float64(100))))
}

// Label uses Font-Awesome 5 glyphs to create a visual label for the battery
func (b Battery) Label() string {
	if b.Charging {
		return ""
	}
	if b.Percent <= 0.1 {
		return ""
	}
	if b.Percent <= 0.3 {
		return ""
	}
	if b.Percent <= 0.5 {
		return ""
	}
	if b.Percent <= 0.8 {
		return ""
	}
	return ""
}

// NotificationLevel outputs notify-send levels according to remaining battery
func (b Battery) NotificationLevel() string {
	if b.Percent <= 0.2 {
		return "critical"
	}
	if b.Percent <= 0.3 {
		return "normal"
	}
	if b.Percent <= 0.4 {
		return "low"
	}
	return ""
}

// Color outputs the color of the i3 block based on `Battery.NotificationLevel`
func (b Battery) Color() string {
	switch b.NotificationLevel() {
	case "critical":
		return "#EF3340"
	case "normal":
		return "#FF8000"
	case "low":
		return "#EFEFEF"
	}
	return ""
}

// I3Block outputs the 4-lines data (http://vivien.github.io/i3blocks/#_format)
func (b Battery) I3Block() string {
	s := fmt.Sprintf("%s %s", b.Label(), b)
	return fmt.Sprintf("%s\n%s\n%s", s, s, b.Color())
}

// Runner to run commands such as `acpi` and `notify-send`
type Runner interface {
	Run(cmd []string) string
}

// AcpiRunner wraps `acpi` calls
type AcpiRunner struct {
}

// Run calls `acpi` with extra argumens `cmd`
func (x AcpiRunner) Run(cmd []string) string {
	out, err := exec.Command("acpi", cmd...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", out)
}

// GetBattery creates a `Battery` from an `acpi` output (`n` is the battery number)
func GetBattery(r Runner) (b Battery, err error) {
	out := r.Run([]string{})
	re := regexp.MustCompile(`Battery (\d+): [\w ]+, (\d+)%`)
	for _, m := range re.FindAllStringSubmatch(out, -1) {
		i, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		b.Percent = float64(i) / float64(100)
		b.Charging = strings.Contains(out, "Charging")
		if b.Percent != 0.0 {
			return b, err
		}
	}

	err = fmt.Errorf("Cannot find battery in:\n%s", out)
	return b, err
}

// Prints the i3block 4-lines
func main() {
	b, err := GetBattery(AcpiRunner{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.I3Block())
}
