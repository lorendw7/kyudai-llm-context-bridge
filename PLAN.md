# 作業計画 / Plan

更新 / Updated: 2026-11-12 (Asia/Tokyo)。本書だけが作業・進捗の一覧です。他の文書は設計・学習資料です。 / This is the sole action and progress checklist. Other documents contain design and learning references.

## 方針 / Decisions

- A案（Go + Ollama + Mem0 OSS）は2026-11-06に利用者が選択。教員承認、玄界account・計算枠は利用者の報告です。 / Route A was selected by the user on 2026-11-06. Supervisor approval and Genkai account/allocation are user-reported.
- 2026-11-08：小型の自作研究UIにReact + TypeScript + Vite + Mantineを選択。独自CSSは初期範囲に含めない。SQLite正本保存を提案。正確なpackage versionsは未固定。 / React + TypeScript + Vite + Mantine were selected for the small owned research UI. Handwritten CSS is outside the initial scope. SQLite authoritative storage is proposed. Exact package versions are not frozen.
- 利用者が研究用実装コード・設定・script・移行・runnerを書く。assistantはテスト、合成fixture、test-only fake、指導・レビュー、依頼された日英文書を担当。2026-11-12の明示的再割当により公開 `website/` と美術素材もassistantが担当する。研究UI `web/` は本人担当のまま。 / The user writes research implementation code, configuration, scripts, migrations and runner. The assistant owns tests, synthetic fixtures, test-only fakes, teaching/review and requested bilingual documentation. The explicit reassignment on 2026-11-12 also gives the assistant responsibility for public `website/` and visual assets; research UI `web/` remains user-owned.
- OSSとローカル実行が必須。導入・download・有料API・遠隔実行は許可範囲に従う。要件・提案・vendor claim・実測を区別。 / OSS and a local route are required. Installation, downloads, paid APIs and remote runs follow authorized scope. Separate requirements, proposals, vendor claims and measurements.
- 公開は現在資料のみ。履歴置換は今回の作業ではない。発表archiveはignored `presentations/` に保持。 / Publish current material only. History replacement is not part of this task. Presentation archives remain under ignored `presentations/`.

## 唯一の作業一覧 / Sole action checklist

