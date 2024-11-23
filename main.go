package main

import (
	"fmt"
	"math"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	width    = 250
	height   = 100
	paused   = "paused"
	stopped  = "stopped"
	running  = "running"
	finished = "finished"
	sixty    = float32(60)
)

func main() {
	rl.InitWindow(width, height, "pmdr")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	sound := rl.LoadSound("jingles_NES13.ogg")
	rl.SetTargetFPS(60)

	cv := float32(0.0)
	max := float32(time.Second * 3 / time.Second)
	done := make(chan bool)

	state := stopped
	halfsize := float32(width / 2)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		gui.ProgressBar(rl.Rectangle{0, 0, width, 45}, "", "", cv, 0, max)
		rl.DrawText(timeFormat(cv), int32(halfsize-25), 10, 20, rl.Black)

		switch state {
		case running:
			if (gui.Button(rl.Rectangle{0, 45, halfsize, 55}, "Pause")) {
				state = paused
				done <- true
			}
			if (gui.Button(rl.Rectangle{halfsize, 45, halfsize, 55}, "Stop")) {
				state = stopped
				done <- true
				cv = 0
			}
		case stopped:
			if (gui.Button(rl.Rectangle{halfsize, 45, halfsize, 55}, "Run 30m")) {
				max = float32(time.Minute * 30 / time.Second)
				state = running
				go tickLoop(done, &cv)
			}
			if (gui.Button(rl.Rectangle{0, 45, halfsize, 55}, "Run 1h")) {
				max = float32(time.Minute * 60 / time.Second)
				state = running
				go tickLoop(done, &cv)
			}
		case paused:
			if (gui.Button(rl.Rectangle{0, 45, halfsize, 55}, "Run")) {
				state = running
				go tickLoop(done, &cv)
			}
			if (gui.Button(rl.Rectangle{halfsize, 45, halfsize, 55}, "Stop")) {
				state = stopped
				cv = float32(0)
			}
		case finished:
			if !rl.IsSoundPlaying(sound) {
				rl.PlaySound(sound)
			}
			if (gui.Button(rl.Rectangle{0, 45, width, 55}, "Off")) {
				state = stopped
			}
		}

		if cv >= max {
			done <- true
			state = finished
			cv = float32(0)
		}

		rl.ClearBackground(rl.RayWhite)
		rl.EndDrawing()
	}

	// if user hits the close window button, clean up the running for loop
	if state == running {
		done <- true
	}
	rl.UnloadSound(sound)
	rl.CloseAudioDevice()
}

func timeFormat(cv float32) string {
	minutes := math.Floor(float64(cv / sixty))
	seconds := cv - float32(minutes*60)
	return fmt.Sprintf("%d:%.2f", int(minutes), seconds)
}

func tickLoop(done chan bool, currentValue *float32) {
	ticker := time.NewTicker(25 * time.Millisecond)
	tick := float32(.025)
	for {
		select {
		case <-done:
			ticker.Stop()
		case <-ticker.C:
			cv := currentValue
			ti := &tick
			*currentValue = *cv + *ti
			// this contiue is very important as it keeps the for loop from
			// hitting the break statement until the done channel is fired
			continue
		}
		break
	}
}
