package raft

import (
	"time"
	"strconv"
)

//
// This defines what a candidates do
// 1. Increment current term
// 2. vote for itself
// 3. Request votes from all others servers
//
func (rf *Raft) doCandidateStuff() {
	rf.currentTerm ++
	rf.init()

	rf.votedFor = rf.me
	rf.voteReceived = 1

	rf.broadcastRequestVote()

	timeout := GetRandTimeoutDuration()

	select {
		case <-rf.chanWinElection:
		// println("Empty BecomeLeaderCH")
		default:
		// println("BecomeLeaderCH is empty, carry on.")
	}

	// Waiting for results
	select {
	case win := <-rf.chanWinElection: // Elected as leader
		if win {
			println("Server " + strconv.Itoa(rf.me) + "@term " + strconv.Itoa(rf.currentTerm) + " Win election!")
		}
	case <-time.After(timeout): // Election time out for this round
	}

}

func (rf *Raft) broadcastRequestVote() {
	println("Server " + strconv.Itoa(rf.me) + "@term " + strconv.Itoa(rf.currentTerm) + " start election!")
	args := &RequestVoteArgs{
		Term: rf.currentTerm,
		CandidateId: rf.me,
	}
	reply := &RequestVoteReply{
		Term: rf.currentTerm,
	}

	// send request vote to all other servers
	go func() {
		for serverIndex := 0; serverIndex < len(rf.peers); serverIndex++ {
			if serverIndex != rf.me {
				rf.sendRequestVote(serverIndex, *args, reply)
				if reply.Term > rf.currentTerm {
					rf.mu.Lock()
					if reply.Term > rf.currentTerm {
						rf.state = FOLLOWER // term outdated. Revert back to follower
						rf.init()
					}
					rf.mu.Unlock()
					break
				}
			}
			if rf.state != LEADER {
				break
			}
		}
	}()
}


//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with &&RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// returns true if labrpc says the RPC was delivered.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args RequestVoteArgs, reply *RequestVoteReply) {
	for {
		ok := rf.peers[server].Call("Raft.RequestVoteRPC", args, reply)
		if ok {
			println("Server " + strconv.Itoa(rf.me) + "@term " + strconv.Itoa(rf.currentTerm) + " receive vote of server " + strconv.Itoa(server))
			if reply.Vote {
				rf.mu.Lock()
				if (reply.Vote && reply.Term == rf.currentTerm) {
					rf.voteReceived ++
				}
				rf.voteReceived ++
				rf.mu.Unlock()
				if rf.state == CANDIDATE && rf.voteReceived > len(rf.peers) / 2 {
					rf.state = LEADER
					rf.chanWinElection <- true
				}

			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
}
