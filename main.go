package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"

	deep "github.com/patrikeh/go-deep"
)

const (
	screenWidth  = 1200
	screenHeight = 700
	actionSpace  = 3
)

var (
	carImage, courseImage *ebiten.Image
	car                   *Car

	course *Course
	debug  bool // print lines for collisons detection
	hit    bool

	rays   []line
	state  []float64
	action int

	dqn *deep.Neural
)

func input(action int) {

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || action == 0 {
		car.turnLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || action == 2 {
		car.turnRight()
	}
	if action == 1 {

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
	action := -1

	if !debug {
		actions := dqn.Predict(state)
		_, action = max(actions)
	}

	input(action)
	car.acc()
	car.update()

	//CAR move
	opCar := &ebiten.DrawImageOptions{}
	// draw Car
	opCar.GeoM.Translate(-car.width/2, -car.height/2)
	opCar.GeoM.Rotate(car.rad())
	opCar.GeoM.Translate(car.x, car.y)

	carBox := carBoundings(car.x, car.y, car.width, car.height, car.rad())

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

	// Set Rays
	rays[0] = newRay(car.x, car.y, 1000, car.rad()-210*math.Pi/180)
	rays[1] = newRay(car.x, car.y, 1000, car.rad()-180*math.Pi/180)
	rays[2] = newRay(car.x, car.y, 1000, car.rad()-150*math.Pi/180)
	rays[3] = newRay(car.x, car.y, 1000, car.rad()-120*math.Pi/180)
	rays[4] = newRay(car.x, car.y, 1000, car.rad()-90*math.Pi/180)
	rays[5] = newRay(car.x, car.y, 1000, car.rad()-60*math.Pi/180)
	rays[6] = newRay(car.x, car.y, 1000, car.rad()-30*math.Pi/180)
	rays[7] = newRay(car.x, car.y, 1000, car.rad()-0*math.Pi/180)
	rays[8] = newRay(car.x, car.y, 1000, car.rad()+30*math.Pi/180)
	rays[9] = newRay(car.x, car.y, 1000, car.rad()-270*math.Pi/180)

	//RAYS check distance to Track
	for i, r := range rays {
		dMin := 1000.
		for _, w := range course.walls {
			x, y, intersect := intersection(w, r)
			if intersect {
				d := math.Sqrt(math.Pow(car.x-x, 2) + math.Pow(car.y-y, 2))
				if d < dMin {
					dMin = d
				}
			}
		}
		state[i] = dMin
	}

	if hit {
		resetCar()
		resetTrack()
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// DRAW
	screen.DrawImage(courseImage, nil)
	screen.DrawImage(carImage, opCar)
	// PRINT Informationen

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
		for _, r := range rays {
			ebitenutil.DrawLine(screen, r.X1, r.Y1, r.X2, r.Y2, color.RGBA{255, 255, 0, 150})
		}
	}

	return nil
}

func resetCar() {
	car = NewCar(620, 600, -80, 15, 25)
}

func resetTrack() {
	course.resetSections()
}

func max(actions []float64) (float64, int) {
	max := actions[0]
	index := 0
	for i, a := range actions {
		if a > max {
			max = a
			index = i
		}
	}
	return max, index
}

// init NN
func init() {
	data, _ := ioutil.ReadFile("./NNdump")
	dqn, _ = deep.Unmarshal(data)
}

// init Car
func init() {
	var err error
	carImage, _, err = ebitenutil.NewImageFromFile("assets/car.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	rays = make([]line, 10)
	state = make([]float64, 10)
	resetCar()
}

// init Course
func init() {
	course = NewCourse("assets/Sections.svg", "assets/Course.svg")
	var err error
	courseImage, _, err = ebitenutil.NewImageFromFile("assets/course.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	resetTrack()
}

func main() {

	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "CAR AI"); err != nil {
		log.Fatal(err)
	}
}
