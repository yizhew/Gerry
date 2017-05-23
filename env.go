package gerry

const test = "test"
const production = "production"
const development = "development"

var env string

func init() {
	SetTestEnv()
}

func IsTestEnv() bool {
	return env == test
}

func IsProductionEnv() bool {
	return env == production
}

func IsDevelopmentEnv() bool {
	return env == development
}

func SetTestEnv() {
	env = test
}

func SetProductionEnv() {
	env = production
}

func SetDevelopmentEnv() {
	env = development
}

func GetEnv() string {
	if IsTestEnv() {
		return test
	}

	if IsDevelopmentEnv() {
		return development
	}

	return production
}
