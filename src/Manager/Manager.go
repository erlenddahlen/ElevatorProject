package Manager

import (
	"../DataStructures"
	"../FSM"
	"../network/bcast"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var latestState DataStructures.GlobalState
var StateUpdate DataStructures.GlobalState

func SpamGlobalState(chMan DataStructures.ManagerChannels) {
	ticker := time.Tick(1000 * time.Millisecond)
	transmitNet := make(chan DataStructures.GlobalState)
	go bcast.Transmitter(16700, transmitNet)

	for {
		select {
		case <-ticker:
			transmitNet <- latestState

		case newUpdate := <-chMan.InternalState:
			latestState = newUpdate

			//Write to log
			var file, err1 = os.Create(DataStructures.Backupfilename)
			isError(err1)
			GStatejason, _ := json.MarshalIndent(latestState, "", "")
			_ = ioutil.WriteFile(DataStructures.Backupfilename, GStatejason, 0644)
			file.Close()
		}
	}
}

func UpdateGlobalState(chMan DataStructures.ManagerChannels, chFSM DataStructures.FSMChannels, id string, GState DataStructures.GlobalState) {
	distribute := false
	var hallOrder DataStructures.ButtonEvent

	for {
		chMan.InternalState <- GState

		if distribute {
			keyBestElevator := ChooseElevator(GState.Map, hallOrder, GState)
			if keyBestElevator == GState.Id {
				var x = GState.Map[GState.Id]
				x.Queue[hallOrder.Floor][hallOrder.Button] = true
				GState.Map[GState.Id] = x
				FSMchan.UpdateFromManager <- GState
			}
			distribute = false
		}

		if DataStructures.HasBackup {
			DataStructures.HasBackup = false
			go func(FSMchan DataStructures.FSMChannels) {
				time.Sleep(3 * time.Second)
				FSMchan.UpdateFromManager <- GState
			}(FSMchan)
		}

		select {
		case Update := <-chMan.ExternalState:
			OutsideElev := Update.Map[Update.Id]
			GState.Map[Update.Id] = OutsideElev

			if OutsideElev.State == DataStructures.DoorOpen {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					GState.HallRequests[OutsideElev.Floor][button] = false
				}
			}

			newOrder := DataStructures.ButtonEvent{}
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if Update.HallRequests[floor][button] && !GState.HallRequests[floor][button] {
						GState.HallRequests[floor][button] = true
						newOrder = DataStructures.ButtonEvent{floor, DataStructures.ButtonType(button)}
						keyBestElevator := ChooseElevator(GState.Map, newOrder, GState)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							FSMchan.UpdateFromManager <- GState
						}
					}
				}
			}
			Lights(GState, GState.Map[id], id)

		case Id := <-chMan.LostElev:
			lostElev := GState.Map[Id]
			delete(GState.Map, Id)
			newOrder := DataStructures.ButtonEvent{}

			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if lostElev.Queue[floor][button] {
						newOrder = DataStructures.ButtonEvent{floor, DataStructures.ButtonType(button)}
						keyBestElevator := ChooseElevator(GState.Map, newOrder, GState)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							FSMchan.UpdateFromManager <- GState
						}
					}
				}
			}

		case hallOrder = <-chMan.AddHallOrder:
			GState.HallRequests[hallOrder.Floor][hallOrder.Button] = true
			distribute = true

		case cabOrderFloor := <-FSMchan.AddCabOrderManager:
			var x = GState.Map[GState.Id]
			x.Queue[cabOrderFloor][2] = true
			GState.Map[GState.Id] = x
			FSMchan.UpdateFromManager <- GState

		case update := <-FSMchan.UpdateFromFSM:
			tempQueue := GState.Map[GState.Id].Queue
			GState.Map[GState.Id] = update
			var x = GState.Map[GState.Id]
			x.Queue = tempQueue
			GState.Map[GState.Id] = x

			if update.State == DataStructures.DoorOpen {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					GState.HallRequests[update.Floor][button] = false
				}
				var x = GState.Map[GState.Id]
				x.Queue[update.Floor][0] = false
				x.Queue[update.Floor][1] = false
				x.Queue[update.Floor][2] = false
				GState.Map[GState.Id] = x
			}
			Lights(GState, GState.Map[id], id)

		case <-chMan.Watchdog:
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if GState.HallRequests[floor][button] {
						var x = GState.Map[GState.Id]
						x.Queue[floor][button] = true
						GState.Map[GState.Id] = x
					}
				}
			}
			FSMchan.UpdateFromManager <- GState
		}
	}
}

func NetworkState(chMan DataStructures.ManagerChannels) {
	seen := make(map[string]time.Time)
	timeOut := DataStructures.Timeout * time.Second
	receiveNet := make(chan DataStructures.GlobalState)
	go bcast.Receiver(16700, receiveNet)

	for {
		select {
		case StateUpdate = <-receiveNet:
			if latestState.Id != StateUpdate.Id {
				seen[StateUpdate.Id] = time.Now()
				chMan.ExternalState <- StateUpdate
				chMan.UpdatefromSpam <- StateUpdate
			}

		default:
			for k, v := range seen {
				t := time.Now().Sub(v)
				if t > timeOut {
					chMan.LostElev <- k
					delete(seen, k)
				}
			}
		}
	}
}

func Watchdog(chMan DataStructures.ManagerChannels, GState DataStructures.GlobalState) {
	watchdogTimer := time.NewTimer(6 * time.Second)
	var changedState bool
	var hallreqExist bool
	var prevState DataStructures.GlobalState
	prevState = GState

	for {
		select {
		case newState := <-chMan.UpdatefromSpam:
			changedState = false
			hallreqExist = false

			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if newState.HallRequests[floor][button] {
						hallreqExist = true
						break
					}
				}
			}

			for key, _ := range newState.Map {
				if newState.Map[key].State != prevState.Map[key].State || newState.Map[key].Floor != prevState.Map[key].Floor {
					changedState = true
					break
				}
			}

			if changedState || !(hallreqExist) {
				watchdogTimer = time.NewTimer(6 * time.Second)
			}
			prevState = newState

		case <-watchdogTimer.C:
			chMan.Watchdog <- 1
		}
	}
}
