// CAS BACnet Stack Golang Server Example
// https://github.com/chipkin/BACnetServerExampleGolang

package main

/*
#cgo LDFLAGS: -L. -lCASBACnetStack_x64_Debug

#include "CASBACnetStackDLL.h"
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
	"unsafe"
)

const GOLANG_EXAMPLE_VERSION = "0.0.1"

var udpConn *net.UDPConn

// Implement Callbacks
// ====================================================
// General Callbacks

//export goCallbackReceiveMessage
func goCallbackReceiveMessage(message unsafe.Pointer, maxMessageLength C.uint16_t, receivedConnectionString unsafe.Pointer, maxConnectionStringLength C.uint8_t, receivedConnectionStringLength unsafe.Pointer, networkType unsafe.Pointer) C.uint16_t {
	buf := make([]byte, uint16(maxMessageLength))
	msgLen, addr, err := udpConn.ReadFromUDP(buf)

	if os.IsTimeout(err) {
		// No message for us, return 0
		return 0
	} else if err != nil {
		fmt.Printf("Error in ReceiveMessage: %v\n", err)
		return 0
	}

	// Convert the received address to the CAS BACnet Stack connection string format
	portBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(portBytes, uint16(addr.Port))
	(*[CONNECTION_STRING_LENGTH]byte)(receivedConnectionString)[0] = addr.IP[12]
	(*[CONNECTION_STRING_LENGTH]byte)(receivedConnectionString)[1] = addr.IP[13]
	(*[CONNECTION_STRING_LENGTH]byte)(receivedConnectionString)[2] = addr.IP[14]
	(*[CONNECTION_STRING_LENGTH]byte)(receivedConnectionString)[3] = addr.IP[15]

	// Take note of endianess
	(*[CONNECTION_STRING_LENGTH]byte)(receivedConnectionString)[4] = portBytes[1]
	(*[CONNECTION_STRING_LENGTH]byte)(receivedConnectionString)[5] = portBytes[0]

	*(*uint8)(receivedConnectionStringLength) = uint8(CONNECTION_STRING_LENGTH)

	// Convert the received data to a format that CASBACnet Stack can process
	for i := 0; i < msgLen && i < int(maxMessageLength); i++ {
		(*[MAX_PACKET_BUFFER_LENGTH]byte)(message)[i] = buf[i]
	}

	// Set the network type
	*(*uint8)(networkType) = 0

	return C.uint16_t(msgLen)
}

//export goCallbackSendMessage
func goCallbackSendMessage(message unsafe.Pointer, messageLength C.uint16_t, connectionString unsafe.Pointer, connectionStringLength C.uint8_t, networkType C.uint8_t, broadcast C.bool) C.uint16_t {
	if networkType != 0 {
		fmt.Println("Error: Unsupported network type.")
		return 0
	}

	// Extract the Connection String from CAS BACnet Stack into an IP address and port
	var destAddr net.UDPAddr
	destAddr.Port = int((*[CONNECTION_STRING_LENGTH]byte)(connectionString)[4])*256 + int((*[CONNECTION_STRING_LENGTH]byte)(connectionString)[5])
	destAddr.IP = make([]byte, 4)
	destAddr.IP[0] = (*[CONNECTION_STRING_LENGTH]byte)(connectionString)[0]
	destAddr.IP[1] = (*[CONNECTION_STRING_LENGTH]byte)(connectionString)[1]
	destAddr.IP[2] = (*[CONNECTION_STRING_LENGTH]byte)(connectionString)[2]
	destAddr.IP[3] = (*[CONNECTION_STRING_LENGTH]byte)(connectionString)[3]

	// TODO: Handle broadcast messages

	// Extract the message from CAS BACnet Stack to a byte slice
	messageBytes := make([]byte, messageLength)
	for i := 0; i < int(messageLength); i++ {
		messageBytes[i] = (*[MAX_PACKET_BUFFER_LENGTH]byte)(message)[i]
	}

	// Send the message
	_, err := udpConn.WriteToUDP(messageBytes, &destAddr)
	if err != nil {
		fmt.Printf("Error writing to UDP: %v\n", err)
		return 0
	}

	return C.uint16_t(messageLength)
}

//export goCallbackGetSystemTime
func goCallbackGetSystemTime() C.time_t {
	return C.longlong(time.Now().Unix())
	// For systems with 32 bit time:
	// return C.long(time.Now().Unix())
}

//export goCallbackLogDebugMessage
func goCallbackLogDebugMessage(message unsafe.Pointer, messageLength C.uint16_t, messageType C.uint8_t) {
	goMessage := (*[MAX_DEBUG_MESSAGE_LENGTH]byte)(message)

	for i := 0; C.uint16_t(i) < messageLength && i < MAX_DEBUG_MESSAGE_LENGTH; i++ {
		fmt.Print(string(goMessage[i]))
	}
}

// Get Property Callbacks

//export goCallbackGetPropertyCharString
func goCallbackGetPropertyCharString(deviceInstance C.uint32_t, objectType C.uint16_t, objectInstance C.uint32_t, propertyIdentifier C.uint32_t, value unsafe.Pointer, valueElementCount unsafe.Pointer, maxElementCount C.uint32_t, encodingType unsafe.Pointer, useArrayIndex C.bool, propertyArrayIndex C.uint32_t) C.bool {
	if uint32(deviceInstance) == device.Instance {
		if uint32(propertyIdentifier) == uint32(BACnet_propertyIdentifier["objectName"]) {
			if uint16(objectType) == uint16(BACnet_objectType["analogInput"]) && uint32(objectInstance) == analogInput.Instance {
				objectNameByte := []byte(analogInput.ObjectName)
				for i := 0; i < len(objectNameByte) && i < int(maxElementCount); i++ {
					(*[MAX_CHARACTER_STRING_SIZE]byte)(value)[i] = objectNameByte[i]
				}
				if len(objectNameByte) > int(maxElementCount) {
					*(*uint32)(valueElementCount) = uint32(maxElementCount)
				} else {
					*(*uint32)(valueElementCount) = uint32(len(objectNameByte))
				}
				return C.bool(true)
			}
		}
		if uint32(propertyIdentifier) == uint32(BACnet_propertyIdentifier["objectName"]) && uint16(objectType) == uint16(BACnet_objectType["device"]) {
			deviceNameByte := []byte(device.ObjectName)
			for i := 0; i < len(deviceNameByte) && i < int(maxElementCount); i++ {
				(*[MAX_CHARACTER_STRING_SIZE]byte)(value)[i] = deviceNameByte[i]
			}
			if len(deviceNameByte) > int(maxElementCount) {
				*(*uint32)(valueElementCount) = uint32(maxElementCount)
			} else {
				*(*uint32)(valueElementCount) = uint32(len(deviceNameByte))
			}
			return C.bool(true)
		}
	}
	return C.bool(false)
}

//export goCallbackGetPropertyReal
func goCallbackGetPropertyReal(deviceInstance C.uint32_t, objectType C.uint16_t, objectInstance C.uint32_t, propertyIdentifier C.uint32_t, value unsafe.Pointer, useArrayIndex C.bool, propertyArrayIndex C.uint32_t) C.bool {
	if uint32(deviceInstance) == device.Instance {
		if uint32(propertyIdentifier) == uint32(BACnet_propertyIdentifier["presentValue"]) {
			if uint16(objectType) == uint16(BACnet_objectType["analogInput"]) && uint32(objectInstance) == analogInput.Instance {
				*(*float32)(value) = analogInput.PresentValue
				return C.bool(true)
			}
		}
	}

	return C.bool(false)
}

func main() {
	fmt.Println("FYI: CAS BACnet Stack Golang Server Example v" + GOLANG_EXAMPLE_VERSION)
	fmt.Println("FYI: https://github.com/chipkin/BACnetServerExampleGolang")

	majorVersion := C.BACnetStack_GetAPIMajorVersion()
	minorVersion := C.BACnetStack_GetAPIMinorVersion()
	patchVersion := C.BACnetStack_GetAPIPatchVersion()
	buildVersion := C.BACnetStack_GetAPIBuildVersion()
	fmt.Printf("FYI: BACnet Version: %v.%v.%v.%v\n", majorVersion, minorVersion, patchVersion, buildVersion)

	// 1. Setup UDP resource
	var udpAddr net.UDPAddr
	udpAddr.Port = 47808
	var err error
	udpConn, err = net.ListenUDP("udp", &udpAddr)
	if err != nil {
		fmt.Printf("Error: cannot create UDP resource - %v\n", err)
		os.Exit(1)
	}

	// 2. Setup the callbacks
	fmt.Print("Setting up callbacks... ")
	RegisterBACnetStackCallbacks()
	fmt.Println("Done.")

	// 3. Add devices
	fmt.Print("Setting up server device... ")
	if !C.BACnetStack_AddDevice(C.uint(device.Instance)) {
		fmt.Println("Error: Failed to add Device")
		os.Exit(1)

	}
	fmt.Println("Done.")

	// 4. Add objects
	fmt.Print("Adding AnalogInput... ")
	if !C.BACnetStack_AddObject(C.uint(device.Instance), C.ushort(BACnet_objectType["analogInput"]), C.uint(analogInput.Instance)) {
		fmt.Println("Error: Failed to add analogInput")
		os.Exit(1)
	}

	// Enable optional properties
	if !C.BACnetStack_SetPropertyEnabled(C.uint(device.Instance), C.ushort(BACnet_objectType["analogInput"]), C.uint(analogInput.Instance), C.uint(BACnet_propertyIdentifier["reliability"]), C.bool(true)) {
		fmt.Println("Error: Failed to enable reliability for analogInput")
		os.Exit(1)
	}
	fmt.Println("Done.")

	// 5. Start the main loop
	fmt.Println("Entering main loop...")
	lastUpdatedTime := time.Now()
	for {
		// Call the DLLs loop function which checks for messages and processes them.
		// Update the UDP read timeout deadline each tick
		udpConn.SetDeadline(time.Now().Add(time.Millisecond * 500))
		C.BACnetStack_Tick()

		// Sleep between loops. Give some time to the other application
		time.Sleep(10 * time.Millisecond)

		// Increment AnalogInput presentValue by 0.1 every 3 seconds
		if lastUpdatedTime.Add(time.Second * 3).Before(time.Now()) {
			analogInput.PresentValue += 0.01
			lastUpdatedTime = time.Now()
			C.BACnetStack_ValueUpdated(C.uint(device.Instance), C.ushort(BACnet_objectType["analogInput"]), C.uint(analogInput.Instance), C.uint(BACnet_propertyIdentifier["presentValue"]))
			fmt.Println("FYI: Updating AnalogInput (0) PresentValue: ", analogInput.PresentValue)
		}
	}
}
