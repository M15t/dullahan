package session

import (
	"fmt"
	"math"
	"time"

	"github.com/golang-module/carbon/v2"
)

func generateQuarters(startDate, endDate time.Time) []int64 {
	quarters := []int64{1}

	diff := endDate.Sub(startDate)

	for i := 1; i < int(math.Round(diff.Hours()/24/30/3)); i++ {
		quarters = append(quarters, int64(i*3))
	}

	return quarters
}

func generateMonths(startDate, endDate time.Time) []int64 {
	months := []int64{}

	diff := endDate.Sub(startDate)

	for i := 0; i < int(math.Round(diff.Hours()/24/30)); i++ {
		months = append(months, int64(i))
	}

	return months
}

func cutoffFloat64(fa []float64, f float64) (fb []float64) {
	for _, a := range fa {
		if a == f {
			fb = append(fb, 0)
			break
		}
		fb = append(fb, a)
	}
	return fb
}

func quarterOf(month int) int {
	quarter := math.Ceil(float64(month) / 3)
	return int(quarter)
}

func getQuarter(d time.Time, i int64) string {
	var year, month, quarter int
	nextDate := d.AddDate(0, int(i), 0)
	year = nextDate.Year()
	month = int(nextDate.Month())
	quarter = quarterOf(month)

	return fmt.Sprintf("%dQ%d", year, quarter)
}

func getMonth(d time.Time, i int64) string {
	c := carbon.FromStdTime(d)
	c = c.SetDay(26)

	nextDate := c.AddMonths(int(i))
	year := nextDate.Year()
	month := nextDate.ToMonthString()[0:3]

	// year%1e2

	return fmt.Sprintf("%s-%d", month, year)
}

func getMonthAndYear(d time.Time, i int64) string {
	c := carbon.FromStdTime(d)
	c = c.SetDay(26)

	nextDate := c.AddMonths(int(i))
	year := nextDate.Year()
	month := nextDate.ToMonthString()[0:3]

	return fmt.Sprintf("%s %d", month, year)
}

func roundFloat(num float64) float64 {
	output := math.Pow(10, float64(2))
	return float64(int(num*output)) / output
}
