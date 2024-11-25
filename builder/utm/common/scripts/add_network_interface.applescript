on run argv
  set vmId to item 1 of argv # Id of the VM
  set modeVal to item 2 of argv # Mode of the network interface
  tell application "UTM"
    set vm to virtual machine id vmId
    set config to configuration of vm

    -- Existing network interfaces
    set networkInterfaces to network interfaces of config

    -- Create a new network interface configuration with given mode
    -- New network interface properties
    -- except mode all are default values
    set newNetworkInterfaceVal to { mode: modeVal}

    -- add the new network interface to the existing network interfaces
    copy newNetworkInterfaceVal to the end of networkInterfaces

    -- Update the VM configuration with the new network interface
    update configuration of vm with config
  end tell
end run