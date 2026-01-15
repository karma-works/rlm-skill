#!/bin/bash

set -e

REPO_BASE="https://raw.githubusercontent.com/BowTiedSwan/rlm-skill/main"
CLAUDE_DIR="$HOME/.claude/skills"
SKILL_DIR="$CLAUDE_DIR/rlm"

GREEN='\033[0;32m'
GRAY='\033[0;90m'
NC='\033[0m'

echo ""
echo -e "${GRAY}Detecting environment...${NC}"

if [ -d "$HOME/.claude" ]; then
    echo -e "${GREEN}✓ Claude Code detected${NC}"
    mkdir -p "$SKILL_DIR"
    
    echo -e "${GRAY}Downloading skill files...${NC}"
    curl -sSL "$REPO_BASE/SKILL.md" -o "$SKILL_DIR/SKILL.md"
    curl -sSL "$REPO_BASE/rlm.ts" -o "$SKILL_DIR/rlm.ts"
    curl -sSL "$REPO_BASE/package.json" -o "$SKILL_DIR/package.json"
    curl -sSL "$REPO_BASE/tsconfig.json" -o "$SKILL_DIR/tsconfig.json"
    
    echo -e "${GRAY}Installing dependencies...${NC}"
    cd "$SKILL_DIR"
    npm install
    
    echo ""
    echo -e "${GREEN}> /rlm installed successfully${NC}"
    echo -e "${GRAY}  Skill: $SKILL_DIR/SKILL.md${NC}"
    echo -e "${GRAY}  Engine: $SKILL_DIR/rlm.ts${NC}"
    echo ""
    exit 0
else
    echo "Claude Code directory (~/.claude) not found."
    echo "Creating directory anyway..."
    mkdir -p "$SKILL_DIR"
    curl -sSL "$REPO_BASE/SKILL.md" -o "$SKILL_DIR/SKILL.md"
    curl -sSL "$REPO_BASE/rlm.ts" -o "$SKILL_DIR/rlm.ts"
    curl -sSL "$REPO_BASE/package.json" -o "$SKILL_DIR/package.json"
    curl -sSL "$REPO_BASE/tsconfig.json" -o "$SKILL_DIR/tsconfig.json"
    cd "$SKILL_DIR"
    npm install
    echo -e "${GREEN}> /rlm installed${NC}"
    exit 0
fi
