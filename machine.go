package main

import (
	"os"

	dev "github.com/uroshercog/sic-machine/devices"
	"github.com/uroshercog/sic-machine/memory"
	"github.com/uroshercog/sic-machine/obj"
	"bufio"
	"github.com/uroshercog/sic-machine/processor"
	"github.com/uroshercog/sic-machine/ui"
	"fmt"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", r)
		}
	}()

	/* 1. Preberi ime datoteke iz command line argumentov */

	if len(os.Args) < 2 {
		panic("No filename provided")
	}
	objectFileName := string(os.Args[1])

	uix := &ui.UI{}
	devices := dev.New()
	RAM := memory.New()

	CPU := processor.NewCPU(RAM, devices)
	CPU.OnStart = append(CPU.OnStart, func() {
		uix.RenderStatusWidget("started")
		uix.RenderRegistersWidget(CPU.GetRegisters())
	})

	CPU.OnStop = append(CPU.OnStop, func() {
		uix.RenderStatusWidget("stopped")
	})

	CPU.OnExec = append(CPU.OnExec, func(cmd string) {
		uix.RenderExecutingCommand(cmd)
		uix.RenderRegistersWidget(CPU.GetRegisters())
		uix.RenderRAMWidget(RAM.GetRaw())
		uix.RenderScreenWidget(RAM.GetRaw())
	})

	uix.Handle(ui.PAUSE, CPU.Stop)
	uix.Handle(ui.CONTINUE, CPU.Start)
	uix.Handle(ui.STEP, CPU.Step)

	/*
		2. Nalozi cel podan fajl v RAM
			 - ime fajla je podano preko argumentov
	*/
	objectCode := parseObjectCode(objectFileName)
	RAM.Load(objectCode)
	CPU.SetStart(objectCode.StartAddr)
	uix.Run(RAM.GetRaw(), CPU.GetRegisters())
}

func parseObjectCode(filename string) *obj.ObjectCode {
	objCode := &obj.ObjectCode{}

	if f, err := os.Open(filename); err != nil {
		panic(err)
	} else {
		// Lets assume that no line can be more than 4096 bytes long
		reader := bufio.NewReader(f)
		for {
			if line, _, err := reader.ReadLine(); err == nil {
				objCode.Load(line)
			} else {
				break
			}
		}
	}

	return objCode
}
