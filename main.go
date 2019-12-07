package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 1200
	screenHeight = 700
)

var (
	carImage, courseImage *ebiten.Image
	car                   *Car
	course                *Course
	episode               int
	debug                 bool // print lines for collisons detection
	hit                   bool
)

func input() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		car.acc()
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		car.br()
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		car.turnLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		car.turnRight()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		resetCar()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if debug == true {
			debug = false
		} else {
			debug = true
		}
	}

}

// Game Loop
func update(screen *ebiten.Image) error {

	screen.DrawImage(courseImage, nil)

	//CAR
	input()
	car.update()
	opCar := &ebiten.DrawImageOptions{}
	opCar.GeoM.Translate(-car.width/2, -car.height/2)
	opCar.GeoM.Rotate(car.drawAngle())
	opCar.GeoM.Translate(car.x, car.y)
	screen.DrawImage(carImage, opCar)

	carBox := carBoundings(opCar.GeoM.Element(0, 2), opCar.GeoM.Element(1, 2), car.width, car.height, car.angle)
	hit = false

	//CAR BOUNDINGBOX
	for _, b := range carBox {
		// Check for collisions with Walls
		for _, w := range course.walls {
			if _, _, ok := intersection(b, w); ok {
				hit = true
			}
		}

		// Check Lap Time (collisions with sections)
		for i, cs := range course.courseSections {
			if _, _, ok := intersection(b, cs); ok {

				if i == 0 && course.finish(course.getLapTime()) {
					resetTrack()
				}
				course.sectionCleared(i)
			}
		}
	}

	if hit {
		resetCar()
		resetTrack()
	}

	// PRINT Informationen
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Episode: %v", episode), 20, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lap time %v", course.getLapTime()), 20, 40)

	for i, v := range course.lapTimes {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v : %v", i+1, v), screenWidth-180, 80+i*20)
	}

	//COURSE
	if debug {
		for _, w := range course.walls {
			ebitenutil.DrawLine(screen, w.X1, w.Y1, w.X2, w.Y2, color.RGBA{255, 0, 0, 255})
		}
		for _, b := range carBox {
			ebitenutil.DrawLine(screen, b.X1, b.Y1, b.X2, b.Y2, color.RGBA{0, 230, 64, 255})
		}
		for _, s := range course.courseSections {
			ebitenutil.DrawLine(screen, s.X1, s.Y1, s.X2, s.Y2, color.RGBA{0, 230, 64, 255})
		}
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	return nil
}

func resetCar() {
	episode++
	w, h := carImage.Size()
	car = NewCar(800, 560, -90, float64(w), float64(h))
}

func resetTrack() {
	course.resetSections()
}

// init Car
func init() {
	var err error
	carImage, _, err = ebitenutil.NewImageFromFile("assets/car.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	resetCar()
}

// init Course
func init() {

	var cS []line
	cS = append(cS,
		line{587, 613.5, 587, 508.5},
		line{272, 604, 283, 471.5},
		line{364.5, 270, 346.5, 133.5},
		line{631.5, 270, 835.5, 294.5},
		line{981, 489, 1146, 469.5})

	course = NewCourse(len(cS), cS, "assets/course.svg")

	var err error
	courseImage, _, err = ebitenutil.NewImageFromFile("assets/course.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	resetTrack()
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "CAR AI"); err != nil {
		log.Fatal(err)
	}
}
