package main

import (
	"flag"
	"fmt"
	"github.com/mattn/go-gimei"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func main() {
	inputSeed := flag.Int64("seed", -1, "seed of random")
	flag.Parse()

	seed := NewSeed(*inputSeed)
	config := CreateConfig(seed)
	names, err := Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	name := names.GetName(config)
	first := name.First
	last := name.Last
	output := fmt.Sprintf("%s %s %s %s %s %s", first.Kanji(), first.Katakana(), first.Hiragana(), last.Kanji(), last.Katakana(), last.Hiragana())
	fmt.Println(output)
}

type Config struct {
	Random *rand.Rand
}

func (config *Config) Select(bound int) int {
	return config.Random.Intn(bound)
}

type Sex func(firstName FirstName) gimei.Item

func (config *Config) SelectSex() Sex {
	sex := config.Select(2)
	if sex%2 == 0 {
		return func(firstName FirstName) gimei.Item {
			female := firstName.Female
			return female[config.Select(len(female))]
		}
	} else {
		return func(firstName FirstName) gimei.Item {
			male := firstName.Male
			return male[config.Select(len(male))]
		}
	}
}

func Load() (*Names, error) {
	file, err := gimei.Assets.Open("/data/names.yml")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var names Names
	err = yaml.Unmarshal(bytes, &names)
	if err != nil {
		return nil, err
	}
	return &names, nil
}

type FirstName struct {
	Male   []gimei.Item `yaml:"male"`
	Female []gimei.Item `yaml:"female"`
	Animal []gimei.Item `yaml:"animal"`
}

type Names struct {
	FirstName   FirstName    `yaml:"first_name"`
	LastName    []gimei.Item `yaml:"last_name"`
	LastNameDog []gimei.Item `yaml:"last_name_dog"`
	LastNameCat []gimei.Item `yaml:"last_name_cat"`
}

func (names *Names) GetName(config Config) gimei.Name {
	sex := config.SelectSex()
	lastName := names.LastName
	return gimei.Name{
		First: sex(names.FirstName),
		Last:  lastName[config.Select(len(lastName))],
	}
}

type Seed int64

func NewSeed(input int64) Seed {
	if input < 0 {
		return Seed(time.Now().UnixNano())
	}
	return Seed(input)
}

func CreateConfig(seed Seed) Config {
	source := rand.NewSource(int64(seed))
	random := rand.New(source)
	return Config{
		Random: random,
	}
}
