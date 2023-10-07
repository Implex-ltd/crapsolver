package main

import (
	"fmt"
	"time"

	"github.com/Implex-ltd/crapsolver/crapsolver"
)

func main() {
	Crap := crapsolver.NewSolver()
	Crap.SetWaitTime(time.Second * 3) // check for complete task every 3s (reduce our load + make less req on your side..)

	/**
	 * Use the function "Crap.SolveUntil(config, max_retry...)" to retry if error spawn (leave 0 = infinite) / return list of spawned errors
	 */
	token, err := Crap.Solve(&crapsolver.TaskConfig{
		SiteKey:  "4c672d35-0701-42b2-88c3-78380b0db560",
		Domain:   "discord.com",
		TaskType: crapsolver.TASKTYPE_ENTERPRISE,
		A11YTfe:  true,
		Turbo:    true,
		TurboSt:  3200,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("solved:", token)
}
