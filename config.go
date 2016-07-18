package main

type RepoConfig struct {
	ApiKey string `yaml:"api_key"`
	Name   string `yaml:"name"`
	Script string `yaml:"script"`
	Tag    string `yaml:"tag"`
}
