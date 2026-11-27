# ADR-001: 保守可能な研究原型 / Maintainable research prototype

更新 / Updated: 2026-11-27。状態 / Status: A案とMantineは選択済み。2026-11-27の方針で、まずCLIとAPIで検証し、検証の後に研究UIを加えます。同日、安全設計を加えました。2026-11-13に、既存のdomain・turn契約に加えて、認証付きsession metadata APIと許可model ID一覧を受け入れました。以下はturn実行、package分割、adapter、保存、CLI・UI境界、安全対策の設計提案です。 / Route A and Mantine are selected. Under the 2026-11-27 direction, verification through the CLI and API comes first and the research UI is added afterwards. A security design was added the same day. On 2026-11-13, acceptance extended from the domain and turn contracts to the authenticated session-metadata API and configured model ID list. Turn execution, package boundaries, adapters, storage, the CLI and UI boundaries and security controls below are design proposals.

決定者は利用者です。実装と導入は別途の受入・許可範囲に従います。進捗は[PLAN.md](../PLAN.md)のみで管理します。

The user is the decision maker. Implementation and installation follow their acceptance and authorized scope. Track progress only in [PLAN.md](../PLAN.md).

2026-11-18に利用者は段階的なモジュール化を選択しました。まず現在のGo process内でL4の切替と失敗時の動作を受け入れます。その後、実際の外部I/Oや独立したtest境界に合わせてpackageを分けます。核心codeとその再配置は本人、外周codeはassistantが担当します（[2026-11-19の分担](../AGENTS.md)）。

On 2026-11-18 the user selected incremental modularization. First accept L4 switching and failure behavior in the current Go process. Then extract packages as external I/O and independent test boundaries emerge. The user owns core code and its moves, and the assistant owns surrounding code ([division of 2026-11-19](../AGENTS.md)).

<a id="revision-2026-11-27"></a>

## 2026-11-27の見直し / Revision of 2026-11-27

利用者の報告によると、教員の方針は「まずcommandlineとAPIで検証し、検証の後にUIを加える」です。同じ日の最初の報告は「UI画面は作らない」でしたが、利用者が訂正しました。あわせて利用者は、他の人が書いたprogramからも呼べること、安全対策を強めること、学べる知識を最大にすることを求めました。

According to the user, the supervisor's direction is to verify through the command line and API first and add a UI after verification. The first report that day said no UI screen would be built; the user corrected it. The user also asked that programs written by other people can call the API, that security be strengthened and that the route teach as much as possible.

