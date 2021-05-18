package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "extensionCallback.h"
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var extensionCallbackFnc C.extensionCallback

func runExtensionCallback(name *C.char, function *C.char, data *C.char) C.int {
	return C.runExtensionCallback(extensionCallbackFnc, name, function, data)
}

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("Version 1.0")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtensionArgs
func goRVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) int {
	var offset = unsafe.Sizeof(uintptr(0))
	var out []string
	for index := C.int(0); index < argc; index++ {
		out = append(out, C.GoString(*argv))
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
	}
	temp := fmt.Sprintf("Function: %s nb params: %d params: %s!", C.GoString(input), argc, out)

	// Return a result to Arma
	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
	return 1
}

func callBackExample() {
	name := C.CString("arma")
	defer C.free(unsafe.Pointer(name))
	function := C.CString("funcToExecute")
	defer C.free(unsafe.Pointer(function))
	// Make a callback to Arma
	for i := 0; i < 3; i++ {
		time.Sleep(2 * time.Second)
		param := C.CString(fmt.Sprintf("Loop: %d", i))
		defer C.free(unsafe.Pointer(param))
		runExtensionCallback(name, function, param)
	}
}

func sendToInflux(data string) {

	fields := strings.Split(data, ",")

	host := fields[0]
	token := fields[1]
	org := fields[2]
	bucket := fields[3]
	profile := fields[4]
	locality := fields[5]
	metric := fields[6]
	value := fields[7]

	int_value, err := strconv.Atoi(value)
	client := influxdb2.NewClient(host, token)
	writeAPI := client.WriteAPI(org, bucket)

	p := influxdb2.NewPoint(metric,
		map[string]string{"profile": profile, "locality": locality},
		map[string]interface{}{"count": int_value},
		time.Now())

	// write point asynchronously
	writeAPI.WritePoint(p)

	// Flush writes
	writeAPI.Flush()

	defer client.Close()

	f, err := os.OpenFile("a3metrics.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	//logger := log.New(f, "", log.LstdFlags)
	//logger.Println(err)

}

//export goRVExtension
func goRVExtension(output *C.char, outputsize C.size_t, input *C.char) {
	// Return by default through ExtensionCallback arma handler the result
	if extensionCallbackFnc != nil {
		go callBackExample()
	} else {
		// Return a result through callextension Arma call
		temp := fmt.Sprintf("Cavmetrics: %s", C.GoString(input))
		result := C.CString(temp)
		defer C.free(unsafe.Pointer(result))
		var size = C.strlen(result) + 1
		if size > outputsize {
			size = outputsize
		}

		go sendToInflux(C.GoString(input))

		C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
	}
}

//export goRVExtensionRegisterCallback
func goRVExtensionRegisterCallback(fnc C.extensionCallback) {
	extensionCallbackFnc = fnc
}

func main() {}
