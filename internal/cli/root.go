package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cgary/js-analyzer-cli/internal/extractor"
	"github.com/cgary/js-analyzer-cli/internal/segmentation"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jsa [archivo_de_entrada]",
	Short: "Analizador de dependencias JS/TS a Markdown",
	Long:  `js-analyzer es una herramienta CLI ultrarrápida...`,

	Args: cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		entry, _ := cmd.Flags().GetString("entry")
		outdir, _ := cmd.Flags().GetString("outdir")
		chunkSize, _ := cmd.Flags().GetInt("chunk-size")

		if len(args) > 0 {
			entry = args[0]
		}

		baseName := filepath.Base(entry)
		ext := filepath.Ext(baseName)
		entryCleanName := strings.TrimSuffix(baseName, ext)
		if outdir == "" {
			outdir = "."
		}

		fmt.Printf("🚀 Iniciando análisis desde: %s\n", entry)

		// 1. Extraer dependencias con esbuild
		files, err := extractor.GetDependencies(entry)
		if err != nil {
			fmt.Printf("❌ Error de Extracción: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Se extrajeron %d archivos locales.\n", len(files))

		// 2. Procesar, segmentar y guardar
		err = segmentation.ProcessAndSave(files, outdir, chunkSize, entryCleanName)
		if err != nil {
			fmt.Printf("❌ Error de Segmentación: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🎉 ¡Proceso completado con éxito! Revisa la carpeta '%s'\n", outdir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("entry", "e", "index.js", "Ruta al archivo de entrada principal")
	rootCmd.Flags().StringP("outdir", "o", "", "Directorio de salida (Por defecto: nombre del archivo de entrada)")
	rootCmd.Flags().IntP("chunk-size", "c", 1000000, "Tamaño máximo en caracteres por archivo")
}
