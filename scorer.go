package main

import (
	"log"
	"math"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

type Scorer interface {
	Score() float64
}

type DefaultScorer struct {
	S *Server
}

func (ds *DefaultScorer) Score() float64 {
	ping := float64(ds.S.Stat.AverageLatency())
	lost := ds.S.Stat.DropRate()
	return math.Pow((1-lost), 50) / ping

}

func NewDefaultScorer(server *Server) *DefaultScorer {
	return &DefaultScorer{
		S: server,
	}
}

type LuaScorer struct {
	L *lua.LState
	S *Server
}

func NewLuaScorer(server *Server, filename string) (*LuaScorer, error) {
	ls := new(LuaScorer)
	ls.S = server
	ls.L = lua.NewState()
	err := ls.L.DoFile(filename)
	if err != nil {
		return nil, err
	}
	latency := func(L *lua.LState) int {
		L.Push(lua.LNumber(ls.S.Stat.AverageLatency() / 1000000))
		return 1
	}
	droprate := func(L *lua.LState) int {

		L.Push(lua.LNumber(ls.S.Stat.DropRate()))
		return 1
	}
	address := func(L *lua.LState) int {

		L.Push(lua.LString(ls.S.Address))
		return 1
	}
	ls.L.SetGlobal("averagelatency", ls.L.NewFunction(latency))
	ls.L.SetGlobal("droprate", ls.L.NewFunction(droprate))
	ls.L.SetGlobal("address", ls.L.NewFunction(address))
	return ls, err

}
func (ls *LuaScorer) Score() float64 {
	err := ls.L.CallByParam(lua.P{
		Fn:      ls.L.GetGlobal("score"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	ret := ls.L.Get(-1)
	ls.L.Pop(1)
	if ret.Type() != lua.LTNumber {
		log.Fatal("scorer return value should be a Number")
	}
	f, err := strconv.ParseFloat(ret.String(), 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
