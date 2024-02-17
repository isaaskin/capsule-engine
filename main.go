package main

import (
	"fmt"
	"time"

	"github.com/isaaskin/capsule-engine/engine"
	"github.com/isaaskin/capsule-engine/registryhandler"
)

func main() {
	capsuleTemplates, _ := registryhandler.GetRepositoryList()

	fmt.Println(capsuleTemplates)

	eng := engine.CreateEngine()

	// _, err := eng.CreateCapsule(models.CapsuleCreateRequest{
	// 	CapsuleTemplate: models.CapsuleTemplate{
	// 		Name:      "capsule-template-go",
	// 		Namespace: "isaaskin",
	// 	},
	// 	Name:       "GoCapsule",
	// 	WorkingDir: "/ben",
	// })

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	_, err := eng.ListCapsules()

	if err != nil {
		panic("Error: " + err.Error())
	}
	// fmt.Println(capsules)

	cM, _ := eng.StartEvent()

	go func() {
		for {
			m := <-cM
			fmt.Println(m)
		}
	}()

	// // engine.PullImage("isaaskin/capsule:main")

	// r, err := eng.CreateCapsule("isaaskin/capsule", "/ii")

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println(r)

	// defer func() {
	// 	eng.StopCapsule("53806266638e00349aca19a9485a17bdf489ab1dcd51c7445f5dd689e37d82fa")
	// }()

	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	select {}
}
