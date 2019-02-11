package elevator

func main() {
    command := make(chan int)
    changeState := make(chan state)

}


func run() {
  for{
    select {
    case targetFloor := <-command:
      switch(state) {
        case IDLE:
          changtoMOVING()
        case DOOR_OPEN:

        case MOVING:
        }
      case floor := <-onFloor:
    }
  }
}
func changtoMOVING() {
  find_dir()
  go floorfinder()

}


func floorfinder() {
  if(on floor)
  <-onFloor
}