- [x] S0 — 日英文書を案内・構成・契約・UI・評価・資料台帳に整理し、実装済み／提案を区別。 / Organize paired bilingual overview, architecture, contracts, UI, evaluation and source register; distinguish implementation from proposals.
- [x] 初期scaffoldの確認（2026-11-08）：health、Message、当時の2つのGoテスト。現在の授業受入は下記を参照。 / Initial scaffold check (2026-11-08): health, Message and the two tests present then. See the lesson acceptance record below for subsequent work.
- [x] S0追補（2026-11-12）— 既存構成を調査し、READMEの案内、UIの情報階層・状態表示、再利用可能な説明用SVGを整理。[視覚設計](docs/design.md)に用途・原本・日英の代替説明を記録。日英対応、実装／提案の区別、local link・SVG構造・描画確認、既存研究codeと未追跡作業の保持を受入。 / S0 follow-up (2026-11-12): reviewed the structure, clarified README navigation and UI hierarchy/states, and created reusable explanatory SVGs. The [visual design reference](docs/design.md) records uses, sources and JA/EN text equivalents. Accepted paired content, implementation/proposal labels, local-link/SVG/render checks and preservation of research code and untracked work.
- [x] S0設計追補（2026-11-12明示依頼）— modular monolithの依存規則、state/store所有者、段階的package分割、UI境界と変更時の検証を日英で整理。状態名・revision意味・session全体のexport/deleteを整合。研究機能codeの再配置・実装は未実施。 / S0 architecture follow-up (explicitly requested 2026-11-12): documented modular-monolith dependency rules, state/store ownership, incremental package extraction, UI boundaries and change-focused validation in JA/EN. Aligned state names, revision semantics and session-wide export/delete. Research functional code was not relocated or implemented.
- [x] main公開（2026-11-12明示依頼）— 文書・website・素材の14ファイルをcommit `eee6aca` でpush。対応するPages deploy成功と公開日英ページ・素材・三段階例を確認。作業中の認証3ファイルと研究機能codeを除外。証拠は下記。 / Main publication (explicitly requested 2026-11-12): pushed the 14 documentation/site/asset files in commit `eee6aca`; verified the matching successful Pages deployment, public JA/EN pages, assets and three-stage example. Excluded the three authentication work-in-progress files and research functional code. Evidence below.
- [ ] S1 — 契約とfake切替（L1–L4）：session、所有者、A → B → A、重複・timeoutを依存追加なしで説明・検証。 / Contracts and fake switching: explain and verify sessions, ownership, A → B → A, duplicates and timeout without new dependencies.
- [ ] S2 — 実モデルと正本保存（L5–L8）：依存・artifact監査、許可後のOllama、SQLite、再起動・再試行。受入は同じsessionの実モデル切替と保存応答の回復。 / Real models and storage: dependency/artifact audit, authorized Ollama and SQLite, restart/retry; accept real switching in one session and saved-response recovery.
- [ ] S3 — 記憶と比較（L9–L12）：Python契約、Mem0、訂正／出典／隔離、部分失敗回復、export/delete、3方式。受入は外向き通信遮断下の訂正保持、隔離、削除検証。 / Memory and controls: Python contract, Mem0, corrections/provenance/isolation, partial-failure recovery, export/delete and three policies; accept correction retention, isolation and deletion checks with outbound traffic blocked.
- [ ] S4 — React研究UI（L13–L16）：本人がUI実装、assistantが画面・API・E2Eテスト。受入は日英、keyboard、狭幅、再読込、二重送信、失敗表示、削除後cache除去。 / React UI: user implementation, assistant UI/API/E2E tests; accept JA/EN, keyboard/narrow layout, reload, duplicate submit, failures and deletion-cache removal.
- [ ] S5 — 再現性と更新（L17–L20）：runner、stage計測、DB移行、API互換性、依存更新、rollback、開発pilot。受入は再現manifestと旧版からの復元試験。 / Reproduction/upgrades: runner, stage metrics, DB migrations, API compatibility, dependency upgrades, rollback and pilot; accept a reproduction manifest and restore test from the previous version.
- [ ] S6 — 玄界・最終評価・報告（L21–L24）：許可済み環境調査とbounded job、評価固定、最終実行、日英報告。欠測・失敗・未達を保持。 / Genkai/evaluation/report: authorized audit and bounded job, freeze, final execution and bilingual report; retain missing data, failures and unmet targets.
- [x] 公開資料のsource同期（2026-11-12再割当）：assistantが日英websiteを更新し、合成例・現在の受入範囲・説明図を揃え、文言・リンク・keyboard・狭幅／広幅の表示を検証。最初にlocal成果物を受入し、その後の明示依頼で上記main公開を実施。 / Public-site source synchronization (reassigned 2026-11-12): the assistant updated both languages, aligned synthetic examples, accepted scope and explanatory diagrams, and verified wording, links, keyboard and narrow/wide presentation. Accepted locally first, then published to main following the explicit request above.
- [ ] 任意U1 — 第三者UI：必要時のみ互換adapterとsession bindingを設計・検証。主評価はnative API。 / Optional U1: design/test a compatibility adapter and session binding only when needed; primary evaluation uses the native API.

## 授業の対応表 / Lesson mapping

各回は説明→一つの課題→テスト→レビューです。表はS1–S6との対応であり、別の進捗一覧ではありません。次段階は前段階の受入後に説明します。 / Each lesson is explanation → one exercise → tests → review. The table maps to S1–S6; it is not another progress tracker. Explain the next stage after acceptance of the current one.

