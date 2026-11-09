# 学習と貢献 / Learning and contributing

[AGENTS.md](AGENTS.md)、[README.md](README.md)、[PLAN.md](PLAN.md)を先に読みます。作業と進捗はPLAN.mdだけで管理し、設計の正本は[docs/architecture.md](docs/architecture.md)、契約は[docs/contracts.md](docs/contracts.md)です。

Read AGENTS.md, README.md and PLAN.md first. Track work/progress only in PLAN.md; architecture and contract references are in docs/architecture.md and docs/contracts.md.

## 分担 / Ownership

利用者は全実装codeを書きます。backend、worker、React、部品library設定、dependency manifest、schema、migration、起動script、runnerを含みます。assistantはtests、synthetic fixtures、test-only fakes、verification harnessesを実装し、指導・code review・依頼された日英文書を担当します。testが失敗したら原因と修正方針を説明し、本人がproduction codeを直します。

The user writes all implementation code: backend, worker, React, component-library configuration, dependency manifests, schemas, migrations, startup scripts and runner. The assistant implements tests, synthetic fixtures, test-only fakes and verification harnesses, and provides teaching, code review and requested bilingual docs. When tests fail, explain the cause and repair direction; the user fixes production code.

## 一回の授業 / One lesson

一つの小さい成果物を扱います。目的とdata flowを説明し、入出力、失敗、完了条件を先に定義します。本人が説明を言い直せることを確認して実装し、assistantの契約testと観測結果でreviewします。完成実装の配布ではなく、hintと設計理由を使います。

Work on one small deliverable. Explain purpose/data flow and define input, output, failure and completion conditions first. The user should restate the design, implement it and review it against assistant-authored contract tests and observations. Teach with hints and design reasons rather than complete implementations.

毎回の記録は「対象S/L、変更契約、本人実装、実行testと結果、残る疑問」です。完了記号や次の作業はPLAN.mdだけに記録します。契約が変わった場合は本人と期待動作を明確にしてからtestを変更し、単に現実装に合わせてtestを弱めません。

Each lesson records the S/L stage, changed contract, user implementation, executed tests/results and open questions. Completion markers and next actions belong only in PLAN.md. Clarify intended behavior before changing tests for a revised contract; do not weaken tests just to match implementation.

授業はdomain契約からinterface、API、切替へ進みます。現在の完了状況と次の授業は[PLAN.md](PLAN.md)だけで管理します。UI授業はAPI・保存・記憶の受入後に行います。授業中の説明は中国語／英語、正式資料とUIは日英です。

Lessons proceed from domain contracts to interfaces, APIs and switching. Track completion and the next lesson only in [PLAN.md](PLAN.md). UI lessons follow acceptance of the API, storage and memory stages. Teaching uses Chinese/English in chat; official material and UI use Japanese/English.

## 検証の分け方 / Verification levels

| 種類 / Level | 対象 / Scope | 証明しないもの / Does not establish |
| --- | --- | --- |
| Unit / 単体 | 純粋domain、policy budget、状態遷移 / Pure domain, policy budget and state transitions | 実model品質 / Real-model quality |
| Contract / 契約 | Go HTTP、worker、UI DTO、fake外部応答 / Go HTTP, worker, UI DTOs and fake external responses | vendor互換性の実測 / Measured vendor compatibility |
| Integration / 連携 | 許可済みSQLite、Ollama、Mem0、遮断network / Authorized SQLite/Ollama/Mem0 and blocked outbound network | 任意環境での再現 / Reproduction on every environment |
| UI / 画面 | component、利用者操作、日英、focus、IME、失敗 / Components, user actions, JA/EN, focus, IME and failures | backend認可の正しさ / Correct backend authorization |
| E2E | login相当、同session切替、訂正、回復、削除 / Authentication, same-session switching, correction, recovery and deletion | 最終研究結果 / Final research results |
| Evaluation / 研究評価 | frozen cases、rubric、cost accounting / Frozen cases, rubric and cost accounting | 未試験集団への一般化 / Generalization to untested populations |

失敗・未実行・skipを明記します。test依存とbrowser binaryも導入許可の対象です。開発fixtureと最終評価caseを混ぜません。state classやcomponent内部構造より、利用者に見える動作と契約を検証します。

Report failures, unexecuted checks and skips. Test dependencies/browser binaries also require installation authorization. Keep development fixtures separate from final cases. Test observable behavior/contracts instead of private component structure or styling classes.

## 現在使える確認 / Available checks

```sh
python scripts/check_repository.py
node --check website/assets/demo.js
gofmt -l cmd/bridge
go test ./... -v -count=1
```

現在のGo scaffoldにmodelやmemory packageは不要です。上記はReact typecheck/build/E2Eを実行しません。React projectはまだ存在せず、将来のchecksは本人の実装と依存許可後に追加します。

The current Go scaffold requires no model/memory packages. These commands do not run React type checks, builds or E2E. The React project does not exist yet; add those checks after user implementation and dependency authorization.

## 文書・公開 / Documentation and publication

日本語と英語を節単位で対応させ、片方だけの要約にしません。識別子とcommandsは共通、standard LICENSE原文は保持します。architecture、contract、UI、evaluation、sourceの説明を各正本文書に置き、READMEは案内にします。別のroadmapやhistorical report集を増やしません。

Pair Japanese and English by section, rather than providing only a shortened translation. Identifiers/commands are shared; preserve the standard LICENSE text. Keep architecture, contracts, UI, evaluation and sources in their reference documents; README is the entry point. Do not add parallel roadmaps or historical-report collections.

通常の公開はfocused branchとdiff確認を使い、明示されたmain更新・履歴置換の指示がある場合はそれに従います。無関係の変更を保持し、実行したcheckだけを報告します。現在のwebsiteはbuild不要の説明用siteで、研究UIとは別です。code変更は本人、assistantは依頼範囲の文言・link・keyboard・狭幅／広幅・合成3段階を検証します。

For ordinary publication, use a focused branch and review diffs; explicit main-update/history-replacement instructions take precedence. Preserve unrelated changes and report only executed checks. The existing buildless public website is separate from the research UI. The user owns code changes; the assistant verifies requested wording, links, keyboard/narrow/wide behavior and the three synthetic stages.

private chats、credential、weights、DB、raw output、発表archiveをGitへ入れません。ignoreだけに頼らずstage内容を確認します。projectの原著作物は[MIT](LICENSE)、第三者資料は元licenseに従い、大学公認を示唆しません。

Keep private chats, credentials, weights, databases, raw output and presentation archives out of Git. Review staged content rather than relying on ignore rules. Original project material uses MIT; preserve third-party terms and do not imply university endorsement.
