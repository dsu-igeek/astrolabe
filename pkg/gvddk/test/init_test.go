package main

import (
	"fmt"
	"github.com/vmware/gvddk/gDiskLib"
	gvddk_high "github.com/vmware/gvddk/gvddk-high"
	"testing"
)

func TestPos1(t *testing.T) {
	// Set up
	res := gDiskLib.Init(6, 7, "/usr/lib/vmware-vix-disklib")
	if res != nil {
		t.Errorf("Init failed, got error code: %d, error message: %s.", res.VixErrorCode(), res.Error())
	}
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	err1 := gDiskLib.PrepareForAccess(params)
	if err1 != nil {
		t.Errorf("Prepare for access failed. Error code: %d. Error message: %s.", err1.VixErrorCode(), err1.Error())
	}
	conn, err2 := gDiskLib.ConnectEx(params)
	if err2 != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Connect to vixdisk lib failed. Error code: %d. Error message: %s.", err2.VixErrorCode(), err2.Error())
	}

	// Disk operation
	dli, err3 := gDiskLib.Open(conn, params)
	if err3 != nil {
		gDiskLib.Disconnect(conn)
		gDiskLib.EndAccess(params)
		t.Errorf("Open a virtual disk failed. Error code: %d. Error message: %s.", err3.VixErrorCode(), err3.Error())
	}
	fmt.Printf("ReadAt test\n")
	buffer := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	diskHandle := gvddk_high.NewDiskHandle(dli, conn, params)
	n, err4 := diskHandle.ReadAt(buffer, 0)
	if err4 != nil {
		t.Errorf("Read from virtual disk failed.")
	}
	fmt.Printf("Read byte n = %d\n", n)
	gDiskLib.Close(dli)
	gDiskLib.Disconnect(conn)
	gDiskLib.EndAccess(params)
}

//func TestPos2(t *testing.T) {
//	// Set up
//	res := gDiskLib.Init(6, 7, "/usr/lib/vmware-vix-disklib")
//	if res != nil {
//		t.Errorf("Init failed, got error code: %d, error message: %s.", res.VixErrorCode(), res.Error())
//	}
//	params := gDiskLib.NewConnectParams("", "10.161.99.58", "31:E1:D5:67:34:50:30:30:0B:8A:96:C8:F0:D1:3F:D4:FD:6A:46:43", "administrator@vsphere.local",
//		"Admin!23", "cf29221a-381b-4036-825a-56bf8294ed38", "datastore-58", "ecb7fa78-cef9-4459-b898-17a39f582d9b", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_READ_ONLY,true, "nbd")
//	err1 := gDiskLib.PrepareForAccess(params)
//	if err1 != nil {
//		t.Errorf("Prepare for access failed. Error code: %d. Error message: %s.", err1.VixErrorCode(), err1.Error())
//	}
//	conn, err2 := gDiskLib.ConnectEx(params)
//	if err2 != nil {
//		gDiskLib.EndAccess(params)
//		t.Errorf("Connect to vixdisk lib failed. Error code: %d. Error message: %s.", err2.VixErrorCode(), err2.Error())
//	}
//	//var filePath string = "./test.vmbk"
//	//createParams := gDiskLib.NewCreateParams(1, 2, 4, 1024)
//	//err3 := gDiskLib.Create(conn, filePath, createParams, "")
//	//if err3 != nil {
//	//	gDiskLib.Disconnect(conn)
//	//	gDiskLib.EndAccess(params)
//	//	t.Errorf("Create a local virtual disk failed. Error code: %d. Error message: %s.", err3.VixErrorCode(), err3.Error())
//	//}
//	gDiskLib.Disconnect(conn)
//	gDiskLib.EndAccess(params)
//}
