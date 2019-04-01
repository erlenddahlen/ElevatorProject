package DataStructures

// Scaleable declaration of #floors and #elevators
const (
	NumFloors 		= 4
	NumButtons    = 3
	Timeout       = 5
	TravelTime    = 5
	DoorOpenTime 	= 3
	NumOfElev     = 3
)

// Name of backupfile
var Backupfilename string
var HasBackup bool

type ElevState int

const (
	Unknown	ElevState = iota - 1
	Idle              = 0
	DoorOpen          = 1
	Moving            = 2
)

type Direction int

const (
	DirDown Direction = -1
	DirStop Direction = 0
	DirUp   Direction = 1
)

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type Elev struct {
	State ElevState
	Dir   MotorDirection
	Floor int
	Queue [NumFloors][NumButtons]bool
}

type GlobalState struct {
	Map          map[string]Elev
	HallRequests [4][2]bool
	Id           string
}

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type FSMChannels struct {
	AtFloor     chan int
	UpdateFromFSM chan Elev
	UpdateFromManger      chan GlobalState
	ButtonPushed     chan ButtonEvent
	AddCabOrder      chan int
	AddCabOrderManager   chan int
}

type ManagerChannels struct { //decalred in config file
	InternalState  chan GlobalState
	ExternalState  chan GlobalState
	LostElev       chan string
	AddHallOrder   chan ButtonEvent
	UpdatefromSpam chan GlobalState
	Watchdog       chan int
}
