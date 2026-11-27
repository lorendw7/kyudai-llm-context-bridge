# Kyudai LLM Context Bridge

**モデルを切り替えても、文脈をつなぐ。 / Change the model. Carry the context.**

[公開サイト / Website](https://lorendw7.github.io/kyudai-llm-context-bridge/) · [作業計画 / Plan](PLAN.md) · [学習方法 / Learning guide](CONTRIBUTING.md) · [MIT](LICENSE)

九州大学RAプロジェクトに由来する独立した研究リポジトリです。LLM切替時の文脈継続、記憶の正確さ、処理負担を調べます。大学公式サイト・公認サービスではありません。共有対象は推論資源であり、会話・記憶・権限は利用者ごとに分離します。

This independent repository originates in a Kyushu University RA project. It studies context continuity, memory correctness and processing overhead when switching LLMs. It is not an official university website or endorsed service. Inference capacity may be shared; conversations, memories and permissions remain separate.

## 受入済みの範囲 / Accepted implementation

2026-11-13現在、localhost health、Session／Messageとturn request／ModelClient契約に加え、Bearer認証、所有sessionの作成・取得・一覧、許可model ID一覧、環境設定からの起動を検証済みです。sessionはmemoryのみで、応答はmetadataです。turn処理、切替、model adapter、記憶、永続保存、CLI、runnerは未受入です。2026-11-27に研究UIの計画を取り下げ、人が使う入口を[CLI](docs/cli.md)にしました。品質・token削減の実測結果はありません。検証証拠・進捗は[PLAN.md](PLAN.md)を参照します。

As of 2026-11-13, verified implementation includes localhost health, Session/Message and turn-request/ModelClient contracts, Bearer authentication, owned-session creation/retrieval/listing, configured model IDs and environment-based startup. Sessions are in memory and responses contain metadata. Turn processing, switching, model adapters, memory, persistence, the CLI and the runner remain unaccepted. On 2026-11-27 the research UI plan was withdrawn and the [CLI](docs/cli.md) became the human entry point. No quality or token-saving results have been measured. See [PLAN.md](PLAN.md) for verification evidence and progress.

利用者は研究の核心codeを書きます。対象は状態・並行実行・保存の原子操作・回復・context policy・記憶の調停・CLIの送信後の待機と回復・再生と集計など、設計判断を含む部分です。assistantはその判断をcodeへ写す外周を書きます。対象はHTTP handlerとDTO変換、error envelope、起動と停止、設定読込、migration実行の枠組み、workerのHTTP外殻、CLIのcommandの枠組みと表示、runnerのCLI、job wrapperです。assistantはこのほかテスト・合成fixture・test-only fake、会話内の中国語／英語による指導・レビュー、依頼された日英文書、公開サイト（`website/`）と美術素材を担当します。2026-11-19に利用者が決めた分担です。CLIへの当てはめは2026-11-27のassistantの提案で、各回の担当は[作業計画](PLAN.md#teaching-route)にあります。

The user writes the core research code: the parts that carry design decisions, such as state, concurrency, atomic storage operations, recovery, context policies, memory reconciliation, waiting and recovery after a CLI send, replay and aggregation. The assistant writes the surrounding code that translates those decisions: HTTP handlers and DTO mapping, the error envelope, startup and shutdown, configuration loading, the migration runner, the worker HTTP shell, the CLI command framework and output, the runner CLI and job wrappers. The assistant also owns tests, synthetic fixtures and test-only fakes, Chinese/English teaching and review in chat, requested Japanese/English documentation, the public website (`website/`) and visual assets. The user set this division on 2026-11-19. Its application to the CLI is the assistant's proposal of 2026-11-27, and the [plan](PLAN.md#teaching-route) names the author for each lesson.

<a id="product-scope"></a>

## 何ができる製品か / What the finished product lets you do

研究者がAに「締切は12日」と伝え、同じ会話でBに「15日に訂正」と伝えてから、再びAに締切を尋ねる。製品はAに渡した文脈と実際の回答を記録する。15日を答えられたかは評価の対象で、12日と答えた場合も結果に残す。

A researcher tells A the deadline is the 12th, corrects it to the 15th with B in the same conversation, then asks A for the deadline. The bridge records which sources were supplied to A and what A actually answered. An answer of “the 12th” stays in the results so the researcher can see where continuity failed.

| 想定する成果物 / Proposed delivery | 利用者ができること / User outcome | 完成を示す例 / Demonstration |
| --- | --- | --- |
| 最初のローカルAPI版 / First usable local API (S1–S2) | 全履歴で許可済みの実モデルを切り替え、回答を取得する / Switch authorized local models using full history and retrieve replies | 同じsessionでA → B → A。再起動後も保存回答が元の要求IDで返り、再送で生成が増えない / Switch A → B → A, restart, then retrieve the saved reply by its original request ID and replay without another generation |
| CLIで使える研究版 / CLI research version (S3–S4) | CLIだけで会話を再開し、次のモデルを選び、訂正の出典を表示し、所有sessionをexport・削除する / Using only the CLI, resume a conversation, choose the next model, show a correction’s source, export and delete an owned session | Bの回答後、`--model` を省いた送信はBを使う。研究者は訂正元の発言を確認でき、削除完了後は会話と派生記憶を取得できない / After B replies, a send without `--model` uses B. The researcher can inspect the correction message, then verify that conversation and derived memory are inaccessible after deletion |
| 再現可能な研究評価 / Reproducible evaluation (S5–S6) | 同じ固定シナリオを全履歴・逐次要約・Mem0で再生し、品質・token・時間・GPU使用を比較する / Replay frozen scripts with full history, rolling summary and Mem0 and compare quality, tokens, time and GPU use | manifestからモデル・設定・caseを特定し、JSONL/CSVと日英報告で失敗・欠測も追跡する / Trace models, configuration and cases from the manifest through JSONL/CSV to the bilingual report, including failures and missing observations |

S1は合成モデルで通信・履歴・失敗を検証する段階です。実モデルの利用には配布物・依存の監査と実行許可が必要です。各段階の受入は[作業計画](PLAN.md#teaching-route)にあります。

S1 exercises transport, history and failures with synthetic models. Real-model use depends on artifact and dependency review and execution authorization. See the [plan](PLAN.md#teaching-route) for each stage’s acceptance criteria.

## 読む順序 / Reading order

L3の続きは[Bの回答をいつ「保存済み」と表示できるか](CONTRIBUTING.md#l4-state-lesson)の一問から始めます。今回書くコードと受入条件は[L4の作業計画](PLAN.md#l4-entry)にあります。 / To resume after L3, start with [When can B’s reply appear as saved?](CONTRIBUTING.md#l4-state-lesson). The [L4 plan](PLAN.md#l4-entry) defines the first code change and when to stop.

| 文書 / Document | 内容 / Purpose |
| --- | --- |
| [PLAN.md](PLAN.md) | 唯一の作業一覧・段階課題・受入基準 / Sole action checklist, stage projects and acceptance gates |
| [CONTRIBUTING.md](CONTRIBUTING.md) | 一回の授業、分担、検証 / Lesson method, ownership and validation |
| [docs/architecture.md](docs/architecture.md) | 技術選定、モジュール境界、保存と更新 / Technology choices, module boundaries, persistence and upgrades |
| [docs/contracts.md](docs/contracts.md) | データ、API、失敗と再試行 / Data, API, failures and retries |
| [docs/cli.md](docs/cli.md) | CLIのcommand、送信と回復、出力と終了code / CLI commands, sending and recovery, output and exit codes |
| [docs/design.md](docs/design.md) | 視覚設計、説明用素材、公開サイトとの対応 / Visual design, explanatory assets and their public-site use |
| [docs/evaluation.md](docs/evaluation.md) | 比較方式、計測、玄界、報告 / Controls, metrics, Genkai and reporting |
| [docs/sources.md](docs/sources.md) | 一次資料、ライセンス、未確認事項 / Primary sources, licenses and unresolved verification |
| [AGENTS.md](AGENTS.md) / [SECURITY.md](SECURITY.md) | 協働規則・機密情報 / Collaboration and sensitive data |

文書は節ごとに日英を対応させます。識別子・コマンドは共通です。`LICENSE` は標準MIT原文を保持し、日本語の説明を[資料台帳](docs/sources.md)に置きます。進捗はPLAN.mdだけで更新します。

Documentation pairs Japanese and English within each section. Identifiers and commands are shared. `LICENSE` retains the standard MIT text; a Japanese explanation is in the [source register](docs/sources.md). Progress is tracked only in PLAN.md.

## 配置と責務 / Repository map

| 場所 / Location | 内容と状態 / Content and status |
| --- | --- |
| `cmd/bridge/` | 現在のGo実装と隣接test / Existing Go implementation and adjacent tests |
| `docs/` | 構成・契約・UI・評価・出典の正本文書 / Reference documents for architecture, contracts, UI, evaluation and sources |
| `website/` | 既存の公開説明サイト。固定の合成例 / Existing public information site with fixed synthetic examples |
| `website/assets/` | 公開サイトと文書で共有する説明図・サイト素材 / Explanatory artwork shared by the site and docs, plus site assets |
| `scripts/` | 既存のrepository検証 / Existing repository verification |
| `internal/` | 未作成の提案配置。module分割と、CLI・runnerが共有するclient / Proposed, not yet created: module boundaries and the client shared by the CLI and runner |
| `presentations/` | ignoredのlocal発表資料。stage／upload対象外 / Ignored local presentation material; excluded from staging/upload |

<a id="route-a-software-architecture"></a>

## 推奨構成 / Recommended architecture

![同じsessionと出典を保ちながらモデルAからB、Aへ切り替える設計案 / Proposed model A to B to A handoff preserving the same session and sources](website/assets/context-continuity.svg)

説明用の合成図です。切替・記憶の動作や効果を実測した結果ではありません。 / This is a synthetic design illustration, not a measured demonstration of switching or memory effectiveness.

選択済みのA案はGo + Ollama + Mem0 OSSです。Goが会話・権限・切替を管理し、Python workerが記憶を扱い、Ollamaがlocal推論を提供する設計です。SQLiteによる正本保存とworker内のQdrant local modeは提案です。人はCLIから `/v1` のAPIを使い、画面は作りません（2026-11-27）。導入・version固定・依存監査は未実施です。

Selected Route A uses Go + Ollama + Mem0 OSS. The design assigns conversations, authorization and switching to Go, memory to a Python worker, and local inference to Ollama. SQLite authoritative storage and Qdrant local mode inside the worker are proposals. People use the `/v1` API through a CLI, and no screen is built (2026-11-27). Installation, version pinning and dependency audit remain outstanding.

CLIとrunnerは同じAPIと同じclientを使います。責務・保存・回復・更新の詳細は[構成設計](docs/architecture.md)、比較条件は[評価設計](docs/evaluation.md)に集約します。

The CLI and runner use the same API through the same client. See [architecture](docs/architecture.md) for responsibilities, persistence, recovery and upgrades, and [evaluation](docs/evaluation.md) for comparison conditions.

保守の基本は `HTTP → app → domain` の依存方向と、app側の小さなinterfaceを実装するadapterです。状態判断はapp、原子的な正本保存はSQLite adapter、派生記憶はworkerが所有します。一括分割せず、連携とtest境界が必要になった時にpackageを抽出します。CLIとrunnerは `client` だけに依存し、server側とはHTTPでつながります。公開websiteはこれらと独立です。

Maintainability rests on `HTTP → app → domain` dependencies and adapters implementing small app-owned interfaces. App owns state decisions, the SQLite adapter owns atomic source persistence, and the worker owns derived memory. Extract packages as integration and test boundaries require them. The CLI and runner depend only on `client` and meet the server side over HTTP. The public website is independent of both.

## 研究範囲 / Research scope

一般向けSaaS、学内SSO、多人数同時運用、ファイルRAG、agent tools、商用API、fine-tuningは初期範囲外です。研究用の画面と第三者UI互換adapterは2026-11-27に範囲から外しました。

A public SaaS, campus SSO, multi-user production operations, file RAG, agent tools, commercial APIs and fine-tuning are outside initial scope. A research screen and a third-party UI adapter were removed from scope on 2026-11-27.

## 現在のコードを確認する / Check existing code

Goが必要です。go.modはGo 1.25.0を宣言し、検証環境はWindows Go 1.26.0です。model／memory packageは不要です。Python標準libraryによる起動testは合成環境設定を子processだけに渡し、一時buildとserverを後処理します。127.0.0.1:8080を空けて実行してください。

Go is required: go.mod declares Go 1.25.0; verification used Windows Go 1.26.0. No model/memory packages are needed. The Python standard-library startup test passes synthetic configuration only to child processes and cleans up its temporary build and server. Port 127.0.0.1:8080 must be free.

```sh
go test ./... -count=1
go vet ./...
go test -race ./... -count=1
python tests/test_bridge_startup.py
```

手動の `go run ./cmd/bridge` には `BRIDGE_TOKEN`、`BRIDGE_LAB_ID`、`BRIDGE_USER_ID`、`BRIDGE_MODEL_IDS` が必要です。model IDsは空白なしのcomma区切りで、空値・重複を拒否します。任意の `BRIDGE_LISTEN_ADDRESS` はloopback IPの `host:port` で、未設定なら `127.0.0.1:8080` です。Ctrl+Cで受付を止め、進行中の要求を待ってから終了します。credentialはlocal環境で設定しcommitしません。healthだけが公開で、`/v1/` routeは `Authorization: Bearer <token>` を要求します。model一覧は設定値であり推論準備状態ではありません。

Manual `go run ./cmd/bridge` requires `BRIDGE_TOKEN`, `BRIDGE_LAB_ID`, `BRIDGE_USER_ID` and `BRIDGE_MODEL_IDS`. Model IDs are comma-separated without whitespace; empty/duplicate entries are rejected. Optional `BRIDGE_LISTEN_ADDRESS` is a loopback-IP `host:port`, defaulting to `127.0.0.1:8080`. Ctrl+C stops accepting, waits for in-flight requests and exits. Set credentials locally and never commit them. Only health is public; `/v1/` routes require `Authorization: Bearer <token>`. The model list represents configuration, not inference readiness.

```sh
python scripts/check_repository.py
node --check website/assets/demo.js
```

依存導入、大型モデル取得、有料API、遠隔実行は許可範囲内でのみ行います。私的会話、重み、DB、実験raw出力、発表archiveは公開しません。既存 `website/` は広報用の静的siteで、研究APIには接続しません。

Installation, large downloads, paid APIs and remote execution require authorized scope. Private conversations, weights, databases, raw output and presentation archives are not published. Existing `website/` is the static public information site and does not connect to the research API.
