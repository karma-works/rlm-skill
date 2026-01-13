package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
)

//go:embed SKILL.md
var skillContent string

type RLMContext struct {
	Root      string
	Index     map[string]string
	ChunkSize int
}

func NewRLMContext(root string) *RLMContext {
	return &RLMContext{
		Root:      root,
		Index:     make(map[string]string),
		ChunkSize: 5000,
	}
}

func (c *RLMContext) LoadContext(pattern string, recursive bool) string {
	loadedCount := 0
	totalSize := 0

	err := filepath.WalkDir(c.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "__pycache__" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		content, err := os.ReadFile(path)
		if err == nil {
			c.Index[path] = string(content)
			loadedCount++
			totalSize += len(content)
		}
		return nil
	})

	if err != nil {
		return fmt.Sprintf("RLM: Error walking directory: %v", err)
	}

	return fmt.Sprintf("RLM: Loaded %d files into hidden context. Total size: %d chars.", loadedCount, totalSize)
}

func (c *RLMContext) Peek(query string, contextWindow int) []string {
	var results []string
	for path, content := range c.Index {
		if strings.Contains(content, query) {
			start := 0
			for {
				idx := strings.Index(content[start:], query)
				if idx == -1 {
					break
				}
				absIdx := start + idx

				snippetStart := absIdx - contextWindow
				if snippetStart < 0 {
					snippetStart = 0
				}

				snippetEnd := absIdx + len(query) + contextWindow
				if snippetEnd > len(content) {
					snippetEnd = len(content)
				}

				snippet := content[snippetStart:snippetEnd]
				results = append(results, fmt.Sprintf("[%s]: ...%s...", path, snippet))

				start = absIdx + 1
				if start >= len(content) {
					break
				}

				if len(results) >= 20 {
					return results
				}
			}
		}
	}
	return results
}

type Chunk struct {
	Source  string `json:"source"`
	ChunkID int    `json:"chunk_id"`
	Content string `json:"content"`
}

func (c *RLMContext) GetChunks(filePattern string) []Chunk {
	var chunks []Chunk
	for path, content := range c.Index {
		if filePattern == "" || strings.Contains(path, filePattern) {
			totalChunks := int(math.Ceil(float64(len(content)) / float64(c.ChunkSize)))
			for i := 0; i < totalChunks; i++ {
				start := i * c.ChunkSize
				end := (i + 1) * c.ChunkSize
				if end > len(content) {
					end = len(content)
				}
				chunks = append(chunks, Chunk{
					Source:  path,
					ChunkID: i,
					Content: content[start:end],
				})
			}
		}
	}
	return chunks
}

func installSkill() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	targets := []string{}
	claudePath := filepath.Join(home, ".claude")
	copilotPath := filepath.Join(home, ".copilot")

	if _, err := os.Stat(claudePath); err == nil {
		targets = append(targets, filepath.Join(claudePath, "skills", "rlm"))
	}
	if _, err := os.Stat(copilotPath); err == nil {
		targets = append(targets, filepath.Join(copilotPath, "skills", "rlm"))
	}

	if len(targets) == 0 {
		fmt.Println("No AI agents (Claude Code or Copilot) detected in home directory.")
		fmt.Println("Installing to default ~/.claude/skills/rlm...")
		targets = append(targets, filepath.Join(claudePath, "skills", "rlm"))
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	for _, target := range targets {
		fmt.Printf("Installing to %s...\n", target)
		if err := os.MkdirAll(target, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", target, err)
			continue
		}

		// Write skill definition
		skillFile := filepath.Join(target, "SKILL.md")
		if err := os.WriteFile(skillFile, []byte(skillContent), 0644); err != nil {
			fmt.Printf("Error writing SKILL.md to %s: %v\n", target, err)
			continue
		}

		// Copy executable
		dstName := "rlm"
		if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") || filepath.Ext(execPath) == ".exe" {
			dstName = "rlm.exe"
		}
		dstPath := filepath.Join(target, dstName)

		if err := copyFile(execPath, dstPath); err != nil {
			fmt.Printf("Error copying executable to %s: %v\n", dstPath, err)
			continue
		}
		if err := os.Chmod(dstPath, 0755); err != nil {
			fmt.Printf("Error setting permissions on %s: %v\n", dstPath, err)
		}
		fmt.Printf("✓ Installed successfully to %s\n", target)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'scan', 'peek', 'chunk' or 'install' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		installSkill()
		return

	case "scan":
		ctx := NewRLMContext(".")
		ctx.LoadContext("**/*", true)
		scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
		pathPtr := scanCmd.String("path", ".", "path to scan")
		scanCmd.Parse(os.Args[2:])

		if *pathPtr != "." {
			ctx = NewRLMContext(*pathPtr)
			fmt.Println(ctx.LoadContext("**/*", true))
		} else {
			fmt.Println("RLM: Loaded " + fmt.Sprint(len(ctx.Index)) + " files into hidden context.")
		}

	case "peek":
		ctx := NewRLMContext(".")
		ctx.LoadContext("**/*", true)
		peekCmd := flag.NewFlagSet("peek", flag.ExitOnError)
		peekCmd.Parse(os.Args[2:])
		if peekCmd.NArg() < 1 {
			fmt.Println("peek requires a query")
			os.Exit(1)
		}
		query := peekCmd.Arg(0)
		results := ctx.Peek(query, 200)
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))

	case "chunk":
		ctx := NewRLMContext(".")
		ctx.LoadContext("**/*", true)
		chunkCmd := flag.NewFlagSet("chunk", flag.ExitOnError)
		patternPtr := chunkCmd.String("pattern", "", "file pattern")
		chunkCmd.Parse(os.Args[2:])

		chunks := ctx.GetChunks(*patternPtr)
		data, _ := json.Marshal(chunks)
		fmt.Println(string(data))

	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
