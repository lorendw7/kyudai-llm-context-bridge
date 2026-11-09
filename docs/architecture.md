# ADR-001: 保守可能な研究原型 / Maintainable research prototype

日付 / Date: 2026-11-08。状態 / Status: A案は選択済み、以下の分割・SQLite・UI詳細は設計提案。 / Route A is selected; the module split, SQLite and UI details below are design proposals.

決定者は利用者です。実装と導入は別途の受入・許可範囲に従います。進捗は[PLAN.md](../PLAN.md)のみで管理します。

The user is the decision maker. Implementation and installation follow their acceptance and authorized scope. Track progress only in [PLAN.md](../PLAN.md).

## 目的と選択肢 / Context and options

一人が学びながら保守し、Windowsで開発し、許可後に玄界で再現実験を行う原型です。主題はモデル変更時の文脈と記憶であり、新しい推論engineの開発ではありません。既存OSSを再利用し、交換箇所を限定します。根拠は[資料台帳](sources.md)に集約します。

One learner maintains the prototype, develops on Windows and later runs authorized reproducible experiments on Genkai. The research concerns context and memory during model changes, not a new inference engine. Reuse existing OSS and confine replacement points. Evidence is in the [source register](sources.md).

| 選択肢 / Option | 利点 / Benefit | 負担・判断 / Cost and decision |
| --- | --- | --- |
| Go一体型application + Python worker / Modular Go application + Python worker | Go学習とMem0再利用を両立 / Keeps Go learning and Mem0 reuse | HTTP境界を1つ維持。推奨 / One worker HTTP boundary; recommended |
| Python単一application / Python-only application | processと型変換が少ない / Fewer processes and translations | Go学習・選択済みA案から外れる / Departs from Go learning and selected Route A |
| 多数のmicroservices / Many microservices | 個別scale / Independent scaling | 配布、認証、障害処理が増える。現段階では不要 / More deployment, auth and failure handling; unnecessary now |
| 完成chat製品のfork / Fork a complete chat product | 画面を再利用 / Reuses a UI | 独自session・記憶・更新の責務が衝突 / Competing session/memory ownership and upgrade burden |

## 技術スタック / Technology stack

| 層 / Layer | 推奨 / Recommendation | 固定・交換条件 / Pinning and replacement |
| --- | --- | --- |
| Application | Go、標準 `net/http`, `encoding/json`, `context`, `database/sql` / Go standard packages | 現在の `go.mod` を出発点にし、ORM・web frameworkを初期追加しない / Start from current go.mod; no initial ORM/web framework |
| 正本保存 / Authoritative storage | SQLite、driver候補 `modernc.org/sqlite` / SQLite with modernc.org/sqlite candidate | Go対応・推移依存・LICENSEを確認後に固定 / Pin after Go compatibility, dependency and license review |
| 推論 / Inference | Ollama HTTP adapter | 特定model API型をapplicationへ流さない / Keep provider types out of application |
| 記憶 / Memory | Python worker + Mem0 OSS `Memory` | `MemoryClient`はcloud版ではない。local設定を明示 / Do not use the cloud client; configure all local providers explicitly |
| 派生保存 / Derived storage | Qdrant local mode + Mem0 history + source mapping | workerが単独所有。正本DBと別ファイル / One worker owns these; separate from source DB |
| UI | React + TypeScript + Vite + Mantine（選択済み） / React + TypeScript + Vite + Mantine (selected) | library stylesとtheme props。自作CSSなし / Library styles and theme props; no custom CSS |
| 検証 / Verification | Go testing、Python unittest；UIにVitest + Testing Library、E2EにPlaywrightを候補 / Go testing, Python unittest; proposed Vitest/Testing Library and Playwright | 現在はGo検証のみ導入済み。追加test依存も監査・導入許可対象 / Only current Go checks exist; new test dependencies also require audit and installation scope |
| 実験 / Experiments | Go runner、JSONL/CSV metrics、versioned JSON設定 / Go runner, JSONL/CSV metrics, versioned JSON config | runnerから同じAPIを呼び、UIに依存しない / Use the same API independently of UI |

SQLiteをJSON snapshotより推奨する理由は、message・turn・request ID・更新待ち記録をtransactionで一緒に保存できることです。DB driver追加とschema migrationの学習は必要です。最初の契約演習はメモリ内repositoryで行い、その実装も本人が担当します。test-only fakeだけをassistantが実装します。

SQLite is preferred over JSON snapshots because messages, turns, request IDs and pending updates can be stored together transactionally. It adds a driver and schema-migration learning. Initial contract exercises use a user-written in-memory repository; only test-only fakes are assistant-owned.

## モジュール境界 / Module boundaries

すべて以下は予定構成です。現コードを自動移動・書換えしません。小さい間はpackageを増やしすぎず、利用側に小さなinterfaceを定義します。application/domainはHTTP、SQLite、Mem0、Reactに依存しません。