| 回 / Lesson | 学習・本人の成果物 / Learning and user deliverable | assistantの受入証拠 / Assistant acceptance evidence |
| --- | --- | --- |
| L1 / S1 | Message検証、Sessionと所有者 / Message validation, Session and ownership | role、空値、session分離の合成例 / Role, empty values and session-isolation examples |
| L2 / S1 | interface、API DTO、依存方向 / Interfaces, API DTOs and dependency direction | 正常・異常契約fixture / Valid and invalid contract fixtures |
| L3 / S1 | 認証、session API、model allowlist / Authentication, session API and model allowlist | 偽identity、他user/lab拒否 / Reject forged identity and other users/labs |
| L4 / S1 | turn状態、fake A → B → A / Turn state and fake A → B → A | timeout、取消、同時・重複要求 / Timeout, cancellation, concurrent/duplicate requests |
| L5 / S2 | artifact台帳、設定schema / Artifact manifest and config schema | version、license、hash、秘密分離 / Version, license, hash and secret separation |
| L6 / S2 | Ollama adapter、予算・usage / Ollama adapter, budget and usage | 許可済み実応答、エラー、token検証 / Authorized response, errors and token validation |
| L7 / S2 | SQLite repositoryとmigration / SQLite repository and migration | 制約、transaction、restart / Constraints, transaction and restart |
| L8 / S2 | 重複・不確定失敗の回復 / Duplicate/uncertain-failure recovery | 保存回答再返却、無条件再生成なし / Saved-answer replay; no blind regeneration |
| L9 / S3 | Python workerとversioned契約 / Python worker and versioned contract | fake worker契約、認証・scope / Fake-worker contract, auth and scope |
| L10 / S3 | local Mem0設定と出典 / Local Mem0 config and provenance | 訂正A → B → A、外向き通信遮断 / Correction across A → B → A; outbound traffic blocked |
| L11 / S3 | memory回復、export/delete、revision / Memory recovery, export/delete and revisions | 遅延write、履歴・vector・checkpoint削除 / Late writes and history/vector/checkpoint deletion |
| L12 / S3 | full-history、summary、Mem0 policy / Full-history, summary and Mem0 policies | 同じ入力・予算、明示overflow / Matched input/budget, explicit overflow |
| L13 / S4 | React/TS、Vite、Mantine、日英辞書 / React/TS, Vite, Mantine and JA/EN dictionaries | fixture画面、library styles、型確認 / Fixture views, library styles and type checks |
| L14 / S4 | API clientとchat画面 / API client and chat view | 切替、IME、二重送信・再読込 / Switching, IME, duplicate submit and reload |
| L15 / S4 | memory/source閲覧と削除 / Memory/source inspection and deletion | 所有者、訂正、注入文字列、cache / Ownership, corrections, injection text and cache |
| L16 / S4 | UI状態と実API連携 / UI states and real API integration | keyboard、狭幅、部分失敗、E2E / Keyboard, narrow layout, partial failure and E2E |
| L17 / S5 | runnerとcase/config schema / Runner and case/config schema | dataset分離、決定的fake再生 / Dataset separation and deterministic fake replay |
| L18 / S5 | stage metricsと結果表示 / Stage metrics and result display | 計算fixture、欠測、load分離 / Accounting fixtures, missing data and load separation |
| L19 / S5 | API/DB/worker/UI更新手順 / API/DB/worker/UI upgrades | migration、互換性、backup復元 / Migration, compatibility and backup restore |
| L20 / S5 | 開発pilotと再現package / Development pilot and reproduction package | 失敗、実費、hash、起動停止 / Failures, measured cost, hashes and lifecycle |
| L21 / S6 | 玄界環境・job packaging / Genkai environment and job packaging | 許可済みaudit、bounded smoke、後処理 / Authorized audit, bounded smoke and cleanup |
| L22 / S6 | 評価仕様と独立case固定 / Freeze evaluation protocol and independent cases | rubric、seed、manifest凍結 / Rubric, seed and manifest freeze |
| L23 / S6 | 固定条件で最終実行 / Final execution under frozen conditions | 324 runs案、欠測・隔離suite / Proposed 324 runs, missing data and isolation suite |
| L24 / S6 | 日英報告・再現ガイド / Bilingual report/reproduction guide | 主張と証拠の対応、未達記載 / Claims mapped to evidence, unmet targets reported |

L1–L2は2026-11-10までに完了しました（2/24、残り22）。L1ではSession、所有者、Message検証を、L2では本人実装のversioned `CreateTurnRequest`、project固有の生成入出力、`context.Context`を受ける`ModelClient`をassistantの正常・異常・取消fixtureで検証しました。`go test ./... -count=1` と `go vet ./...` が成功しました。次はL3の認証、session API、model allowlistです。S1全体はL1–L4の受入後に完了します。

L1–L2 were completed by 2026-11-10 (2/24 complete, 22 remaining). L1 verified Session, ownership and Message validation. L2 used assistant-authored valid, invalid and cancellation fixtures to verify the user-written versioned `CreateTurnRequest`, project-specific generation input/output and `ModelClient` accepting `context.Context`. `go test ./... -count=1` and `go vet ./...` passed. L3 is next: authentication, the session API and the model allowlist. S1 completes after acceptance of L1–L4.

## 今回の資料・website検証 / Documentation and website verification

