package main

import (
	"flag"
	"fmt"
	"github.com/blang/gosqm"
	"io"
	"os"
)

var (
	output = flag.String("output", "", "file to write slotlist to, empty to print to stdout")
	input  = flag.String("input", "mission.sqm", "mission.sqm to read slotlist from")
)

func init() {
	flag.Parse()
}

func main() {
	if *input == "" {
		fmt.Println("Specify path to mission.sqm")
		return
	}
	f, err := os.Open(*input)
	if err != nil {
		fmt.Printf("Can't open mission.sqm: %s\n", err)
		return
	}
	defer f.Close()
	dec := gosqm.NewDecoder(f)
	missionFile, err := dec.Decode()
	if err != nil {
		fmt.Printf("Error while reading mission.sqm: %s", err)
		return
	}
	var out io.Writer

	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Printf("Can't write to output: %s\n", err)
			return
		}
		defer f.Close()
		out = io.Writer(f)
	} else {
		out = os.Stdout
	}

	if missionFile.Mission != nil {
		for _, g := range missionFile.Mission.Groups {
			groupHasSlots := false
			for _, u := range g.Units {
				if u.Player != "" {
					if !groupHasSlots {
						fmt.Fprintf(out, "--Group--\n")
						groupHasSlots = true
					}
					if u.Description != "" {
						fmt.Fprintf(out, "%s\n", u.Description)
					} else {
						fmt.Fprintf(out, "Playable without description%s\n", u.Description)
					}
				}
			}
		}

	} else {
		fmt.Println("No Mission found")
		return
	}
}
