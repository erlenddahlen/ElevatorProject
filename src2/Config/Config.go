package Config

// Scaleable declaration of #floors and #elevators
const (
	NumFloors    = 4
//	NumElevators = 3
	NumButtons   = 3
    //timeOut      = 1 second
    TRAVEL_TIME = 5
    DOOR_OPEN_TIME = 3
)

type ElevState int
const (
	UNKNOWN ElevState = iota -1
	IDLE
	DOOR_OPEN
	MOVING
)

// type Dir int
// const (
// 	UP   Dir = 0
// 	DOWN     = 1
// 	STOP     = 2
// )

type Direction int
const (
	DirDown Direction = iota - 1
	DirStop
	DirUp
)

// type MotorDirection int
//
// const (
// 	MD_Up   MotorDirection = 1
// 	MD_Down                = -1
// 	MD_Stop                = 0
// )

type Elev struct {
	State ElevState
	Dir   Direction
	Floor int
	Queue [NumFloors][NumButtons]bool
}

type GlobalState struct {
	Map map[string]Elev
	HallRequests [4][2]int
	Id int
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