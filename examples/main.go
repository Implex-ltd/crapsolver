package main

import (
	"log"
	"time"

	"github.com/Implex-ltd/crapsolver/crapsolver"
)

var (
	KEY    = "4c672d35-0701-42b2-88c3-78380b0db560"
	DOMAIN = "discord.com"
)

func main() {
	Crap, err := crapsolver.NewSolver("id:superapikey")
	if err != nil {
		panic(err)
	}

	Crap.SetWaitTime(time.Second * 3) // check for complete task every 3s (reduce our load + make less req on your side..)

	// get restriction for the current sitekey
	restrictions, err := crapsolver.GetRestrictions(KEY)
	if err != nil {
		panic(err)
	}

	log.Println(restrictions)

	/**
	 * Use the function "Crap.SolveUntil(config, max_retry...)" to retry if error spawn (leave 0 = infinite) / return list of spawned errors
	 */
	token, err := Crap.Solve(&crapsolver.TaskConfig{
		SiteKey:  KEY,
		Domain:   DOMAIN,
		TaskType: crapsolver.TASKTYPE_ENTERPRISE,
		A11YTfe:  true,
		Turbo:    true,
		TurboSt:  3200,
	})

	if err != nil {
		panic(err)
	}

	log.Println("solved:", token)
}
