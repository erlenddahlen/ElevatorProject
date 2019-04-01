package Manager

//
import (
	"../DataStructures"
	"../FSM"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func managerInit(GState DataStructures.GlobalState, id string) DataStructures.GlobalState {
	DataStructures.HasBackup = false
	var _, err1 = os.Stat(DataStructures.BackupFilename)

	if os.IsNotExist(err1) {

		//Set up elevator in unknown state
		Queue1 := [4][3]bool{{false, false, false}, {false, false, false}, {false, false, false}, {false, false, false}}
		ElevState := DataStructures.Elev{DataStructures.Unknown, DataStructures.MotorDirUp, 0, Queue1}
		GState.Map = make(map[string]DataStructures.Elev)
		GState.HallRequests = [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
		GState.Id = id
		GState.Map[GState.Id] = ElevState

		// Create backupfile
		var file, err1 = os.Create(DataStructures.BackupFilename)
		isError(err1)
		defer file.Close()
		GStatejason, _ := json.MarshalIndent(GState, "", "")
		_ = ioutil.WriteFile(DataStructures.BackupFilename, GStatejason, 0644)
		return GState
	} else {
		DataStructures.HasBackup = true

		//Read backup into GlobalState
		GStateByte, err := ioutil.ReadFile(DataStructures.BackupFilename)
		if err != nil {
			fmt.Print(err)
		}
		error := json.Unmarshal(GStateByte, &GState)
		if error != nil {
			fmt.Println("error:", error)
		}

		//Set the state to Unknown
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
		tempElevator.Dir = FSM.getNextDir(tempElevator)
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
		if FSM.shouldStop(tempElevator) {
			if tempElevator.Floor == button.Floor {
				return timeUsed
			}
			tempElevator.Queue[tempElevator.Floor][DataStructures.Cab] = false
			if tempElevator.Dir == DataStructures.MotorDirUp || !FSM.orderAbove(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][DataStructures.HallUp] = false
			}
			if tempElevator.Dir == DataStructures.MotorDirDown || !FSM.orderBelow(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][DataStructures.HallDown] = false
			}
			timeUsed += DataStructures.DoorOpenTime
			tempElevator.Dir = FSM.getNextDir(tempElevator)
		}
		tempElevator.Floor += int(tempElevator.Dir)
		timeUsed += DataStructures.TravelTime
		count = count + 1
	}
}
