# 学習と貢献 / Learning and contributing

[AGENTS.md](AGENTS.md)、[README.md](README.md)、[PLAN.md](PLAN.md)を先に読みます。作業と進捗はPLAN.mdだけで管理し、設計の正本は[docs/architecture.md](docs/architecture.md)、契約は[docs/contracts.md](docs/contracts.md)です。

Read AGENTS.md, README.md and PLAN.md first. Track work/progress only in PLAN.md; architecture and contract references are in docs/architecture.md and docs/contracts.md.

## 分担 / Ownership

利用者は研究用の実装codeを書きます。backend、worker、React、部品library設定、dependency manifest、schema、migration、起動script、runnerを含みます。assistantはtests、synthetic fixtures、test-only fakes、verification harnessesを実装し、指導・code review・依頼された日英文書を担当します。testが失敗したら原因と修正方針を説明し、本人が研究用production codeを直します。

The user writes research implementation code: backend, worker, React, component-library configuration, dependency manifests, schemas, migrations, startup scripts and runner. The assistant implements tests, synthetic fixtures, test-only fakes and verification harnesses, and provides teaching, code review and requested bilingual docs. When tests fail, explain the cause and repair direction; the user fixes research production code.

2026-11-12の例外：利用者が公開 `website/` と美術素材をassistantへ再割当しました。HTML・CSS・JavaScript、説明図、表示検証を含みます。研究UI `web/` の本人実装とMantine・独自CSSなしの方針は維持します。

Exception dated 2026-11-12: the user reassigned public `website/` and visual assets to the assistant, including HTML, CSS, JavaScript, explanatory graphics and display verification. Research UI `web/` remains user-written with Mantine and no handwritten CSS.

<a id="lesson-method"></a>

## 授業の進め方 / Lesson method

各回は一つの問いから始めます。関連するsourceまたは合成fixtureを読み、本人が実行前に結果と理由を予測します。assistantは一度に一つのhintと設計理由を伝えます。確認後に条件を一つ変え、本人が同じ原則で説明できるかを確かめます。

Start each session with one question. Read the relevant source or synthetic fixture, and have the user predict the result and explain why before running it. The assistant gives one hint at a time and explains the design reason. Then change one condition and ask the user to explain it using the same principle.

観測で検証できる内容は既存test、trace、手計算で扱います。たとえば[session storeの並行test](cmd/bridge/session_store_test.go)を読み、同じIDの並行要求で成功が1件になる理由を説明します。次に同じrequest IDで本文だけを変えたときの結果を予測し、L4の冪等性に応用します。

Use existing tests, traces or manual calculations when observation can answer the question. For example, read the [session-store concurrency test](cmd/bridge/session_store_test.go) and explain why one request wins when concurrent inserts use the same ID. Then apply that reasoning to L4 idempotency by predicting what happens when a request ID is reused with changed text.

機能を実装する回は、入出力・失敗・受入条件を先に決め、本人が既存codeを再利用して最小の変更を書きます。assistantが用意した契約testと観測結果でreviewします。契約が変わる場合は期待動作を先に確認し、testはその契約を反映します。

For an implementation session, agree on inputs, outputs, failure behavior and acceptance first. The user writes the smallest change that reuses existing code, then reviews it with assistant-authored contract tests and observations. If the contract changes, clarify the expected behavior first and update tests to reflect that contract.

記録は対象段階、説明できた設計判断、変更した契約、実行testと結果、残る疑問に絞ります。進捗と次回の演習は[PLAN.mdの学習ルート](PLAN.md#teaching-route)で管理します。授業中の説明は中国語／英語、正式資料とUIは日英です。

Keep the record to the stage, a design decision the user can explain, any contract change, checks run and their results, and open questions. Track progress and the next exercise in the [learning route in PLAN.md](PLAN.md#teaching-route). Teaching uses Chinese/English in chat. Official material and the UI use Japanese/English.

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

通常の公開はfocused branchとdiff確認を使い、明示されたmain更新・履歴置換の指示がある場合はそれに従います。無関係の変更を保持し、実行したcheckだけを報告します。現在のwebsiteはbuild不要の説明用siteで、研究UIとは別です。2026-11-12以降はassistantがwebsiteのcode・素材・文言・link・keyboard・狭幅／広幅・合成3段階を担当します。

For ordinary publication, use a focused branch and review diffs; explicit main-update/history-replacement instructions take precedence. Preserve unrelated changes and report only executed checks. The existing buildless public website is separate from the research UI. Since 2026-11-12, the assistant owns website code, assets, wording, links, keyboard/narrow/wide behavior and the three synthetic stages.

private chats、credential、weights、DB、raw output、発表archiveをGitへ入れません。ignoreだけに頼らずstage内容を確認します。projectの原著作物は[MIT](LICENSE)、第三者資料は元licenseに従い、大学公認を示唆しません。

Keep private chats, credentials, weights, databases, raw output and presentation archives out of Git. Review staged content rather than relying on ignore rules. Original project material uses MIT; preserve third-party terms and do not imply university endorsement.
