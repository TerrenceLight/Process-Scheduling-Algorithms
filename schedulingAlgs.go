//Terrence Light te965355
//COP 4600 Operating Systems
//Project 1: Scheduling Algorithms
package main

import (
	"bufio"
	"log"
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"
)

//Struct to store the vital information of a simulated process
type proc struct {
	name string
	arrival int
	burst int
	originalBurst int
	finishTime int
	tat int
	wait int
}

func fcfs(runfor int, procCount int, procs []proc, output *bufio.Writer) {
	//Title output
	fmt.Fprintf(output, "Using First-Come First-Served\n")
	
	//We will create the queue and an int to count the time clock and run the scheduler in a loop
	var schedQueue []proc
	selected := -1
	isRunning := false
	numScheduled := 0
	clock := 0
	
	//Scheduler runs in a loop
	for clock = 0; clock < runfor; clock++ {
		//We must first check if any process has arrived
		for j := 0; j < procCount; j++ {
			
			//If it's the appropriate arrival time for a process, add it to the queue
			if procs[j].arrival == clock {
				schedQueue = append(schedQueue, procs[j])
				numScheduled++
				fmt.Fprintf(output, "Time %3d : %s arrived\n", clock, procs[j].name)
			}
		}
		
		//If an algorithm is running, update the burst time remaining
		if isRunning == true {
			//If the selected process finishes, update to be no longer running a process
			schedQueue[selected].burst--
			if schedQueue[selected].burst == 0 {
				isRunning = false
				schedQueue[selected].finishTime = clock
				schedQueue[selected].tat = clock - schedQueue[selected].arrival
				schedQueue[selected].wait = schedQueue[selected].tat - schedQueue[selected].originalBurst
				fmt.Fprintf(output, "Time %3d : %s finished\n", clock, schedQueue[selected].name)
			}
		} 
		
		if isRunning == false {
			// If an algorithm isn't running, check the queue to see if we can run one
			if numScheduled > (selected + 1) {
				selected++
				isRunning = true
				fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", clock, schedQueue[selected].name, schedQueue[selected].burst)
			} else {
				//If there is no program available to run, then we're idle
				fmt.Fprintf(output, "Time %3d : Idle\n", clock)
			}
		}
	}
	
	//Enter wait, and turnaround time data into the process array
	//Our school server did not have golang updated, so I didn't have access to slice sorting
	for i := 0; i < procCount; i++ {
		for curProc := 0; curProc < procCount; curProc++ {
			if procs[curProc].name == schedQueue[i].name {
				procs[curProc].tat = schedQueue[i].tat
				procs[curProc].wait = schedQueue[i].wait
			}
		}
	}
	
	//Scheduling has ended, time to report final data
	fmt.Fprintf(output, "Finished at time %3d\n", clock)
	timeReport(procCount, procs, output)
}

