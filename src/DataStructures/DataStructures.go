package DataStructures

// Config constants
const (
	NumFloors                = 4
	NumButtons               = 3
	Timeout                  = 5
	TravelTime               = 5
	DoorOpenTime             = 3
	NumOfElev                = 3
	SpamTime                 = 1000 //ms
	PortNumber               = 16700
	MotorstopWatchdogTimeout = 6 //s
)

var BackupFilename string
var HasBackup bool

type ElevState int

const (
	Unknown  ElevState = iota - 1
	Idle               = 0
	DoorOpen           = 1
	Moving             = 2
)

type MotorDirection int

const (
	MotorDirUp   MotorDirection = 1
	MotorDirDown                = -1
	MotorDirStop                = 0
)

/*
Structure of Queue
Floor 			Button
	1.	Up	|		-		|	Cab
	2.	Up	|	Down	|	Cab
	3.	Up	|	Down	|	Cab
	4.	-		|	Down	|	Cab
*/
type Elev struct {
	State ElevState
	Dir   MotorDirection
	Floor int
	Queue [NumFloors][NumButtons]bool
}

/*
Structure of HallRequests
Floor 	Button
	1.	Up	|	-
	2.	Up	|	Down
	3.	Up	|	Down
	4.	-		|	Down
*/
type GlobalState struct {
	Map          map[string]Elev
	HallRequests [4][2]bool
	Id           string
}

type ButtonType int

const (
	HallUp   ButtonType = 0
	HallDown            = 1
	Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type FSMChannels struct {
	AtFloor            chan int
	UpdateFromFSM      chan Elev
	UpdateFromManager  chan GlobalState
	ButtonPushed       chan ButtonEvent
	AddCabOrder        chan int
	AddCabOrderManager chan int
}

type ManagerChannels struct { //decalred in config file
	InternalState     chan GlobalState
	ExternalState     chan GlobalState
	LostElev          chan string
	AddHallOrder      chan ButtonEvent
	UpdatefromSpam    chan GlobalState
	MotorstopWatchdog chan int
}