Everything below is a target structure. Existing code is not automatically moved or rewritten. Keep packages modest and define small interfaces at their consumers. Application/domain code must not depend on HTTP, SQLite, Mem0 or React.

| Module / 予定場所 | 責務と境界 / Responsibility and boundary | 交換・検証 / Replacement and verification |
| --- | --- | --- |
| Domain / `internal/domain` | Session、Message、Turn、Scope、revision。外部I/Oなし / Core records, scope and revision; no external I/O | 不変条件・訂正・削除のunit tests / Invariant, correction and deletion unit tests |
| Application / `internal/app` | session作成、turn進行、model切替、export/delete、回復を調整 / Coordinates sessions, turns, switching, export/delete and recovery | 下記portのfakeで状態遷移を検証 / State transitions with port fakes |
| Identity / `internal/auth` | credential→Principal、session所有権。clientのidentity指定を信用しない / Credential→Principal and ownership; never trust client identity fields | 他user/lab、期限切れ、偽scopeを拒否 / Deny other users/labs, expiry and forged scope |
| HTTP / `internal/transport/httpapi` | DTO、入力制限、status/error変換、UI静的配信 / DTOs, input limits, status/error mapping and UI static serving | schemaとHTTP契約。business ruleを置かない / Schema/HTTP tests; no business rules |
| Policies / `internal/contextpolicy` | full history、rolling summary、Mem0 + recent。token予算内のContextPlan / Three policies produce a budgeted ContextPlan | 同じfixture・budget・overflow条件 / Same fixtures, budget and overflow conditions |
| Models / `internal/adapters/ollama` | `ModelClient`、capabilities、usage変換、取消 / ModelClient, capabilities, usage mapping and cancellation | 記録した合成HTTP応答と許可済み実連携 / Synthetic recorded HTTP responses and authorized integration |
| Token budget / `internal/tokenbudget` | template込み入力と出力予約。推定／実測を区別 / Templated input plus output reserve; distinguish estimates/measurements | tokenizer/runtime一致fixture。不明なら評価対象外 / Tokenizer/runtime fixtures; unresolved counts block evaluation |
| Source repository / `internal/adapters/sqlite` | `SessionRepository`、transaction、turnとjob記録 / Repository, transactions, turns and pending-operation records | restart、重複、migration、restore / Restart, duplicate, migration and restore tests |
| Memory adapter / `internal/adapters/memoryhttp` | versioned worker API、timeout、operation ID、error変換 / Versioned worker API, timeout, operation ID and errors | worker契約suiteを固定 / Stable worker contract suite |
| Worker / `memory-worker` | Mem0呼出し、scope、出典mapping、Qdrant/history、idempotent operation状態 / Mem0, scope, provenance mapping, Qdrant/history and idempotent operation state | Mem0更新時は隔離・訂正・削除suiteを再実行 / Rerun isolation/correction/deletion suite for Mem0 upgrades |
| Metrics / `internal/metrics` | stage、時間、token、欠測理由。本文をlogしない / Stage/time/tokens/missing reasons; no content logging | 合成計算fixture、schema version / Synthetic accounting fixtures and schema version |
| Runner / `cmd/evaluate` | frozen case・config読込、API再生、checkpoint、結果集計 / Frozen cases/config, API replay, checkpoints and aggregation | fake API再生、途中再開、dataset分離 / Fake replay, resume and dataset separation |
| UI / `web/src` | backend状態の表示と操作。記憶処理・採点なし / Displays/operates backend state; no memory policy or scoring | component/API/E2E tests。詳細は[UI設計](ui.md) / Component/API/E2E tests; see UI design |
| Composition / `cmd/bridge`, `configs`, `scripts` | 設定検証、portの組立、readiness、停止 / Config validation, wiring, readiness and shutdown | 不正設定、起動失敗、子process終了 / Invalid config, startup failure and child-process cleanup |

主要portは `ModelClient`, `MemoryClient`, `SessionRepository`, `ContextPolicy`, `TokenCounter`, `MetricsSink`。domain型を使い、vendor SDK objectを返しません。UIはversioned DTOのみを参照します。契約の詳細は[contracts.md](contracts.md)です。

Core ports are ModelClient, MemoryClient, SessionRepository, ContextPolicy, TokenCounter and MetricsSink. They use domain types and do not return vendor SDK objects. The UI consumes versioned DTOs only. See [contracts.md](contracts.md).

## データと失敗 / Data and failures

原messageを正本とします。summary、memory、vector、browser表示は派生情報です。memoryはsource IDs、revision、抽出versionを持ち、assistant回答を検証済み事実へ自動昇格しません。modelを切り替えてもsession ID、owner、policy、sourceは維持し、helper/embedding modelは実験内で固定します。

