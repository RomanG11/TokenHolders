package enviroment

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	TestEnv    = "unit-test"
	DockerEnv  = "docker"
	DefaultEnv = "default"
)

func VCSDir() (string, error) {
	wd, _ := os.Getwd()
	return findVCS(wd, 10)
}

func findVCS(dir string, nestLevel int) (string, error) {
	if nestLevel == 0 {
		return "", os.ErrNotExist
	}
	_, err := os.Stat(fmt.Sprintf("%s/.git", dir))
	if err != nil {
		if os.IsNotExist(err) {
			return findVCS(fmt.Sprintf("%s/..", dir), nestLevel-1)
		}
		return "", os.ErrNotExist
	}
	return dir, nil
}

func GetEnv() map[string]string {
	env := os.Getenv("ENV")
	if env == "" {
		env = DefaultEnv
	}
	var myEnv = make(map[string]string, 0)

	if env != DockerEnv {
		wd, err := VCSDir()
		if err != nil {
			log.Panic().Msgf("Can't find project root %s", wd)
		}
		file := fmt.Sprintf("%s/configs/.%s", wd, env)
		godotenv.Load()
		err = godotenv.Load(file)
		if err != nil {
			log.Panic().Msgf("Can't load enviroment file, create %s/.%s file in configs folder", wd, env)
		}
		myEnv, err = godotenv.Read(file)
		if err != nil {
			log.Panic().Msgf("Can't read .%s file", env)
		}

	} else {
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			myEnv[pair[0]] = pair[1]
		}
	}

	return myEnv
}
