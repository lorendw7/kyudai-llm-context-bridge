# 作業計画 / Plan

更新 / Updated: 2026-11-10 (Asia/Tokyo)。本書だけが作業・進捗の一覧です。他の文書は設計・学習資料です。 / This is the sole action and progress checklist. Other documents contain design and learning references.

## 方針 / Decisions

- A案（Go + Ollama + Mem0 OSS）は2026-11-06に利用者が選択。教員承認、玄界account・計算枠は利用者の報告です。 / Route A was selected by the user on 2026-11-06. Supervisor approval and Genkai account/allocation are user-reported.
- 2026-11-08：小型の自作研究UIにReact + TypeScript + Vite + Mantineを選択。独自CSSは初期範囲に含めない。SQLite正本保存を提案。正確なpackage versionsは未固定。 / React + TypeScript + Vite + Mantine were selected for the small owned research UI. Handwritten CSS is outside the initial scope. SQLite authoritative storage is proposed. Exact package versions are not frozen.
- 利用者が全実装コード・設定・script・移行・runnerを書く。assistantはテスト、合成fixture、test-only fake、指導・レビュー、依頼された日英文書を担当。 / The user writes all implementation code, configuration, scripts, migrations and runner. The assistant owns tests, synthetic fixtures, test-only fakes, teaching/review and requested bilingual documentation.
- OSSとローカル実行が必須。導入・download・有料API・遠隔実行は許可範囲に従う。要件・提案・vendor claim・実測を区別。 / OSS and a local route are required. Installation, downloads, paid APIs and remote runs follow authorized scope. Separate requirements, proposals, vendor claims and measurements.
- 公開は現在資料のみ。履歴置換は今回の作業ではない。発表archiveはignored `presentations/` に保持。 / Publish current material only. History replacement is not part of this task. Presentation archives remain under ignored `presentations/`.

## 唯一の作業一覧 / Sole action checklist

- [x] S0 — 日英文書を案内・構成・契約・UI・評価・資料台帳に整理し、実装済み／提案を区別。 / Organize paired bilingual overview, architecture, contracts, UI, evaluation and source register; distinguish implementation from proposals.
- [x] 既存scaffold：health、Message、2つのGoテストを確認。 / Existing scaffold: health, Message and two Go tests checked.
- [ ] S1 — 契約とfake切替（L1–L4）：session、所有者、A → B → A、重複・timeoutを依存追加なしで説明・検証。 / Contracts and fake switching: explain and verify sessions, ownership, A → B → A, duplicates and timeout without new dependencies.
- [ ] S2 — 実モデルと正本保存（L5–L8）：依存・artifact監査、許可後のOllama、SQLite、再起動・再試行。受入は同じsessionの実モデル切替と保存応答の回復。 / Real models and storage: dependency/artifact audit, authorized Ollama and SQLite, restart/retry; accept real switching in one session and saved-response recovery.
- [ ] S3 — 記憶と比較（L9–L12）：Python契約、Mem0、訂正／出典／隔離、部分失敗回復、export/delete、3方式。受入は外向き通信遮断下の訂正保持、隔離、削除検証。 / Memory and controls: Python contract, Mem0, corrections/provenance/isolation, partial-failure recovery, export/delete and three policies; accept correction retention, isolation and deletion checks with outbound traffic blocked.
- [ ] S4 — React研究UI（L13–L16）：本人がUI実装、assistantが画面・API・E2Eテスト。受入は日英、keyboard、狭幅、再読込、二重送信、失敗表示、削除後cache除去。 / React UI: user implementation, assistant UI/API/E2E tests; accept JA/EN, keyboard/narrow layout, reload, duplicate submit, failures and deletion-cache removal.
- [ ] S5 — 再現性と更新（L17–L20）：runner、stage計測、DB移行、API互換性、依存更新、rollback、開発pilot。受入は再現manifestと旧版からの復元試験。 / Reproduction/upgrades: runner, stage metrics, DB migrations, API compatibility, dependency upgrades, rollback and pilot; accept a reproduction manifest and restore test from the previous version.
- [ ] S6 — 玄界・最終評価・報告（L21–L24）：許可済み環境調査とbounded job、評価固定、最終実行、日英報告。欠測・失敗・未達を保持。 / Genkai/evaluation/report: authorized audit and bounded job, freeze, final execution and bilingual report; retain missing data, failures and unmet targets.
- [ ] 公開資料同期：本人が必要なwebsite実装変更を行い、assistantが文言・リンク・表示を検証。 / Public-site synchronization: the user makes required website implementation changes; the assistant verifies wording, links and presentation.
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
