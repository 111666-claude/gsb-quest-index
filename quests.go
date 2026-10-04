// Package quests 是任务需求索引：按物品需求建倒排索引，玩家获得物品时只触碰
// 等待该物品的任务，产出本次新可完成的任务。
package quests

import "sort"

type playerState struct {
	owned   map[string]int
	ready   map[string]bool
	pending map[string]map[string]int
	byItem  map[string]map[string]bool
	inited  bool
}

// Board 是任务需求台账。
type Board struct {
	quests  map[string]map[string]int
	players map[string]*playerState
	scanned int
}

// New 建空台账。
func New() *Board {
	return &Board{quests: map[string]map[string]int{}, players: map[string]*playerState{}}
}

// AddQuest 登记一个任务及所需物品数，返回是否登记成功。
// 同一个任务 id 重复登记返回 false，且不覆盖已有需求。
func (b *Board) AddQuest(id string, need map[string]int) bool {
	if _, dup := b.quests[id]; dup {
		return false
	}
	left := map[string]int{}
	for item, count := range need {
		if count <= 0 {
			continue
		}
		left[item] = count
	}
	b.quests[id] = left
	for _, state := range b.players {
		if state.inited {
			b.offer(state, id, left)
		}
	}
	return true
}

// offer 把一个新任务计入已初始化玩家的待满足集合，已满足的直接进入可完成集合。
func (b *Board) offer(state *playerState, id string, need map[string]int) {
	remaining := map[string]int{}
	for item, count := range need {
		if left := count - state.owned[item]; left > 0 {
			remaining[item] = left
		}
	}
	if len(remaining) == 0 {
		state.ready[id] = true
		return
	}
	state.pending[id] = remaining
	for item := range remaining {
		waiting := state.byItem[item]
		if waiting == nil {
			waiting = map[string]bool{}
			state.byItem[item] = waiting
		}
		waiting[id] = true
	}
}

// initPlayer 首次见到玩家时按当前拥有数建立待满足集合，返回直接可完成的任务。
func (b *Board) initPlayer(state *playerState) []string {
	state.inited = true
	out := []string{}
	for id, need := range b.quests {
		b.scanned++
		before := len(state.ready)
		b.offer(state, id, need)
		if len(state.ready) != before {
			out = append(out, id)
		}
	}
	return out
}

// Acquire 记录该玩家获得 count 件物品，返回本次新可完成的任务 id（升序）。
func (b *Board) Acquire(player, item string, count int) []string {
	state := b.players[player]
	if state == nil {
		state = &playerState{
			owned:   map[string]int{},
			ready:   map[string]bool{},
			pending: map[string]map[string]int{},
			byItem:  map[string]map[string]bool{},
		}
		b.players[player] = state
	}
	state.owned[item] += count
	if !state.inited {
		out := b.initPlayer(state)
		sort.Strings(out)
		return out
	}
	out := []string{}
	for id := range state.byItem[item] {
		b.scanned++
		remaining := state.pending[id]
		remaining[item] -= count
		if remaining[item] <= 0 {
			delete(remaining, item)
			delete(state.byItem[item], id)
		}
		if len(remaining) > 0 {
			continue
		}
		state.ready[id] = true
		delete(state.pending, id)
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// Ready 返回该玩家当前可完成的任务 id（升序）。
func (b *Board) Ready(player string) []string {
	out := []string{}
	state := b.players[player]
	if state == nil {
		return out
	}
	for id := range state.ready {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// Scanned 是累计扫描的任务数（规模观测）。
func (b *Board) Scanned() int { return b.scanned }
