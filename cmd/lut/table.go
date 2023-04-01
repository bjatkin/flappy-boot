package main

import (
	"math"
	"os"
)

type sinTable struct {
	file   string
	Values [1 << 6]int
}

func newSinTable(file string) *sinTable {
	// generate a full lookup table with all 256 entries
	entries := 1 << 8
	lut := make([]int, entries)
	for i := 0; i < entries; i++ {
		s := sin(float64(i) / float64(entries))
		fix := floatToFix(s)
		lut[i] = int(fix)
	}

	// compress the table down so the generated code is smaller
	var values [1 << 6]int
	for i := len(lut) - 1; i >= 0; i-- {
		value := lut[i]
		key := i & 0x00FC
		key >>= 2

		values[key] = value
	}

	return &sinTable{
		file:   file,
		Values: values,
	}
}

func (s *sinTable) execute() error {
	f, err := os.Create(s.file)
	if err != nil {
		return err
	}

	err = goTemplates.ExecuteTemplate(f, "sin.go.tmpl", s)
	if err != nil {
		return err
	}

	return nil
}

func sin(f float64) float64 {
	return math.Sin(2 * math.Pi * f)
}
