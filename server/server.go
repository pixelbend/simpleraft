package server

import (
	"github.com/hashicorp/raft"
	"github.com/teapartydev/simpleraft/server/rafthandlers"
	"github.com/teapartydev/simpleraft/server/storehandlers"
)

type Srv struct {
	listenAddress string
	raft          *raft.Raft
	app           *fiber.App
}

func (s Srv) Start() error {
	return s.app.Listen(s.listenAddress)
}

// New return new server
func New(listenAddr string, badgerDB *badger.DB, r *raft.Raft) *Srv {
	app := fiber.New()

	// Raft server
	raftHandler := rafthandlers.New(r)
	app.Post("/raft/join", raftHandler.JoinRaftHandler)
	app.Post("/raft/remove", raftHandler.RemoveRaftHandler)
	app.Get("/raft/stats", raftHandler.StatsRaftHandler)

	// Store server
	storeHandler := storehandlers.New(r, badgerDB)
	app.Post("/store", storeHandler.Insert)
	app.Get("/store/:key", storeHandler.Get)
	app.Delete("/store/:key", storeHandler.Delete)

	return &Srv{
		listenAddress: listenAddr,
		app:           app,
		raft:          r,
	}
}
