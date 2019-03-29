package Config

// Scaleable declaration of #floors and #elevators
const (
	NumFloors = 4
	//	NumElevators = 3
	NumButtons     = 3
	TIMEOUT        = 5
	TRAVEL_TIME    = 5
	DOOR_OPEN_TIME = 3
	NumOfElev      = 3
)

// Name of backupfile
var Backupfilename string
var HasBackup bool

type ElevState int

const (
	UNKNOWN   ElevState = iota - 1
	IDLE                = 0
	DOOR_OPEN           = 1
	MOVING              = 2
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
	CurrentFloor     chan int
	LocalStateUpdate chan Elev
	PingFromGov      chan GlobalState
	ButtonPushed     chan ButtonEvent
	AddCabOrder      chan int
	AddCabOrderGov   chan int
}

type GovernorChannels struct { //decalred in config file
	InternalState  chan GlobalState
	ExternalState  chan GlobalState
	LostElev       chan string
	AddHallOrder   chan ButtonEvent
	UpdatefromSpam chan GlobalState
	Watchdog       chan int
	TakeBackup     chan GlobalState
}
