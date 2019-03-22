
package main

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
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

	// Sending Log to backupfile
	log := Log {}
	filename := "backupfile.txt"
	// Delete file(because we want to replace it with new content/log)
	var err = os.Remove(filename)
	if isError(err) {
			return
	}

	// create file again
	var _, err1 = os.Stat(filename)
	if os.IsNotExist(err1) {
			var file, err1 = os.Create(filename)
			if isError(err1) {
					return
			}
			defer file.Close()
	}

	fmt.Println("File Created Successfully", filename)

	// Write Log to file in json-format
	file, _ := json.MarshalIndent(log, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

	// Read file
	var text = make([]byte, 1024)
 for {
		 _, err = file.Read(text)

		 // Break if finally arrived at end of file
		 if err == io.EOF {
				 break
		 }

		 // Break if error occured
		 if err != nil && err != io.EOF {
				 isError(err)
				 break
		 }
 }

 fmt.Println("Reading from file.")
 fmt.Println(string(text))

	// Undo json
		var new = []byte(string(file))
		fmt.Println(new)
		var tests Log
		error := json.Unmarshal(new, &tests)
		if error != nil {
			fmt.Println("error:", error)
		}
		fmt.Println(tests)
	}

	//for range time.Tick(500*time.Millisecond) {
		//file, err := os.Open(filename) // For read access.
		//if err != nil {
			//fmt.Println(err)
		//}
		//b := make([]byte, 4)
    	//	file.Read(b)
		//i := binary.BigEndian.Uint32(b)
		//fmt.Println(i)
		//i += 1
		//file.WriteString(fmt.Sprintf("%d\n", i))
		//file.Close()
	//}
//}
