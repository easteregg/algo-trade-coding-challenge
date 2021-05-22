package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

const TWO_HUNDRED = 200
const FIFTY = 50

func RandFloat(min, max float64) float64 {
	var int_part, int_min, int_max int
	int_min = int(min)
	int_max = int(max)
	rand.Seed(time.Now().UnixNano())
	diff := int_max - int_min

	if diff != 0 {
		int_part = rand.Intn(diff) + int_min
	} else {
		int_part = int_min
	}

	float_part := rand.Float64()

	return float64(int_part) + float_part
	return 0.0
}

func RandFloats(min, max float64, n int) []float64 {
	list := make([]float64, n)

	for i := range list {
		list[i] = RandFloat(min, max)
	}

	return list
}

func getDates(how_many int) []int {
	list := make([]int, how_many)

	// Jan 01 2021
	start_time := 1609488000

	for i := range list {
		list[i] = start_time + (i * 86400)
	}

	return list
}

func writeToFile(output string) {
	ioutil.WriteFile("1.txt", []byte(output), 0644)
}

func getOHLC(start, min, max float64) (float64, float64, float64, float64) {
	open := RandFloat(min, max)
	if start != 0 {
		open = start
	}

	close := RandFloat(min, max)
	high_low := RandFloats(min, max, 2)
	sort.Float64s(high_low)
	if high_low[1] < open {
		high_low[1] = open
	}
	return open, high_low[1], high_low[0], close
}

type tAverageDay struct {
	date        string
	day_average float64
}

type tMovingAverage struct {
	average_day         tAverageDay
	fifty_day_avg       float64
	two_hundred_day_avg float64
}

/*
 *  This is the actual trade params
 */
func tradeWithMovingAverage() {
	content, err := ioutil.ReadFile("1.txt")
	if err != nil {
		panic("Could'nt reach file")
	}

	string_content := fmt.Sprintf("%s", content)
	days := strings.Split(string_content, "\n")
	var average_days []tAverageDay = make([]tAverageDay, len(days))
	var moving_average_days []tMovingAverage = make([]tMovingAverage, len(days))

	// Calculate moving average for the day
	for i := range days {
		if days[i] == "" {
			continue
		}
		day := strings.Split(days[i], " ")
		timestamp := day[0]
		open, _ := strconv.ParseFloat(day[1], 64)
		high, _ := strconv.ParseFloat(day[2], 64)
		low, _ := strconv.ParseFloat(day[3], 64)
		var in_date tAverageDay = tAverageDay{
			day_average: (open + high + low) / 3,
			date:        timestamp,
		}

		average_days[i] = in_date

		var in_date_moving_average tMovingAverage = tMovingAverage{
			average_day: in_date,
		}
		// Calculate 50 day moving average after 50 day
		if i >= FIFTY {
			var fifty_day_sum float64 = 0.0
			for j := (i - FIFTY); j < i; j++ {
				fifty_day_sum += average_days[j].day_average
			}
			in_date_moving_average.fifty_day_avg = fifty_day_sum / FIFTY
		} else {
			in_date_moving_average.fifty_day_avg = 0
		}

		if i >= TWO_HUNDRED {
			two_hundred_sum := 0.0
			for j := (i - TWO_HUNDRED); j < i; j++ {
				two_hundred_sum += in_date.day_average
			}
			in_date_moving_average.two_hundred_day_avg = two_hundred_sum / TWO_HUNDRED
		} else {
			in_date_moving_average.two_hundred_day_avg = 0
		}

		moving_average_days[i] = in_date_moving_average
	}

	has_bought := false
	var bought_at tAverageDay
	for i, output := range moving_average_days {

		if output.fifty_day_avg > output.two_hundred_day_avg && !has_bought {
			fmt.Println("index: ", i)
			fmt.Printf("BUY %s %.2f\n", output.average_day.date, output.average_day.day_average)
			bought_at = output.average_day
			has_bought = true
		}
		if output.fifty_day_avg < output.two_hundred_day_avg && has_bought {
			fmt.Printf("SELL %s %.2f\n", output.average_day.date, output.average_day.day_average)
			income := output.average_day.day_average*200 - bought_at.day_average*200
			fmt.Printf("Income: %.2f\n\n", income)
			has_bought = false
		}

		// fmt.Printf(
		// 	"Day: %s\n 50 day moving average: %.2f\n 200 day moving average: %.2f\n+-------+\n",
		// 	output.average_day.date,
		// 	output.fifty_day_avg,
		// 	output.two_hundred_day_avg,
		// )
	}

}

func main() {
	dates := getDates(2000)
	var output string
	var o, h, l, c float64
	// var test_cases []string
	o, h, l, c = getOHLC(0, 93, 140)
	for i := range dates {
		o, h, l, c = getOHLC(o, 93, 140)
		output += fmt.Sprintf("%d %.2f %.2f %.2f %.2f\n", dates[i], o, h, l, c)
		o = c
	}
	writeToFile(output)

	tradeWithMovingAverage()
}
