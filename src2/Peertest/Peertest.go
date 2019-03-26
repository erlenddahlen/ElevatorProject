package Peertest
import (
    //"../Governor"
    "../Config"
    //"fmt"
)


func GetNextDir(peer Config.Elev)Config.Direction  {

  if peer.Dir == Config.DirStop{
      if OrderAbove(peer){
          return Config.DirUp
      } else if OrderBelow(peer){
          return Config.DirDown
      } else {
          return Config.DirStop
      }
  } else if peer.Dir == Config.DirUp {
      if OrderAbove(peer){
          return Config.DirUp
      } else if OrderBelow(peer) {
          return Config.DirDown
      } else {
          return Config.DirStop
      }
  } else {
      if OrderBelow(peer){
          return Config.DirDown
      } else if OrderAbove(peer){
          return Config.DirUp
      } else {
          return Config.DirStop
      }
  }
  return Config.DirStop
}

func OrderAbove(peer Config.Elev)bool {
  for floor:= peer.Floor + 1; floor < Config.NumFloors; floor++{
    for button:= 0; button < Config.NumButtons-1; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func OrderBelow(peer Config.Elev)bool {
  for floor:= 0; floor < peer.Floor; floor++{
    for button:= 0; button < Config.NumButtons-1; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func ShouldStop(peer Config.Elev)bool {
    //fmt.Println("peer.Dir: ", peer.Dir)
    //fmt.Println("peer.Floor: ", peer.Floor)
    switch peer.Dir {
    case Config.DirUp:
        //fmt.Println("dirUp: ", peer.Queue[peer.Floor][Config.BT_HallUp], peer.Queue[peer.Floor][Config.BT_Cab], !OrderAbove(peer))
        return(peer.Queue[peer.Floor][Config.BT_HallUp] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderAbove(peer))
    case Config.DirDown:
        //fmt.Println("dirDown: ", peer.Queue[peer.Floor][Config.BT_HallDown],  peer.Queue[peer.Floor][Config.BT_Cab], !OrderBelow(peer))
        return (peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderBelow(peer))
    case Config.DirStop:
        // fmt.Println("dirstop: ", peer.Queue[peer.Floor][Config.BT_HallUp ], peer.Queue[peer.Floor][Config.BT_HallDown], peer.Queue[peer.Floor][Config.BT_Cab])
        // if(peer.Queue[peer.Floor][Config.BT_HallUp ] || peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab]){
        //     return true
        // }
    default:
    }
    return false
}
  // if peer.Dir == Config.DirUp{
	// 	if(peer.Queue[peer.Floor][Config.BT_HallUp ] || peer.Queue[peer.Floor][Config.BT_Cab] || !orderAbove(peer.Floor)){
	// 		return true
	// 	}
	// } else if peer.Dir == Config.DirDown {
	// 	if(peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab] || !orderBelow(peer.Floor)){
	// 		return true
	// 	}
	// } else if peer.Dir == Config.DirStop  {
	// 	if(peer.Queue[peer.Floor][Config.BT_HallUp ] || peer.Queue[peer.Floor][Config.BT_HallDown|| peer.Queue[peer.Floor][Config.BT_Cab]){
	// 		return true
	// 	}

//  if (peer.Dir == Config.DirDown && peer.Queue[peer.Floor][Config.BT_HallUp ] && !queue_order_below(peer.Floor))
//	{
	//		return true
//	} else if (peer.Dir == Config.DirUp && peer.Queue[peer.Floor][Config.BT_HallDown] && !queue_order_above(peer.Floor)){
	//		return true
//	}
