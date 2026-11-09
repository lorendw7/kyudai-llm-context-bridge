# 評価・玄界・報告 / Evaluation, Genkai and reporting

更新 / Updated: 2026-11-08。すべて未実行の評価設計です。数値は固定前の提案で、project測定結果ではありません。実行の順序と進捗は[PLAN.md](../PLAN.md)にあります。

Updated 2026-11-08. This is an unexecuted evaluation design. Numbers are proposals pending freeze, not project measurements. Execution order and progress are in [PLAN.md](../PLAN.md).

## 研究上の問い / Research question

モデルA → B → Aで作業を続けるとき、全履歴・逐次要約・Mem0のどれが訂正・出典・制約を保持でき、どの追加負担を伴うかを調べます。Mem0はmemory framework、Ollamaはruntime、Qwenはmodel familyであり、同じ種類の比較対象ではありません。

When continuing work across A → B → A, which of full history, rolling summary and Mem0 retains corrections, provenance and constraints, and at what additional cost? Mem0 is a memory framework, Ollama a runtime and Qwen a model family; they are not interchangeable comparison categories.

## 比較条件 / Comparison conditions

| Policy | 固定する処理 / Fixed behavior |
| --- | --- |
| Full history / 全履歴 | system + committed source history + current user turn。暗黙truncateなし / System, committed source history and current user turn; no silent truncation |
| Rolling summary / 逐次要約 | versioned summary + recent messages。更新timingとpromptを固定 / Versioned summary plus recent messages; fixed update timing/prompt |
| Mem0 | scoped retrieved memories + recent messages。検索数、予算、extract/update設定を固定 / Scoped retrieved memory plus recent messages; fixed retrieval count, budget and extraction/update config |

answer modelの切替順、model/tokenizer/template、helper/embedding model、generation設定、case、seed計画を揃えます。policyをsession内で切り替えず、独立namespaceで比較します。条件順は事前固定のbalanced orderにし、load/warm-upを記録します。

Match answer-model switch order, model/tokenizer/template, helper/embedder, generation configuration, cases and seed schedule. Compare policies in independent namespaces without changing policy within a session. Use a predeclared balanced condition order and record loading/warm-up.

主評価は固定user scriptでの実session会話です。各方式の生成回答はその方式の履歴に入るため、後続の入力は同一にはなりません。この会話全体の差を主評価し、入力完全一致による因果効果とは主張しません。同じsource prefixを使うcontext-only診断は別の補助実験とし、集計を混ぜません。

The primary evaluation uses real sessions with fixed user scripts. Generated answers enter each policy's own history, so later inputs are not identical. Evaluate the end-to-end conversation difference without claiming an effect under identical prompts. A context-only diagnostic using identical source prefixes is a separate secondary experiment and must not be pooled.

## ケースと採点 / Cases and scoring

6 family（recall、correction、temporal order、provenance、abstention、continuation）×日英2言語を想定します。開発は各組合せ1件で12件、最終は独立3件で36件です。翻訳や固有名詞置換だけで開発と最終を分けず、最終は異なる事実・経過・出典を持たせます。

Use six families (recall, correction, temporal order, provenance, abstention, continuation) × two languages (JA/EN). One development case per combination gives 12; three independently authored final cases give 36. Do not create the held-out split merely by translating or renaming development cases; use different facts, trajectories and sources.

各caseは24 user turns、8・16turn後にmodel切替、事前指定の6 probesを提案します。3方式×3反復×36casesで324 conversation runs、完遂時は1,944 probe responsesです。seedを固定してもruntimeの完全決定性は仮定しません。Qwen3-4B/8Bなら同一family内の切替であり、異familyへの一般化を主張しません。

Each proposed case has 24 user turns, switches after turns 8/16 and six predefined probes. Three policies × three repetitions × 36 cases yield 324 conversation runs and 1,944 probe responses if all complete. Fixed seeds do not imply deterministic runtime behavior. Qwen3-4B/8B represent within-family switching, not evidence of cross-family generalization.

rubricは期待事実、許容表現、必須source、古い事実、捏造、未回答の扱いをfreeze前に定義します。probeは事前の二値pass/failを基本とし、訂正前情報を現行として答えた場合は失敗です。人手評価ではpolicy名を隠し、一部を二重採点して不一致を記録します。未校正LLM judgeを唯一の採点者にしません。

Before freezing, define expected facts, acceptable wording, required sources, stale facts, fabrication and unanswered probes. Use predefined binary pass/fail per probe; treating superseded information as current is a failure. Blind human scorers to policy and double-score a subset, retaining disagreements. Do not use an uncalibrated LLM judge as the sole scorer.

## 予算・計測 / Budgets and measurement

出発案はcontext 16,384 tokens、出力予約512、recent 4 message pairs、summary上限1,024、memory最大5件かつ1,024 tokensです。system、template、special tokensも入力予算に含めます。exact countの手段を固定runtimeで確認できるまで最終評価を始めません。推定値をexactと表示しません。

Starting proposals: 16,384 context tokens, 512 reserved output tokens, four recent message pairs, summary limit 1,024, up to five memories within 1,024 tokens. Include system/template/special tokens in the budget. Do not begin final evaluation until exact counting is validated against the pinned runtime. Never label estimates as exact.

主cohortは全履歴がwindow内に収まるケースにし、overflowは別cohortで報告します。最終評価中のoverflowを都合よく除外せず、失敗・欠測として残します。事前countとruntimeの `prompt_eval_count` の不一致はtruncateやtemplate差を調査する理由です。

