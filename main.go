package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
)

// Battery handles its label and notification level based on its remaining energy
type Battery struct {
	Percent float64
}

// Represents the battery in a human-readable strung such as 42%
func (b Battery) String() string {
	return fmt.Sprintf("%d%%", int64(math.Round(b.Percent*float64(100))))
}

// Label uses Font-Awesome 5 glyphs to create a visual label for the battery
func (b Battery) Label() string {
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
	if b.Percent <= 0.1 {
		return "critical"
	}
	if b.Percent <= 0.2 {
		return "normal"
	}
	if b.Percent <= 0.3 {
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
func GetBattery(r Runner, n int) Battery {
	re := regexp.MustCompile(`Battery (\d+): [\w ]+, (\d+)%`)
	out := r.Run([]string{})
	s := strconv.Itoa(n)
	b := Battery{0.0}

	for _, m := range re.FindAllStringSubmatch(out, -1) {
		if m[1] != s {
			continue
		}

		i, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}
		b.Percent = float64(i) / float64(100)
	}

	if b.Percent == 0.0 {
		log.Fatalf("Cannot find battery %d in:\n%s", n, out)
	}
	return b
}

// Prints the i3block 4-lines
func main() {
	fmt.Println(GetBattery(AcpiRunner{}, 1).I3Block())
}
