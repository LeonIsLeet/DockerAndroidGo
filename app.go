package main

import (
	"androidauto/DockerManager"
	"fmt"
	"time"
)

func checkError(err error) {
	if err != nil {
		// fmt.Println("Error occured", err)
		panic(err)
	}
}

func main() {
	docker, err := DockerManager.NewDockerManager()
	container, err := docker.CreateContainer("budtmo/docker-android:emulator_11.0", "TestContainer", nil)
	checkError(err)
	fmt.Println("container started")
	tenSeconds := 10 * time.Second
	docker.StopContainer(container, &tenSeconds)

}
