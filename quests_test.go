package quests

import (
	"fmt"
	"testing"
)

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

func TestAcquireCountsByCount(t *testing.T) {
	board := New()
	board.AddQuest("q1", map[string]int{"ore": 2})
	if got := board.Acquire("p", "ore", 2); len(got) != 1 || got[0] != "q1" {
		t.Fatalf("一次获得 count 件应按 count 累加，得到 %v", got)
	}

	board = New()
	board.AddQuest("q2", map[string]int{"ore": 3})
	if got := board.Acquire("p", "ore", 2); len(got) != 0 {
		t.Fatalf("累计 2/3 时不该可完成，得到 %v", got)
	}
	if got := board.Acquire("p", "ore", 1); len(got) != 1 || got[0] != "q2" {
		t.Fatalf("累计到 3 件应该可完成，得到 %v", got)
	}
}

func TestReadyIsIncremental(t *testing.T) {
	board := New()
	board.AddQuest("q1", map[string]int{"ore": 1})
	board.Acquire("p", "ore", 1)
	if got := board.Acquire("p", "ore", 1); len(got) != 0 {
		t.Fatalf("已经可完成的任务不该重复返回，得到 %v", got)
	}
	if got := board.Ready("p"); len(got) != 1 || got[0] != "q1" {
		t.Fatalf("已可完成的任务应留在集合里，得到 %v", got)
	}
}

func TestDuplicateAddQuestIsRejected(t *testing.T) {
	board := New()
	if !board.AddQuest("q1", map[string]int{"ore": 1}) {
		t.Fatal("首次登记应该成功")
	}
	if board.AddQuest("q1", map[string]int{"ore": 2}) {
		t.Fatal("重复登记应该返回 false")
	}
	if got := board.Acquire("p", "ore", 1); len(got) != 1 || got[0] != "q1" {
		t.Fatalf("重复登记不该覆盖已有需求，得到 %v", got)
	}
}

func TestReadyResultsAreSorted(t *testing.T) {
	board := New()
	board.AddQuest("q2", map[string]int{"ore": 1})
	board.AddQuest("q1", map[string]int{"ore": 1})
	got := board.Acquire("p", "ore", 1)
	if len(got) != 2 || got[0] != "q1" || got[1] != "q2" {
		t.Fatalf("返回结果应按任务 id 升序，得到 %v", got)
	}
}

func TestScannedDoesNotGrowWithCalls(t *testing.T) {
	board := New()
	for index := 0; index < 3000; index++ {
		board.AddQuest(fmt.Sprintf("q-%d", index), map[string]int{"ore": 1})
	}
	for index := 0; index < 3000; index++ {
		board.Acquire("p", "ore", 1)
	}
	if got := board.Scanned(); got > 6000 {
		t.Fatalf("Scanned 不该随任务数乘调用次数增长，得到 %d", got)
	}
}
