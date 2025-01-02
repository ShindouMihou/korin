package main

const GO_PATH = "{$ENV:GOPATH}" // +k:env

const (
	JAVA_HOME   = "{$ENV:JAVA_HOME}" // +k:env
	JAVA_HOME_2 = "{$ENV:JAVA_HOME}" // +k:env
)

var GO_PATH_2 = "{$ENV:GOPATH}" // +k:env

var (
	JAVA_HOME_3 = "{$ENV:JAVA_HOME}" // +k:env
	JAVA_HOME_4 = "{$ENV:JAVA_HOME}" // +k:env
)

const PORT = "{$ENV:PORT}"               // +k:env(int)
const MODIFIER = "{$ENV:MODIFIER}"       // +k:env(float64)
const ENABLE_TEST = "{$ENV:ENABLE_TEST}" // +k:env(bool)
