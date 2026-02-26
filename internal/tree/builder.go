package tree

import (
	"path/filepath"
	"sort"
	"strings"
)

// Node representa un archivo o carpeta en el árbol
type Node struct {
	Name     string
	Children map[string]*Node
}

// BuildASCII convierte una lista de rutas absolutas en un árbol de texto visual
func BuildASCII(paths []string, basePath string) string {
	root := &Node{Name: ".", Children: make(map[string]*Node)}

	// Construimos la estructura de datos en forma de árbol
	for _, p := range paths {
		rel, err := filepath.Rel(basePath, p)
		if err != nil {
			rel = filepath.Base(p) // Fallback de seguridad
		}

		parts := strings.Split(rel, string(filepath.Separator))
		current := root
		for _, part := range parts {
			if current.Children[part] == nil {
				current.Children[part] = &Node{Name: part, Children: make(map[string]*Node)}
			}
			current = current.Children[part]
		}
	}

	// Dibujamos el árbol en un string
	var builder strings.Builder
	builder.WriteString(root.Name + "\n")

	// Ordenamos las llaves de la raíz para que la salida sea consistente
	var keys []string
	for k := range root.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		printNode(&builder, root.Children[k], "", i == len(keys)-1)
	}

	return builder.String()
}

// printNode dibuja recursivamente las ramas del árbol ASCII
func printNode(b *strings.Builder, node *Node, prefix string, isLast bool) {
	b.WriteString(prefix)
	if isLast {
		b.WriteString("└── ")
		prefix += "    "
	} else {
		b.WriteString("├── ")
		prefix += "│   "
	}
	b.WriteString(node.Name + "\n")

	var keys []string
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Orden alfabético

	for i, k := range keys {
		printNode(b, node.Children[k], prefix, i == len(keys)-1)
	}
}
