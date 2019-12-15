package main

import (
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Course the course of the game
type Course struct {
	walls           []line
	courseSections  []line
	lapTimes        []time.Duration
	sections        int
	sectionsCounter []bool
	sTime           time.Time
	startTime       time.Time
}

// NewCourse constructor for Course returns *Course
func NewCourse(sectionsFilePath, filePath string) *Course {
	sections := sectionsFromSVG(sectionsFilePath)

	c := &Course{
		walls:           []line{},                    // to check for collisons
		courseSections:  sections,                    // sections line check for collisons
		lapTimes:        []time.Duration{},           // Leaderboard
		sections:        len(sections),               // count of section
		sectionsCounter: make([]bool, len(sections)), // check all section have passed No cheating :D
	}

	c.initWithSVG(filePath)
	return c
}

func (c *Course) finish(lapTime time.Duration) bool {
	// check for all sections
	for _, s := range c.sectionsCounter {
		if s != true {
			return false
		}
	}

	// Update Leaderboard
	c.lapTimes = append(c.lapTimes, lapTime)
	if len(c.lapTimes) > 1 {
		sort.Slice(c.lapTimes, func(i, j int) bool { return c.lapTimes[i] < c.lapTimes[j] })
	}
	if len(c.lapTimes) > 5 {
		c.lapTimes = c.lapTimes[:5]
	}

	return true
}

// section passed
func (c *Course) sectionCleared(section int) bool {
	if section == 0 {
		c.startTime = time.Now()
	}

	if c.sectionsCounter[section] == false {
		c.sectionsCounter[section] = true
		return true
	}

	return false
}

// getLapTime returns the time since Lap start
func (c *Course) getLapTime() time.Duration {
	if c.sectionsCounter[0] {
		return time.Now().Sub(c.startTime).Truncate(10000000)
	}
	return 0
}

// reset Sections
func (c *Course) resetSections() {
	c.sectionsCounter = make([]bool, c.sections)

}

func (c *Course) initWithSVG(filePath string) {

	var lines []line
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Refactor :D
	for _, v := range strings.Split(string(file), "\n") {
		if strings.Contains(v, "polygon") {
			// Remove everything beside values in points="..." somehow...
			points := strings.Split(strings.Split(strings.TrimLeft(strings.TrimRight(v, "></polygon>"), "<polygon "), "points=")[1], " ")
			points[0] = strings.TrimLeft(points[0], "\"")
			points[len(points)-1] = strings.TrimRight(points[len(points)-1], "\"")

			// number of Lines = Points in Polygon / 2
			for i := 0; i < len(points)/2; i++ {
				if i == len(points)/2-1 {
					// last line end point = first line start point
					lines = append(lines, getPointsForLine(i+i, i+i+1, 0, 1, &points))
				} else {
					lines = append(lines, getPointsForLine(i+i, i+i+1, i+i+2, i+i+3, &points))
				}
			}
		}
	}
	c.walls = append(c.walls, lines...)
}

func sectionsFromSVG(filePath string) []line {
	var lines []line
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range strings.Split(string(file), "\n") {
		if strings.Contains(v, "line") {
			lineStr := strings.TrimLeft(strings.TrimRight(v, "></line>"), "<line ")
			points := strings.Split(lineStr, " ")
			pointBuff := make([]float64, 4)
			for i, v := range points {
				if strings.Contains(v, "x") || strings.Contains(v, "y") {
					pointBuff[i], _ = strconv.ParseFloat(strings.TrimLeft(strings.TrimRight(v[4:], "\""), "\""), 64)

				}
			}
			lines = append(lines, line{pointBuff[0], pointBuff[1], pointBuff[2], pointBuff[3]})
		}
	}
	return lines
}

// getPointsForLine helper for initWithSVG()
func getPointsForLine(x1, y1, x2, y2 int, points *[]string) line {
	X1, _ := strconv.ParseFloat(string((*points)[x1]), 64)
	Y1, _ := strconv.ParseFloat(string((*points)[y1]), 64)
	X2, _ := strconv.ParseFloat(string((*points)[x2]), 64)
	Y2, _ := strconv.ParseFloat(string((*points)[y2]), 64)
	return line{X1, Y1, X2, Y2}
}
