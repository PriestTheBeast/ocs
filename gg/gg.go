// Functions that Go should already have implemented but doesn't.
package gg

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Random generator already with source
// var RandGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

const float64EqualityThreshold = 1e-5

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func AlmostEqualOrBigger(a, b float64) bool {
	return AlmostEqual(a, b) || a > b
}

func AlmostEqualOrLower(a, b float64) bool {
	return AlmostEqual(a, b) || a < b
}

// Rand_Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [min,max)
// from the default Source.
// It panics if max-min <= 0 (except both min and max == 0).
func Rand_Intn(min, max int) int {
	if min == 0 && (max == 1 || max == 0) {
		return 0
	}
	return rand.Intn(max-min) + min
}

// Rand_Choose randomly chooses from given slice
func Rand_Choose[V interface{}](slice []V) *V {
	return &slice[Rand_Intn(0, len(slice))]
}

func Rand_Shuffle[X any](s []X) []X {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

// Remove removes an item from a slice. It changes slice order.
func Remove[X any](s []X, i int) []X {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Remove removes an item from a slice. It keeps slice order.
func RemoveKeepOrder[X any](slice []X, s int) []X {
	return append(slice[:s], slice[s+1:]...)
}

// Gets the index of the biggest element
func GetIndexOfBiggest(slice []float64) int {
	max := 0.0
	index := 0

	for i := range slice {
		if slice[i] > max {
			max = slice[i]
			index = i
		}
	}

	return index
}

func TimeTrack(start time.Time, doSomethingFunc func(time.Duration)) {
	elapsed := time.Since(start)
	doSomethingFunc(elapsed)
}

func ContextWithSigTerm(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// Start a goroutine to listen for the interrupt signal
	go func() {
		<-c
		fmt.Println("\nCtrl+C pressed. Cancelling...")
		// Cancel the context when the interrupt signal is received
		cancel()
	}()

	return ctx, cancel
}

// SumSquareError returns the squared error of the elems in the array, all summed up
func SumSquareError(actual []float64, expected []float64) float64 {
	total_error := 0.0
	for i := range actual {
		value := actual[i] - expected[i]
		total_error += value * value
	}

	return total_error
}

// ParseToFloat64 is just an util function to parse a string to Float64
// without having to deal with error handling
func ParseToFloat64(value string) float64 {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Parsing string to float, error: %v", err)
	}

	return number
}

func RoundToDecimalPlaces(number float64, decimalPlaces int) float64 {
	conversionNumber := math.Pow10(decimalPlaces)
	return math.Round(number*conversionNumber) / conversionNumber
}

// Map function
func MapXtoY[X any, Y any](collection []X, fn func(elem X) Y) []Y {
	var result []Y
	for _, item := range collection {
		result = append(result, fn(item))
	}
	return result
}

// MapCompreensionMakeValues makes a map with the result of the Map function
func MapCompreensionMakeValues[X comparable, Y any](collection []X, getValue func(elem X) Y) map[X]Y {
	mapResult := make(map[X]Y)

	for i := range collection {
		elem := collection[i]
		mapResult[elem] = getValue(elem)
	}

	return mapResult
}

// MapCompreensionMakeValues makes a map with the result of the Map function
func MapCompreensionMakeKeys[X any, Y comparable](collection []X, getKey func(elem X) Y) map[Y]X {
	mapResult := make(map[Y]X)

	for i := range collection {
		elem := collection[i]
		mapResult[getKey(elem)] = elem
	}

	return mapResult
}

// Reduce function
// func ReduceXtoY[X any, Y any](collection []X, init Y, fn func(memo Y, elem X) Y) Y {
// 	result := init
// 	for _, item := range collection {
// 		result = fn(result, item)
// 	}
// 	return result
// }

// Filter function
// func FilterXs[X any, Y any](collection []X, fn func(elem X) bool) []X {
// 	var result []X
// 	for _, item := range collection {
// 		if fn(item) {
// 			result = append(result, item)
// 		}
// 	}
// 	return result
// }
