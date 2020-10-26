// Copyright (c) 2019-2020 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/reconditematter/geomys"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// cmd-popcen {npar} -- prints top 10 central locations from the standard input using `npar`
// concurrent goroutines to compute the population-weighted sum of geographic distances.
func main() {
	npar, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	if npar < 1 {
		npar = 1
	}
	if npar > 256 {
		npar = 256
	}
	//
	ps := read()
	fmt.Println("nrec", len(ps), "npar", npar)
	//
	T := time.Now()
	computeall(npar, ps)
	fmt.Println("time", time.Since(T).Milliseconds(), "ms")
	//
	sort.Sort(locs(ps))
	fmt.Println("central locations (top 10)")
	for k := range ps {
		if k == 10 {
			break
		}
		fmt.Println(ps[k])
	}

}

type loc struct {
	id       string
	lat, lon float64
	weight   float64
	cen      float64
}

type locs []loc

func (s locs) Len() int           { return len(s) }
func (s locs) Less(i, j int) bool { return s[i].cen < s[j].cen }
func (s locs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var sph geomys.Spheroid = geomys.WGS1984()

// computeone -- computes the weighted sum of distances from ps[k]
// to all the points in `ps`.
func computeone(k int, ps []loc) {
	q := geomys.Geo(ps[k].lat, ps[k].lon)
	sum := 0.0
	for i := range ps {
		p := geomys.Geo(ps[i].lat, ps[i].lon)
		sum += ps[i].weight * geomys.Andoyer(sph, q, p)
	}
	ps[k].cen = sum
}

// computeall -- computes the weighted sum of distances for all the
// points in `ps` by calling `computeone` using `npar` concurrent
// goroutines.
func computeall(npar int, ps []loc) {
	var wg sync.WaitGroup
	wg.Add(npar)
	for par := 0; par < npar; par++ {
		first := par * len(ps) / npar
		limit := (par + 1) * len(ps) / npar
		go func(first, limit int) {
			for k := first; k < limit; k++ {
				computeone(k, ps)
			}
			wg.Done()
		}(first, limit)
	}
	wg.Wait()
}

// read -- reads the standard input records and returns an array of locations
// with computed weights.
// The input:
//	id,population,latitude,longitude
func read() []loc {
	result := make([]loc, 0)
	rdr := csv.NewReader(os.Stdin)
	W := 0.0
	for {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		p := loc{id: record[0], lat: mustf64(record[2]), lon: mustf64(record[3]), weight: mustf64(record[1]), cen: 0}
		W += p.weight
		result = append(result, p)
	}
	//
	for k := range result {
		result[k].weight /= W
	}
	//
	return result
}

// mustf64 -- converts `s` to a float64 number, or panics on a format error.
func mustf64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