Keep the primary cohort within the full-history window and report overflow separately. Preserve unexpected final-evaluation overflow as failures/missing data rather than selectively excluding it. Disagreement between preflight counts and runtime prompt_eval_count triggers investigation of truncation/template differences.

| 計測 / Metric | 意味 / Meaning |
| --- | --- |
| Quality / 品質 | probe pass rate、family・言語・answer model別の失敗 / Pass rate and failures by family/language/answer model |
| Answer input / 回答入力 | runtimeが報告する回答用入力token。helper分は別 / Runtime-reported answer input tokens; helpers separate |
| Other tokens / 他token | summary、抽出、更新、embedding、再試行の入出力 / Summary, extraction, update, embedding and retry input/output |
| Time / 時間 | queue、model load、context、検索、生成、memory更新、全turn / Queue, load, context, search, generation, memory update and total turn |
| Resources / 資源 | RAM/VRAM peak、割当GPU-hours、実際の計上points / Peak RAM/VRAM, allocated GPU-hours and actual charged points |
| Reliability / 信頼性 | restart、重複、timeout、部分失敗、削除後検索 / Restart, duplicates, timeout, partial failure and post-deletion retrieval |

回答入力token削減率は `1 - sum(policy answer input) / sum(baseline answer input)` を同じcase・反復の対応集合で計算し、全履歴・要約の各baselineを明示します。memory overheadを含むtotal costと分けます。小さいcontextだけでは総費用削減を証明できません。

Compute answer-input reduction as 1 - sum(policy answer input) / sum(baseline answer input) on matched case/repetition sets, naming full-history or summary as the baseline. Separate it from total cost including memory overhead. Smaller context alone does not prove total savings.

品質は予定probeを分母とし、技術失敗で未到達のprobeも未成功として計上して欠測理由を併記します。token比較は共通完了集合と除外数を報告し、完了例だけの品質を主結果にしません。反復を独立caseと誤認せず、case単位のpaired uncertaintyと全失敗を報告します。

Quality uses scheduled probes as its denominator; probes not reached due to technical failure count as unsuccessful with missing reasons reported. Token comparisons report the common completed set and excluded counts; complete-case quality is not the primary result. Treat repeated runs as clustered by case and report paired case-level uncertainty and all failures.

目標案は品質≥90%、全履歴からの低下≤5 percentage points、回答入力token削減≥30%です。隔離／削除は独立した必須suiteで、試験した範囲で未許可操作・削除後検索が0件であることを要求します。有限のtest成功を普遍的安全性と表現しません。

Proposed targets are quality ≥90%, degradation ≤5 percentage points versus full history and answer-input reduction ≥30%. Isolation/deletion is a separate mandatory suite requiring zero unauthorized operations/post-deletion retrievals in tested cases. Passing finite tests is not universal security assurance.

## 玄界 / Genkai

許可後にaccount、modules、driver、private storage、GPU基盤のOSS境界、network/job隔離を確認します。公式資料はPJMとSingularityを説明しています。実際の配布方法を環境に合わせ、runtime install/downloadを避けます。[PJM](https://www.cc.kyushu-u.ac.jp/scp/usage/job/)、[Singularity](https://www.cc.kyushu-u.ac.jp/scp/system/genkai/software/Singularity/)（2026-11-08参照）。

After authorization, verify entitlement, modules, drivers, private storage, the GPU platform's OSS boundary and network/job isolation. Official documentation describes PJM and Singularity. Match packaging to the actual environment and avoid runtime installs/downloads. [PJM](https://www.cc.kyushu-u.ac.jp/scp/usage/job/), [Singularity](https://www.cc.kyushu-u.ac.jp/scp/system/genkai/software/Singularity/) (accessed 2026-11-08).

Go、worker、Ollama、runnerを一つの割当job内で実行し、研究UIは不要です。最初は1 GPU、smoke 30分・開発pilot 60分を上限案とし、本評価枠は実測から決めます。loopbackだけでは同一hostの他jobを隔離できず、raw Ollama endpointの保護も必要です。

Run Go, worker, Ollama and runner in one allocated job; the UI is unnecessary. Proposed initial caps are one GPU, 30-minute smoke and 60-minute development pilot; size final allocation from measurements. Loopback alone does not isolate other jobs on a host, and raw Ollama endpoints also require protection.

node-local DBは整合したbackupをprivate永続storageへ保存します。稼働中DB fileの単純copyをcheckpointにせず、job終了時のlocal storage消失を前提に途中保存します。job wrapperはreadiness、timeout、失敗記録、自分の子processだけの終了を管理します。

Save consistent backups of node-local databases to private persistent storage. Do not checkpoint by copying a live DB file blindly; save progress before local storage disappears at job end. The wrapper manages readiness, timeouts, failure records and termination of its own children only.

## 公開する証拠 / Reporting evidence

freezeしたmanifest、case/rubric hash、software/model version、実行条件、集計結果、失敗・未達、再現手順を日英で報告します。私的raw output、weights、DB、credentialは公開しません。引用benchmarkはvendor/published claimと表示し、project測定と分離します。

Report the frozen manifest, case/rubric hashes, software/model versions, execution conditions, aggregates, failures/unmet targets and reproduction instructions in JA/EN. Do not publish private raw output, weights, databases or credentials. Label cited benchmarks as vendor/published claims, separate from project measurements.
