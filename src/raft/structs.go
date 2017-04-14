package raft

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

//
// Hold information for each log entry
//
type LogEntry struct {
	Command interface{}
	Term    int // term when entry was received by leader
}

//
// Also used as heartbeat messages
//
type AppendEntriesArgs struct {
	Term     int
	LeaderId int
}

type AppendEntriesReply struct {
	Term int
	Outdated bool  // Add here in case a leader's term is out of date
}


type RequestVoteArgs struct {
	Term        int // candidate's term
	CandidateId int
}

type RequestVoteReply struct {
	Term int  // currentTerm
	Vote bool // Yes/No
}


