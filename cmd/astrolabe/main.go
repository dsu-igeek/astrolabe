package main

import (
	"flag"
	"fmt"
	"github.com/vmware-tanzu/astrolabe/gen/client"
	"github.com/vmware-tanzu/astrolabe/gen/client/operations"
)

func main() {
	host := flag.String("host", "localhost:1323", "Astrolabe server")
	insecure := flag.Bool("insecure", false, "Only use HTTP")
	flag.Parse()

	transport := client.DefaultTransportConfig()
	transport.Host = *host
	if *insecure {
		transport.Schemes = []string{"http"}
	}

	client := client.NewHTTPClientWithConfig(nil, transport)

	results, err := client.Operations.ListServices(operations.NewListServicesParams())
	if err != nil {
		fmt.Errorf("ListServices failed with err %v\n", err)
	}

	fmt.Println("Services:")
	for num, curService:= range results.Payload.Services {
		fmt.Printf("%d: %s\n", num, curService)
	}

	lpeParams := operations.NewListProtectedEntitiesParams()
	lpeParams.SetService("ivd")
	lpeResults, err := client.Operations.ListProtectedEntities(lpeParams)
	if err != nil {
		fmt.Errorf("ListProtectedEntities failed with err %v\n", err)
	}

	for num, curPE := range lpeResults.Payload.List {
		fmt.Printf("%d: %v\n", num, curPE)

	}
}
