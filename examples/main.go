package main

import (
	"fmt"
	"time"

	"github.com/Implex-ltd/crapsolver/crapsolver"
)

func main() {
	Crap := crapsolver.NewSolver()
	Crap.SetWaitTime(time.Second)

	token, err := Crap.Solve(&crapsolver.TaskConfig{
		SiteKey:  "4c672d35-0701-42b2-88c3-78380b0db560",
		Domain:   "discord.com",
		TaskType: crapsolver.TASKTYPE_ENTERPRISE,
		A11YTfe:  true,
		Turbo:    true,
		TurboSt:  2900,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("solved:", token)
}
