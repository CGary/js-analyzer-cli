package identifier

import gonanoid "github.com/matoous/go-nanoid"

// GenerateSuffix crea un ID corto, seguro y alfanumérico.
func GenerateSuffix() string {
	// Usamos un alfabeto seguro (sin caracteres confusos) y longitud 5
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 5)
	if err != nil {
		// Fallback de seguridad poco probable
		return "abcde"
	}
	return id
}
