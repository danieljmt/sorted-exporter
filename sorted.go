package sorted

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/danieljmt/sorted-exporter/paprika"
	"github.com/google/uuid"
)

type Engine struct {
	destination string
	c           *Client
	format      Format
}

const (
	sortedEndpoint = "https://cook.sorted.club"
	defaultFormat  = Paprika
)

func New(user, pass, dest string) (*Engine, error) {
	e := &Engine{
		format:      defaultFormat,
		destination: dest,
	}
	var err error
	e.c, err = NewClient(sortedEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}
	token, err := e.Auth(user, pass)
	if err != nil {
		return nil, err
	}
	e.c.SetHeader("Authorization", fmt.Sprintf("Token %s", token))
	return e, nil
}

func (e *Engine) Run() error {
	packs, err := e.GetPacks()
	if err != nil {
		return err
	}

	var fullPacks []*Pack
	for _, p := range packs {
		pack, err := e.GetPack(p.ID)
		if err != nil {
			return err
		}
		fullPacks = append(fullPacks, pack)
	}
	if err := e.Export(fullPacks); err != nil {
		return err
	}

	return nil
}

type Format string

const (
	Paprika Format = "paprika"
)

func (e *Engine) Export(packs []*Pack) error {
	switch e.format {
	case Paprika:
		collection := paprika.NewCollection()
		for _, p := range packs {
			recipe, err := e.formatPaprika(p)
			if err != nil {
				return fmt.Errorf("failed encoding pack `%d`: %w", p.ID, err)
			}
			collection.Add(recipe...)
		}

		return collection.Export(e.destination)
	}
	return nil
}

func (e *Engine) formatPaprika(p *Pack) ([]paprika.Recipe, error) {
	var recipes []paprika.Recipe
	for _, sr := range p.Pack.Recipes {
		imParts := strings.Split(sr.AltImages.Thumbnail, ".")
		uid := uuid.NewMD5(uuid.Nil, []byte(sr.Title))

		resp, err := http.Get(sr.AltImages.Thumbnail)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		imgBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		imgData := base64.StdEncoding.EncodeToString(imgBytes)

		pr := paprika.Recipe{
			UID:             uid.String(),
			Created:         time.Now().Format("2006-01-02 15:04:05"), // 2020-08-31 20:51:49
			Name:            sr.Title,
			Description:     "",
			Ingredients:     sr.Ingredients.String(),
			Directions:      sr.Method.String(),
			Notes:           "",
			NutritionalInfo: "",
			PrepTime:        "",
			CookTime:        "",
			TotalTime:       sr.CookingTime.String(),
			Difficulty:      "",
			Servings:        "2",
			Rating:          0,
			Source:          fmt.Sprintf("Sorted Packs - %s", p.Pack.Name),
			SourceURL:       "",
			Photo:           fmt.Sprintf("%s.%s", uid, imParts[len(imParts)-1]),
			PhotoLarge:      "",
			PhotoHash:       "",
			ImageURL:        sr.AltImages.Thumbnail,
			Categories:      []string{fmt.Sprintf("Sorted - %s", strings.Title(p.Pack.Name))},
			PhotoData:       imgData,
		}
		recipes = append(recipes, pr)
	}
	return recipes, nil
}

func (e *Engine) Auth(user, pass string) (string, error) {
	creds := Credentials{
		Username: user,
		Password: pass,
	}
	var authResp AuthResponse
	if err := e.c.Do(context.Background(), "POST", "/auth/login/", creds, &authResp); err != nil {
		return "", fmt.Errorf("failed to auth: %w", err)
	}
	return authResp.Key, nil
}

func (e *Engine) GetPacks() (Packs, error) {
	var packs Packs
	if err := e.c.Do(context.Background(), "GET", "/api/v1/user/packs/", nil, &packs); err != nil {
		return nil, fmt.Errorf("failed to list packs: %w", err)
	}
	return packs, nil
}

func (e *Engine) GetPack(id int) (*Pack, error) {
	var pack Pack
	if err := e.c.Do(context.Background(), "GET", fmt.Sprintf("/api/v1/user/packs/%d", id), nil, &pack); err != nil {
		return nil, fmt.Errorf("failed to get pack `%d`: %w", id, err)
	}
	return &pack, nil
}
