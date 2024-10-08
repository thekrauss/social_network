package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func UploadImages(w http.ResponseWriter, r *http.Request, filePath string) (string, error) {
	// Limite de taille du formulaire multipart (20 MB)
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		return "", fmt.Errorf("error parsing form: %v", err)
	}

	// Récupère le fichier image (si fourni)
	file, handler, err := r.FormFile("image")
	if err != nil {
		// Si aucun fichier n'est fourni, on ne lève pas d'erreur ici, mais on retourne une chaîne vide
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", fmt.Errorf("error retrieving file from form: %v", err)
	}
	defer file.Close()

	// Lire les données du fichier
	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file data: %v", err)
	}

	// Vérification de la taille du fichier
	fileSize := int64(len(fileData))
	maxFileSize := int64(20 << 20)
	if fileSize > maxFileSize {
		return "", fmt.Errorf("file size too large: %d bytes (max allowed: %d bytes)", fileSize, maxFileSize)
	}

	// Création et écriture du fichier sur le serveur
	f, err := os.OpenFile(filepath.Join(filePath, handler.Filename), os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return "", fmt.Errorf("error creating file on server: %v", err)
	}
	defer f.Close()

	_, err = f.Write(fileData)
	if err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}

	return path.Join(".", filePath, handler.Filename), nil
}

func IsValidImageExtension(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	fileExt := strings.ToLower(filepath.Ext(filename))
	for _, ext := range validExtensions {
		if ext == fileExt {
			return true
		}
	}
	return false
}
