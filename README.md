# quest-index

任务需求索引：按物品需求索引任务，玩家获得物品时产出新可完成的任务，只用 Go 标准库。

```
go run ./cmd/quest-index --sample=exact
go run ./cmd/quest-index --sample=again
go run ./cmd/quest-index --sample=dup
go run ./cmd/quest-index --sample=multi
go run ./cmd/quest-index --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **需求计数**：一次获得 `count` 件就按 `count` 累加，玩家拥有数不小于需求数才算满足。
- **可完成**：任务的每项需求都满足才进入可完成集合；返回结果按任务 id 升序。
- **增量**：`Acquire` 只返回本次新进入可完成集合的任务，已经可完成的不再重复返回。
- **幂等**：同一个任务 id 重复登记返回 false，且不覆盖已有需求。
- **代价**：`Acquire` 与 `Ready` 不许每次全扫任务表，`Scanned` 不随任务数乘调用次数增长；
  10 万个任务、每秒 10 万次获得。

## 输出契约（不改格式）

```
ready=..
second=..
second=.. ready=..
scanned=..
```

## 目录

```
quests.go            需求计数、可完成集合与幂等登记
quests_test.go       单元测试
cmd/quest-index/     命令行入口
```
