# BACnet Server Example Golang

A basic BACnet IP server example written in Golang using the [CAS BACnet Stack](https://store.chipkin.com/services/stacks/bacnet-stack).

## Compile and Run

Download the source and place the [CAS BACnet Stack DLL/so](https://store.chipkin.com/services/stacks/bacnet-stack) and [CASBACnetStackDLL.h](https://store.chipkin.com/services/stacks/bacnet-stack) in the source directory. 

Linux: Replace backslash with slash in `go.mod`

Run `go run .` in the source directory to run the example. Golang v1.9 or above is required.

Pre-configured with the following example BACnet device and objects:

- **Device**: 390000 (Example Device Yellow)
  - analog_input: 0 (Analog Input White)

## Example Output

```txt
FYI: CAS BACnet Stack Golang Server Example v0.0.1
FYI: https://github.com/chipkin/BACnetServerExampleGolang
FYI: BACnet Version: 3.28.1.1980
Setting up callbacks... Done.
Setting up server device... Done.
Adding AnalogInput... Done.
Entering main loop...
```