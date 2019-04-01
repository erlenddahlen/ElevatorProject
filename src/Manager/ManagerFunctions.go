package Manager

//
import (
	"../DataStructures"
	"../FSM"
	"../elevio"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"sort"
	"strconv"
)

func ManagerInit(GState DataStructures.GlobalState, id string) DataStructures.GlobalState {
	DataStructures.HasBackup = false
	var _, err1 = os.Stat(DataStructures.Backupfilename)

	if os.IsNotExist(err1) {

		//Set up elevator in unknown state
		Queue1 := [4][3]bool{{false, false, false}, {false, false, false}, {false, false, false}, {false, false, false}}
		ElevState := DataStructures.Elev{DataStructures.Unknown, DataStructures.MotorDirUp, 0, Queue1}
		GState.Map = make(map[string]DataStructures.Elev)
		GState.HallRequests = [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
		GState.Id = id
		GState.Map[GState.Id] = ElevState

		// Create backupfile
		var file, err1 = os.Create(DataStructures.Backupfilename)
		isError(err1)
		defer file.Close()
		GStatejason, _ := json.MarshalIndent(GState, "", "")
		_ = ioutil.WriteFile(DataStructures.Backupfilename, GStatejason, 0644)
		return GState
	} else {
		DataStructures.HasBackup = true
		GStateByte, err := ioutil.ReadFile(DataStructures.Backupfilename)
		if err != nil {
			fmt.Print(err)
		}
		error := json.Unmarshal(GStateByte, &GState)
		if error != nil {
			fmt.Println("error:", error)
		}
		GState.Map[GState.Id] = DataStructures.Elev{DataStructures.Unknown, DataStructures.MotorDirUp, 0, GState.Map[GState.Id].Queue}
		return GState
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func chooseElevator(elevators map[string]DataStructures.Elev, NewOrder DataStructures.ButtonEvent, GState DataStructures.GlobalState) string {
	times := make([]int, len(elevators))
	indexBestElevator := 0
	keyBestElevator := GState.Id

	i := 0
	keys := make([]int, 0)
	for key, _ := range elevators {
		keyInt, error := strconv.Atoi(key)
		if error != nil {
			fmt.Println(error.Error())
		}
		keys = append(keys, keyInt)
	}
	sort.Ints(keys)
	for _, v := range keys {
		keyString := strconv.Itoa(v)
		times[i] = timeEstimateToServeOrder(elevators[keyString], NewOrder)
		if times[i] <= times[indexBestElevator] {
			indexBestElevator = i
			keyBestElevator = keyString
		}
		i++
	}
	return keyBestElevator
}

func timeEstimateToServeOrder(e DataStructures.Elev, button DataStructures.ButtonEvent) int {
	tempElevator := e
	tempElevator.Queue[button.Floor][button.Button] = true

	timeUsed := 0

	switch tempElevator.State {
	case DataStructures.Idle:
		tempElevator.Dir = FSM.GetNextDir(tempElevator)
		if tempElevator.Dir == DataStructures.MotorDirStop {
			return timeUsed
		}
	case DataStructures.Moving:
		timeUsed += DataStructures.TravelTime / 2
		tempElevator.Floor += int(tempElevator.Dir)
	case DataStructures.DoorOpen:
		timeUsed += DataStructures.DoorOpenTime / 2
	}
	count := 0
	for {
		if FSM.ShouldStop(tempElevator) {
			if tempElevator.Floor == button.Floor {
				return timeUsed
			}
			tempElevator.Queue[tempElevator.Floor][DataStructures.BT_Cab] = false
			if tempElevator.Dir == DataStructures.MotorDirUp || !FSM.OrderAbove(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][DataStructures.BT_HallUp] = false
			}
			if tempElevator.Dir == DataStructures.MotorDirDown || !FSM.OrderBelow(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][DataStructures.BT_HallDown] = false
			}
			timeUsed += DataStructures.DoorOpenTime
			tempElevator.Dir = FSM.GetNextDir(tempElevator)
		}
		tempElevator.Floor += int(tempElevator.Dir)
		timeUsed += DataStructures.TravelTime
		count = count + 1
	}
}
