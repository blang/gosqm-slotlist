package main

import (
	"flag"
	"fmt"
	"github.com/blang/gosqm"
	"io"
	"os"
)

const LF = "\r\n"

var (
	output = flag.String("output", "", "file to write slotlist to, empty to print to stdout")
	input  = flag.String("input", "mission.sqm", "mission.sqm to read slotlist from")
)

func init() {
	flag.Parse()
}

func main() {
	if *input == "" {
		fmt.Print("Specify path to mission.sqm" + LF)
		return
	}
	f, err := os.Open(*input)
	if err != nil {
		fmt.Printf("Can't open mission.sqm: %s"+LF, err)
		return
	}
	defer f.Close()
	dec := gosqm.NewDecoder(f)
	missionFile, err := dec.Decode()
	if err != nil {
		fmt.Printf("Error while reading mission.sqm: %s"+LF, err)
		return
	}
	var out io.Writer

	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Printf("Can't write to output: %s"+LF, err)
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
						fmt.Fprint(out, "--Group--"+LF)
						groupHasSlots = true
					}
					if u.Description != "" {
						fmt.Fprint(out, u.Description+LF)
					} else {
						fmt.Fprint(out, "Playable without description: "+u.Description+LF)
					}
				}
			}
		}

	} else {
		fmt.Println("No Mission found")
		return
	}
}
