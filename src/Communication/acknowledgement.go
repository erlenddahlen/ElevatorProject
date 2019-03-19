//ACKING:
//Use prefix to determine if hash or msg
//Hashes should be directly sent to ack-thread
//Use global datastruct
//Use channelID/counter
//Universal timer used to resend messages from datastruct

//All incoming chans in NH --> send ack
//All outgoing chans in NH --> list of ACK, resend



package main

import (
	"fmt"
	"crypto/md5"
	"encoding/json"
)

type ElevPos struct {
	Floor  int
	Id     int
}

type ack struct {
  chanNum int
  msg string
}



var m map[string]ack

func ackhandler(){
  //declaring variables
  m = make(map[string]ack)

  for{
      select{
        case chanNum, msg <- messages from NH(in json.Marshall format):
        //out, err:= json.Marshal(test)
      	//if err != nil {
        //      panic (err)
        //}
        h := md5.New()
      	h.Write([]byte(msg))
        newack := ack{1, string(out)}
        hash := hex.EncodeToString(h.Sum(nil))

        //add item to map
        m[hash] = newack

      case chanNum, hash <- hashes from NH:
        //need to compare hash AND chanNum!
        for k, v := range m {
		        if (v.chanNum == chanNum) && (k == hash){
                delete(m, k)
                break //quit for loop
            }
	      }

      case <- timer:
        for k, v := range m {
		        fmt.Println(k, v)
            msg := []byte(v.msg)
            switch v.chanNum {
            case 1:
              channel <- msg

            case 2:
            case 3:
            case 4:
            //and so on
            //should unmarshal in respective channels to be sure to achieve datastructs
            }
	      }
        //reset timer

      }

  }

}

///running in go online

package main

import (
	"fmt"
	"crypto/md5"
	"encoding/json"
	"encoding/hex"
)

type ack struct {
  chanNum int
  msg string
}

type ElevPos struct {
	Floor  int
	Id     int
}

var m map[string]ack

func main() {
	m = make(map[string]ack)

	test := ElevPos{1,1}
	out, err:= json.Marshal(test)
	if err != nil {
        panic (err)
  	}

	h := md5.New()
	h.Write([]byte(out))


	newack := ack{1, string(out)}
	hash := hex.EncodeToString(h.Sum(nil))


	var new = []byte(string(out))
	fmt.Println("Out:", string(out))
	fmt.Println("new:", new)

	var tests ElevPos
	error := json.Unmarshal(new, &tests)
	if error != nil {
		fmt.Println("error:", error)
	}
	fmt.Printf("%+v", tests)

	m[hash] = newack

}
