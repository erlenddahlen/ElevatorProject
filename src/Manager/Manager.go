package Manager

import (
	"../DataStructures"
	"../network/bcast"
	"../FSM"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var latestState DataStructures.GlobalState
var stateUpdate DataStructures.GlobalState

func SpamGlobalState(chMan DataStructures.ManagerChannels) {

	spamTicker := time.Tick(DataStructures.SpamTime * time.Millisecond)

	transmitNet := make(chan DataStructures.GlobalState)
	go bcast.Transmitter(DataStructures.PortNumber, transmitNet)

	for {
		select {
		case <-spamTicker:
			transmitNet <- latestState

		case newUpdate := <-chMan.InternalState:
			latestState = newUpdate

			//Write global state to log
			var file, err1 = os.Create(DataStructures.BackupFilename)
			isError(err1)
			GStatejason, _ := json.MarshalIndent(latestState, "", "")
			_ = ioutil.WriteFile(DataStructures.BackupFilename, GStatejason, 0644)
			file.Close()
		}
	}
}

func UpdateGlobalState(chMan DataStructures.ManagerChannels, chFSM DataStructures.FSMChannels, id string, GState DataStructures.GlobalState) {
	newOrders := false

	var hallOrder DataStructures.ButtonEvent
	for {
		chMan.InternalState <- GState

		//Ckeck if this elevator should take a new order
		if newOrders {
			keyBestElevator := chooseElevator(GState.Map, hallOrder, GState)
			if keyBestElevator == GState.Id {
				var x = GState.Map[GState.Id]
				x.Queue[hallOrder.Floor][hallOrder.Button] = true
				GState.Map[GState.Id] = x
				chFSM.UpdateFromManager <- GState
			}
			newOrders = false
		}

		//After initialization with backup, make sure FSM take the old orders
		if DataStructures.HasBackup {
			DataStructures.HasBackup = false
			go func(chFSM DataStructures.FSMChannels) {
				time.Sleep(3 * time.Second)
				chFSM.UpdateFromManager <- GState
			}(chFSM)
		}

		select {
		case Update := <-chMan.ExternalState:
			otherElev := Update.Map[Update.Id]
			GState.Map[Update.Id] = otherElev

			// Turn off floor light where the other elevator is
			if otherElev.State == DataStructures.DoorOpen {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					GState.HallRequests[otherElev.Floor][button] = false
				}
			}

			// Check for new hall orders and distribute responsibility
			newOrder := DataStructures.ButtonEvent{}
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if Update.HallRequests[floor][button] && !GState.HallRequests[floor][button] {
						GState.HallRequests[floor][button] = true
						newOrder = DataStructures.ButtonEvent{floor, DataStructures.ButtonType(button)}
						keyBestElevator := chooseElevator(GState.Map, newOrder, GState)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							chFSM.UpdateFromManager <- GState
						}
					}
				}
			}
			FSM.SetHallAndCabLights(GState, GState.Map[id], id)

		case Id := <-chMan.LostElev:
			lostElev := GState.Map[Id]
			delete(GState.Map, Id)
			newOrder := DataStructures.ButtonEvent{}

			// Redistribute for lost elevator's orders
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if lostElev.Queue[floor][button] {
						newOrder = DataStructures.ButtonEvent{floor, DataStructures.ButtonType(button)}
						keyBestElevator := chooseElevator(GState.Map, newOrder, GState)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							chFSM.UpdateFromManager <- GState
						}
					}
				}
			}

		case hallOrder = <-chMan.AddHallOrder:
			GState.HallRequests[hallOrder.Floor][hallOrder.Button] = true
			newOrders = true

		case cabOrderFloor := <-chFSM.AddCabOrderManager:
			var x = GState.Map[GState.Id]
			x.Queue[cabOrderFloor][2] = true
			GState.Map[GState.Id] = x
			chFSM.UpdateFromManager <- GState

		case update := <-chFSM.UpdateFromFSM:
			tempQueue := GState.Map[GState.Id].Queue
			GState.Map[GState.Id] = update
			var x = GState.Map[GState.Id]
			x.Queue = tempQueue
			GState.Map[GState.Id] = x

			// Turn off floor light where this elevator is
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
			FSM.SetHallAndCabLights(GState, GState.Map[id], id)

		case <-chMan.MotorstopWatchdog:
			// Take responsibility for all hall orders
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if GState.HallRequests[floor][button] {
						var x = GState.Map[GState.Id]
						x.Queue[floor][button] = true
						GState.Map[GState.Id] = x
					}
				}
			}
			chFSM.UpdateFromManager <- GState
		}
	}
}

func UpdateNetworkPeers(chMan DataStructures.ManagerChannels) {
	seen := make(map[string]time.Time)
	timeOut := DataStructures.Timeout * time.Second
	receiveNet := make(chan DataStructures.GlobalState)
	go bcast.Receiver(DataStructures.PortNumber, receiveNet)

	for {
		select {
		case stateUpdate = <-receiveNet:
			if latestState.Id != stateUpdate.Id {
				seen[stateUpdate.Id] = time.Now()
				chMan.ExternalState <- stateUpdate
				chMan.UpdatefromSpam <- stateUpdate
			}

		default:
			// Ckeck timeout for other other elevators
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

func MotorstopWatchdog(chMan DataStructures.ManagerChannels, GState DataStructures.GlobalState) {
	watchdogTimer := time.NewTimer(DataStructures.MotorstopWatchdogTimeout * time.Second)
	var changedState bool
	var hallreqExists bool
	prevState := GState

	for {
		select {
		case newState := <-chMan.UpdatefromSpam:
			changedState = false
			hallreqExists = false

			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if newState.HallRequests[floor][button] {
						hallreqExists = true
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

			if changedState || !(hallreqExists) {
				watchdogTimer = time.NewTimer(DataStructures.MotorstopWatchdogTimeout * time.Second)
			}
			prevState = newState

		case <-watchdogTimer.C:
			chMan.MotorstopWatchdog <- 1
		}
	}
}
