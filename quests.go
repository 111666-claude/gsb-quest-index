// Package quests 是任务需求索引：按物品需求建倒排索引，玩家获得物品时只触碰
// 相关任务，产出本次新可完成的任务；需求按获得数量累加，重复登记幂等。
package quests

import "sort"

type quest struct {
	id    string
	needs map[string]int
}

type playerState struct {
	owned     map[string]int
	ready     map[string]bool
	remaining map[string]map[string]int // 任务 id -> 物品 -> 还差几件
	unmet     map[string]int            // 任务 id -> 尚未满足的需求种数
	pending   map[string][]string       // 物品 -> 该玩家尚未可完成的任务 id
	seen      map[string]int            // 物品 -> 已并入 pending 的倒排索引长度
}

// Board 是任务需求台账。
type Board struct {
	quests  map[string]*quest
	index   map[string][]string // 物品 -> 需要它的任务 id（只增不删）
	players map[string]*playerState
	scanned int
}

// New 建空台账。
func New() *Board {
	return &Board{
		quests:  map[string]*quest{},
		index:   map[string][]string{},
		players: map[string]*playerState{},
	}
}

// AddQuest 登记一个任务及所需物品数，同一个任务 id 重复登记返回 false 且不覆盖已有需求。
func (b *Board) AddQuest(id string, need map[string]int) bool {
	if _, ok := b.quests[id]; ok {
		return false
	}
	needs := map[string]int{}
	for item, count := range need {
		if count <= 0 {
			continue
		}
		needs[item] = count
	}
	b.quests[id] = &quest{id: id, needs: needs}
	for item := range needs {
		b.index[item] = append(b.index[item], id)
	}
	return true
}

// Acquire 记录该玩家获得 count 件物品，返回本次新可完成的任务 id（升序）。
func (b *Board) Acquire(player, item string, count int) []string {
	state := b.players[player]
	if state == nil {
		state = &playerState{
			owned:     map[string]int{},
			ready:     map[string]bool{},
			remaining: map[string]map[string]int{},
			unmet:     map[string]int{},
			pending:   map[string][]string{},
			seen:      map[string]int{},
		}
		b.players[player] = state
	}
	state.owned[item] += count
	if full := b.index[item]; state.seen[item] < len(full) {
		state.pending[item] = append(state.pending[item], full[state.seen[item]:]...)
		state.seen[item] = len(full)
	}
	pending := state.pending[item]
	if len(pending) == 0 {
		return nil
	}
	b.scanned += len(pending)
	out := []string{}
	kept := pending[:0]
	for _, id := range pending {
		remaining, ok := state.remaining[id]
		if !ok {
			remaining = map[string]int{}
			for needItem, needCount := range b.quests[id].needs {
				remaining[needItem] = needCount
			}
			state.remaining[id] = remaining
			state.unmet[id] = len(remaining)
		}
		if left, ok := remaining[item]; ok {
			left -= count
			if left <= 0 {
				delete(remaining, item)
				state.unmet[id]--
			} else {
				remaining[item] = left
			}
		}
		if state.unmet[id] == 0 {
			state.ready[id] = true
			delete(state.remaining, id)
			delete(state.unmet, id)
			out = append(out, id)
		} else {
			kept = append(kept, id)
		}
	}
	state.pending[item] = kept
	sort.Strings(out)
	return out
}

// Ready 返回该玩家当前可完成的任务 id（升序）。
func (b *Board) Ready(player string) []string {
	state := b.players[player]
	if state == nil {
		return nil
	}
	out := make([]string, 0, len(state.ready))
	for id := range state.ready {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// Scanned 是累计扫描的任务数（规模观测）。
func (b *Board) Scanned() int { return b.scanned }