2026-11-12、基点commit `cc7fe99fd93ab4681ed65e781da1591d0316ecdb` に対する未commitの文書・site・素材変更を確認しました。Python 3.13.12で `python scripts/check_repository.py`（35 files）、Node v24.14.0で `node --check website/assets/demo.js`、`git diff --check` が成功しました。2つの新SVGのXML・title/desc・外部参照／script不在と、日英siteのruntime素材がlocal fileのみであることを確認しました。

On 2026-11-12, reviewed the uncommitted documentation/site/asset changes against base commit `cc7fe99fd93ab4681ed65e781da1591d0316ecdb`. `python scripts/check_repository.py` passed for 35 files with Python 3.13.12; `node --check website/assets/demo.js` passed with Node v24.14.0, as did `git diff --check`. Checked both new SVGs for valid XML, titles/descriptions and absence of scripts/external references, and verified that both site languages use only local runtime assets.

本機previewで日英の3段階切替、締切12→15→15、単一の選択状態、Enter／Space、Tab focus、skip linkを確認しました。1280pxの広幅、390pxの狭幅、日英320pxで表示を確認し、検査した幅で横方向のoverflowはありませんでした。2図の描画とwebsiteへの画像読込を確認。研究codeは変更せず、開始時の未追跡 `auth.go`、`auth_test.go`、`models_test.go` はSHA-256一致を確認しました。

Local browser review verified both languages' three stages, the 12 → 15 → 15 deadline, a single selected state, Enter/Space, Tab focus and the skip link. Reviewed wide 1280px, narrow 390px and both languages at 320px, with no horizontal overflow at the inspected widths. Verified both diagrams' rendering and site image loading. Research code was unchanged; SHA-256 hashes confirmed preservation of the initially untracked `auth.go`, `auth_test.go` and `models_test.go`.

本検証は文書と静的websiteの受入であり、研究評価ではありません。Go testの再実行、React実装／E2E、実model評価、外部URLの可用性、公開deployは今回の確認対象外です。新依存・modelの導入や取得はありません。L3以降の受入状態は変更しません。

This accepts documentation and a static website, not research results. Go tests were not rerun; React implementation/E2E, real-model evaluation, external URL availability and public deployment were not verified in this task. No dependencies or models were installed/downloaded. Acceptance of L3 and later lessons is unchanged.

### 架構整理と公開の受入 / Architecture revision and publication acceptance

2026-11-12の追加依頼で[構成](docs/architecture.md)、[契約](docs/contracts.md)、[UI](docs/ui.md)を再整理しました。Mantine選択を維持し、導入・version固定は行っていません。runtime呼出しとimport依存を分け、appが状態決定、adapterが原子保存、workerが派生storeを担当する設計です。既存Go codeの移動はせず、受入済みL1–L2を進めたとは扱いません。

The additional 2026-11-12 request refined [architecture](docs/architecture.md), [contracts](docs/contracts.md) and [UI](docs/ui.md). Mantine remains selected, without installation or version pinning. The design separates runtime calls from imports: app decides state, adapters commit atomically, and the worker owns derived stores. Existing Go code was not moved and this does not advance accepted L1–L2 work.

公開対象はこの整理の文書・website・SVGのみです。`cmd/bridge/auth.go`、`auth_test.go`、`models_test.go` はstage対象外とし、開始時のSHA-256と照合します。差分・local link・一般的credential pattern・JavaScript構文を確認後、履歴を置換せずmainへpushします。Pages workflowの成功と公開日英HTML・SVGを実際に確認するまで公開完了としません。

Publication includes only this revision's documentation, website and SVGs. Exclude `cmd/bridge/auth.go`, `auth_test.go` and `models_test.go` from staging and compare their SHA-256 hashes against the starting snapshot. After reviewing diffs, local links, common credential patterns and JavaScript syntax, push to main without rewriting history. Do not mark publication complete until the Pages workflow succeeds and public JA/EN HTML and SVG are verified.

公開後に主要link・三段階例・素材読込が壊れた場合はsiteだけの修正を優先します。復旧できない場合は公開site変更を取り消す新commitで戻し、本人の作業treeや認証codeをresetしません。DB移行・model実行・依存導入は今回の公開に含みません。

If primary links, the three-stage example or assets fail after publication, prefer a focused site fix. If recovery is not possible, restore the public site with a new reverting commit, without resetting the user's worktree or authentication code. This release includes no database migration, model execution or dependency installation.

