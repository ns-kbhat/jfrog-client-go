package cliutils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

func IsArtifactoryConfExists() bool {
	return readConf().Artifactory != nil
}

func IsBintrayConfExists() bool {
	return readConf().Bintray != nil
}

func ReadArtifactoryConf() *ArtifactoryDetails {
	details := readConf().Artifactory
	if details == nil {
		return new(ArtifactoryDetails)
	}
	return details
}

func ReadBintrayConf() *BintrayDetails {
	details := readConf().Bintray
	if details == nil {
		return new(BintrayDetails)
	}
	return details
}

func SaveArtifactoryConf(details *ArtifactoryDetails) {
	config := readConf()
	config.Artifactory = details
	saveConfig(config)
}

func SaveBintrayConf(details *BintrayDetails) {
	config := readConf()
	config.Bintray = details
	saveConfig(config)
}

func saveConfig(config *Config) {
	b, err := json.Marshal(&config)
	CheckError(err)
	var content bytes.Buffer
	err = json.Indent(&content, b, "", "  ")
	CheckError(err)
	ioutil.WriteFile(getConFilePath(), []byte(content.String()), 0x777)
}

func readConf() *Config {
	confFilePath := getConFilePath()
	config := new(Config)
	if !IsFileExists(confFilePath) {
		return config
	}
	content := ReadFile(confFilePath)
	json.Unmarshal(content, &config)

	return config
}

func getConFilePath() string {
	userDir := GetHomeDir()
	if userDir == "" {
		Exit(ExitCodeError, "Couldn't find home directory. Make sure your HOME environment variable is set.")
	}
	confPath := userDir + "/.jfrog/"
	os.MkdirAll(confPath, 0777)
	return confPath + "jfrog-cli.conf"
}

type Config struct {
	Artifactory *ArtifactoryDetails `json:"artifactory,omitempty"`
	Bintray     *BintrayDetails     `json:"bintray,omitempty"`
}

type ArtifactoryDetails struct {
	Url            string            `json:"url,omitempty"`
	User           string            `json:"user,omitempty"`
	Password       string            `json:"password,omitempty"`
	SshKeyPath     string            `json:"sshKeyPath,omitempty"`
	SshAuthHeaders map[string]string `json:"-"`
}

type BintrayDetails struct {
	ApiUrl             string `json:"-"`
	DownloadServerUrl  string `json:"-"`
	User               string `json:"user,omitempty"`
	Key                string `json:"key,omitempty"`
	DefPackageLicenses string `json:"defPackageLicense,omitempty"`
}
