# Learning and contributing

Read [AGENTS.md](AGENTS.md), [README.md](README.md) and [PLAN.md](PLAN.md). PLAN.md is the sole action checklist. Keep current design in README; do not add parallel roadmaps or historical report collections.

## Teaching sequence

Route A is selected: **Go bridge + Ollama + Mem0 OSS**. The user writes all functional code, including Go, the Python memory worker, operational scripts, and the experiment runner. The assistant teaches in Chinese and English in chat, sets a small exercise, explains acceptance checks and reviews the user's solution. The assistant writes and maintains test code, synthetic fixtures, test-only fakes, and verification harnesses. It must not implement or rewrite functional code unless the user explicitly changes this arrangement. A failing test leads to an explanation and a user-written fix, not an assistant-written production patch. Official documentation and website material use Japanese and English.

Each lesson follows: concept explanation → one coding exercise → assistant-authored tests and observable verification → review. Agree on expected behavior before implementation; write tests against that contract rather than copying implementation details. Start with hints and reasons; avoid handing over a complete functional implementation. Move on when the acceptance checks pass and the user can explain both behavior and design. Report skipped or blocked integration tests explicitly. Installation, model downloads, and remote runs still require their authorized scope.

The next exercise is to describe `Session`, `Message` and `ModelClient`: what each owns, its inputs/outputs, and what remains unchanged when the answer model changes. The existing health endpoint is the starting point. Later topics are ownership checks, fake adapters, Ollama, scoped Mem0, persistence/recovery and fixed evaluation. Track completion and effort only in PLAN.md.

本人が全機能コードを実装し、assistantはテストコード・合成fixtureを担当し、会話内で中国語・英語による説明・課題・レビューを提供します。正式文書は日英です。一度に一つの成果物を扱い、受入確認に合格し本人が動作と理由を説明できてから次へ進みます。授業数・受入基準・進捗はPLAN.mdだけで管理します。

## Local checks

```sh
python scripts/check_repository.py
node --check website/assets/demo.js
gofmt -l cmd/bridge
go test ./... -v -count=1
```

No model or memory package is required for the current Go scaffold. Do not install dependencies or download models as part of routine validation. Exact versions, licenses and execution scope must be resolved before integration.

```sh
python -m http.server 8000 --bind 127.0.0.1 --directory website
```

For website edits, check Japanese/English, narrow/wide layouts, keyboard focus, navigation and all three synthetic demonstration stages. The site is buildless. GitHub Pages publishes `website/` from `main`.

## Repository workflow

Use a focused branch and review the diff before ordinary publication; an explicitly authorized main update or history replacement takes precedence. Preserve unrelated work. Record sources, access dates, versions, assumptions and checks actually executed. Separate requirements, design, published claims and measurements. Freeze evaluation conditions before final testing and retain failures.

Keep presentation archives, generated artifacts, private conversations, credentials, model weights, databases and local history backups out of Git. `presentations/` is local-only. Ignore rules do not remove already-tracked files; inspect the staged tree. Contributions use the project's [MIT license](LICENSE); preserve third-party terms. Do not imply university endorsement.
