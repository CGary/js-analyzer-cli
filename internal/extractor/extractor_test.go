package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetDependencies verifica que esbuild rastrea correctamente las dependencias locales
// ignorando módulos que no existen en nuestro entorno de prueba.
func TestGetDependencies(t *testing.T) {
	// 1. t.TempDir() es mágico en Go: crea una carpeta temporal en tu RAM/disco
	// y la destruye automáticamente cuando termina el test. ¡Cero basura!
	tempDir := t.TempDir()

	// 2. Simulamos un pequeño proyecto JavaScript
	utilsCode := `export const saludar = () => "Hola Mundo";`
	indexCode := `import { saludar } from './utils.js'; console.log(saludar());`

	utilsPath := filepath.Join(tempDir, "utils.js")
	indexPath := filepath.Join(tempDir, "index.js")

	// Escribimos los archivos virtuales
	err := os.WriteFile(utilsPath, []byte(utilsCode), 0644)
	if err != nil {
		t.Fatalf("No se pudo crear utils.js: %v", err)
	}

	err = os.WriteFile(indexPath, []byte(indexCode), 0644)
	if err != nil {
		t.Fatalf("No se pudo crear index.js: %v", err)
	}

	// 3. Ejecutamos nuestro adaptador
	files, err := GetDependencies(indexPath)

	// 4. Aserciones (Verificaciones)
	if err != nil {
		t.Errorf("GetDependencies falló inesperadamente: %v", err)
	}

	// Esperamos encontrar exactamente 2 archivos (index.js y utils.js)
	if len(files) != 2 {
		t.Errorf("Se esperaban 2 archivos, pero se encontraron %d", len(files))
	}

	// Verificamos que las rutas extraídas coincidan con las reales
	encontradoIndex := false
	encontradoUtils := false

	for _, f := range files {
		if f == indexPath {
			encontradoIndex = true
		}
		if f == utilsPath {
			encontradoUtils = true
		}
	}

	if !encontradoIndex {
		t.Errorf("El árbol de dependencias no incluyó el archivo principal: %s", indexPath)
	}
	if !encontradoUtils {
		t.Errorf("El árbol de dependencias no incluyó la dependencia: %s", utilsPath)
	}
}
