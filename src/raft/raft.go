package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import (
	"labrpc"
	"time"
	"strconv"
)

// import "bytes"
// import "encoding/gob"

const (
	LEADER = iota
	FOLLOWER
	CANDIDATE
)


//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu           sync.Mutex
	peers        []*labrpc.ClientEnd
	persister    *Persister
	me           int // index into peers[]


	state        int

	currentTerm  int
	votedFor     int
	logEntry     []*LogEntry

									 // Volatile state on all servers
	commitIndex  int
	lastApplied  int

									 // Volatile state on leaders
	nextIndex    []int
	matchIndex   []int

									 // Number of votes received currently
	voteReceived int

	// Channels
	chanWinElection chan bool
	chanHeartbeat chan bool
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	return rf.currentTerm, rf.state == LEADER
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here.
	// Example:
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here.
	// Example:
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.xxx)
	// d.Decode(&rf.yyy)
}



//
// RequestVote RPC handler. Be initiate by rf.sendRequestVote()
// When a vote request received:
//   1. Approve the first request with valid term
//	 2. Reject all the following requests
//
func (rf *Raft) RequestVoteRPC(args RequestVoteArgs, reply *RequestVoteReply) {
	println("SERVER " + strconv.Itoa(rf.me) + " @term" + strconv.Itoa(rf.currentTerm) + ": Vote request from server" + strconv.Itoa(args.CandidateId) + " @term" + strconv.Itoa(args.Term))
	if rf.currentTerm > args.Term || (rf.currentTerm == args.Term && rf.votedFor != -1){
		// Outdated candidates. Simply reject. or
		// Already voted in this round of election
		reply.Term = rf.currentTerm
		reply.Vote = false
	} else {
		// First valid request received. Approve this request

		if rf.state == LEADER {
			println("Server " + strconv.Itoa(rf.me) + "@term" + strconv.Itoa(rf.currentTerm) + " is no longer a leader 2")
		}
		rf.mu.Lock()
		rf.state = FOLLOWER
		rf.voteReceived = 0
		rf.votedFor = args.CandidateId
		rf.currentTerm = args.Term
		rf.mu.Unlock()
		reply.Vote = true
	}
}

//
// AppendEntries RPC Handler. Called by rf.sendAppendEntries()
// When an appendEntries received:
//   1. Reject outdated message
//   2. Move to follower state if is a candidate
//   3. Confirm receive of heartbeat message and update term
//
func (rf *Raft) AppendEntriesRPC(args AppendEntriesArgs, reply *AppendEntriesReply) {
	if args.Term < rf.currentTerm {
		// Outdated heartbeat message
		reply.Term = rf.currentTerm
		reply.Outdated = true
		return
	}
	reply.Outdated = false
	rf.currentTerm = args.Term

	if rf.state == CANDIDATE {
		rf.state = FOLLOWER
		rf.init()
	}

	rf.chanHeartbeat <- true
	rf.currentTerm = args.Term
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//
// The background go routine
//
func (rf *Raft) run() {
	for {
		// the state read can be in a race condition
		// we have to check again and do the right thing in each rule.
		switch rf.state {
		case FOLLOWER:
			select {
				case <-rf.chanHeartbeat:  // Receive heartbeat message
					DPrintf("Received hearbeat message", rf.me)
				case <-time.After(TIMEOUT_DURATION):
						rf.mu.Lock()
						if rf.state != LEADER {
							rf.state = CANDIDATE
						}
						rf.mu.Unlock()
				}
		case CANDIDATE:
			rf.doCandidateStuff()
		case LEADER:
			rf.doLeaderStuff()
		}
	}
}

func (rf *Raft) init() {
	rf.votedFor = -1;
	rf.voteReceived = 0;
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//


	func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {

		rf := &Raft{
			peers:     peers,
			persister: persister,
			me:        me,

			state: FOLLOWER,

			// persistent state on all servers
			currentTerm: 0,
			votedFor:    -1, // initialize as -1 for not voted yet
			logEntry:         []*LogEntry{{Term: 0}},

			// volatile state on all servers
			commitIndex: 0,
			lastApplied: 0,

			// volatile state on leaders
			nextIndex:  make([]int, len(peers)),
			matchIndex: make([]int, len(peers)),

			// volatile state on candidates
			voteReceived: 0,

			// Channels
			chanWinElection: make(chan bool),
			chanHeartbeat: make(chan bool),
		}

		// TODO: Find out what persistent means..

		go rf.run() // Start background go routine

		// initialize from state persisted before a crash
		rf.readPersist(persister.ReadRaftState())


		return rf
	}
