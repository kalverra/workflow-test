package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func main() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "kalverra/parrot:latest",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor:   wait.ForHealthCheck(),
	}
	parrot, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer parrot.Terminate(ctx)
	fmt.Println("Parrot is running!")

	parrotURL, err := parrot.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	resty.New().R().SetContext(ctx).Get(parrotURL + "/health")
	fmt.Println("Parrot is healthy!")
	logReader, err := parrot.Logs(ctx)
	if err != nil {
		panic(err)
	}
	defer logReader.Close()

	fmt.Println("Parrot logs:")
	scanner := bufio.NewScanner(logReader)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Print each line of the log
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
