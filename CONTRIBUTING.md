# 学習と貢献 / Learning and contributing

L3の続きは[Bの回答とturn状態の演習](#l4-state-lesson)から始めます。編集前に[協働規則](AGENTS.md)と[作業計画](PLAN.md)を確認します。設計資料は[READMEの案内](README.md)から参照できます。

To resume after L3, start with the [exercise on B’s reply and turn state](#l4-state-lesson). Before editing, read the [collaboration guidelines](AGENTS.md) and [plan](PLAN.md). The [README](README.md) links to the design references.

## 分担 / Ownership

利用者は研究の核心codeを書きます。対象は状態・並行実行・保存の原子操作・回復・context policy・記憶の調停・UIの状態処理・再生と集計など、設計判断を含む部分です。assistantはその判断をcodeへ写す外周を書きます。対象はHTTP handlerとDTO変換、error envelope、起動と停止、設定読込、migration実行の枠組み、workerのHTTP外殻、UIの画面構成と文言、runnerのCLI、job wrapperです。assistantはこのほかtests、synthetic fixtures、test-only fakes、verification harnessesを実装し、指導・code review・依頼された日英文書を担当します。testが失敗したら、そのcodeの担当者が直します。本人のcodeならassistantが原因と修正方針を説明し、本人が直します。依存の選定と導入許可は本人が行います。2026-11-19に利用者が決めた分担です。

The user writes the core research code: the parts that carry design decisions, such as state, concurrency, atomic storage operations, recovery, context policies, memory reconciliation, UI state handling, replay and aggregation. The assistant writes the surrounding code that translates those decisions: HTTP handlers and DTO mapping, the error envelope, startup and shutdown, configuration loading, the migration runner, the worker HTTP shell, UI layout and copy, the runner CLI and job wrappers. The assistant also implements tests, synthetic fixtures, test-only fakes and verification harnesses, and provides teaching, code review and requested bilingual docs. When tests fail, the author of that code fixes it. For user-written code, the assistant explains the cause and repair direction and the user makes the fix. The user selects dependencies and approves installation. The user set this division on 2026-11-19.

公開 `website/` と美術素材は2026-11-12からassistantの担当で、HTML・CSS・JavaScript、説明図、表示検証を含みます。研究UI `web/` はMantineを使い独自CSSを書かない方針を維持します。

The assistant has owned public `website/` and visual assets since 2026-11-12, including HTML, CSS, JavaScript, explanatory graphics and display verification. Research UI `web/` keeps Mantine without handwritten CSS.

<a id="lesson-method"></a>

## 授業の進め方 / Lesson method

各回は一つの問いから始めます。関連するsourceまたは合成fixtureを読み、本人が実行前に結果と理由を予測します。assistantは一度に一つのhintと設計理由を伝えます。assistantが外周codeを書く回は、本人が先に動作を予測してからそのcodeを読み、条件を一つ変えた場合を説明します。確認後に条件を一つ変え、本人が同じ原則で説明できるかを確かめます。その回に[未決の設計事項](docs/architecture.md#open-questions)があれば、本人が選択肢を比べて決めてから実装します。段階ごとの実験は既存codeを一時的に変えて観測するもので、終わったら元に戻します。

Start each session with one question. Read the relevant source or synthetic fixture, and have the user predict the result and explain why before running it. The assistant gives one hint at a time and explains the design reason. In lessons where the assistant writes surrounding code, the user predicts its behavior first, reads the code, then explains one changed condition. Then change one condition and ask the user to explain it using the same principle. When the lesson has an [open design question](docs/architecture.md#open-questions), the user compares options and decides before implementing. Stage experiments change existing code temporarily to observe behavior, then restore it.

機能を実装する回は、入出力・失敗・受入条件を先に決め、本人が既存codeを再利用して最小の変更を書きます。assistantが用意した契約testと観測結果でreviewします。契約が変わる場合は期待動作を先に確認し、testはその契約を反映します。

For an implementation session, agree on inputs, outputs, failure behavior and acceptance first. The user writes the smallest change that reuses existing code, then reviews it with assistant-authored contract tests and observations. If the contract changes, clarify the expected behavior first and update tests to reflect that contract.

記録は対象段階、説明できた設計判断、変更した契約、実行testと結果、残る疑問に絞ります。進捗と次回の演習は[PLAN.mdの学習ルート](PLAN.md#teaching-route)で管理します。授業中の説明は中国語／英語、正式資料とUIは日英です。

Keep the record to the stage, a design decision the user can explain, any contract change, checks run and their results, and open questions. Track progress and the next exercise in the [learning route in PLAN.md](PLAN.md#teaching-route). Teaching uses Chinese/English in chat. Official material and the UI use Japanese/English.

<a id="l4-state-lesson"></a>

### L4の一問：Bの回答をいつ「保存済み」と表示できるか / L4 exercise · When can B’s reply appear as saved?

Bへ送った締切の訂正を例にします。サーバーが受理した直後の状態は `accepted` です。`answer_saved` は契約にある名前です。それだけで画面に「回答保存済み」と表示できるでしょうか。先に理由を一文で答えてください。説明できれば下の答え合わせへ、迷ったら既存のテストへ進みます。

Take the deadline correction sent to B. The server has just accepted it. The turn is `accepted`, and `answer_saved` is a name in the contract. Is that enough to display “Answer saved”? Give a one-sentence reason. Go straight to the answer if you can explain it. If unsure, try the existing test.

**手掛かり / Hint.** 既存テストの `unknown role` は `MessageRole("system")` を作って `NewMessage` に渡します。この値をGoで作れても、コンストラクタは受理するでしょうか。予測してからリポジトリ直下で実行します。[テスト](cmd/bridge/domain_test.go)・[コンストラクタ](cmd/bridge/domain.go)

The existing `unknown role` test constructs `MessageRole("system")` and passes it to `NewMessage`. Go permits constructing the value. Predict whether the constructor accepts it, then run this from the repository root. [Test](cmd/bridge/domain_test.go) · [Constructor](cmd/bridge/domain.go)

```sh
go test ./cmd/bridge -run '^TestNewMessageRejectsInvalidInput$/^unknown_role$' -v -count=1
```

表示される `unknown_role` の PASS は、未知のroleに対してコンストラクタがerrorを返したことを示します。 / A PASS on the named `unknown_role` subtest means the constructor returned an error for that input.

値の判定をturnにも応用します。`" generating "` と `memory_pending` のうち契約にある名前はどちらでしょうか。その名前は、全履歴でBが正常に回答する時にも現れるでしょうか。

Apply the same check to a turn. Which of `" generating "` and `memory_pending` is a valid state name? Would you see that state on B’s successful full-history turn?

<details>
<summary>予測した後の答え合わせ / Check after predicting</summary>

状態名の検証は「この値を記録できるか」、遷移の検証は「今、この変更を確定できるか」を問います。全履歴turnは `accepted → generating → answer_saved → complete` と進みます。Bの回答とモデル・usageを保存し、サーバーが `answer_saved` またはその後の `complete` を返して初めて、画面に保存済みと表示できます。

Value validation asks whether the value belongs to the contract. Transition validation asks whether this change can happen now. A full-history turn follows `accepted → generating → answer_saved → complete`. The UI can show a saved reply once the server has committed B’s answer, model and usage, and reported `answer_saved` or the subsequent `complete`.

`memory_pending` は8つの正本値の一つですが、派生記憶を更新しない全履歴の正常経路では使いません。`" generating "` は空白を含むため、値の検証で拒否します。この2例を説明できれば今回の読解は終了です。[状態の契約](docs/contracts.md#turn-states)

`memory_pending` is one of the eight valid names. A successful full-history turn skips it because it has no derived memory to update. The string `" generating "` fails value validation because its spaces are significant. Once you can explain both, the reading exercise is complete. [Turn-state contract](docs/contracts.md#turn-states)

</details>

本人が書く最初のコードと、その受入条件は[PLANのL4](PLAN.md#l4-entry)にあります。 / Find the first user-written change and its acceptance criteria in the [L4 plan](PLAN.md#l4-entry).

<a id="handoff-acceptance"></a>

### Aに戻る前に何が確定したか / What was committed before switching back to A?

合成の全履歴sessionを一つ使います。要求 `r1` はAに「締切は12日」、`r2` はBに「15日に訂正」、`r3` はAに「締切はいつ？」。まずBの回答保存の直前で止め、次に保存直後で止めたとき、Aへの入力に何が残るか予測します。以下は予定の受入条件です。

Use one synthetic full-history session. In request `r1`, tell A “The deadline is the 12th.” In `r2`, tell B “Correction: the 15th.” In `r3`, ask A “When is the deadline?” Predict what reaches A if execution stops just before B’s answer is saved, then just after. These are proposed acceptance cases.

| 段階・操作 / Stage and action | 観測できたら合格 / Pass when observed |
| --- | --- |
| S1 · 3回とも成功、同じ時刻 / All three succeed with tied timestamps | fakeの呼出し順はA,B,A。最後の入力は `[user1,answer1,user2,answer2,user3]` とID・順序・件数が一致。各要素は同じsession由来。別sessionに置いた目印は含まれない / Fake calls are A,B,A. The final input matches exactly those five message IDs in that order, once each, from this session. A marker placed in another session is absent |
| S1 · Bの生成中に `r2` と `r3` を送る / Replay r2 and submit r3 during B | 同内容のr2は同じturn、同IDでmodelをAに変えると409 conflict、r3は409 busy。Bの呼出し1回、訂正1件、r3の記録0件。B完了後、明示送信したr3が一度通る / Identical r2 returns its turn, changing its model to A yields 409 conflict, and r3 yields 409 busy. Count one B call, one correction and zero r3 records. Explicit r3 after B completes is admitted once |
| S1 · 保存前にBがtimeout／取消 / Timeout or cancel B before save | r2はそれぞれ `failed`／`cancelled`、回答0件。同じr2の再送で追加生成0回。直後にAへ尋ねた場合の入力は `[user1,answer1,user3]`。別の試行で訂正を新IDで再送し完了させると、その成功対だけがAの入力に加わる / r2 is failed or cancelled respectively, with zero answers. Replay adds zero calls. Asking A immediately gives exactly those three inputs. In a separate run, explicitly retry the correction with a fresh ID and complete it. Only that successful pair is added to A’s input |
| S1 · 遅いBが新turnの途中で戻る / Late B returns during the next active turn | r2の回答は0件のまま、新turnの状態は変わらず、さらに別のIDは409 busy。旧呼出しの終了で新しい実行枠を空けない / r2 still has zero answers, the new turn keeps its state, and a third distinct request receives 409 busy |
| S1 · B保存後にHTTP応答が消失 / Lose the HTTP response after B is saved | r2の照会・同内容再送で同じ回答ID、Bへの追加呼出し0回。r3には訂正対が一度だけ入る / Lookup and identical replay return the same answer ID with zero extra B calls. r3 receives the correction pair exactly once |
| S2 · Bの途中でprocess再起動 / Restart the process during B | 未送出の `accepted` なら回復後の生成1回。保存回答のない `generating` なら自動生成0回、照合中はr3がbusy。所有者が未確定のr2を回答なしで閉じたら新IDを受理でき、r2への遅い保存は失敗する。`answer_saved` なら元の回答IDで完了、追加生成0回 / Resume an undispatched accepted turn with one call. For generating with no committed answer, make zero automatic calls and keep r3 busy during reconciliation. After owner-confirmed closure, admit a fresh ID and reject a late save to r2. For answer_saved, complete with the original answer ID and zero new calls |
| S3 · 同じ例をMem0用sessionで再生し、Bの回答保存後に記憶更新だけ失敗 / Replay the example in a Mem0 session and fail the memory update after B’s answer is saved | 保存回答を表示し、`memory_pending` の間はr3がbusy。元の更新を照合・回復した後、この1事実のfixtureでは有効な訂正が1件、Bへの追加生成0回。Aに渡す訂正の出典IDは `user2` / Show the saved answer and keep r3 busy while memory_pending. Reconcile the original update. For this one-fact fixture, expect one current correction, zero additional B generations and message ID `user2` as the source of the correction supplied to A |
| S4 · 取消確認後に古い照会結果が到着 / Deliver an old poll after confirmed cancellation | Bの取消済み表示と次のモデル選択を保持。再読込を伴う場合も生成を自動送信しない / B stays cancelled, the composer retains the next-model selection, and reloading adds zero submissions |

各行は同じ初期状態から独立に試します。記録するのは要求ID・呼出し先モデル、fakeに渡したmessage ID列、保存回答数、最終状態です。最後の回答がたまたま「15日」でも、入力列に訂正元がなければ履歴の受入は失敗です。[turn契約](docs/contracts.md#l4-turn-design)・[再起動時の判定](docs/contracts.md#restart-cut-points)

Run each row independently from the same starting conditions. Record request IDs, called model, ordered input message IDs, saved-answer count and final state. Even if the final answer happens to say “the 15th”, missing correction history fails the input check. [Turn contract](docs/contracts.md#l4-turn-design) · [Restart decisions](docs/contracts.md#restart-cut-points)

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
