package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const MAX_DRIVETIME_MINUTES = float64(12 * 60)

type Route struct {
	Number int64
	StartX float64
	StartY float64
	EndX   float64
	EndY   float64
}

type Driver struct {
	Routes []Route
}

func main() {

	if len(os.Args) < 1 {
		panic(errors.New("please provide a path to a valid load file"))
	}
	arg := os.Args[1]
	f, err := os.Open(arg)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var routes []Route

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "loadNumber") {
			//ignore the header. later we may look at it if columns are dynamic
		} else {

			route := parseRoute(line)
			routes = append(routes, route)
		}

	}

	drivers := make([]Driver, 1)

	// Determine routes
	for _, route := range routes {

		for i, driver := range drivers {
			// fmt.Println("driver ", i)
			if driver.canTakeRoute(route) {
				drivers[i].acceptRoute(route)
				break
			} else if i == len(drivers)-1 {
				t := Driver{}
				t.acceptRoute(route)
				drivers = append(drivers, t)
			}
		}

	}

	// Output results
	for _, driver := range drivers {

		var sb strings.Builder
		sb.WriteString("[")
		for i, route := range driver.Routes {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strconv.FormatInt(route.Number, 10))
		}
		sb.WriteString("]")
		fmt.Println(sb.String())
	}

}

func (d Driver) canTakeRoute(r Route) bool {

	// When the driver currently has no routes/loads
	if len(d.Routes) < 1 && r.calculateRouteTime() < MAX_DRIVETIME_MINUTES {
		return true
	}

	currentTime := d.getTotalRoutesTime(false)

	// The time of the new load pickups -> delivery + the time going from the previous dropoff to the proposed pickup + the time from the proposed dropoff back to the home base
	proposedAddedTime := r.calculateRouteTime() + calculateRouteTime(d.Routes[len(d.Routes)-1].EndX, d.Routes[len(d.Routes)-1].EndY, r.StartX, r.StartY) + calculateRouteTime(r.EndX, r.EndY, 0.0, 0.0)

	return (proposedAddedTime + currentTime) < MAX_DRIVETIME_MINUTES
}

func (d Driver) getTotalRoutesTime(includeReturnTrip bool) float64 {

	currentTime := 0.0

	for i, route := range d.Routes {

		//if first add time from origin
		if i == 0 {
			currentTime += calculateRouteTime(0.0, 0.0, route.StartX, route.StartY)
		}

		// add load time
		currentTime += route.calculateRouteTime()

		// if end, add time from last dropoff to origin
		if i == len(d.Routes)-1 && includeReturnTrip {
			currentTime += calculateRouteTime(route.StartX, route.StartY, 0.0, 0.0)
		} else if i != len(d.Routes)-1 && len(d.Routes) > 1 {
			currentTime += calculateRouteTime(route.EndX, route.EndY, d.Routes[i+1].StartX, d.Routes[i+1].StartY)
			//add time to next load
		}
	}

	return currentTime
}

func (d *Driver) acceptRoute(r Route) {

	d.Routes = append(d.Routes, r)

}

func (r Route) calculateRouteTime() float64 {

	return calculateRouteTime(r.StartX, r.StartY, r.EndX, r.EndY)

}

func calculateRouteTime(startX, startY, endX, endY float64) float64 {

	xDiff := endX - startX
	yDiff := endY - startY
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseRoute(line string) Route {
	route := Route{}
	var err error

	tokens := strings.Split(line, " ")

	if len(tokens) != 3 {
		panic(errors.New("invalid file format"))
	}

	route.Number, err = strconv.ParseInt(tokens[0], 10, 64)
	check(err)
	pickup := strings.Split(tokens[1], ",")

	conv, err := strconv.ParseFloat(strings.Trim(pickup[0], "()"), 64)
	check(err)
	route.StartX = conv

	conv, err = strconv.ParseFloat(strings.Trim(pickup[1], "()"), 64)
	check(err)
	route.StartY = conv

	dropoff := strings.Split(tokens[2], ",")

	conv, err = strconv.ParseFloat(strings.Trim(dropoff[0], "()"), 64)
	check(err)
	route.EndX = conv

	conv, err = strconv.ParseFloat(strings.Trim(dropoff[1], "()"), 64)
	check(err)
	route.EndY = conv

	return route
}
