# Recursive Language Model (RLM) Skill

> **"Context is an external resource, not a local variable."**

This skill equips Claude Code (and compatible agents) with the **Recursive Language Model (RLM)** pattern described in the research paper:
**[Recursive Language Modeling (ArXiv:2512.24601)](https://arxiv.org/pdf/2512.24601)**.

It enables the agent to handle massive codebases (100+ files, millions of lines) by treating the filesystem as a database and using parallel background agents to process information recursively, eliminating "context rot".

## 📦 Installation

Run this one-liner in your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/BowTiedSwan/rlm-skill/main/install.sh | bash
```

Auto-detects Claude Code and installs the skill.

## 🚀 Usage

Once installed, simply ask Claude to handle a large task:

> "Use RLM to analyze the entire codebase for security vulnerabilities."
> "Scan all 500 files and find where UserID is defined."

The skill triggers automatically on keywords like:
- "analyze codebase"
- "scan all files"
- "large repository"
- "RLM"

## 🧠 How It Works

The skill operates in two distinct modes to eliminate "context rot":

1.  **Native Mode (Default)**: Optimized for **Zero-Shot Filtering**. It uses high-speed filesystem tools like `grep` and `find` for rapid codebase traversal and pattern discovery. Best for mapping project structure and locating specific definitions.
2.  **Strict Mode (Paper Implementation)**: Optimized for **Dense Data Processing**. It uses the `rlm.py` engine to perform **Programmatic Slicing (Chunking)**. By loading data into memory and serving it in atomic chunks, it allows precise analysis of massive logs, monorepos, and CSVs that exceed standard context limits.

### The Pipeline
1.  **Index**: The agent scans your file structure using `find` or `ls`.
2.  **Filter**: It uses `grep` / `ripgrep` to narrow down candidate files (Zero-Shot Filtering).
3.  **Map**: It spawns multiple **parallel background agents**. Each sub-agent reads *one* file and answers *one* question.
4.  **Reduce**: The main agent collects the structured outputs and synthesizes the final answer.

## 📜 Credits

- **Research Paper**: [Recursive Language Modeling](https://arxiv.org/pdf/2512.24601)
- **Skill Author**: [Bowtiedswan](https://x.com/Bowtiedswan)
