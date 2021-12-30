package testmain

import (
	"os"
	"testing"

	testspace "github.com/Jiu2015/gotestspace"
)

var myTestSpace testspace.Space

func TestMain(m *testing.M) {
	var err error

	myTestSpace, err = testspace.Create(
		testspace.WithPathOption("testspace-*"),
		testspace.WithCleanersOption(func() error {
			// Add custom clean functions
			return nil
		}),
		testspace.WithShellOption(`
			git config --global core.abbrev 10 &&
			git config --global init.defaultBranch master &&
			git init --bare repo.git &&
			git clone repo.git workdir &&
			(
				cd workdir &&
				printf "A\n" >A &&
				git add A &&
				test_tick &&
				git commit -m "A" &&
				printf "B\n" >B &&
				git add B &&
				test_tick &&
				git commit -m "B" &&
				git push -u origin HEAD
			)
		`),
	)
	if err != nil {
		panic(err)
	}
	defer myTestSpace.Cleanup()

	// You could register the other cleaner, it will be called while myTestSpace.Cleanup() running
	myTestSpace.RegistrationCustomCleaner(func() error {
		// Other custom action
		return nil
	})

	res := m.Run()
	os.Exit(res)
}
