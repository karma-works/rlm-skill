# Recursive Language Model (RLM) Skill

> **"Context is an external resource, not a local variable."**

This skill equips Claude Code (and compatible agents like github copilot) with the **Recursive Language Model (RLM)** pattern described in the research paper:
**[Recursive Language Modeling (ArXiv:2512.24601)](https://arxiv.org/pdf/2512.24601)**.

It enables the agent to handle massive codebases (100+ files, millions of lines) by treating the filesystem as a database and using parallel background agents to process information recursively, eliminating "context rot".

## 📦 Installation

This project is built with Go. No external scripts are required.

```bash
# Clone the repository
git clone https://github.com/BowTiedSwan/rlm-skill.git
cd rlm-skill

# Build and install the skill
go build -o rlm rlm.go
./rlm install
```

The `install` command auto-detects **Claude Code** and **GitHub Copilot** and installs the skill locally using embedded resources.

## 📜 Credits & Inspiration

- **Original Inspiration**: This project was inspired by the work of **[Bowtiedswan](https://x.com/Bowtiedswan)**, who first prototyped the RLM pattern for AI agents.
- **Research Paper**: [Recursive Language Modeling (ArXiv:2512.24601)](https://arxiv.org/pdf/2512.24601)

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

1.  **Index**: The agent scans your file structure using `find` or `ls`.
2.  **Filter**: It uses `grep` / `ripgrep` to narrow down candidate files (Zero-Shot Filtering).
3.  **Map**: It spawns multiple **parallel background agents**. Each sub-agent reads *one* file and answers *one* question.
4.  **Reduce**: The main agent collects the structured outputs and synthesizes the final answer.