公開証拠（2026-11-12）：[commit eee6aca](https://github.com/lorendw7/kyudai-llm-context-bridge/commit/eee6aca676f71c042c538a76edf97687b95d08aa) をmainへpushし、[Pages run 43155372738](https://github.com/lorendw7/kyudai-llm-context-bridge/actions/runs/43155372738) がsuccessで完了しました。[公開site](https://lorendw7.github.io/kyudai-llm-context-bridge/) の `index.html`、`ja.html`、`assets/style.css`、`assets/context-continuity.svg` はHTTP 200で、公開内容がrelease fileとbyte単位で一致しました。公開browserでも日英の12→15→15、単一選択状態、画像読込を確認し、そのsite tabのerror/warning logは空でした。

Publication evidence (2026-11-12): pushed [commit eee6aca](https://github.com/lorendw7/kyudai-llm-context-bridge/commit/eee6aca676f71c042c538a76edf97687b95d08aa) to main; [Pages run 43155372738](https://github.com/lorendw7/kyudai-llm-context-bridge/actions/runs/43155372738) completed successfully. The [public site](https://lorendw7.github.io/kyudai-llm-context-bridge/) returned HTTP 200 for index.html, ja.html, assets/style.css and assets/context-continuity.svg, with byte-for-byte equality to released files. Public-browser checks confirmed both languages' 12 → 15 → 15 sequence, single selected state and loaded images; that site's tab had no error/warning logs.

commit前の14-file allowlistとstage内容一致を検証し、commit treeに認証3ファイルがないことを確認しました。本人の3ファイルは開始時とSHA-256が一致し、未追跡のまま保持しました。最終の記録更新は文書のみで、site・研究codeを変更しません。

Verified the exact 14-file staging allowlist and its equality with reviewed content, and confirmed the three authentication files were absent from the commit tree. Their SHA-256 hashes matched the starting snapshot and they remained untracked. The final publication-record update changes documentation only, not site or research code.

## 工数 / Effort

週20時間で学びながら本人実装する場合の推定です。通常の授業は3–6時間、重い連携課題は分割します。S4はS3受入後に進めます。 / Estimates assume user implementation while learning, 20 hours/week. Lessons usually take 3–6 hours; heavier integration work may be split. S4 follows S3 acceptance.

| 段階 / Stage | 追加工数 / Incremental effort | 成果物 / Deliverable |
| --- | --- | --- |
| S1 | 12–20 h | 依存不要の切替契約 / Dependency-free switching contract |
| S2 | 16–28 h | 実モデル＋保存 / Real models and storage |
| S3 | 20–32 h | 記憶・隔離・削除・3方式 / Memory, isolation, deletion and policies |
| S4 | 16–28 h | React研究UI / React research UI |
| S5 | 16–28 h | 再現・更新可能な原型 / Reproducible, upgradeable prototype |
| S6 | 24–48 h | 最終評価と日英報告 / Final evaluation and bilingual report |

累積：実モデルAPI demo 28–48h、記憶原型48–80h、UI付き64–108h、報告まで104–184h（約6–10週間）。学習、互換性、許可、queue待ちで変動します。2授業・1週間の実績で見直し、納期保証にはしません。

Cumulative: real-model API demo 28–48h, memory prototype 48–80h, UI included 64–108h, report completed 104–184h (about 6–10 weeks). Learning, compatibility, authorization and queue waits affect duration. Re-estimate after two lessons/one week; these are not delivery guarantees.

## 評価・実行ゲート / Evaluation and execution gates

開発12／最終36ケース、3方式×3反復、24 turns・6 probes、324 conversation runsを出発案とします。品質90%以上・全履歴比5 percentage points以内の低下・回答入力token30%以上削減は評価前に確定する目標案です。[評価仕様](docs/evaluation.md)を参照します。

The starting proposal is 12 development/36 final cases, three policies × three repetitions, 24 turns/six probes and 324 conversation runs. Quality ≥90%, degradation ≤5 percentage points versus full history and ≥30% answer-input-token reduction are proposed targets to freeze before evaluation. See the [evaluation design](docs/evaluation.md).

玄界のaccount、GPU基盤のOSS適用範囲、配布物、実行権限を確認してからjobを扱います。初期案は1 GPU、smoke最大30分、開発pilot最大60分です。全評価枠はpilot実測から算出します。未実行検証を成功と記録しません。

Verify Genkai entitlement, the GPU platform's OSS boundary, artifacts and execution authorization before jobs. Initial caps are one GPU, 30 minutes for smoke and 60 minutes for a development pilot. Size final allocation from the pilot. Never record unexecuted checks as passed.
