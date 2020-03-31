package main

import "errors"

// usableMemory is the amount of memory that the user has subscribed
var usableMemory = 1024
// memoryUsed is the amount of memory that is used so far by the user.
var memoryUsed =0
var Cloud[nPlusM][]byte

func cloudRead(index int) ([]byte,error) {
	if index < nPlusM {
		return Cloud[index], nil
	} else {
		return []byte("0"), errors.New("access beyond nPlusM")
	}
}

func cloudWrite(index int, data []byte) error {
	if index < nPlusM {
		//update the total usable memory by adding present length and subtracting previous length of the data
		memoryUsed = memoryUsed - len(Cloud[index]) + len(data)
		if memoryUsed > usableMemory{ panic("usage quota exceeding thus data not written")}
		Cloud[index]=data
		return nil
	} else {
		return errors.New("access beyond user space")
	}

}