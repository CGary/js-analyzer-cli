# js-analyzer-cli

**js-analyzer-cli** (ejecutable: `jsa`) es una herramienta de interfaz de línea de comandos (CLI) desarrollada en Go. Su función principal es analizar el punto de entrada de aplicaciones JavaScript o TypeScript, extraer su árbol de dependencias interno utilizando `esbuild` y consolidar el código fuente en archivos Markdown debidamente estructurados y segmentados.

## Características Principales

* **Rastreo Preciso de Dependencias**: Utiliza `esbuild` para resolver módulos estáticos, garantizando que solo el código que forma parte de la aplicación sea extraído.
* **Filtrado de Activos**: Ignora automáticamente archivos estáticos binarios o multimedia (imágenes, fuentes, vectores) para evitar procesamientos innecesarios y mantener el formato de texto limpio.
* **Generación de Árbol ASCII**: Incluye una representación visual en formato de texto del árbol de directorios y dependencias al inicio de la exportación.
* **Segmentación Automática**: Divide el contenido exportado en múltiples archivos Markdown si el tamaño total excede un límite de caracteres configurable, evitando la saturación de los lectores de documentos.
* **Soporte Multilenguaje**: Identifica y aplica el formato de sintaxis correcto en Markdown para extensiones como `.js`, `.ts`, `.jsx`, `.tsx`, `.css`, `.json`, entre otras.

---

## Requisitos Previos

* [Go](https://go.dev/) 1.25.0 o superior.

---

## Instalación

Puede compilar la herramienta directamente desde el código fuente. Clone el repositorio y ejecute el siguiente comando en el directorio raíz:

```bash
go build -o jsa ./cmd/jsa/main.go

```

Para instalarlo globalmente en su sistema (asegúrese de que su ruta `$GOPATH/bin` esté en su variable de entorno `PATH`):

```bash
go install ./cmd/jsa

```

---

## Uso

La sintaxis básica del comando es la siguiente:

```bash
jsa [archivo_de_entrada] [opciones]

```

### Ejemplos de Ejecución

1. **Ejecución básica**:
Analiza `index.js` (valor por defecto) y genera los archivos Markdown en el directorio actual.
```bash
jsa

```


2. **Especificar un punto de entrada y directorio de salida**:
```bash
jsa src/main.ts --outdir ./documentacion

```


3. **Ajustar el límite de segmentación**:
Establece el tamaño máximo del archivo a 500,000 caracteres.
```bash
jsa src/app.tsx -c 500000

```



### Opciones y Banderas

| Bandera Corta | Bandera Larga | Descripción | Valor por Defecto |
| --- | --- | --- | --- |
| `-e` | `--entry` | Ruta al archivo de entrada principal. | `index.js` |
| `-o` | `--outdir` | Directorio de salida para los archivos generados. | Directorio actual (`.`) |
| `-c` | `--chunk-size` | Tamaño máximo en caracteres por archivo Markdown. | `1000000` |

---

## Estructura del Proyecto

El proyecto sigue una arquitectura estándar de Go, separando la interfaz de línea de comandos de la lógica de negocio:

* **`cmd/jsa/`**: Contiene el punto de entrada de la aplicación (`main.go`).
* **`internal/cli/`**: Define la configuración de comandos y banderas utilizando la biblioteca `Cobra`.
* **`internal/extractor/`**: Contiene la lógica de integración con `esbuild` para el análisis y extracción de los archivos de código fuente.
* **`internal/segmentation/`**: Procesa la lectura de los archivos extraídos, aplica el formato Markdown y maneja la escritura en fragmentos (chunks) limitados por tamaño.
* **`internal/tree/`**: Construye la representación en ASCII de la estructura de dependencias relativas.
* **`internal/identifier/`**: Genera identificadores alfanuméricos cortos (mediante `nanoid`) para evitar colisiones en los nombres de los archivos de salida.