func sjf(runfor int, procCount int, procs []proc, output *bufio.Writer) {
	//Title output
	fmt.Fprintf(output, "Using preemptive Shortest Job First\n")
	
	//We will create the queue and an int to count the time clock and run the scheduler in a loop
	var schedQueue []proc
	selected := -1
	oldSelection := -1
	clock := 0
	numScheduled := 0
	numCompleted := 0
	firstSchedFlag := true
	isRunning := false
	
	//Scheduler runs in a loop
	for clock = 0; clock < runfor; clock++ {
		//We must first check if any process has arrived
		for j := 0; j < procCount; j++ {
			
			//If it's the appropriate arrival time for a process, add it to the queue
			if procs[j].arrival == clock {
				schedQueue = append(schedQueue, procs[j])
				numScheduled++
				fmt.Fprintf(output, "Time %3d : %s arrived\n", clock, procs[j].name)
			}
		}
		
		//Kick off the process scheduling with the first arriving process
		if ((numScheduled == 1) && (firstSchedFlag == true)) {
			selected = 0
			oldSelection = 0
			fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", clock, schedQueue[selected].name, schedQueue[selected].burst)
			firstSchedFlag = false
			isRunning = true
			continue
		}
		
		//Run the currently selected algorithm
		if isRunning == true {
			schedQueue[selected].burst--
			
			//If the scheduled process finishes, we'll set it to math.MaxInt32
			if schedQueue[selected].burst == 0 {
				numCompleted++
				schedQueue[selected].burst = math.MaxInt32
				isRunning = false
				schedQueue[selected].finishTime = clock
				schedQueue[selected].tat = clock - schedQueue[selected].arrival
				schedQueue[selected].wait = schedQueue[selected].tat - schedQueue[selected].originalBurst
				fmt.Fprintf(output, "Time %3d : %s finished\n", clock, schedQueue[selected].name)
			}
		}
		
		//If the number of completed processes == the number of scheduled processes, then we're idle until another arrives, or the clock ends
		if numCompleted == numScheduled {
			fmt.Fprintf(output, "Time %3d : Idle\n", clock)
			continue
		}
		
		//Iterate through the schedule queue to select the current shortest job
		for i := 0; i < numScheduled; i++ {
			//If a process has a shorter burst time than our current process, run that one
			if schedQueue[i].burst < schedQueue[selected].burst {
				selected = i
			}
			
			//Once we've checked all of the processes and found the shortest one, select it (assuming it's a different process)
			if (((i + 1) == numScheduled) && (oldSelection != selected)) {
				oldSelection = selected
				isRunning = true
				fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", clock, schedQueue[selected].name, schedQueue[selected].burst)
			}
		}
	}
	
	//Enter wait, and turnaround time data into the process array
	//Our school server did not have golang updated, so I didn't have access to slice sorting
	for i := 0; i < procCount; i++ {
		for curProc := 0; curProc < procCount; curProc++ {
			if procs[curProc].name == schedQueue[i].name {
				procs[curProc].tat = schedQueue[i].tat
				procs[curProc].wait = schedQueue[i].wait
			}
		}
	}
	
	//Scheduling has ended, time to report final data
	fmt.Fprintf(output, "Finished at time %3d\n", clock)
	timeReport(procCount, procs, output)
}

func rr(runfor int, procCount int, procs []proc, quantum int, output *bufio.Writer) {
	//Title output
	fmt.Fprintf(output, "Using Round-Robin\n")
	fmt.Fprintf(output, "Quantum %3d\n\n", quantum)
	
	//Instead of making a queue of processes, we will be working with a queue of ints, and directly editing the procs slice
	var schedQueue []int
	clock := 0
	selected := 0
	oldSelection := 0
	numScheduled := 0
	numCompleted := 0
	quantCnt := 0
	quantLock := false
	firstSchedFlag := true
	
	//Scheduler runs in a loop
	for clock = 0; clock < runfor; clock++ {
		//We must first check if any process has arrived
		for j := 0; j < procCount; j++ {
			
			//If it's the appropriate arrival time for a process, add it to the queue
			if procs[j].arrival == clock {
				schedQueue = append(schedQueue, j)
				numScheduled++
				fmt.Fprintf(output, "Time %3d : %s arrived\n", clock, procs[j].name)
			}
		}
		
		//Kick off the process scheduling with the first arriving process
		if ((numScheduled == 1) && (firstSchedFlag == true)) {
			selected = schedQueue[0]
			oldSelection = schedQueue[0]
			fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", clock, procs[selected].name, procs[selected].burst)
			firstSchedFlag = false
			quantLock = true
			continue
		}
		
		//Run the currently selected program until it either completes, or it's quantum has ended
		if quantLock == true {
			procs[selected].burst--
			quantCnt++
			
			//If the process completes running, do not insert it back into the queue
			if procs[selected].burst == 0 {
				numCompleted++
				quantLock = false
				quantCnt = 0
				procs[selected].finishTime = clock
				procs[selected].tat = clock - procs[selected].arrival
				procs[selected].wait = procs[selected].tat - procs[selected].originalBurst
				fmt.Fprintf(output, "Time %3d : %s finished\n", clock, procs[selected].name)
				schedQueue = append(schedQueue[:0], schedQueue[1:]...)
				numScheduled--
			}
			
			//If a processes time quantum is up, but not completed put it at the end of the queue
			if quantCnt == quantum {
				quantLock = false
				quantCnt = 0
				oldSelection = selected
				schedQueue = append(schedQueue[:0], schedQueue[1:]...)
				schedQueue = append(schedQueue, oldSelection)
			}
		}
		
		//Select a process if one isn't currently locked in
		//Since all of our manipulations happen with the queue, schedQueue[0] will always be the next to select
		if quantLock == false {
			//We can only select a process if one is scheduled
			if numScheduled > 0 {
				selected = schedQueue[0]
				quantLock = true
				fmt.Fprintf(output, "Time %3d : %s selected (burst %3d)\n", clock, procs[selected].name, procs[selected].burst)
			} else {
				//If there isn't a process that we're able to schedule, we're idle
				fmt.Fprintf(output, "Time %3d : Idle\n", clock)
			}
		}
	}
	
	//Scheduling has ended, time to report final data
	fmt.Fprintf(output, "Finished at time %3d\n", clock)
	timeReport(procCount, procs, output)
}

