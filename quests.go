// Package quests 是任务需求索引：按物品需求索引任务，玩家获得物品时产出新可完成的任务。
// 缺陷：没有倒排索引（每次全扫任务表）、需求按获得次数减一而不是按数量累加、
// 已经可完成的任务会被重复返回、重复登记会覆盖已有需求、查询同样全扫任务表。
package quests

import "sort"

type quest struct {
	id   string
	left map[string]int
}

type playerState struct {
	owned map[string]int
	ready map[string]bool
}

// Board 是任务需求台账。
type Board struct {
	quests  map[string]*quest
	order   []string
	players map[string]*playerState
	scanned int
}

// New 建空台账。
func New() *Board {
	return &Board{quests: map[string]*quest{}, players: map[string]*playerState{}}
}

// AddQuest 登记一个任务及所需物品数，返回是否登记成功。
func (b *Board) AddQuest(id string, need map[string]int) bool {
	left := map[string]int{}
	for item, count := range need {
		if count <= 0 {
			continue
		}
		left[item] = count
	}
	b.quests[id] = &quest{id: id, left: left}
	b.order = append(b.order, id)
	return true
}

// Acquire 记录该玩家获得物品，返回本次新可完成的任务 id（升序）。
func (b *Board) Acquire(player, item string, count int) []string {
	b.scanned += len(b.order)
	state := b.players[player]
	if state == nil {
		state = &playerState{owned: map[string]int{}, ready: map[string]bool{}}
		b.players[player] = state
	}
	state.owned[item] += count
	out := []string{}
	for _, id := range b.order {
		entry := b.quests[id]
		if state.ready[id] {
			out = append(out, id)
			continue
		}
		if _, ok := entry.left[item]; ok {
			entry.left[item]--
		}
		done := true
		for _, left := range entry.left {
			if left > 0 {
				done = false
				break
			}
		}
		if done {
			state.ready[id] = true
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// Ready 返回该玩家当前可完成的任务 id（升序）。
func (b *Board) Ready(player string) []string {
	b.scanned += len(b.order)
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
