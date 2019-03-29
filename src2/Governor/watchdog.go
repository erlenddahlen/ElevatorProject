

//Watchdog should be own go routine
//In SpamGlobalState, send latestState to Watchdog on own channel
//Watchdog should be a case in for/select in UpdateGlobalState
//When activated, put all Hallreq in own queue


func Watchdog(gchan Config.GovernorChannels){
	latestState
	ticker := time.Tick(6 * time.Second)
	var change bool
	var hallreqs bool

	for{
		select{
			case newUpdate := <- gchan.UpdatefromSpam
			change = false
			hallreqs = false

			//Iterate through Hallreqs
				//If exists orders hallreq = true

			//Iterate throug and check if chang in elev flor and state between latestState and newUpdate
				//If exists change = true

			if change OR !(hallreqs){
				reset ticker
			}
			case <- ticker.C
			gchan.pingfromwatch <- 1
		}

	}

}