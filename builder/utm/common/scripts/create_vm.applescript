---
-- create_vm.applescript
-- This script creates a new VM with the specified properties.
-- Usage: osascript create_vm.applescript --name <VM_NAME> --backend <BACKEND> --arch <ARCH> --iso <ISO_PATH> --size <DISK_SIZE>
-- Example: osascript create_vm.applescript --name "MyVM" --backend "QeMu" --arch "aarch64" --iso "/path/to/image.iso" --size 65536
on run argv
    -- Initialize variables
    set vmName to ""
    set vmBackend to ""
    set vmArch to ""

    -- Parse arguments
    repeat with i from 1 to (count argv)
        set currentArg to item i of argv
        if currentArg is "--name" then
            set vmName to item (i + 1) of argv
        else if currentArg is "--backend" then
            set vmBackend to item (i + 1) of argv as string
        else if currentArg is "--arch" then
            set vmArch to item (i + 1) of argv
        else if currentArg is "--iso" then
            set isoPath to POSIX file (POSIX path of item (i + 1) of argv)
        else if currentArg is "--size" then
            set diskSize to item (i + 1) of argv
        end if
    end repeat

    -- Create a new VM with the specified properties
    tell application "UTM"
        set vm to make new virtual machine with properties �
          { backend:vmBackend, �
            configuration:{ �
              name:vmName, �
              architecture:vmArch, �
              drives:{  �
                {removable:true, source:isoPath}, �
                {guest size:diskSize} �
              } �
            } �
          }
    end tell
end run