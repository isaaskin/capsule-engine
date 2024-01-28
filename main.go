package main

import (
	"fmt"
	"time"

	"github.com/isaaskin/capsule-engine/engine"
)

func main() {
	engine := engine.CreateEngine()

	capsules, err := engine.ListCapsules()

	if err != nil {
		panic("Error: " + err.Error())
	}
	fmt.Println(capsules)

	cM, _ := engine.StartEvent()

	go func() {
		for {
			select {
			case m := <-cM:
				fmt.Println(m)
			}
		}
	}()

	// for _, container := range controller.ListContainers() {
	// 	err := controller.StopContainer(container)
	// 	if err != nil {
	// 		log.Fatalln(err.Error())
	// 	}
	// }

	//myDevEnv := dockercontroller.MyDevEnv{®
	//res, err := controller.CreateContainer(&myDevEnv)

	//if err != nil {
	//	//log.Fatalf(err.Error())
	//}

	// err := controller.StartContainer("53806266638e00349aca19a9485a17bdf489ab1dcd51c7445f5dd689e37d82fa")

	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }

	//fmt.Println(fmt.Sprintf("Container has been created: %s", res.ID))

	//controller.Attach()

	// controller.Exec()

	//controller.StartEvent()

	// _, err := controller.CreateContainer()

	// if err != nil {
	// 	fmt.Println(err)
	// }

	defer func() {
		engine.StopCapsule("53806266638e00349aca19a9485a17bdf489ab1dcd51c7445f5dd689e37d82fa")
	}()

	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	select {}
}
