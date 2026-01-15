import * as fs from 'fs-extra';
import * as path from 'path';
import { glob } from 'glob';
import { Command } from 'commander';

interface Chunk {
    source: string;
    chunk_id: number;
    content: string;
}

class RLMContext {
    private root: string;
    private index: Record<string, string> = {};
    private chunkSize: number = 5000;

    constructor(rootDir: string = ".") {
        this.root = path.resolve(rootDir);
    }

    async loadContext(pattern: string = "**/*", recursive: boolean = true): Promise<string> {
        const files = await glob(pattern, {
            cwd: this.root,
            nodir: true,
            absolute: true,
            ignore: ['**/.git/**', '**/__pycache__/**', '**/node_modules/**']
        });

        let loadedCount = 0;
        for (const f of files) {
            try {
                const content = await fs.readFile(f, 'utf8');
                this.index[f] = content;
                loadedCount++;
            } catch (error) {
                // Ignore files that can't be read (e.g., binary files)
            }
        }

        const totalChars = Object.values(this.index).reduce((sum, content) => sum + content.length, 0);
        return `RLM: Loaded ${loadedCount} files into hidden context. Total size: ${totalChars} chars.`;
    }

    peek(query: string, contextWindow: number = 200): string[] {
        const results: string[] = [];
        for (const [filePath, content] of Object.entries(this.index)) {
            let start = 0;
            while (true) {
                const idx = content.indexOf(query, start);
                if (idx === -1) break;

                const snippetStart = Math.max(0, idx - contextWindow);
                const snippetEnd = Math.min(content.length, idx + query.length + contextWindow);
                const snippet = content.slice(snippetStart, snippetEnd);
                results.push(`[${path.relative(this.root, filePath)}]: ...${snippet}...`);

                if (results.length >= 20) return results;
                start = idx + 1;
            }
        }
        return results;
    }

    getChunks(filePattern: string | null = null): Chunk[] {
        const chunks: Chunk[] = [];
        const targets = Object.keys(this.index).filter(f => !filePattern || f.includes(filePattern));

        for (const filePath of targets) {
            const content = this.index[filePath];
            const totalChunks = Math.ceil(content.length / this.chunkSize);
            for (let i = 0; i < totalChunks; i++) {
                const start = i * this.chunkSize;
                const end = Math.min((i + 1) * this.chunkSize, content.length);
                chunks.push({
                    source: path.relative(this.root, filePath),
                    chunk_id: i,
                    content: content.slice(start, end)
                });
            }
        }
        return chunks;
    }
}

async function main() {
    const program = new Command();
    const ctx = new RLMContext();

    program
        .name('rlm')
        .description('RLM Engine in TypeScript')
        .version('1.0.0');

    program.command('scan')
        .description('Scan the directory for files')
        .option('--path <path>', 'Path to scan', '.')
        .action(async (options) => {
            const scanner = new RLMContext(options.path);
            const result = await scanner.loadContext();
            console.log(result);
        });

    program.command('peek')
        .description('Peek into the context for a query')
        .argument('<query>', 'Query string')
        .action(async (query) => {
            await ctx.loadContext();
            const results = ctx.peek(query);
            console.log(JSON.stringify(results, null, 2));
        });

    program.command('chunk')
        .description('Get chunks of the context')
        .option('--pattern <pattern>', 'File pattern to filter')
        .action(async (options) => {
            await ctx.loadContext();
            const chunks = ctx.getChunks(options.pattern);
            console.log(JSON.stringify(chunks));
        });

    await program.parseAsync(process.argv);
}

main().catch(console.error);
