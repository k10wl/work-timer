package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/andybrewer/mack"
)

const (
	DefaultDuration = 45

	Full  = "█"
	Empty = "─"

	ProgressFallbackSize = 40

	ClearLine = "\033[1A\033[K"

	SayText = "You are doing great, now its time to rest! Go stretch your bums."
)

func main() {
	d := duration()
	clearFn := clearConsole()
	go countdown(d, clearFn)
	time.Sleep(time.Minute * time.Duration(d))
	notify()
}

func duration() int {
	args := os.Args
	if len(args) < 2 {
		return DefaultDuration
	}
	d := args[1]
	dur, err := strconv.Atoi(d)
	if err != nil {
		panic("Cannot parse duration")
	}
	return dur
}

func countdown(d int, clearFn func()) {
	upTo := d * 60
	fmt.Printf("\nConcentration time!\n\n")
	clearFn()
	for i := 0; i <= upTo; i++ {
		(func(cur int) {
			if cur > 0 {
				fmt.Print(strings.Repeat(ClearLine, 4))
				clearFn()
			}
			percentage := getPercentage(cur, upTo)
			fmt.Println("Time remaining:", formatTime(time.Second*time.Duration(upTo-cur)))
			fmt.Println(drawProgress(percentage))
			fmt.Println(drawProgress(percentage))
			fmt.Println()
			time.Sleep(time.Second)
		})(i)
	}
}

func notify() {
	mack.Say(SayText)
}

func getPercentage(num int, total int) float64 {
	percentage := float64(num) / float64(total) * 100
	return math.Round(percentage*100) / 100
}

func formatPercentage(num float64) string {
	return fmt.Sprintf("%.2f%%", num)
}

func drawProgress(percentage float64) string {
	progress := "["
	timerProgress := percentage / 100
	width := consoleWidth() - 2
	for i := 0; i < width; i++ {
		(func(cur int) {
			drawProgress := float64(cur) / float64(width)
			if drawProgress > timerProgress {
				progress += Empty
				return
			}
			progress += Full
		})(i)
	}
	progress += "]"
	return progress
}

func formatTime(duration time.Duration) string {
	formater := "04:05"
	if duration.Hours() > 1 {
		formater = "15:04:05"
	}
	time := time.Unix(0, 0).UTC().Add(duration).Format(formater)
	return time

	return duration.String()
}

func consoleWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Cannot read stty, returning fallback size")
		return ProgressFallbackSize
	}

	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")

	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		fmt.Println("Cannot convert width to int, returning fallback size")
		return ProgressFallbackSize
	}
	return width
}

func clearConsole() func() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	return func() {
		c.Run()
	}
}
