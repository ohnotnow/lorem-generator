# lorem-generator

A tiny Go CLI that generates **contextual lorem** text using OpenAI’s Responses API.

* **Single file, no deps** – just the Go standard library
* **Terminal-friendly** – prompt via args or STDIN

---

## Quick start

```bash
# 1) Download a release (recommended) OR build from source

# 2) Set your API key
export OPENAI_API_KEY=sk-your-key

# 3) Generate text
lorem Two paragraphs about the Edinburgh Festival 2025
```

> Windows (PowerShell): `setx OPENAI_API_KEY "sk-your-key"` then restart your shell.

---

## Installation

### Option A: Download a release

Grab the latest assets from **Releases**:

* `lorem-<version>-linux-amd64.tar.gz`
* `lorem-<version>-linux-arm64.tar.gz`
* `lorem-<version>-darwin-arm64.tar.gz` (Apple Silicon)
* `lorem-<version>-windows-amd64.zip`
* `lorem-<version>-windows-arm64.zip`

Unpack and place the binary on your `PATH`.

### Option B: Build from source

Requires Go ≥ 1.21.

```bash
git clone https://github.com/ohnotnow/lorem-generator
cd lorem-generator
make build
# or: go build -o dist/lorem ./cmd/lorem
./dist/lorem "Short blurb about our checkout onboarding"
```

---

## Usage

```text
lorem [flags] <prompt>

Flags:
  -model string       OpenAI model (default "gpt-5-mini")
  -effort string      Reasoning effort: minimal|medium|high (default "minimal")
  -verbosity string   Verbosity: low|medium|high (default "low")
```

### Examples

```bash
# Plain usage (args)
lorem Two paragraphs about the Edinburgh Festival 2025

# Pipe from STDIN
cat topic.txt | lorem

# Make it chattier (but still fast)
lorem -verbosity=medium "Feature announcement placeholder copy"

# Ask for deeper reasoning when you want richer structure
lorem -effort=high -verbosity=high "Outline the agenda for a 1-day developer offsite"
```

---

## What it does

This CLI POSTs a minimal JSON payload to `POST /v1/responses` and prints the `output_text` field:

* `model` – defaults to `gpt-5-mini`
* `text.verbosity` – defaults to `"low"`
* `reasoning.effort` – defaults to `"minimal"`

That’s it.

---

## Environment

* **OPENAI_API_KEY** (required): your API key
* **Proxy support**: standard env vars like `HTTPS_PROXY` are honoured if set

Request timeout: 60 seconds.

---

## Tips for better placeholder copy

* Be explicit about **length/shape**: e.g., “2 paragraphs, 3–5 sentences each.”
* Call out **tone** if you care: “newsy”, “marketing”, “straight-to-the-point”.
* For UK-style output, mention it: “Use British English spelling.”

Example:

```bash
lorem "Write 2 paragraphs (3–5 sentences each) using British English about the 2025 programme announcement. Tone: newsy."
```

---

## Exit codes

* `0` – success (text printed)
* `1` – runtime/API error
* `2` – usage error (e.g., missing prompt)

---

## CI/CD (Releases)

A GitHub Actions workflow is included that:

* triggers on tag pushes like `v1.0.0`
* cross-compiles for Linux (amd64/arm64), macOS (arm64), Windows (amd64/arm64)
* uploads archives to the tag’s Release

Publish a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Then check your repo’s **Releases** page for assets.

---

## Project layout

```
.
├── cmd/
│   └── lorem/
│       └── main.go        # single-file CLI (no external deps)
├── Makefile               # tiny build helpers
└── README.md
```

---

## Security & privacy

* Prompts are sent to OpenAI’s API. Avoid secrets and production PII.
* Consider separate API keys per environment/project.

---

## Troubleshooting

* **`error: OPENAI_API_KEY not set`** – Export your key and restart the shell.
* **`API error: 401 Unauthorized`** – Invalid/expired key or wrong org/project.
* **`no text in response`** – Retry; if persistent, try `-effort=medium` and check proxy/network.
* **Windows SmartScreen** – Unsigned binaries may need one-time approval.

---

## License

MIT (see `LICENSE`).

---

## Contributing

Issues and PRs welcome—especially improvements to error messages, extra flags like `-paras`, or optional UK-style defaults.

