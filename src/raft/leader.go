package raft

import (
	"time"
	"strconv"
)

// Simply send heartbeat messages for now(lab2-A)
func (rf *Raft) doLeaderStuff() {
	time.Sleep(20*time.Millisecond)  // Avoid too many hearbeat message in a minute
	rf.BroadcastAppendEntriesRPC()  // Start broadcast heartbeat messages
}


// BroadcastAppendEntries to all the servers
func (rf *Raft) BroadcastAppendEntriesRPC() {
	args := &AppendEntriesArgs{
		Term: rf.currentTerm,
		LeaderId: rf.me,
	}
	reply := &AppendEntriesReply{}

	for serverIndex := 0; serverIndex < len(rf.peers); serverIndex++ {
		// A common mistake !!!
		// By adding serverIndex as a parameter to the closure,
		// serverIndex is evaluated at each iteration and placed on the stack for the goroutine,
		// so each slice element is available to the goroutine when it is eventually executed.
		// Same as javascript
		go func(serverIndex int) {
			rf.sendAppendEntriesRPC(serverIndex, *args, reply)
		}(serverIndex)
	}
}

func (rf *Raft) sendAppendEntriesRPC(server int, args AppendEntriesArgs, reply *AppendEntriesReply) {
	for {
		ok := rf.peers[server].Call("Raft.AppendEntriesRPC", args, reply)
		if ok {
			if reply.Outdated{
				if rf.state != LEADER {
					return
				}
				println("Server " + strconv.Itoa(rf.me) + "@term" + strconv.Itoa(rf.currentTerm) + " is no longer a leader")
				println("The up-to-date term is " + strconv.Itoa(reply.Term) + " from server " + strconv.Itoa(server))
				rf.mu.Lock()
				if rf.state == LEADER {
					rf.state = FOLLOWER
					rf.init()
					rf.currentTerm = reply.Term
				}
				rf.mu.Unlock()
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}
