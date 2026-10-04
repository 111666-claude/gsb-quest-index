// Command quest-index 跑任务需求索引样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/quests"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("quest-index", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "exact", "exact / again / dup / multi / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "exact":
		board := quests.New()
		board.AddQuest("q1", map[string]int{"ore": 2})
		fmt.Fprintf(stdout, "ready=%d\n", len(board.Acquire("p", "ore", 2)))
	case "again":
		board := quests.New()
		board.AddQuest("q1", map[string]int{"ore": 1})
		board.Acquire("p", "ore", 1)
		fmt.Fprintf(stdout, "second=%d\n", len(board.Acquire("p", "ore", 1)))
	case "dup":
		board := quests.New()
		board.AddQuest("q1", map[string]int{"ore": 1})
		second := board.AddQuest("q1", map[string]int{"ore": 2})
		board.Acquire("p", "ore", 1)
		fmt.Fprintf(stdout, "second=%v ready=%d\n", second, len(board.Ready("p")))
	case "multi":
		board := quests.New()
		board.AddQuest("q1", map[string]int{"ore": 1})
		board.AddQuest("q2", map[string]int{"ore": 1})
		fmt.Fprintf(stdout, "ready=%d\n", len(board.Acquire("p", "ore", 1)))
	case "scan":
		board := quests.New()
		for index := 0; index < 3000; index++ {
			board.AddQuest(fmt.Sprintf("q-%d", index), map[string]int{"ore": 1})
		}
		for index := 0; index < 3000; index++ {
			board.Acquire("p-0", "ore", 1)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", board.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample exact|again|dup|multi|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
