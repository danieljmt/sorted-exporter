package paprika

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Collection struct {
	recipes []Recipe
}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) Add(r ...Recipe) {
	c.recipes = append(c.recipes, r...)
}

func (c *Collection) Export(dest string) error {
	archive, err := os.Create(filepath.Join(dest, "Recipes.paprikarecipes"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := archive.Close(); err != nil {
			log.Println("failed to close archive writer: %w", err)
		}
	}()
	zw := zip.NewWriter(archive)
	defer func() {
		if err := zw.Close(); err != nil {
			log.Println("failed to close zip writer: %w", err)
		}
	}()

	for _, r := range c.recipes {
		data, err := r.compress()
		if err != nil {
			return fmt.Errorf("failed compressing recipe: %w", err)
		}

		illegalChars := regexp.MustCompile("[/\\?%*:|\"<>]")
		fileName := illegalChars.ReplaceAllString(r.Name, "_")

		iow, err := zw.Create(filepath.Clean(fmt.Sprintf("%s.paprikarecipe", fileName)))
		if err != nil {
			return fmt.Errorf("failed creating recipe in archive: %w", err)
		}
		if _, err := iow.Write(data); err != nil {
			return fmt.Errorf("failed writing recipe to archive: %w", err)
		}
	}

	return nil
}

type Recipe struct {
	UID             string   `json:"uid"`
	Created         string   `json:"created"`
	Hash            string   `json:"hash"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Ingredients     string   `json:"ingredients"`
	Directions      string   `json:"directions"`
	Notes           string   `json:"notes"`
	NutritionalInfo string   `json:"nutritional_info"`
	PrepTime        string   `json:"prep_time"`
	CookTime        string   `json:"cook_time"`
	TotalTime       string   `json:"total_time"`
	Difficulty      string   `json:"difficulty"`
	Servings        string   `json:"servings"`
	Rating          int      `json:"rating"`
	Source          string   `json:"source"`
	SourceURL       string   `json:"source_url"`
	Photo           string   `json:"photo"`
	PhotoLarge      string   `json:"photo_large"`
	PhotoHash       string   `json:"photo_hash"`
	ImageURL        string   `json:"image_url"`
	Categories      []string `json:"categories"`
	PhotoData       string   `json:"photo_data"`
}

func (r *Recipe) compress() ([]byte, error) {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("failed marshaling recipe `%s`: %w", r.Name, err)
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
