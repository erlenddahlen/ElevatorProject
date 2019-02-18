numbFloors = 4

type Dir int
struct elevator {
  dir Dir = 0
  floor = 0
  dist_order outsideOrder
}
struct outsideOrder {
  floor int = 0
  dir Dir = 0
}

var elevators [3]elevator

func shortestPath(elevators, order) {
  for elev in elevators {
    if elev.outsideOrder.floor
  }
}
func dist(fromFloor int, fromDir Dir, toFloor int, toDir Dir) {
  toFloorIndex := 0
  fromFloorINdex := 0
  if(toDir == UP) {
    toFloorIndex = toFloor + numbFloors
  } else {
    toFloorIndex = numbFloors - toFloor
  }
  return (toFloor- fromFloor)%(numbFloors-1)*2
}
