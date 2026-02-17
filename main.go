package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/HMasataka/connect/connect"
)

func main() {
	var (
		inPlace bool
		outDir  string
		clean   bool
		config  string
		dryRun  bool
		verbose bool
		help    bool
	)

	flag.BoolVar(&inPlace, "in-place", true, "元ファイルに直接注入する")
	flag.StringVar(&outDir, "out-dir", "", "コピーを作成してそちらに注入する")
	flag.BoolVar(&clean, "clean", false, "注入済みスクリプトを除去する")
	flag.StringVar(&config, "config", "", "設定ファイルのパス")
	flag.BoolVar(&dryRun, "dry-run", false, "変更内容を表示するが、ファイルは書き換えない")
	flag.BoolVar(&verbose, "verbose", false, "処理の詳細を表示する")
	flag.BoolVar(&help, "help", false, "ヘルプを表示する")
	flag.Parse()

	if help {
		printUsage()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: designs-dir is required")
		printUsage()
		os.Exit(1)
	}

	designsDir := flag.Arg(0)
	if err := run(designsDir, outDir, config, clean, dryRun, verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: connect <designs-dir> [options]

Arguments:
  designs-dir          デザインディレクトリのパス

Options:
  --in-place           元ファイルに直接注入する（デフォルト）
  --out-dir <path>     コピーを作成してそちらに注入する
  --clean              注入済みスクリプトを除去する
  --config <path>      設定ファイルのパス（デフォルト: <designs-dir>/connect.json）
  --dry-run            変更内容を表示するが、ファイルは書き換えない
  --verbose            処理の詳細を表示する
  --help               ヘルプを表示する`)
}

func run(designsDir, outDir, configPath string, clean, dryRun, verbose bool) error {
	if configPath == "" {
		configPath = connect.ConfigPath(designsDir)
	}

	cfg, err := connect.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	pages, err := connect.DetectPages(designsDir, cfg.Ignore)
	if err != nil {
		return fmt.Errorf("failed to detect pages: %w", err)
	}

	if verbose {
		fmt.Printf("Detected pages: %v\n", pages)
		for _, ignored := range cfg.Ignore {
			fmt.Printf("Skipping: %s (ignored)\n", ignored)
		}
	}

	targetDir := designsDir
	if outDir != "" {
		if err := copyDir(designsDir, outDir); err != nil {
			return fmt.Errorf("failed to copy directory: %w", err)
		}
		targetDir = outDir
	}

	script := connect.GenerateScript(pages, cfg.Selectors, cfg.Mapping, cfg.Toolbar)

	var processed int
	for _, page := range pages {
		htmlPath := filepath.Join(targetDir, page, "index.html")

		if clean {
			cleaned, err := connect.CleanScript(htmlPath, dryRun)
			if err != nil {
				return fmt.Errorf("failed to clean %s: %w", htmlPath, err)
			}
			if verbose && cleaned {
				fmt.Printf("%s: removed <script data-connect>\n", htmlPath)
			}
		} else {
			if err := connect.InjectScript(htmlPath, script, dryRun); err != nil {
				return fmt.Errorf("failed to inject into %s: %w", htmlPath, err)
			}
			if verbose {
				fmt.Printf("%s: inject <script data-connect> (pages: %d entries)\n", htmlPath, len(pages))
			}
		}
		processed++
	}

	fmt.Printf("%d files processed.\n", processed)
	return nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
