package Governor
import (
	"../Config"
	 "fmt"
	 "time"
)

//Watchdog should be own go routine
//In SpamGlobalState, send latestState to Watchdog on own channel
//Watchdog should be a case in for/select in UpdateGlobalState
//When activated, put all Hallreq in own queue

func Watchdog(gchan Config.GovernorChannels, GState Config.GlobalState){
	watchdogTimer := time.NewTimer(6 * time.Second)
	var changedState bool
	var hallreqExist bool
	var prevState Config.GlobalState
	prevState = GState


	for{
		select{
			case newState := <- gchan.UpdatefromSpam
			changedState = false
			hallreqExist = false

			// Iterate through Hallreqs, if exists orders hallreq = true
			for floor:= 0; floor < Config.NumFloors; floor++ {
				for button := 0; button < Config.NumButtons-1; button++ {
					if newState.Hallreq[floor][button] {
						hallreqExist = true
						break
					}
				}
			}

			// Iterate throug and check if change in elev floor and state between latestState
			// and newState, if exists change = true

			for key, _:= range newState.Map{
				if newState.Map[key].State != prevState.Map[key].State ||  newState.Map[key].State != prevState.Map[key].State{
					changedState = true
					break
				}
			}

			if changedState || !(hallreqExist){
					watchdogTimer = time.NewTimer(6 * time.Second)
			}
			// This could be directly set in governor. The governor takes in the timer
			case <- watchdogTimer.C
			gchan.pingfromwatch <- 1

		}
	}
}
