// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step shuts down the machine.
//
// Uses:
//
//	communicator packersdk.Communicator
//	driver Driver
//	ui     packersdk.Ui
//	vmName string
//
// Produces:
//
//	<nothing>
type StepShutdown struct {
	Command         string
	Timeout         time.Duration
	Delay           time.Duration
	DisableShutdown bool
}

func (s *StepShutdown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	comm := state.Get("communicator").(packersdk.Communicator)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)
	vmId := state.Get("vmId").(string)

	if !s.DisableShutdown {
		if s.Command != "" {
			ui.Say("Gracefully halting virtual machine...")
			log.Printf("Executing shutdown command: %s", s.Command)
			cmd := &packersdk.RemoteCmd{Command: s.Command}
			if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
				err := fmt.Errorf("failed to send shutdown command: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}

		} else {
			ui.Say("Halting the virtual machine...")
			if err := driver.Stop(vmId); err != nil {
				err := fmt.Errorf("error stopping VM: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}
		}
	} else {
		ui.Say("Automatic shutdown disabled. Please shutdown virtual machine.")
	}

	// Wait for the machine to actually shut down
	log.Printf("Waiting max %s for shutdown to complete", s.Timeout)
	shutdownTimer := time.After(s.Timeout)
	for {
		running, _ := driver.IsRunning(vmId)
		if !running {

			if s.Delay.Nanoseconds() > 0 {
				log.Printf("Delay for %s after shutdown to allow locks to clear...", s.Delay)
				time.Sleep(s.Delay)
			}

			break
		}

		select {
		case <-shutdownTimer:
			err := errors.New("timeout while waiting for machine to shutdown")
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}

	log.Println("VM shut down.")
	return multistep.ActionContinue
}

func (s *StepShutdown) Cleanup(state multistep.StateBag) {}
