package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variables
var (
	seatingCapacity = 10
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (b *BarberShop) CutHair(barber, client string) {
	color.Green("%s is cutting %s's hair", barber, client)
	time.Sleep(b.HairCutDuration)
	color.Green("%s is finished cutting %s's hair", barber, client)
}

func (b *BarberShop) SendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	b.BarbersDoneChan <- true
}

func (b *BarberShop) AddClient(client string) {
	// print out a message
	color.Green("*** %s arrives!", client)

	if b.Open {
		select {
		case b.ClientsChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves", client)
		}
	} else {
		color.Red("The shop is already closed, so %s leaves!", client)
	}
}

func (b *BarberShop) CloseShopForDay() {
	color.Cyan("Closing shop for the day")
	close(b.ClientsChan)
	b.Open = false

	for i := 1; i <= b.NumberOfBarbers; i++ {
		<-b.BarbersDoneChan
	}

	close(b.BarbersDoneChan)

	color.Green("--------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day, and everyone has gone home")
}

func (b *BarberShop) AddBarber(barber string) {
	b.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		for {
			// if there are no clients, the barber goes to sleep
			if len(b.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			client, shopOpen := <-b.ClientsChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up", client, barber)
					isSleeping = false
				}
				// cut hair
				b.CutHair(barber, client)
			} else {
				// shop is closed, so send the barber home nad close this  goroutine
				b.SendBarberHome(barber)
				return
			}
		}
	}()
}

func main() {
	// seed our random number generator
	rand.Seed(time.Now().UnixNano())

	// print welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels if we need any
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	// add barbers
	shop.AddBarber("Frank")
	shop.AddBarber("Gerard")
	shop.AddBarber("Milton")
	shop.AddBarber("Susan")
	shop.AddBarber("Kelly")
	shop.AddBarber("Pat")

	// start the barber shop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.CloseShopForDay()
		closed <- true
	}()

	// add clients
	i := 1
	go func() {
		for {
			// get a random number with average arrival rate
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.AddClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// block until the barbershop is closed
	<-closed
}
