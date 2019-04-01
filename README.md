TTK4145 - Elevator Project
================================

The system consists of three modules and two helping packages, which all serves specific tasks to accomplish the main goal of creating software for controlling `n` elevators working in parallel across `m` floors. The architecture is constructed on a peer-to-peer concept where all the peers on the network cooperate to execute orders. The idea is that all the peers always have the same and latest information about each other and about the orders awaiting to be executed. With this assumption they can decide which of the other peers that should handle a specific order, by optimizing a cost function based on the state of the elevator of the peer. The chosen peer is then responsible for this order.

In addition, the system also have functionality to handle error. This includes spamming the network to achieve correctness of information, backup files and a watchdog to monitor elevator motor stop.

This code is specifically made for three elevators and four floors.


Module descriptions
---------------------
###  FSM

The FSM module is the Finite State Machine for controlling one elevator. It consists of four states and three main events:

State| Event
------------ | -------------
Unknown| Initial
Idle| Update of Global State from Manager
DoorOpen| At a floor
Moving| Door timer finished



The elevator starts at state Unknown before it finds out which floor it is located, then it returns to state Idle. From Idle it will wait for the event of an update from Manager, update its queue and search it to find which order it shall execute (#1). It finds its direction and goes to state Moving. When the elevator arrives at a floor this event triggers it to check if it should stop(#1). If the elevator is in the floor of the order it is executing it will stop and start a door timer of 3 sec. The elevator goes to DoorOpen state, and waits for the timer to finish. After finishing timer, the elevator goes back to Idle again and waits for a new update from Manager. Through this whole process the elevator sends updates to the Manager of its state, floor, direction and queue.

The FSM module also have a GO routine for checking if buttons are pushed and if the elevator arrives at a new floor.

(#1) The functions in FSMFunctions.go handles these operations.

###  Manager

The Manager module is responsible for distribution of correct and newest data between other peers and its own FSM module. This data is called the GState (GlobalState) which consists of all the states of all the elevators in the network, an id and the common hall orders(#2). It cooperates with the Network module for transmitting and receiving the newest GState every time a change in any of the peers happens(#3).

It is also responsible to decide which peer that should execute an order, and if it results in its own elevator it has to send this information to its FSM.

Since it has all the information about everything in the system it also has the responsibility of monitoring which elevators that are cooperative. It does this by checking if a peer disappears from the network or if it has a motor stop(#4).


(#2) Happens in the UpdateGlobalState function
(#3) Happens in the SpamGlobalState function
(#4) Happens in the functions UpdateNetworkPeers and MotorstopWatchdog

###  Network
The Network module enables communication over the network through UDP broadcasting. The module is used directly from the project resources. Follow the link for more information on its features.
[GitHub: Network Go](https://github.com/TTK4145/Network-go)


###  Extra packages
#### Elevio
Handles the hardware of the system. Initializing hardware, handling motor, buttons, lights, floor sensor etc.


#### DataStructures
Defining constants and types of structs.

Error handling functionality
---------------------------
Error  | Description of solution
------------ | -------------
Packet loss| For this system it is extremely important that all the peers on the network always has the same data. This is important for guaranteeing that they all get the same result of which elevator that should serve a specific order. To be able to guarantee correct and similar information the system will spam the network with the newest information at all times. Even though we may have packet loss, the probability for this gets small when the data is sent frequently.
Motor stop or lost peer| To prevent that orders get lost the system has functionality to monitor weather an elevator has motor stop or if an elevator is disconnected from the system. A watchdog checks updates from all elevators to decide if a motor error occurs. This is done by finding out if an elevator has an order in its queue but don't changes state or floor within 6 sec. If a motor error occurs at another elevator, the elevator that discovered the error takes all the hall orders. The network module sends frequently an updated list of peers in the network, if a peer disconnects all its hall orders are redistributed.
Elevator shut down| If an elevator shuts down the orders in its queue will disappear. To prevent these orders from not being executed a backup log file is frequently updated with the latest hall orders given to the elevator. When an elevator turns on again it starts executing the orders in this file.
