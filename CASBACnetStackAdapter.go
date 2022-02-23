// CAS BACnet Stack Golang Adapter and Example Database
// This file is used to help integration between the CAS BACnet stack DLL/.so files and a Golang application.
// Full documentation on what these function, and their parameters can be found in the CASBACnetStackDLL.h file.
// It also acts as a minimal database, storing device and object values.

package main

/*
#cgo LDFLAGS: -L. -lCASBACnetStack_x64_Debug

#include "CASBACnetStackDLL.h"

// Go Callback Export Declarations
// ====================================================
// General Callbacks
extern uint16_t goCallbackReceiveMessage(uint8_t* message, const uint16_t maxMessageLength, uint8_t* receivedConnectionString, const uint8_t maxConnectionStringLength, uint8_t* receivedConnectionStringLength, uint8_t* networkType);
extern uint16_t goCallbackSendMessage(const uint8_t* message, const uint16_t messageLength, const uint8_t* connectionString, const uint8_t connectionStringLength, const uint8_t networkType, bool broadcast);
extern time_t goCallbackGetSystemTime();
extern void goCallbackLogDebugMessage(const char* message, const uint16_t messageLength, const uint8_t messageType);

// Get Property Callbacks
extern bool goCallbackGetPropertyCharString(const uint32_t deviceInstance, const uint16_t objectType, const uint32_t objectInstance, const uint32_t propertyIdentifier, char* value, uint32_t* valueElementCount, const uint32_t maxElementCount, uint8_t* encodingType, const bool useArrayIndex, const uint32_t propertyArrayIndex);
extern bool goCallbackGetPropertyReal(const uint32_t deviceInstance, const uint16_t objectType, const uint32_t objectInstance, const uint32_t propertyIdentifier, float* value, const bool useArrayIndex, const uint32_t propertyArrayIndex);


// C Callback Declarations
// ====================================================
// General Callbacks
extern uint16_t fpCallbackReceiveMessage(uint8_t* message, const uint16_t maxMessageLength, uint8_t* receivedConnectionString, const uint8_t maxConnectionStringLength, uint8_t* receivedConnectionStringLength, uint8_t* networkType);
extern uint16_t fpCallbackSendMessage(const uint8_t* message, const uint16_t messageLength, const uint8_t* connectionString, const uint8_t connectionStringLength, const uint8_t networkType, bool broadcast);
extern time_t fpCallbackGetSystemTime();
extern void fpCallbackLogDebugMessage(const char* message, const uint16_t messageLength, const uint8_t messageType);

// Get Property Callbacks
extern bool fpCallbackGetPropertyCharString(const uint32_t deviceInstance, const uint16_t objectType, const uint32_t objectInstance, const uint32_t propertyIdentifier, char* value, uint32_t* valueElementCount, const uint32_t maxElementCount, uint8_t* encodingType, const bool useArrayIndex, const uint32_t propertyArrayIndex);
extern bool fpCallbackGetPropertyReal(const uint32_t deviceInstance, const uint16_t objectType, const uint32_t objectInstance, const uint32_t propertyIdentifier, float* value, const bool useArrayIndex, const uint32_t propertyArrayIndex);


// Callback C Wrappers
// ====================================================
// General Callbacks
uint16_t fpCallbackReceiveMessage(uint8_t* message, const uint16_t maxMessageLength, uint8_t* receivedConnectionString, const uint8_t maxConnectionStringLength, uint8_t* receivedConnectionStringLength, uint8_t* networkType) {
	return goCallbackReceiveMessage(message, maxMessageLength, receivedConnectionString, maxConnectionStringLength, receivedConnectionStringLength, networkType);
}
uint16_t fpCallbackSendMessage(const uint8_t* message, const uint16_t messageLength, const uint8_t* connectionString, const uint8_t connectionStringLength, const uint8_t networkType, bool broadcast) {
	return goCallbackSendMessage(message, messageLength, connectionString, connectionStringLength, networkType, broadcast);
}
time_t fpCallbackGetSystemTime() {
	return goCallbackGetSystemTime();
}
void fpCallbackLogDebugMessage(const char* message, const uint16_t messageLength, const uint8_t messageType) {
	goCallbackLogDebugMessage(message, messageLength, messageType);
}

// Get Property Callbacks
bool fpCallbackGetPropertyCharString(const uint32_t deviceInstance, const uint16_t objectType, const uint32_t objectInstance, const uint32_t propertyIdentifier, char* value, uint32_t* valueElementCount, const uint32_t maxElementCount, uint8_t* encodingType, const bool useArrayIndex, const uint32_t propertyArrayIndex) {
	return goCallbackGetPropertyCharString(deviceInstance, objectType, objectInstance, propertyIdentifier, value, valueElementCount, maxElementCount, encodingType, useArrayIndex, propertyArrayIndex);
}
bool fpCallbackGetPropertyReal(const uint32_t deviceInstance, const uint16_t objectType, const uint32_t objectInstance, const uint32_t propertyIdentifier, float* value, const bool useArrayIndex, const uint32_t propertyArrayIndex) {
	return goCallbackGetPropertyReal(deviceInstance, objectType, objectInstance, propertyIdentifier, value, useArrayIndex, propertyArrayIndex);
}
*/
import "C"

// Example DB Structs
type ExampleDBAnalogInput struct {
	ObjectName string
	Instance   uint32

	PresentValue float32
	CovIncurment float32
	Reliability  uint32
	Description  string
}
type ExampleDBDevice struct {
	ObjectName string
	Instance   uint32

	UTCOffset         int
	CurrentTimeOffset int64
	Description       string
	SystemStatus      uint32
}

// BACnet Constants
const MAX_DEBUG_MESSAGE_LENGTH = 512
const MAX_CHARACTER_STRING_SIZE = 256
const MAX_PACKET_BUFFER_LENGTH = 1497
const CONNECTION_STRING_LENGTH = 6

var BACnet_propertyIdentifier = map[string]uint32{
	"objectName":   77,
	"presentValue": 85,
	"reliability":  103,
}

var BACnet_objectType = map[string]uint16{
	"analogInput": 0,
	"device":      8,
}

// DB
var device = ExampleDBDevice{
	ObjectName:        "Example Device Yellow",
	Instance:          390000,
	UTCOffset:         0,
	CurrentTimeOffset: 0,
	Description:       "Golang Server Example Device",
	SystemStatus:      0,
}
var analogInput = ExampleDBAnalogInput{
	ObjectName:   "Analog Input White",
	Instance:     101,
	PresentValue: 1.001,
	CovIncurment: 2.0,
	Reliability:  0,
	Description:  "Golang Server Example Analog Input White",
}

// Uncomment or add callback registerations as needed
func RegisterBACnetStackCallbacks() {
	// General Callbacks
	C.BACnetStack_RegisterCallbackReceiveMessage((*[0]byte)(C.fpCallbackReceiveMessage))
	C.BACnetStack_RegisterCallbackSendMessage((*[0]byte)(C.fpCallbackSendMessage))
	C.BACnetStack_RegisterCallbackGetSystemTime((*[0]byte)(C.fpCallbackGetSystemTime))
	C.BACnetStack_RegisterCallbackLogDebugMessage((*[0]byte)(C.fpCallbackLogDebugMessage))

	// Get Property Callbacks
	C.BACnetStack_RegisterCallbackGetPropertyCharacterString((*[0]byte)(C.fpCallbackGetPropertyCharString))
	C.BACnetStack_RegisterCallbackGetPropertyReal((*[0]byte)(C.fpCallbackGetPropertyReal))
}