| 項目 / Item | 以前 / Before | 現在の案 / Now |
| --- | --- | --- |
| 人が使う入口 / Human entry point | 研究UIだけ / The research UI only | まず[CLI](cli.md)。検証の後に[研究UI](ui.md)を加える / The [CLI](cli.md) first; the [research UI](ui.md) is added after verification |
| APIの利用者 / API consumers | browser、runner | CLI、runner、他の人のprogram、後の研究UI。全てが[同じHTTP契約](contracts.md#other-programs)を呼ぶ / The CLI, the runner, other people's programs and the later research UI, all calling the [same HTTP contract](contracts.md#other-programs) |
| 順序 / Order | UIはS4 / UI in S4 | S4はCLIと外部programからの利用、その最後に安全性の点検。UIはS7 / S4 is the CLI and use by external programs, ending with a security review. The UI is S7 |
| 安全対策 / Security | 各所に分散した要件 / Requirements scattered across documents | [安全設計](#security)に層ごとにまとめ、各段階に授業を置く / Gathered by layer in the [security design](#security), with lessons in every stage |
| 授業数 / Lessons | 42 | 56。安全6回、UI 8回を追加 / 56, adding six security lessons and eight UI lessons |

変えないものは、A案（Go + Ollama + Mem0 OSS）、`/v1` の契約と受入済みのL1–L3、L4のturn設計、SQLiteとworkerの境界、非streamingの202と照会の形、研究UIのReact + TypeScript + Vite + Mantineです。

Unchanged: Route A (Go + Ollama + Mem0 OSS), the `/v1` contract and accepted L1–L3, the L4 turn design, the SQLite and worker boundaries, the non-streaming 202-then-lookup shape, and React + TypeScript + Vite + Mantine for the research UI.

UIを後にする利点は三つあります。APIの動作と失敗時の扱いを、画面の状態と混ぜずに確かめられます。browserが加わる前に、認証・隔離・濫用対策を点検できます。UIは、検証済みの契約とfixtureをそのまま使えます。CLIの最初の5回（S4-1〜S4-5）はS1のrouteだけに依存するので、S1の受入直後に進めることもできます。順序は利用者が決めます。

Deferring the UI has three benefits. API behavior and failure handling are verified without view state mixed in. Authentication, isolation and abuse controls are reviewed before a browser joins. The UI then reuses the verified contract and fixtures as they are. The first five CLI lessons (S4-1 to S4-5) depend only on S1 routes and can follow S1 acceptance directly. The user decides the order.

## 目的と選択肢 / Context and options

一人が学びながら保守し、Windowsで開発し、許可後に玄界で再現実験を行う原型です。主題はモデル変更時の文脈と記憶であり、新しい推論engineの開発ではありません。既存OSSを再利用し、交換箇所を限定します。根拠は[資料台帳](sources.md)に集約します。

One learner maintains the prototype, develops on Windows and later runs authorized reproducible experiments on Genkai. The research concerns context and memory during model changes, not a new inference engine. Reuse existing OSS and confine replacement points. Evidence is in the [source register](sources.md).

| 選択肢 / Option | 利点 / Benefit | 負担・判断 / Cost and decision |
| --- | --- | --- |
| Go一体型application + Python worker / Modular Go application + Python worker | Go学習とMem0再利用を両立 / Keeps Go learning and Mem0 reuse | HTTP境界を1つ維持。推奨 / One worker HTTP boundary; recommended |
| Python単一application / Python-only application | processと型変換が少ない / Fewer processes and translations | Go学習・選択済みA案から外れる / Departs from Go learning and selected Route A |
| 多数のmicroservices / Many microservices | 個別scale / Independent scaling | 配布、認証、障害処理が増える。現段階では不要 / More deployment, auth and failure handling; unnecessary now |
| 完成chat製品のfork / Fork a complete chat product | 画面を再利用 / Reuses a UI | 独自session・記憶・更新の責務が衝突 / Competing session/memory ownership and upgrade burden |

### Goの実装順序の選択 / Choice of Go implementation sequence

| 選択肢 / Option | 利点と負担 / Benefit and cost |
| --- | --- |
| A · 段階的なモジュール化（選択済み） / Incremental modularization (selected) | 現在の `cmd/bridge` で切替を検証し、同じtestで後から境界の抽出を確認できる。境界が必要になるまでは同じpackageで責務を分ける / Verify switching in the current `cmd/bridge` package, then use the same tests to verify extraction. Separate responsibilities within that package until a package boundary is needed |
| B · package境界を先に作る / Establish packages first | import制約を早期に確認できる。turn契約が変化する時期に型の公開・変換・移動も必要になる / Enforce import constraints early, with additional type exposure, mapping and code moves while the turn contract is still evolving |
| C · APIと実行を別processにする / Separate the API and executor processes | 実行を独立して再起動できる。永続的なjobの受理・認領、障害検知、取消・遅延writeの調停が必要 / Restart execution independently, with durable admission/claiming, failure detection, cancellation and coordination of late writes |

## 技術スタック / Technology stack

| 層 / Layer | 推奨 / Recommendation | 固定・交換条件 / Pinning and replacement |
| --- | --- | --- |
| Application | Go、標準 `net/http`, `encoding/json`, `context`, `database/sql` / Go standard packages | 現在の `go.mod` を出発点にし、ORM・web frameworkを初期追加しない / Start from current go.mod; no initial ORM/web framework |
| 正本保存 / Authoritative storage | SQLite、driver候補 `modernc.org/sqlite` / SQLite with modernc.org/sqlite candidate | Go対応・推移依存・LICENSEを確認後に固定 / Pin after Go compatibility, dependency and license review |
| 推論 / Inference | Ollama HTTP adapter | 特定model API型をapplicationへ流さない / Keep provider types out of application |
| 記憶 / Memory | Python worker + Mem0 OSS `Memory` | `MemoryClient`はcloud版ではない。local設定を明示 / Do not use the cloud client; configure all local providers explicitly |
| 派生保存 / Derived storage | Qdrant local mode + Mem0 history + source mapping | workerが単独所有。正本DBと別ファイル / One worker owns these; separate from source DB |
| CLI | Go、標準 `flag` と `net/http` / Go with standard `flag` and `net/http` | CLI frameworkを足さない。出力の形と終了codeを契約として固定 / No CLI framework. Output shape and exit codes are fixed as a contract |
| UI（検証の後、S7） / UI (after verification, S7) | React + TypeScript + Vite + Mantine（選択済み） / React + TypeScript + Vite + Mantine (selected) | library stylesとtheme props。自作CSSなし / Library styles and theme props; no custom CSS |
| 検証 / Verification | Go testing、Python unittest；UIにVitest + Testing Library、E2EにPlaywrightを候補 / Go testing, Python unittest; proposed Vitest/Testing Library and Playwright | 現在はGo検証のみ導入済み。追加test依存も監査・導入許可対象 / Only current Go checks exist; new test dependencies also require audit and installation scope |
| 実験 / Experiments | Go runner、JSONL/CSV metrics、versioned JSON設定 / Go runner, JSONL/CSV metrics, versioned JSON config | runnerはCLIと同じclient packageで同じAPIを呼ぶ / The runner calls the same API through the client package the CLI uses |

SQLiteをJSON snapshotより推奨する理由は、message・turn・request ID・更新待ち記録をtransactionで一緒に保存できることです。DB driver追加とschema migrationの学習は必要です。最初の契約演習はメモリ内repositoryで行い、その実装も本人が担当します。test-only fakeだけをassistantが実装します。

SQLite is preferred over JSON snapshots because messages, turns, request IDs and pending updates can be stored together transactionally. It adds a driver and schema-migration learning. Initial contract exercises use a user-written in-memory repository; only test-only fakes are assistant-owned.

## 実行時の境界 / Runtime boundaries

```mermaid
flowchart LR
    CLI["CLI / 命令行"] -->|"HTTP /v1"| Go["Go bridge / 会話・権限・状態"]
    Other["Other programs / 他の人のprogram"] -->|"HTTP /v1"| Go
    UI["React + Mantine / 研究UI（S7）"] -.->|"HTTP /v1"| Go
    Runner["Evaluation runner / 評価実行"] -->|"HTTP /v1"| Go
    Go --> Source["SQLite / 正本・更新待ち"]
    Go -->|"generation / 推論"| Ollama["Ollama / ローカルモデル"]
    Go -->|"versioned HTTP / 版付き契約"| Worker["Python worker / Mem0 OSS"]
    Worker --> Derived["Qdrant local + history / 派生情報"]
    Worker -->|"helper + embedding / 補助モデル"| Ollama
```

これは予定のruntime呼出し図で、Goのimport図ではありません。常駐するのはGo・Python worker・Ollamaの3 processです。CLIとrunnerはcommandごとに起動して終了し、他の人のprogramと同じnative APIを呼びます。点線の研究UIは検証の後に加え、その時からGoが静的buildを同じoriginで配信します。Node/Viteは開発・build時だけ使います。公開 `website/` は別の説明用siteで、研究APIやDBに接続しません。

This is a proposed runtime-call diagram, not a Go import graph. Three processes stay running: Go, the Python worker and Ollama. The CLI and runner start and exit per command and call the same native API as other people's programs. The dashed research UI is added after verification; from then Go serves its static build from the same origin. Node/Vite are development/build tools only. Public `website/` is a separate information site without access to the research API or databases.

## モジュール境界と依存規則 / Module boundaries and dependency rules

すべて以下は予定構成です。現コードを自動移動・書換えしません。小さい間はpackageを増やしすぎず、利用側に小さなinterfaceを定義します。application/domainはHTTP、SQLite、Mem0、CLI、Reactに依存しません。

Everything below is a target structure. Existing code is not automatically moved or rewritten. Keep packages modest and define small interfaces at their consumers. Application/domain code must not depend on HTTP, SQLite, Mem0, the CLI or React.

| Module / 予定場所 | 責務と境界 / Responsibility and boundary | 変更・検証の接点 / Change and test boundary |
| --- | --- | --- |
| Domain / `internal/domain` | Session、Message、Turn、Principal/Scopeの値と不変条件。外部I/Oなし / Core values and invariants, including identity/scope; no external I/O | 純粋な所有権・role・状態規則のunit tests / Pure ownership, role and state-rule unit tests |
| Application / `internal/app` | session、turn、回復、export/deleteのuse caseと消費側port。state変更を決定 / Use cases and consumer-owned ports; decides state transitions | session/turn/recovery/data操作をまず同packageの別fileに分け、fakeで検証 / Separate responsibility files within one package first; test with fakes |
| HTTP / `internal/transport/httpapi` | versioned DTO、認証結果の受渡し、入力制限、濫用の上限、監査log、HTTP error変換、S7からUI配信と安全header / Versioned DTOs, authenticated identity forwarding, input limits, abuse limits, audit logging, HTTP error mapping and, from S7, UI serving with security headers | domain/use-case型との境界mapping。認可判断をhandlerに重複させない / Boundary mapping; do not duplicate use-case authorization in handlers |
| Adapters / `internal/adapters/{ollama,sqlite,memoryhttp}` | 外部の推論・正本保存・worker通信をそれぞれ隔離 / Isolate inference, source persistence and worker communication separately | 実際の連携時に作成。provider応答、短transaction、timeout・errorの契約 / Create with real integrations; provider responses, short transactions and timeout/error contracts |
| Worker / `memory-worker` | Python内でAPI境界、operation調整、Mem0依存、派生store所有を分ける / Separate API handling, operation coordination, Mem0 dependency and derived-store ownership inside Python | HTTP handlerから直接storeを触らない。source mapping、重複、訂正、削除を検証 / Handlers do not manipulate stores directly; test mapping, duplicates, corrections and deletion |
| Client / `internal/client` | `/v1` を呼ぶHTTP client、client側DTO、応答の検証、error envelopeの読取り / HTTP client for `/v1`, client-side DTOs, response validation and error-envelope decoding | server側packageをimportしない。serverと同じ契約fixtureで検証 / Imports no server-side package; verified with the same contract fixtures as the server |
| CLI / `cmd/bridge` のsubcommand | 引数の解釈、送信後の待機と回復、人向け表示と `--json`、終了code / Argument parsing, waiting and recovery after a send, human output and `--json`, exit codes | clientだけを使う。会話の規則や認可を持たない。詳細は[CLI設計](cli.md) / Uses only the client; holds no conversation rules or authorization. See the [CLI design](cli.md) |
| UI / `web/src`（S7） | appが認証・session・画面間連携を所有。featuresが表示と操作を担当 / App owns auth/session/cross-view coordination; features own views and interaction | `app → features → lib/api → lib/contracts`。詳細は[UI設計](ui.md) / See the [UI design](ui.md) for dependency and test seams |
| Runner / `cmd/evaluate` | frozen case/config、API再生、checkpoint、結果集計 / Frozen input, API replay, checkpoints and aggregation | 同じclientでnative APIだけを利用し、DBやmemoryへ直接書込まない / Use only the native API through the same client; never write directly to DB or memory |
| Composition / `cmd/bridge` の `serve` | 設定検証、具体adapterの生成・注入、起動・readiness・停止 / Validate config, construct/inject concrete adapters, start, report readiness and stop | 構築と終了だけを調整し、会話のbusiness ruleを置かない / Coordinate construction/lifecycle without conversation business rules |

port候補は `ModelClient`, `MemoryClient`, `SessionRepository`, `ContextPolicy`, `TokenCounter`, `MetricsSink`。domain値とproject固有のuse-case入出力型を使い、vendor SDK objectを返しません。現在のuse caseが必要なportだけを導入します。CLI、runner、UIはversioned DTOのみを参照します。契約の詳細は[contracts.md](contracts.md)です。

Candidate ports are ModelClient, MemoryClient, SessionRepository, ContextPolicy, TokenCounter and MetricsSink. They use domain values and project-owned use-case input/output types, not vendor SDK objects. Introduce only the ports a current use case needs. The CLI, runner and UI consume versioned DTOs only. See [contracts.md](contracts.md).

Goの依存方向は `httpapi → app → domain` と `adapters → app/domain` です。client側は `CLI・runner → client` で、server側のpackageとはHTTPだけでつながります。`cmd/bridge` だけが具体adapterとuse caseを組み立てます。domainは他のproject layerをimportせず、appはHTTP DTO・SQL driver・Ollama/Mem0型をimportしません。app側が必要最小限のinterfaceを定義し、adapterが実装します。大きな共通 `ports`／`utils` packageへ型を集めません。

Go dependencies point `httpapi → app → domain` and `adapters → app/domain`. On the client side they point `CLI/runner → client`, which meets the server-side packages only over HTTP. Only `cmd/bridge` wires concrete adapters into use cases. Domain imports no other project layer; app imports no HTTP DTOs, SQL drivers or Ollama/Mem0 types. App defines the small interfaces it consumes and adapters implement them. Do not collect every type into broad `ports` or `utils` packages.

credential→Principalの検証はHTTP境界の責務、Principalと対象sessionの認可はuse caseの責務です。認証providerが独立して変化するときだけ `internal/auth` に抽出します。policy・token予算・metricsは初めはapp内の小さい部品でよく、独立した実装が必要になった時に `internal/contextpolicy` 等へ抽出し、appへ注入します。appと実装packageが互いをimportする構成にしません。

Credential-to-Principal verification belongs at the HTTP boundary; authorization of that Principal for the target session belongs in the use case. Extract `internal/auth` when authentication providers need to change independently. Policies, token budgets and metrics can begin as small app-local components. Extract packages such as `internal/contextpolicy` when independent implementations are needed, then inject them into app. Avoid app/implementation import cycles.

### 段階的な抽出 / Incremental extraction

以下は[PLAN.md](../PLAN.md)の段階との対応であり、追加の作業一覧ではありません。空directoryや全interfaceを先に作らず、一つの責務を抽出するたびに既存契約・失敗動作を維持して確認します。

This maps to the stages in [PLAN.md](../PLAN.md), not a second task list. Do not pre-create empty directories or every interface. Extract one responsibility at a time and verify that existing contracts and failure behavior stay intact.

| 段階 / Stage | 分離する理由と範囲 / Reason and scope |
| --- | --- |
| S1 | 現在の `cmd/bridge` 内でL4を検証。認証済みsessionでfake A → B → A、重複・取消・遅延応答を確認し、後の抽出に同じfixtureを使う / Verify L4 inside the current `cmd/bridge`: fake A → B → A in an authenticated session, replay, cancellation and late results. Reuse the fixtures for later extraction |
| S2 | Ollama呼出しとSQLiteの原子操作を実際のadapter境界にする。HTTP契約を保ち、先のfixtureに実推論での切替と再起動後の照会を加える / Use actual Ollama calls and SQLite atomic operations as adapter boundaries. Keep the HTTP contract and extend the existing fixtures with real switching and lookup after restart |
| S3 | worker境界と3方式の違いが出現した時にmemory adapterとpolicyを分離 / Separate the memory adapter and policies when worker integration and three distinct policies exist |
| S4–S5 | CLIとrunnerをnative API上に構築。CLIで作ったclientをrunnerが再利用する。必要な計測を追加し、再利用と独立更新の必要がある部分だけ抽出 / Build the CLI and runner on the native API, with the runner reusing the client written for the CLI. Add metrics and extract only for real reuse or independent updates |
| S7 | 検証済みのnative API上に `web/` を構築。Go側の追加は静的build配信と安全headerだけ / Build `web/` on the verified native API. The only Go additions are static-build serving and security headers |

境界を抽出する変更と動作を追加する変更は、それぞれ同じ契約を検証できる差分にします。抽出では既存route、所有者判定、状態遷移を維持し、domain/appへのvendor依存や循環importを確認します。単一Go processでHTTP応答と有界な実行を保てる間は、その構成を使います。

Keep boundary extraction and behavior changes independently verifiable against the same contract. During extraction, preserve routes, authorization and state transitions, and check for vendor dependencies in domain/app and import cycles. Use one Go process while it can serve HTTP and keep execution bounded.

## データと失敗 / Data and failures

原messageを正本とします。summary、memory、vector、CLIやbrowserの表示は派生情報です。memoryはsource IDs、revision、抽出versionを持ち、assistant回答を検証済み事実へ自動昇格しません。modelを切り替えてもsession ID、owner、policy、sourceは維持し、helper/embedding modelは実験内で固定します。

Source messages are authoritative. Summaries, memories, vectors and CLI or browser views are derived. Memory retains source IDs, revision and extraction version; assistant output is not automatically promoted to a verified fact. Switching models preserves session ID, owner, policy and sources; helper/embedding models stay fixed within an experiment.

| 所有者 / Owner | 決定・保存する内容 / Decisions and persisted content |
| --- | --- |
| Go app | session/turnの状態遷移、認可、policy/model選択、回復判断 / Session/turn transitions, authorization, policy/model selection and recovery decisions |
| SQLite adapter | 正本message、turn、重複key、更新待ち記録を原子的に保存。制約を実施 / Atomically persist authoritative messages, turns, deduplication keys and pending work; enforce constraints |
| Python worker | 派生記憶、vector、history、出典mapping、worker operation結果。Goのturnを直接更新しない / Derived memory/vectors/history/provenance and worker-operation results; never directly mutate Go turns |
| CLI / runner / UI | APIから受け取ったsnapshotと操作意図。手元に会話の正本を保存せず、正本や権限を独自決定しない / API snapshots and user/run intent. They keep no conversation state locally and never define authority or access themselves |

1. 認証・所有者確認後、request ID／payload、user message、初期turn状態を短transactionで保存し、同sessionのactive turnが一つという制約を守る。 / After authentication/authorization, atomically persist request identity/payload, user message and initial turn state while enforcing one active turn per session.
2. transactionを閉じてcontext組立・推論を行う。answer、usage、派生更新が必要な場合の更新意図を別の短transactionで保存する。 / Close the transaction before context construction/inference. Atomically save the answer, usage and, when needed, derived-update intent in a separate short transaction.
3. transaction外でworkerへ更新を要求し、結果を照合してGoが完了を保存する。派生更新不要のpolicyは直接完了する。 / Request the worker update outside a transaction; Go reconciles the result and records completion. Policies without derived updates complete directly.

Repository portは上記の原子操作を表現します。一般的なCRUDの寄せ集めでuse caseに途中状態を組み立てさせず、SQL transactionやdriver型をappへ公開しません。revisionの意味・増分主体・削除後の遅延write拒否は[契約設計](contracts.md)に集約します。

The repository port expresses these atomic operations. Avoid generic CRUD that forces use cases to assemble partially committed state; expose neither SQL transactions nor driver types to app. Revision meanings, increment owners and rejection of writes after deletion are defined in [contracts.md](contracts.md).

Go側も、要求受理・回答保存・回復完了の各transactionで、sessionが `deleting/deleted` でなく、operation開始時の `session_epoch` と一致することを原子的に確認します。不一致なら正本への書込を拒否します。worker側の検査だけに依存せず、削除後に遅い推論回答が正本へ戻る経路を閉じます。

Go also atomically checks, inside request-acceptance, answer-save and recovery-completion transactions, that the session is neither deleting nor deleted and its `session_epoch` matches the operation's starting epoch. Reject source writes on mismatch. Worker checks alone are insufficient: a late inference result must not recreate authoritative content after deletion.

正本DBとQdrantは同じtransactionにはできません。追加queue serviceは使わず、DBの更新待ち記録を再調停します。answer保存後のmemory失敗は `memory_pending` として可視化し、回復まで次turnを止めます。不確定な推論・writeを無条件に再送せず、exactly-once inferenceは主張しません。

The source DB and Qdrant cannot share a transaction. Reconcile pending DB operations without adding a queue service. A memory failure after answer persistence becomes visible as memory_pending and blocks the next turn until recovery. Do not blindly retry uncertain inference/writes or claim exactly-once inference.

sessionごとに同時実行は1turn、最初のGPU実行も直列です。削除開始時に `session_epoch` を進め、遅延writeを拒否し、全派生store・content-bearing checkpointを削除して検証します。削除APIは物理媒体の完全消去を保証しません。

Allow one active turn per session and initially serialize GPU work. Deletion advances `session_epoch`, rejects late writes, purges all derived stores/content-bearing checkpoints and verifies completion. The deletion API does not promise physical-media erasure.

### 実行の寿命 / Execution lifecycle

202を返す前に受理を記録します。S1はメモリ内、S2以降は永続保存です。appがGo process内の有界な実行・回復を管理し、受理済みjobはHTTP切断後も続行します。取消には明示APIを使います。起動時は[保存したturn状態](contracts.md#restart-cut-points)から、未送出の再開、全履歴回答の保存後の完了、不確定な推論の照合を判断します。

Record acceptance before returning 202, in memory for S1 and durably from S2. App owns bounded execution and recovery inside the Go process. Disconnecting HTTP leaves accepted work running. Cancellation uses the explicit API. Startup uses the [durable turn state](contracts.md#restart-cut-points) to resume undispatched work, complete a saved full-history turn or expose uncertain inference for reconciliation.

S1のL4では、同じGo process内で受理・履歴snapshot・回答確定の境界を検証します。storeは同じsessionのrequest記録・message・turnの更新を原子的に適用し、appはI/Oをその排他区間の外で実行します。これにより、後続のSQLite adapterでは同じ原子操作をtransactionで表現できます。状態と競合の詳しい設計は [L4のturn設計](contracts.md#l4-turn-design) を参照します。再起動回復はS2で検証します。

L4 in S1 verifies the admission, history snapshot and answer-commit boundaries inside one Go process. The store atomically updates a session's request records, messages and turns. The app performs I/O outside that critical section. A later SQLite adapter can express the same atomic operations as transactions. See the [L4 turn design](contracts.md#l4-turn-design) for states and races. S2 verifies restart recovery.

組立入口がprocessの停止signalと期限を渡し、まず新しい受理を止め、進行中処理を確定または回復可能な状態にしてから接続を閉じます。workerもoperation状態を記録します。process自身が起動していないOllama等を暗黙停止しません。livenessと依存readinessを区別し、具体的な期限・上限は実装時に設定契約で固定します。

The composition root supplies shutdown signals/deadlines: stop new acceptance, settle active work or preserve recoverable state, then close resources. The worker records operation state too. Do not implicitly terminate an Ollama process that this process did not start. Separate liveness from dependency readiness; freeze actual deadlines and bounds in the implementation's configuration contract.

## 更新と回復の設計 / Upgrade and recovery design

| 変更 / Change | 安定境界 / Stable boundary | 移行・戻し方 / Migration and rollback |
| --- | --- | --- |
| React/Mantine/Vite（S7） | native `/v1` DTO | lockfile、build ID、直前のUI artifactを保持。API互換テスト / Lockfile/build ID, retain previous UI artifact, API contract tests |
| CLI | native `/v1` DTO、`--json` の形、終了code / Native `/v1` DTOs, the `--json` shape and exit codes | serverと同じsourceからbuildし、build IDを記録。出力を読むscriptのために出力の形を契約testで固定 / Build from the same source as the server and record the build ID. Contract tests pin the output shape for scripts that read it |
| Go内部 / Go internals | domain invariants、ports、API | package単位で変更、unit/contract回帰 / Change by package; unit/contract regression |
| API/worker | `/v1`、`schema_version`、capability情報 / Version and capability information | 旧consumerが無視できる任意response fieldだけを互換追加とする。必須field・enum・意味変更は互換性判定とversion移行 / Only optional response fields older consumers can ignore are compatible additions; required fields, enums and semantics require compatibility review and version migration |
| SQLite schema | repository port | migration IDとbackup。private copyで移行・restore検証後に適用 / Migration IDs/backups; test migration/restore on private copies before use |
| Mem0/Qdrant | worker contract | 新しいnamespace/storeへ再構築して比較。旧storeを直接書換えない / Rebuild into a new namespace/store and compare; do not overwrite the old store directly |
| embedding model/dimension | source IDsとmemory schema / Source IDs and memory schema | vector再生成、新collection、切替前照合 / Re-embed into a new collection and verify before switching |
| answer/helper model | model adapter | weights/tokenizer/template/quantization/hashを固定。新しい実験ID / Pin weights/tokenizer/template/quantization/hashes; use a new experiment ID |
| inference runtime | ModelClient/TokenCounter | capability・usage・timeout・token対応を実測確認 / Verify capabilities, usage, timeouts and token accounting |

一度に一つの依存系統を更新し、productionとtest依存をlockします。前版manifest、artifact、移行前backupを私的に保持し、schema互換性を確認せず古いbinaryを新DBへ戻しません。prototype段階では停止時間を許容し、zero-downtime移行は要求しません。

Upgrade one dependency family at a time and lock production/test dependencies. Retain the previous manifest, artifacts and pre-migration backup privately. Never run an old binary against a new DB without schema compatibility. Prototype upgrades may use downtime; zero-downtime migration is not required.

削除済みデータをbackupから復活させないため、restore時も削除台帳を適用します。backup、WAL、exportの保持・削除範囲を文書化します。Genkaiではlocal storageとprivate永続storageを区別し、共有filesystem上のlive SQLite/WALを前提にしません。

Apply the deletion ledger during restore to avoid resurrecting deleted data. Document retention/deletion of backups, WAL and exports. On Genkai, distinguish local storage from private persistent storage; do not assume live SQLite/WAL is safe on a shared filesystem.

## 変更の範囲と検証 / Change scope and verification

| 変更例 / Change example | 主に変更する場所 / Primary change location | 受入証拠 / Acceptance evidence |
| --- | --- | --- |
| 推論providerの応答形式 / Provider response shape | Ollama adapterと対応fixture / Adapter and corresponding fixtures | ModelClient契約、usage欠測、timeout・取消 / ModelClient contract, missing usage, timeout/cancellation |
| 訂正・turn回復の規則 / Correction or recovery rule | domain/appと契約正本 / Domain/app and contract reference | 同じ合成入力で状態・重複・回復の回帰 / State, duplicate and recovery regression with the same synthetic inputs |
| Mem0更新 / Mem0 upgrade | Worker内部、依存lock、必要ならworker protocol / Worker internals, locks and protocol only if necessary | 出典、scope隔離、削除、operation replay / Provenance, scope isolation, deletion and replay |
| CLIの表示・commandの変更 / CLI output or command change | CLI subcommandと出力のfixture / CLI subcommands and output fixtures | `--json` と終了codeの回帰、日本語の入出力、stdoutとstderrの分離 / Regression of `--json` and exit codes, Japanese input/output and stdout/stderr separation |
| UI library更新（S7） / UI library upgrade (S7) | webのview・provider設定とlock / Web views/provider setup and locks | DTO境界を保ったkeyboard、IME、狭幅、削除後cache検証 / Stable DTO boundaries plus keyboard, IME, narrow layout and deletion-cache checks |
| 安全対策の変更 / Security control change | 該当する層のcodeと[安全設計](#security)の表 / Code in the affected layer and the [security design](#security) table | 安全性の回帰test（S4-9）と点検表の該当項目 / The security regression suite (S4-9) and the matching checklist items |
| DB schema変更 / DB schema change | SQLite adapterとmigration / Adapter and migration | 前版snapshotからの移行・restore、整合制約 / Migration/restore from the prior snapshot and integrity constraints |

package数ではなく、変更理由が一か所に収まり、外部依存をfakeへ置換でき、失敗の意味が境界で維持されることを保守性の基準にします。新しい抽象化には実際の利用者を必要とし、共通wrapper・plugin機構・event busを先に作りません。構造reviewではimport循環・vendor型の漏出・正本への二重書込を確認します。

Maintainability means one place for a reason to change, replaceable external dependencies and preserved failure semantics across boundaries, rather than a larger package count. Require actual consumers for new abstractions; do not prebuild wrapper libraries, plugin systems or event buses. Structure review checks import cycles, leaking vendor types and multiple writers to authoritative data.

利点はprovider・store・CLI・UI更新の影響を閉じ込められることです。一方、境界のDTO変換、fake契約の維持、Go/worker間の部分失敗処理は必要になります。初期版では短い明示的なmappingと少数の契約fixtureを優先し、code generation基盤は契約が安定してから再検討します。

This limits the impact of provider, storage, CLI and UI changes. The cost is explicit DTO mapping, maintaining fake contracts and handling partial failures between Go and the worker. Initially favor small explicit mappings and a few meaningful contract fixtures; reconsider code generation only after contracts stabilize.

<a id="security"></a>

## 安全設計 / Security design

2026-11-27に利用者が求めた追加です。防御を一か所に頼らず層ごとに置き、各対策にtestを付けます。「現在」の列は今のcodeの状態で、「提案」と書いた対策はまだ存在しません。testで確認するまで安全とは表現しません。本原型は第三者の評価を受けたsystemではありません。

The user asked for this addition on 2026-11-27. Defenses are layered rather than resting on one point, and each control has a test. The "now" column states what the code does today; a control marked proposed does not exist yet. Nothing is called secure until a test confirms it. This prototype has had no third-party assessment.

最初に脅威を書き出します（L4-10）。守る資産は会話の本文、派生記憶、token、GPUの実行枠、評価結果です。信頼の境界は、呼出し側とGo、Goとworker、GoとOllama、processとfile、後にはbrowserとGoの間にあります。

Threats are written down first (L4-10). The assets are conversation text, derived memory, tokens, GPU execution capacity and evaluation results. Trust boundaries lie between callers and Go, Go and the worker, Go and Ollama, processes and files, and later the browser and Go.

| 層 / Layer | 脅威 / Threat | 対策 / Control | 現在 / Now | 回 / Lesson |
| --- | --- | --- | --- | --- |
| 認証 / Authentication | tokenの推測、漏えい、使い回し / Guessed, leaked or shared tokens | 呼出し側ごとに `crypto/rand` で作ったtoken。保存はhashだけ。定数時間の比較。失効と更新 / A token per caller from `crypto/rand`, stored only as a hash, compared in constant time, with revocation and rotation | 環境変数の一つのtokenを完全一致で検索 / One token from the environment, looked up by exact match | L4-10 |
| 認可 / Authorization | IDを指定して他人のsessionや記憶を読む / Reading another owner's session or memory by ID | use caseでの所有者確認。他人のobjectは不存在と同じ404。Goとworkerの両方でscopeを検証 / Ownership checks in use cases, the same 404 for foreign and absent objects, and scope validation in both Go and the worker | sessionは実装済み（L3）。turnと記憶は提案 / Implemented for sessions (L3); proposed for turns and memory | L4-7、S3-6 |
| 入力 / Input | 巨大・不正・未知のfieldを含む要求 / Oversized, malformed or unknown-field requests | body上限、厳格なJSON、未知の重要fieldの拒否、model許可一覧 / Body limit, strict JSON, rejection of unknown critical fields and the model allowlist | session作成で実装済み（4096 bytes） / Implemented for session creation (4096 bytes) | L4-8、S4-9 |
| 濫用 / Abuse | 連続要求でGPUや他の利用者の枠を使い切る / Exhausting the GPU or other callers' capacity with a flood | 呼出し側ごとのrate limit、待機数の上限、生成期限、serverの期限 / Per-caller rate limits, a bound on waiting turns, a generation deadline and server timeouts | serverの期限だけ実装済み / Only server timeouts are implemented | L4-5、L4-6、L4-11 |
| 監査 / Audit | 不正な利用に気づけない。logから秘密が漏れる / Misuse goes unnoticed, or logs leak secrets | 構造化した監査log。誰が・何を・結果を記録し、tokenと本文は記録しない。改行の注入を防ぐ / A structured audit log of who, what and outcome, never tokens or message text, hardened against line injection | 提案 / Proposed | L4-11 |
| 通信 / Transport | 平文HTTPからtokenと会話を読まれる / Tokens and conversations read from plain HTTP | loopbackだけで待ち受ける。他の機械からはSSHのport転送。TLSは別に決める / Listen on loopback only, reach from other machines by SSH port forwarding, decide TLS separately | loopback以外を拒否（実装済み） / Non-loopback addresses are refused (implemented) | S4-8 |
| 保存 / Storage | SQL injection、DB・WAL・backup・exportの読取り、削除済みdataの復活 / SQL injection, reading the DB, WAL, backups or exports, and resurrection of deleted data | parameter付きのquery、fileの権限、backupの置き場所、削除台帳 / Parameterized queries, file permissions, backup placement and the deletion ledger | 提案。保存はまだmemory内 / Proposed; storage is still in memory | S2-5、S2-9、S3-7 |
| 供給網 / Supply chain | 悪意ある依存やmodel file、既知の脆弱性 / Malicious dependencies or model files, and known vulnerabilities | 導入前の監査、versionとhashの固定、既知の脆弱性の確認、npmの導入時scriptの確認 / Audit before installation, pinned versions and hashes, checks for known vulnerabilities and review of npm install scripts | Goの依存は標準libraryだけ / Go depends on the standard library only | S2-1、S7-1 |
| LLM固有 / LLM-specific | 発言や記憶に埋めた指示（prompt injection）、記憶の汚染、他のscopeへの漏えい / Instructions planted in messages or memory (prompt injection), memory poisoning and leakage across scopes | 記憶と履歴をdataとして区切って渡す。検索はscope内だけ。出典を保持。モデルの回答を信頼しない / Pass memory and history as delimited data, search only within scope, keep provenance and never trust model output | 提案 / Proposed | S3-6、S3-9 |
| 外向き通信 / Outbound traffic | 会話がcloudやtelemetryへ送られる / Conversations sent to a cloud or telemetry | local providerの明示設定、telemetryの無効化、遮断下での検証 / Explicit local providers, telemetry off and verification with outbound traffic blocked | 提案 / Proposed | S3-8 |
| 出力 / Output | 回答に含まれる制御文字やscriptが端末やbrowserで動く / Control characters or script in a reply executing in a terminal or browser | CLIは制御文字をescape。UIはplain text表示とCSP / The CLI escapes control characters; the UI renders plain text under a CSP | 提案 / Proposed | S4-5、S7-8 |
| browser（S7） | XSS、CSRF、clickjacking、tokenの窃取 / XSS, CSRF, clickjacking and token theft | [browserでの安全対策](ui.md#browser-security)の表 / The [browser security](ui.md#browser-security) table | 提案 / Proposed | S7-7、S7-8 |

S4-9で、ここまでの対策を攻撃する側の視点で点検します。所有者の全組合せ、不正入力のfuzz、上限の確認、logの検査を回帰testにし、点検表の結果を記録します。この点検の合格が、研究UIを加える前の関門です。S7の最後に、browserの項目を加えて同じ点検をもう一度行います。

S4-9 reviews the controls so far from an attacker's point of view. Every ownership combination, fuzzing of invalid input, limit checks and log inspection become regression tests, and the checklist results are recorded. Passing this review is the gate before the research UI is added. At the end of S7 the same review runs again with the browser items included.

根拠にする資料は[参考候補](sources.md#security-references)にあり、まだ開いていません。L4-10で読んで版と参照日を記録します。検査toolの導入には許可が必要です。攻撃の練習は合成dataと自分のlocal環境だけで行い、他人のsystemには行いません。

The intended references are listed as [candidates](sources.md#security-references) and have not been opened yet; L4-10 reads them and records edition and access date. Installing scanners needs authorization. Attack exercises use synthetic data and the user's own local environment only, never anyone else's system.

<a id="open-questions"></a>

## 未決事項と決定時期 / Open questions and decision points

2026-11-18の設計レビューで見つかった未決事項です。右列は提案で、決定者は利用者です。各項目は[学習ルート](../PLAN.md#teaching-route)の該当回で、実装前に本人が選択肢を比較して決めます。決定後は該当する契約へ移し、この表から外します。

The design review on 2026-11-18 found these unresolved points. The right column is a proposal, and the user decides. Each item is settled in the matching lesson of the [learning route](../PLAN.md#teaching-route), where the user compares options before implementing. Move each decision into the relevant contract and remove it from this table.

最初の10項目は2026-11-27に加わりました。HTTP serviceを残すかどうかは、後でUIを加える方針で決着したので表から外しました。 / The first ten items were added on 2026-11-27. Whether the HTTP service stays was settled by the direction to add a UI later, so it is no longer listed.

| 未決事項 / Open question | 現状と影響 / Current state and impact | 提案と決定回 / Proposal and lesson |
| --- | --- | --- |
| binaryの構成 / Binary layout | 現在は引数なしの `bridge` がserverを起動する。CLIを同じbinaryに置くか別にするかが未定 / Today a bare `bridge` starts the server. Whether the CLI lives in the same binary is undecided | 一つの `bridge` に `serve` とclientのsubcommandを置く。配布物が一つで、serverとCLIの版が必ず一致する。引数なしの起動は使い方の表示に変わる。S4-1 / One `bridge` with `serve` and client subcommands: a single artifact whose server and CLI versions always match. A bare invocation then prints usage. S4-1 |
| CLIが手元に持つ状態 / State the CLI keeps locally | 選択中のsession、次のmodel、送信中のrequest IDをどこに置くかが未定 / Where the current session, next model and in-flight request ID live is undecided | 状態fileを作らない。sessionは毎回引数で渡し、次のmodelはserverの記録から求め、request IDは送信前に表示する。S4-3とS4-4 / No state file. The session is an argument every time, the next model is derived from the server record and the request ID is printed before sending. S4-3 and S4-4 |
| 待機中のCtrl+C / Ctrl+C while waiting | 中断がturnの取消を意味するかが未定 / Whether an interrupt cancels the turn is undecided | 待機だけをやめ、turnは続ける。取消は `turn cancel` で明示する。HTTP切断後も続行する契約と同じ扱い。S4-3 / Stop waiting only and let the turn continue. Cancellation is the explicit `turn cancel`, matching the contract that work continues after a disconnect. S4-3 |
| 出力の契約 / Output contract | 終了codeの割当、`--json` の形、messageの言語が未定。scriptとrunnerが依存する / Exit-code assignment, the `--json` shape and message language are undecided, and scripts and the runner depend on them | [終了codeの案](cli.md)を使い、error `code` と短い英語の説明を表示する。`--json` はserverのDTOに沿う。S4-1とS4-5 / Use the [proposed exit codes](cli.md) and print the error `code` with a short English explanation. `--json` follows the server DTOs. S4-1 and S4-5 |
| 呼出し側ごとの認証 / Authentication per caller | 現在は環境変数の一つのtokenが一人の所有者に対応する。他の人のprogramが同じtokenを使うと、会話と記憶が同じ所有者に混ざる / Today one token from the environment maps to one owner. If other people's programs share it, their conversations and memories fall under the same owner | 呼出し側ごとにtokenと `user_id` を発行し、repository外のfileから対応表を読む。隔離は既存の所有者確認がそのまま担う。tokenはhashで保存する。L4-10 / Issue a token and `user_id` per caller and load the table from a file outside the repository, storing tokens as hashes. The existing ownership checks then provide isolation. L4-10 |
| loopbackの外からの到達 / Reach beyond loopback | 現在はloopback以外のaddressを拒否する。他の機械のprogramは届かない。平文HTTPを外へ開くとtokenと会話が読まれる / Non-loopback addresses are refused today, so programs on other machines cannot connect. Plain HTTP opened to a network exposes tokens and conversations | loopbackのままにし、他の機械からはSSHのport転送で使う。TLSでの待受けは必要が確かめられてから別に決める。玄界では同一hostの他jobも考慮する。S4-8 / Stay on loopback and reach it from other machines through SSH port forwarding. Decide on a TLS listener separately once the need is confirmed. On Genkai, account for other jobs on the same host. S4-8 |
| 機械可読のschema / Machine-readable schema | 他言語の呼出し側が頼れるのは文書とGoのfixtureだけ / Callers in other languages have only prose and Go fixtures to rely on | OpenAPI fileを一つ置き、契約fixtureと照合するtestでずれを検出する。言語別のSDKは作らない。[呼出し方](contracts.md#other-programs)。L4-8 / Keep one OpenAPI file and detect drift with a test against the contract fixtures. Build no per-language SDKs. See [how to call](contracts.md#other-programs). L4-8 |
| 濫用の上限 / Abuse limits | 呼出し側ごとの要求数に上限がない。一人のprogramが全員のGPU枠を使い切れる / No per-caller request limit exists, so one program can consume everyone's GPU capacity | 呼出し側ごとのtoken bucketと同時turn数の上限。超過は429と再試行までの時間を返す。値は設定にする。L4-11 / A token bucket and a concurrent-turn cap per caller. Excess returns 429 with a retry time. The values are configuration. L4-11 |
| 保存dataの暗号化 / Encryption of stored data | SQLiteのfileは平文。同じ機械の他の利用者やbackupの持ち出しが脅威になるかは環境による / The SQLite file is plaintext. Whether other users of the machine or copied backups are a threat depends on the environment | まずfileの権限とbackupの置き場所で守り、暗号化は脅威を書き出してから決める。暗号化libraryは新しい依存になる。S2-9 / Protect with file permissions and backup placement first, and decide on encryption after the threat model. An encryption library would be a new dependency. S2-9 |
| browserでのtokenの置き場所 / Where the browser keeps the token | memoryだけに置く案は再読込のたびに再認証が要る。cookieにする案はCSRF対策が要る / Memory-only needs reauthentication on every reload; a cookie needs CSRF protection | memoryだけに置く。再認証の手間を実際に測ってから `HttpOnly` cookieを検討する。S7-8 / Keep it in memory only, and consider an `HttpOnly` cookie after measuring the reauthentication burden. S7-8 |
| 応答の `revision` の意味 / Meaning of the response `revision` | 現行APIは常に1を返す。[契約](contracts.md)は曖昧な単一revisionを避けるとするが、この値がどの版かは未定義 / The current API always returns 1. The [contract](contracts.md) rejects a single ambiguous revision, yet this field's role is undefined | 原文版 `source_revision` として扱い、user・assistant messageの追記ごとに増やす。L4-4 / Treat it as `source_revision` and advance it on every appended user or assistant message. L4-4 |
| 会話順の表現 / Representation of conversation order | 「追記順」と定めるが、Messageにあるのは `CreatedAt` のみ。同時刻で順序が決まらない / Order is defined as append order, but Message carries only `CreatedAt`, which cannot break ties | session内の単調増加の連番をstoreが採番する。時刻は表示用。L4-4 / The store assigns a monotonically increasing per-session sequence number. Timestamps are for display. L4-4 |
| session間の実行枠 / Execution capacity across sessions | 「暗黙queueなし」はsession内の規則。GPUを直列にすると、別sessionの `accepted` は待つことになるが、待ち数の上限と超過時の応答が未定義 / "No implicit queue" is a per-session rule. Serial GPU work makes `accepted` turns in other sessions wait, with no bound or overflow response defined | 待機数に上限を置き、超過は受理せず503と再試行可能の表示を返す。待ち時間はqueue時間として計測。L4-5 / Bound the number of waiting turns. Refuse admission beyond it with 503 marked retryable. Measure the wait as queue time. L4-5 |
| 生成期限と時計 / Generation deadline and clock | 期限の値と起点が未定義。実時間に依存するtestは不安定になる / The deadline value and its starting point are undefined, and tests tied to wall-clock time become flaky | 期限は設定値とし、`generating` への遷移時から数える。時計と期限をtestから差し替えられる形で受け取る。L4-6 / Make the deadline configurable, counted from the move to `generating`. Accept the clock and deadline in a form tests can replace. L4-6 |
| 失敗理由の記録 / Recording failure reasons | Turnの必須項目に失敗理由がなく、期限切れ・model error・回復不能を区別できない / Turn has no failure-reason field, so deadline, model error and unresolved recovery are indistinguishable | 終端状態に固定の理由codeを一つ持たせ、error envelopeの `code` と同じ語彙を使う。L4-8 / Give terminal states one reason code from a fixed vocabulary shared with the error envelope's `code`. L4-8 |
| 認可の置き場所 / Where authorization lives | 設計はuse caseの責務とするが、現行は取得handler内で所有者を確認している / The design assigns it to use cases, while the current GET handler checks ownership itself | turnのuse caseを書く時に、session取得と所有者確認を一つの関数へまとめ、handlerはそれを呼ぶ。L4-7 / When writing the turn use case, combine session lookup and ownership checking in one function that handlers call. L4-7 |
| 停止時の受理済みturn / Admitted turns at shutdown | 2026-11-19に [server](../cmd/bridge/server.go) へHTTPの期限、signalによる停止、設定できるloopback addressを実装済み。期限の値（header 5秒、読取り15秒、書込み30秒、idle 60秒、停止猶予10秒）は提案。停止時に `accepted`／`generating` のturnをどう確定するかは未定 / On 2026-11-19 the [server](../cmd/bridge/server.go) gained HTTP timeouts, signal-driven shutdown and a configurable loopback address. The timeout values (header 5s, read 15s, write 30s, idle 60s, shutdown grace 10s) are proposals. How `accepted` and `generating` turns settle at shutdown is undecided | 本人が期限の値を確認する。停止要求後は新しいturnを受理せず、猶予内に終わらない生成は取消して `failed` にする。L4-9 / The user confirms the timeout values. After a stop request, admit no new turns and cancel generation that outlasts the grace period, marking it `failed`. L4-9 |
| 記憶更新が回復しないturn / A turn whose memory update never recovers | 回答保存済みで記憶更新が永久に失敗する場合の出口がなく、sessionが止まり続ける。現在の出口はsession削除のみ / A saved answer whose memory update fails permanently has no exit, so the session stays blocked. Deleting the session is the only way out | 所有者が「記憶未反映のまま完了」を明示的に選べる操作を設け、そのturnを出典に未反映と記録する。S3-5 / Add an explicit owner action to complete the turn with memory not applied, and record that gap in provenance. S3-5 |
| worker間の認証 / Go-to-worker authentication | 「service認証」とあるが方式が未定義 / "Service authentication" is required without a mechanism | 起動時に生成した共有secretをloopback接続のheaderで渡す。Genkaiでは同一host上の他jobを考慮する。S3-1 / Pass a shared secret generated at startup in a header over loopback. On Genkai, account for other jobs on the same host. S3-1 |
| 統計手法と採点者 / Statistics and raters | 「case単位のpaired uncertainty」の手法、採点者数、一致度の指標が未定義。目標の品質90%がどの方式を指すかも曖昧 / The method behind "paired case-level uncertainty", rater count and agreement measure are undefined. Which policy the 90% quality target applies to is also unclear | case単位のpaired bootstrapで信頼区間を出し、二重採点部分はCohenのκを報告する。90%は各方式に個別に適用する。freeze前のS6-2 / Use a case-level paired bootstrap for intervals and report Cohen's κ on the double-scored subset. Apply 90% to each policy separately. S6-2, before freezing |

## 見直す条件 / Revisit triggers

multi-process書込、複数worker、常時サービスが必要になればPostgreSQLやQdrant serverを別ADRで検討します。Ollamaの能力・throughputが実測で不足すればllama.cpp/vLLM adapterを別実験にします。Mem0とモデルを同じ種類の代替物として比較しません。回答を生成中に少しずつ表示する必要が出たら、streamingのrouteを別ADRで検討します。loopbackの外へ常時公開する必要が出たら、TLS終端・利用者管理・監視を別ADRで検討します。

If multi-process writers, multiple workers or a persistent service become necessary, consider PostgreSQL/Qdrant server in a separate ADR. If measured capabilities/throughput are insufficient, evaluate a llama.cpp/vLLM adapter separately. Do not compare Mem0 and a model as interchangeable alternatives. If replies need to appear incrementally while generating, consider a streaming route in a separate ADR. If the service must be permanently exposed beyond loopback, consider TLS termination, user management and monitoring in a separate ADR.
