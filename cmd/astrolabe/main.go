package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"github.com/vmware-tanzu/astrolabe/pkg/server"
	"io"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag {
			&cli.StringFlag{
			Name: "host",
			Value: "localhost:1323",
			Usage: "Astrolabe server",
			},
			&cli.BoolFlag{
				Name:        "insecure",
				Usage:       "Only use HTTP",
				Required:    false,
				Hidden:      false,
				Value:       false,
			},
			&cli.StringFlag{
				Name: "confDir",
				Usage: "Configuration directory",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:	"types",
				Usage: "shows Protected Entity Types",
				Action: types,
			},
			{
				Name: "ls",
				Usage:	"lists entities for a type",
				Action: ls,
			},
			{
				Name: "lssn",
				Usage: "lists snapshots for a Protected Entity",
				Action: lssn,
				ArgsUsage: "<protected entity id>",
			},
			{
				Name: "snap",
				Usage: "snapshots a Protected Entity",
				Action: snap,
				ArgsUsage: "<protected entity id>",
			},
			{
				Name: "rmsn",
				Usage: "removes a Protected Entity snapshot",
				Action: rmsn,
				ArgsUsage: "<protected entity snapshot id>",
			},
			{
				Name: "cp",
				Usage: "copies a Protected Entity snapshot",
				Action: cp,
				ArgsUsage: "<src> <dest>",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setupProtectedEntityManager(c *cli.Context) (astrolabe.ProtectedEntityManager) {
	confDirStr := c.String("confDir")
	_, pem := server.NewProtectedEntityManager(confDirStr, 0)

	return pem
}
func types (c *cli.Context) error {
	pem := setupProtectedEntityManager(c)
	for _, curPETM:= range pem.ListEntityTypeManagers() {
		fmt.Println(curPETM.GetTypeName())
	}
	return nil
}

func ls (c *cli.Context) error {
	pem := setupProtectedEntityManager(c)
	peType := c.Args().First()
	petm := pem.GetProtectedEntityTypeManager(peType)
	if petm == nil {
		log.Fatalf("Could not find type named %s", peType)
	}
	peIDs, err := petm.GetProtectedEntities(context.TODO())
	if err != nil {
		log.Fatalf("Could not retrieve protected entities for type %s err:%b", peType, err)
	}

	for _, curPEID := range peIDs {
		fmt.Println(curPEID.String())
	}
	return nil
}

func lssn (c *cli.Context) error {
	peIDStr := c.Args().First()
	peID, err := astrolabe.NewProtectedEntityIDFromString(peIDStr)
	if err != nil {
		log.Fatalf("Could not parse protected entity ID %s, err: %v", peIDStr, err)
	}

	pem := setupProtectedEntityManager(c)

	pe, err := pem.GetProtectedEntity(context.TODO(), peID)
	if err != nil {
		log.Fatalf("Could not retrieve protected entity ID %s, err: %v", peIDStr, err)
	}

	snaps, err := pe.ListSnapshots(context.TODO())
	if err != nil {
		log.Fatalf("Could not get snapshots for protected entity ID %s, err: %v", peIDStr, err)
	}

	for _, curSnapshotID := range snaps {
		curPESnapshotID := peID.IDWithSnapshot(curSnapshotID)
		fmt.Println(curPESnapshotID.String())
	}
	return nil
}

func snap (c *cli.Context) error {
	peIDStr := c.Args().First()
	peID, err := astrolabe.NewProtectedEntityIDFromString(peIDStr)
	if err != nil {
		log.Fatalf("Could not parse protected entity ID %s, err: %v", peIDStr, err)
	}

	pem := setupProtectedEntityManager(c)

	pe, err := pem.GetProtectedEntity(context.TODO(), peID)
	if err != nil {
		log.Fatalf("Could not retrieve protected entity ID %s, err: %v", peIDStr, err)
	}
	snap, err := pe.Snapshot(context.TODO())
	if err != nil {
		log.Fatalf("Could not snapshot protected entity ID %s, err: %v", peIDStr, err)
	}
	fmt.Println(snap.String())
	return nil
}

func rmsn (c *cli.Context) error {
	peIDStr := c.Args().First()
	peID, err := astrolabe.NewProtectedEntityIDFromString(peIDStr)
	if err != nil {
		log.Fatalf("Could not parse protected entity ID %s, err: %v", peIDStr, err)
	}
	if !peID.HasSnapshot() {
		log.Fatalf("Protected entity ID %s does not have a snapshot ID", peIDStr)
	}

	pem := setupProtectedEntityManager(c)

	pe, err := pem.GetProtectedEntity(context.TODO(), peID)
	if err != nil {
		log.Fatalf("Could not retrieve protected entity ID %s, err: %v", peIDStr, err)
	}
	success, err := pe.DeleteSnapshot(context.TODO(), peID.GetSnapshotID())
	if err != nil {
		log.Fatalf("Could not remove snapshot ID %s, err: %v", peIDStr, err)
	}
	if success {
		log.Printf("Removed snapshot %s\n", peIDStr)
	}
	return nil
}

func cp (c *cli.Context) error {
	if c.NArg() != 2 {
		log.Fatalf("Expected two arguments for cp, got %d", c.NArg())
	}
	srcStr := c.Args().First()
	destStr := c.Args().Get(1)
	var err error
	var srcPEID, destPEID astrolabe.ProtectedEntityID
	var srcFile, destFile string
	srcPEID, err = astrolabe.NewProtectedEntityIDFromString(srcStr)
	if err != nil {
		srcFile = srcStr
	}
	destPEID, err = astrolabe.NewProtectedEntityIDFromString(destStr)
	if err != nil {
		destFile = destStr
	}
	pem := setupProtectedEntityManager(c)

	var reader io.ReadCloser
	var writer io.WriteCloser
	fmt.Printf("cp from ")
	if srcFile != "" {
		fmt.Printf("file %s", srcFile)
	} else {
		fmt.Printf("pe %s", srcPEID.String())
		srcPE, err := pem.GetProtectedEntity(context.TODO(), srcPEID)
		if err != nil {
			log.Fatalf("Could not retrieve protected entity ID %s, err: %v", srcPEID.String(), err)
		}
		var dw io.WriteCloser
		reader, dw = io.Pipe()
		go zipPE(context.TODO(), srcPE, dw)
	}
	fmt.Printf(" to ")
	if destFile != "" {
		fmt.Printf("file %s", destFile)
		writer, err = os.Create(destFile)
		if err != nil {
			log.Fatalf("Could not create file %s, err: %v", destFile, err)
		}
	} else {
		fmt.Printf("pe %s", destPEID.String())
	}
	fmt.Printf("\n")

	bytesCopied, err := io.Copy(writer, reader)
	if err != nil {
		log.Fatalf("Error copying %v", err)
	}
	fmt.Printf("Copied %d bytes\n", bytesCopied)
	return nil
}

func zipPE(ctx context.Context, pe astrolabe.ProtectedEntity, writer io.WriteCloser) {
	defer writer.Close()
	err := astrolabe.ZipProtectedEntity(ctx, pe, writer)
	if err != nil {
		log.Fatalf("Failed to zip protected entity %s, err = %v", pe.GetID().String(), err)
	}
}