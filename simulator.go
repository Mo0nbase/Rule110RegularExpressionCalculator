package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"
)

var sim [][]uint64

// make sure to initialize with a raw array of 1 and 0
func initialize(initial []uint64, evolutions int, history bool) {
	for {
		if initial[0] == 1 {
			break
		}
		initial = initial[1:]
	}

	loopLength := int(math.Ceil(float64(len(initial)) / 64.0))

	if history == true {
		sim = make([][]uint64, evolutions+1)
		for i := range sim {
			sim[i] = make([]uint64, loopLength+int(math.Ceil(float64(evolutions/64.0))))
		}
	} else {
		sim = make([][]uint64, 2)
		for i := range sim {
			sim[i] = make([]uint64, loopLength+int(math.Ceil(float64(evolutions/64.0))))
		}
	}

	for i, j, k, p := 0, 0, 0, 0; i < 64; i++ {
		for {
			if k > ((len(sim[0]))*64 - len(initial)) {
				if p < len(initial) && initial[p] == 1 {
					sim[0][j] = setBit(sim[0][j], uint64(i))
				}
				p++
			}
			k++
			j++
			if j == len(sim[0]) {
				j = 0
				break
			}
		}
	}
}

func simulate(history bool, evolutions int) time.Duration {
	start := time.Now()
	if sim == nil {
		fmt.Println("Program uninitialized using default...")
		fmt.Println()
		var testSplit = r110Default()
		initialize(testSplit, evolutions, history)
	}

	if history == true {
		for i := 1; i < evolutions+1; i++ {
			for j := 0; j < len(sim[i]); j++ {
				if j == 0 {
					sim[i][j] = ((^(sim[i-1][len(sim[i-1])-1]) << 1) & sim[i-1][j]) | (sim[i-1][j] ^ sim[i-1][j+1])
				} else if j == len(sim[i])-1 {
					sim[i][j] = ((^(sim[i-1][j-1])) & sim[i-1][j]) | (sim[i-1][j] ^ sim[i-1][0]>>1)
				} else {
					sim[i][j] = ((^(sim[i-1][j-1])) & sim[i-1][j]) | (sim[i-1][j] ^ sim[i-1][j+1])
				}
			}
		}
	} else {
		for i, k := 1, 0; i < evolutions+1; i++ {
			for j := 0; j < len(sim[0]); j++ {
				if k == 0 {
					if j == 0 {
						sim[k+1][j] = ((^(sim[k][len(sim[k])-1]) << 1) & sim[k][j]) | (sim[k][j] ^ sim[k][j+1])
					} else if j == len(sim[k])-1 {
						sim[k+1][j] = ((^(sim[k][j-1])) & sim[k][j]) | (sim[k][j] ^ sim[k][0]>>1)
					} else {
						sim[k+1][j] = ((^(sim[k][j-1])) & sim[k][j]) | (sim[k][j] ^ sim[k][j+1])
					}
					if j == len(sim[0])-1 {
						k++
					}
				} else {
					if j == 0 {
						sim[k-1][j] = ((^(sim[k][len(sim[k])-1]) << 1) & sim[k][j]) | (sim[k][j] ^ sim[k][j+1])
					} else if j == len(sim[k])-1 {
						sim[k-1][j] = ((^(sim[k][j-1])) & sim[k][j]) | (sim[k][j] ^ sim[k][0]>>1)
					} else {
						sim[k-1][j] = ((^(sim[k][j-1])) & sim[k][j]) | (sim[k][j] ^ sim[k][j+1])
					}
					if j == len(sim[0])-1 {
						k--
					}
				}
			}
		}
	}

	//duration := time.Since(start)
	//fmt.Println(duration)
	//fmt.Println("------------------------------------------------")
	//PrintMemUsage()

	//fmt.Println("---------------------------------------------------")
	//fmt.Println()
	//displayFancy()
	return time.Since(start)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func readTape() uint64 {
	return 1
}

func writeToFile() {

}

func readFromFile() {

}

// REMEMBER: BINARY NUMBERS READ RIGHT TO LEFT!!!
func displayRaw(layer int) {
	//fmt.Println("REMEMBER: BINARY NUMBERS READ RIGHT TO LEFT!!!")
	for i := 0; i < len(sim[layer]); i++ {
		fmt.Printf("%064d", strconv.FormatInt(int64(sim[layer][i]), 2))
		fmt.Println()
	}
}

func displayFancy() {
	if sim == nil {
		fmt.Println("Array Empty!")
	}
	for i := 0; i < len(sim); i++ {
		for j := 0; j < 64; j++ {
			for k := 0; k < len(sim[i]); k++ {
				if getBit(sim[i][k], uint64(j)) == 0 {
					fmt.Print("□")
				} else {
					fmt.Print("■")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Sets the bit at pos in the integer n.
func setBit(n uint64, pos uint64) uint64 {
	n |= 1 << pos
	return n
}

// Clears the bit at pos in n.
func clearBit(n uint64, pos uint64) uint64 {
	mask := ^(1 << pos)
	n &= uint64(mask)
	return n
}

func getBit(n uint64, pos uint64) int {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}
