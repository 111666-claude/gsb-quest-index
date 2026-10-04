package quests

import "testing"

func TestFreshBoardHasNothingReady(t *testing.T) {
	if got := New().Ready("p"); len(got) != 0 {
		t.Fatalf("空台账不该有可完成任务，得到 %v", got)
	}
}

func TestSingleNeedBecomesReady(t *testing.T) {
	board := New()
	board.AddQuest("q1", map[string]int{"ore": 1})
	if got := board.Acquire("p", "ore", 1); len(got) != 1 || got[0] != "q1" {
		t.Fatalf("需求满足后应该可完成，得到 %v", got)
	}
}

func TestPartialNeedIsNotReady(t *testing.T) {
	board := New()
	board.AddQuest("q1", map[string]int{"ore": 1, "wood": 1})
	if got := board.Acquire("p", "ore", 1); len(got) != 0 {
		t.Fatalf("还差一件时不该可完成，得到 %v", got)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
