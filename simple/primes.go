package main

// Finds the first `n` primes, where `n` is passed to the command line.
import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func isPrime(n int) bool {
	if n%2 == 0 && n != 2 {
		return false
	}

	sqrt_n := int(math.Sqrt(float64(n)))

	for i := 3; i <= sqrt_n; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func isPrimeCached(n int, knownPrimes []int) (bool, []int) {
	if n%2 == 0 && n != 2 {
		return false, knownPrimes
	}

	sqrt_n := int(math.Sqrt(float64(n)))

	i := 0
	for i < len(knownPrimes) && knownPrimes[i] <= sqrt_n {
		if n%knownPrimes[i] == 0 {
			return false, knownPrimes
		}

		i += 1
	}

	return true, append(knownPrimes, n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: ./%s [count]", os.Args[0])
		return
	}

	ceiling, err := strconv.Atoi(os.Args[1]) //, 10, 64)
	if err != nil {
		fmt.Printf("Usage: ./%s [count]", os.Args[0])
		return
	}

	primes := []int{2}
	for i := 3; i <= ceiling; i += 1 {
		if wellIsIt, newPrimes := isPrimeCached(i, primes); wellIsIt {
			fmt.Printf("%d is prime\n", i)
			primes = newPrimes
		}
	}
}