Source messages are authoritative. Summaries, memories, vectors and browser views are derived. Memory retains source IDs, revision and extraction version; assistant output is not automatically promoted to a verified fact. Switching models preserves session ID, owner, policy and sources; helper/embedding models stay fixed within an experiment.

1. 認証・所有者確認→request IDとuser messageを保存→context作成。 / Authenticate/authorize → save request ID and user message → assemble context.
2. 推論→answerとusage、memory更新待ち記録を一つの正本transactionで保存。 / Generate → save answer, usage and pending memory operation in one source transaction.
3. workerで更新→operation結果を保存→turn完了。 / Update through the worker → save operation result → complete the turn.

正本DBとQdrantは同じtransactionにはできません。追加queue serviceは使わず、DBの更新待ち記録を再調停します。answer保存後のmemory失敗は `memory_pending` として可視化し、回復まで次turnを止めます。不確定な推論・writeを無条件に再送せず、exactly-once inferenceは主張しません。

The source DB and Qdrant cannot share a transaction. Reconcile pending DB operations without adding a queue service. A memory failure after answer persistence becomes visible as memory_pending and blocks the next turn until recovery. Do not blindly retry uncertain inference/writes or claim exactly-once inference.

sessionごとに同時実行は1turn、最初のGPU実行も直列です。削除開始時にrevisionを進め、遅延writeを拒否し、全派生store・content-bearing checkpointを削除して検証します。削除APIは物理媒体の完全消去を保証しません。

Allow one active turn per session and initially serialize GPU work. Deletion advances the revision, rejects late writes, purges all derived stores/content-bearing checkpoints and verifies completion. The deletion API does not promise physical-media erasure.

## 更新と回復の設計 / Upgrade and recovery design

| 変更 / Change | 安定境界 / Stable boundary | 移行・戻し方 / Migration and rollback |
| --- | --- | --- |
| React/Mantine/Vite | native `/v1` DTO | lockfile、build ID、直前のUI artifactを保持。API互換テスト / Lockfile/build ID, retain previous UI artifact, API contract tests |
| Go内部 / Go internals | domain invariants、ports、API | package単位で変更、unit/contract回帰 / Change by package; unit/contract regression |
| API/worker | `/v1`、`schema_version`、capability情報 / Version and capability information | 追加fieldは互換、削除・意味変更は新version。非互換なら起動拒否 / Additive fields compatible; removal/semantic change needs new version; reject incompatibility at startup |
| SQLite schema | repository port | migration IDとbackup。private copyで移行・restore検証後に適用 / Migration IDs/backups; test migration/restore on private copies before use |
| Mem0/Qdrant | worker contract | 新しいnamespace/storeへ再構築して比較。旧storeを直接書換えない / Rebuild into a new namespace/store and compare; do not overwrite the old store directly |
| embedding model/dimension | source IDsとmemory schema / Source IDs and memory schema | vector再生成、新collection、切替前照合 / Re-embed into a new collection and verify before switching |
| answer/helper model | model adapter | weights/tokenizer/template/quantization/hashを固定。新しい実験ID / Pin weights/tokenizer/template/quantization/hashes; use a new experiment ID |
| inference runtime | ModelClient/TokenCounter | capability・usage・timeout・token対応を実測確認 / Verify capabilities, usage, timeouts and token accounting |

一度に一つの依存系統を更新し、productionとtest依存をlockします。前版manifest、artifact、移行前backupを私的に保持し、schema互換性を確認せず古いbinaryを新DBへ戻しません。prototype段階では停止時間を許容し、zero-downtime移行は要求しません。

Upgrade one dependency family at a time and lock production/test dependencies. Retain the previous manifest, artifacts and pre-migration backup privately. Never run an old binary against a new DB without schema compatibility. Prototype upgrades may use downtime; zero-downtime migration is not required.

削除済みデータをbackupから復活させないため、restore時も削除台帳を適用します。backup、WAL、exportの保持・削除範囲を文書化します。Genkaiではlocal storageとprivate永続storageを区別し、共有filesystem上のlive SQLite/WALを前提にしません。

Apply the deletion ledger during restore to avoid resurrecting deleted data. Document retention/deletion of backups, WAL and exports. On Genkai, distinguish local storage from private persistent storage; do not assume live SQLite/WAL is safe on a shared filesystem.

## 見直す条件 / Revisit triggers

multi-process書込、複数worker、常時サービスが必要になればPostgreSQLやQdrant serverを別ADRで検討します。Ollamaの能力・throughputが実測で不足すればllama.cpp/vLLM adapterを別実験にします。Mem0とモデルを同じ種類の代替物として比較しません。

If multi-process writers, multiple workers or a persistent service become necessary, consider PostgreSQL/Qdrant server in a separate ADR. If measured capabilities/throughput are insufficient, evaluate a llama.cpp/vLLM adapter separately. Do not compare Mem0 and a model as interchangeable alternatives.
