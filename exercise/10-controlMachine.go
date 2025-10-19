package main

import "fmt"

type Car struct {
	speed   int
	battery int
}

func NewCar(speed, battery int) *Car {
	return &Car{
		speed:   speed,
		battery: battery,
	}
}

func GetSpeed(car *Car) int {
	return car.speed
}

func GetBattery(car *Car) int {
	return car.battery
}

func ChargeCar(car *Car, minutes int) {
	if minutes < 0 {
		return
	}
	charge := minutes / 2
	car.battery += charge
	if car.battery > 100 {
		car.battery = 100
	}
}

func TryFinish(car *Car, distance int) string {
	if distance <= 0 || car.speed <= 0 {
		return ""
	}

	usage := distance / 2

	if usage > car.battery {
		car.battery = 0
		return ""
	}

	car.battery -= usage

	time := float64(distance) / float64(car.speed)

	return fmt.Sprintf("%.2f", time)
}
