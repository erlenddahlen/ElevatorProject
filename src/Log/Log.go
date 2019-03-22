
package main

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"io"
)

// Channels in: newLog, mergeLog, logPos, requestLogInt
// Channels out: fullLog

func isError(err error) bool {
    if err != nil {
        fmt.Println(err.Error())
    }

    return (err != nil)
}


type Log struct {
	HallMat	  [4][2]int
	CabMat	  [4][3]int
}
//func logSystem(chLog CommTest.LogChannels, chComm CommTest.CommunicationChannels, )


func main() {

	// CASE: newLog
	// Sending Log to backupfile
	log := Log {}
	filename := "backupfile.txt"
	// Delete/create file again
	var _, err1 = os.Stat(filename)
	if os.IsNotExist(err1) {
			var file, err1 = os.Create(filename)
			if isError(err1) {
					return
			}
			defer file.Close()
	} else {
			var err = os.Remove(filename)
			var file, err1 = os.Create(filename)
			if isError(err) || isError(err1){
					return
				}
				defer file.Close()
			}

	fmt.Println("File Created Successfully", filename)

	// Write Log to file in json-format
	file, _ := json.MarshalIndent(log, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

// CASE: requestLogInt, out: fullLog
	// Read file
	var readFile, err2 = os.OpenFile(filename, os.O_RDWR, 0644)
			if isError(err2) {
					return
			}
			defer readFile.Close()

	var text = make([]byte, 1024)

 for {
		 _, err3 := readFile.Read(text)

		 // Break if finally arrived at end of file
		if err3 == io.EOF {
				 break
		 }

		 // Break if error occured
		 if err3 != nil && err3 != io.EOF {
				 isError(err3)
				 break
		 }
 }

 fmt.Println("Reading from file.")
 fmt.Println(string(text))

	// Undo json
	var new = []byte(string(file))
	//fmt.Println(new)
	var tests Log
	error := json.Unmarshal(new, &tests)
	if error != nil {
		fmt.Println("error:", error)
		}
		fmt.Println(tests)
}
