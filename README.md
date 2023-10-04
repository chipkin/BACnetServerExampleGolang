# BACnet Server Example Golang

A basic BACnet IP server example written in Golang using the [CAS BACnet Stack](https://store.chipkin.com/services/stacks/bacnet-stack).

## Compile and Run

1. Place the following in the source directory: CASBACnetStack_x64_Debug.dll, CASBACnetStack_x64_Debug.so, CASBACnetStack_x64_Release.dll, CASBACnetStack_x64_Release.so, and [CASBACnetStackDLL.h](https://store.chipkin.com/services/stacks/bacnet-stack). These can be found in the [CAS BACnet Stack](https://store.chipkin.com/services/stacks/bacnet-stack).
2. Linux only: Replace backslash with slash in `go.mod`.
3. Run `go run .` in the source directory to run the example. Golang v1.9 or above is required.

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