func timeReport(procCount int, procs []proc, output *bufio.Writer) {
	fmt.Fprintf(output, "\n")
	for i := 0; i < procCount; i++ {
		fmt.Fprintf(output, "%s wait %3d turnaround %3d\n", procs[i].name, procs[i].wait, procs[i].tat )
	}
}

func main() {	
	//Grab the input file name from the command land
	//Open up the input file 
	inp, err := os.Open(os.Args[1])
	if err != nil{
		log.Fatal(err)
	}
	
	//Grab the output file name from the command line
	//Open up the output file
	out, err := os.Create(os.Args[2])
	if err != nil{
		log.Fatal(err)
	}
	
	defer inp.Close()
	defer out.Close()
	
	//Create a write buffer
	w := bufio.NewWriter(out)
	
	//We will preemptively create a processes counter to track which array slot to insert the process into
	//Declare variables
	procCount := 0
	schedAlg := 0
	numProcs := 0
	runtime := 0
	quant := 0
	var procArr []proc
	
	//Parse the file input
	//For this approach, we're going to read in all of the file before we begin processing
	scanner := bufio.NewScanner(inp)
	for scanner.Scan() {
		curLine := scanner.Text()
		
		//After we scan the line of text, we need to parse it
		parsedLine := strings.Fields(curLine)
		
		//Time to try to deal with each word in the line
		for i := 0; i < len(parsedLine); i++{
			//If we find "processcount" we need to record how many procceses the program will have
			if parsedLine[i] == "processcount" {
				i += 1
				test := parsedLine[i]
				numProcs, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, integer does not follow processcount. Ending")
					os.Exit(3)
				}
				
				_ , err = fmt.Fprintf(w, "%3d processes\n", numProcs)
				
				break
			}
			
			//If we find "runfor" we need to define the runtime for the simulation
			if parsedLine[i] == "runfor" {
				i += 1
				test := parsedLine[i]
				runtime, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, run time must be an integer. Ending")
					os.Exit(3)
				}
				
				break
			}
			
			//If we find "use" we need see which algorithm we're using
			//We'll arbitrarily assign fcfs to 1, sjf to 2, and rr to 3
			if parsedLine[i] == "use" {
				i += 1
				
				//Time to check which scheduling algorithm was found
				if parsedLine[i] == "fcfs" {
					schedAlg = 1
				}
				
				if parsedLine[i] == "sjf" {
					schedAlg = 2
				}
				
				if parsedLine[i] == "rr" {
					schedAlg = 3
				}
				
				break
			}
			
			//If we're running round robin, we need to note the time quantum
			if parsedLine[i] == "quantum" {
				i += 1
				test := parsedLine[i]
				quant, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, integer does not follow quantum. Ending")
					os.Exit(3)
				}
				
				break
			}
			
			//If we find "process" we need to create an new process inside the procArr
			if parsedLine[i] == "process" {
				//Process name will always be 2 words after process
				i += 2	
				procName := parsedLine[i]
				
				//Arrival time is then 2 words after the process name
				i += 2
				test := parsedLine[i]
				arrivalTime, err := strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, arrival time must be an integer. Ending")
					os.Exit(3)
				}
				
				//Burst time is then 2 words after the arrival time
				i += 2
				test = parsedLine[i]
				burstTime, err := strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, burst time must be an integer. Ending")
					os.Exit(3)
				}
				
				//Create the process to insert into the procArr
				testProc := proc{
					name: procName,
					arrival: arrivalTime,
					burst: burstTime,
					originalBurst: burstTime,
				}
				
				procArr = append(procArr, testProc)
				procCount += 1
				
				break
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	//Once we've finished parsing the file input, we must begin scheduling
	//First one is fcfs
	if schedAlg == 1 {
		fcfs(runtime, procCount, procArr, w)
	}
	
	//Second is sjf
	if schedAlg == 2 {
		sjf(runtime, procCount, procArr, w)
	}
	
	//Last is rr
	if schedAlg == 3 {
		rr(runtime, procCount, procArr, quant, w)
	}
	
	w.Flush()
}
