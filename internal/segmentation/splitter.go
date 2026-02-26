package segmentation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cgary/js-analyzer-cli/internal/identifier"
	"github.com/cgary/js-analyzer-cli/internal/tree"
)

func getMarkdownLang(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".ts":
		return "typescript"
	case ".tsx":
		return "tsx"
	case ".js", ".cjs", ".mjs":
		return "javascript"
	case ".jsx":
		return "jsx"
	case ".json":
		return "json"
	default:
		return "javascript" // fallback por defecto
	}
}

func ProcessAndSave(files []string, outDir string, chunkSize int, entryName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo el directorio actual: %v", err)
	}

	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creando directorio de salida: %v", err)
		}
	}

	treeASCII := tree.BuildASCII(files, cwd)
	var currentChunk string = fmt.Sprintf("# Árbol de Dependencias de %s\n```text\n%s```\n\n---\n\n", entryName, treeASCII)

	var chunkIndex int = 1
	var totalBytes int = len(currentChunk)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		relPath, err := filepath.Rel(cwd, file)
		if err != nil {
			relPath = filepath.Base(file)
		}

		lang := getMarkdownLang(file)

		formattedFile := fmt.Sprintf("## Archivo: %s\n```%s\n%s\n```\n\n", relPath, lang, string(content))

		if len(currentChunk)+len(formattedFile) > chunkSize && len(currentChunk) > 0 {
			// NUEVO: Pasamos el entryName
			if err := guardarArchivo(outDir, chunkIndex, currentChunk, entryName); err != nil {
				return err
			}
			chunkIndex++
			currentChunk = ""
		}

		currentChunk += formattedFile
		totalBytes += len(formattedFile)
	}

	if len(currentChunk) > 0 {
		// NUEVO: Pasamos el entryName
		if err := guardarArchivo(outDir, chunkIndex, currentChunk, entryName); err != nil {
			return err
		}
	}

	fmt.Printf("📦 Generación completada: %d archivo(s) consolidados.\n", chunkIndex)
	return nil
}

func guardarArchivo(outDir string, index int, content string, entryName string) error {
	suffix := identifier.GenerateSuffix()

	outName := filepath.Join(outDir, fmt.Sprintf("%s_part%d_%s.md", entryName, index, suffix))

	if err := os.WriteFile(outName, []byte(content), 0644); err != nil {
		return fmt.Errorf("error guardando %s: %v", outName, err)
	}

	fmt.Printf("  -> Creado: %s (%d bytes)\n", outName, len(content))
	return nil
}
