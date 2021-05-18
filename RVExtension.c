#include <stdlib.h>

#include "extensionCallback.h"

extern void goRVExtension(char *output, int outputSize, char *input);
extern void goRVExtensionVersion(char *output, int outputSize);
extern int goRVExtensionArgs(char* output, int outputSize, char* input, char** argv, int argc);
// extern void goRVExtensionRegisterCallback(extensionCallback fnc);

//--- Called by Engine on extension load 
__attribute__((dllexport)) void RVExtensionVersion(char *output, int outputSize);
//--- STRING callExtension STRING
__attribute__((dllexport)) void RVExtension(char *output, int outputSize, char *input);
//--- STRING callExtension ARRAY
__attribute__((dllexport)) int RVExtensionArgs(char *output, int outputSize, char* input, char** argv, int argc);



void RVExtension(char *output, int outputSize, char *input)
{
	goRVExtension(output, outputSize, input);
}

void RVExtensionVersion(char *output, int outputSize)
{
	goRVExtensionVersion(output, outputSize);
}


int RVExtensionArgs(char *output, int outputSize, char* input, char** argv, int argc)
{
	return goRVExtensionArgs(output, outputSize, input, argv, argc);
}


// __declspec(dllexport) void RVExtensionRegisterCallback(extensionCallback fnc) {
// 	goRVExtensionRegisterCallback(fnc);
// }

