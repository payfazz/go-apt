package configsource

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// FromEnvFile return config source from .env file
func FromEnvFile(filename string) ConfigSource {
	configMap := readFile(filename)
	return FromMap(configMap)
}

func readFile(filename string) map[string]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("env file %s not found, config from file not loaded", filename)
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	result := make(map[string]string)

	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		if key != "" && value != "" {
			result[key] = value
		}
	}
	return result
}

func parseLine(line string) (key string, value string) {
	trimmedLine := strings.TrimSpace(line)

	if len(trimmedLine) == 0 || trimmedLine[0] == '#' {
		return
	}

	splitLine := strings.Split(trimmedLine, "=")
	if len(splitLine) < 2 {
		return
	}

	key = splitLine[0]

	value = strings.TrimSpace(strings.Join(splitLine[1:], ""))
	return
}
