package Peertest
import (
    //"../Governor"
    "../Config"
    "fmt"
)


func GetNextDir(peer Config.Elev)Config.Direction  {

  if peer.Dir == Config.DirStop{
      if OrderAbove(peer.Floor, peer){
          return Config.DirUp
      } else if OrderBelow(peer.Floor, peer){
          return Config.DirDown
      } else {
          return Config.DirStop
      }
  } else if peer.Dir == Config.DirUp {
      if OrderAbove(peer.Floor, peer){
          return Config.DirUp
      } else if OrderBelow(peer.Floor , peer) {
          return Config.DirDown
      } else {
          return Config.DirStop
      }
  } else {
      if OrderBelow(peer.Floor, peer){
          return Config.DirDown
      } else if OrderAbove(peer.Floor, peer){
          return Config.DirUp
      } else {
          return Config.DirStop
      }
  }
  return Config.DirStop
}

func OrderAbove(Floor int, peer Config.Elev)bool {
  for floor:= Floor + 1; floor < Config.NumFloors; floor++{
    for button:= 0; button < Config.NumButtons-1; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func OrderBelow(Floor int, peer Config.Elev)bool {
  for floor:= 0; floor < Floor; floor++{
    for button:= 0; button < Config.NumButtons-1; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func ShouldStop(peer Config.Elev)bool {
    fmt.Println("peer.Dir: ", peer.Dir)
    fmt.Println("peer.Floor: ", peer.Floor)
    switch peer.Dir {
    case Config.DirUp:
        fmt.Println("dirUp: ", peer.Queue[peer.Floor][Config.BT_HallUp], peer.Queue[peer.Floor][Config.BT_Cab], !OrderAbove(peer.Floor, peer))
        return(peer.Queue[peer.Floor][Config.BT_HallUp] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderAbove(peer.Floor, peer))
    case Config.DirDown:
        fmt.Println("dirDown: ", peer.Queue[peer.Floor][Config.BT_HallDown],  peer.Queue[peer.Floor][Config.BT_Cab], !OrderBelow(peer.Floor, peer))
        return (peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderBelow(peer.Floor, peer))
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
