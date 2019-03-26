package config

// Scaleable declaration of #floors and #elevators
const (
	NumFloors    = 4
//	NumElevators = 3
	NumButtons   = 3
    //timeOut      = 1 second
)

type ElevState int
const (
	UNKNOWN state = iota -1
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


type Elev struct {
	State ElevState
	Dir   Direction
	Floor int
	Cab [NumFloors]bool  //Queue [NumFloors][NumButtons]bool
}

type globalState struct {
	Map map[int]Elev
	HallRequests [4][2]int
	Id int
}
