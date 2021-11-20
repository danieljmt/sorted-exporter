package sorted

import (
	"fmt"
	"strconv"
	"strings"
)

type Packs []Pack

type Pack struct {
	ID        int       `json:"id"`
	Pack      InnerPack `json:"pack"`
	Active    bool      `json:"active"`
	NumPeople int       `json:"num_people"`
}

type InnerPack struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Recipes []Recipe `json:"recipes"`
	Tags    []string `json:"tags"`
}

type Recipe struct {
	ID          int          `json:"id"`
	CookingTime CookingTimes `json:"cooking_time"`
	Ingredients Ingredients  `json:"ingredients"`
	Method      Method       `json:"method"`
	Title       string       `json:"title"`
	AltImages   struct {
		Thumbnail string `json:"thumbnail"`
	} `json:"alt_images"`
	NumPeople string `json:"num_people"`
}

type (
	CookingTimes []CookingTime
	CookingTime  struct {
		NumPeople int `json:"num_people"`
		Duration  int `json:"duration"`
	}
)

func (cts CookingTimes) String() string {
	for _, ct := range cts {
		if ct.NumPeople != 2 {
			continue
		}
		return fmt.Sprintf("%d mins", ct.Duration)
	}
	return ""
}

type (
	Ingredients []Ingredient
	Ingredient  struct {
		Ingredient struct {
			Name string `json:"name"`
			Type struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}
		} `json:"ingredient"`
		Quantities Quantities `json:"quantities"`
	}
)

func (is Ingredients) String() string {
	var sb strings.Builder
	for _, i := range is {
		if i.Ingredient.Type.Name == "" || i.Ingredient.Type.Name == "Equipment" {
			continue
		}
		sb.WriteString(i.String() + "\n")
	}
	return sb.String()
}

func (i *Ingredient) String() string {
	q := i.Quantities.GetQuantity()
	unit := q.Unit.Name
	if q.Quantity > 1 {
		unit = q.Unit.PluralAbbrev
	}
	return fmt.Sprintf("%s %s %s", strconv.FormatFloat(q.Quantity, 'f', -1, 64), unit, i.Ingredient.Name)
}

type (
	Quantities []Quantity
	Quantity   struct {
		ID   int `json:"id"`
		Unit struct {
			Name         string `json:"name"`
			NameAbbrev   string `json:"name_abbrev"`
			PluralAbbrev string `json:"plural_abbrev"`
			Type         string `json:"type"`
			System       string `json:"system"`
		} `json:"unit"`
		NumPeople int     `json:"num_people"`
		Quantity  float64 `json:"quantity"`
	}
)

func (qs Quantities) GetQuantity() Quantity {
	var quantity Quantity
	for _, q := range qs {
		if q.Unit.System == "si" && q.NumPeople == 2 { // Default to SI units when available
			quantity = q
			break
		}
		if q.NumPeople == 2 {
			quantity = q
		}
	}
	return quantity
}

type Method struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Tasks     []Task `json:"tasks"`
	NumPeople int    `json:"num_people"`
}

type Task struct {
	Method    string `json:"method"`
	NumPeople int    `json:"num_people"`
}

func (m Method) String() string {
	var sb strings.Builder
	for _, s := range m.Steps {
		sb.WriteString(s.String())
	}
	return sb.String()
}

func (s Step) String() string {
	var sb strings.Builder
	if s.NumPeople != 2 {
		return ""
	}
	for _, t := range s.Tasks {
		sb.WriteString(t.Method + "\n")
	}
	return sb.String()
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Key string `json:"key"`
}
