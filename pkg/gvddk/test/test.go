package main

import (
	"fmt"
	"github.com/vmware/gvddk/gDiskLib"
	gvddk_high "github.com/vmware/gvddk/gvddk-high"
)

func main() {
	// Init
	fmt.Println("Init Start")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	res1 := gDiskLib.Init(majorVersion, minorVersion, path)
	fmt.Println(res1)
	fmt.Println("Init End.")

	// Connect
	fmt.Println("Connect start")
	//params := gDiskLib.NewConnectParams("", "10.161.131.94", "D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
	//	"Admin!23", "093e3932-6484-4c25-822b-4cf2ad8bd7b5", "datastore-58", "7b9b2783-274b-44e8-8b7f-ac16e0f25ba0", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ | gDiskLib.VIXDISKLIB_FLAG_OPEN_READ_ONLY,
	//	true, "nbd")
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//gDiskLib.EndAccess(params)
	error := gDiskLib.PrepareForAccess(params)
	fmt.Println(error)
	fmt.Println("Prepare finished.")
	conn, err := gDiskLib.ConnectEx(params)
	//conn, err := gDiskLib.Connect(params, false, "nbd")
	fmt.Println(err)
	fmt.Println("Connect end")

	// Open
	fmt.Println("Open a vmdk file")
	//filePath := "[datastore1] fcd/618e688af55d4039a8da46cd18d4151d.vmdk"
	dli, e := gDiskLib.Open(conn, params)
	//_, e := gDiskLib.Open(conn, filePath, gDiskLib.VIXDISKLIB_FLAG_OPEN_READ_ONLY)
	fmt.Println(e)
	fmt.Println("Open end")
	diskHandle := gvddk_high.NewDiskHandle(dli, conn, params)

	// ReadAt
	fmt.Printf("ReadAt test\n")
	buffer := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	n, err4 := diskHandle.ReadAt(buffer, 0)
	fmt.Printf("Read byte n = %d\n", n)
	fmt.Println(buffer)
	fmt.Println(err4)
	//// WriteAt
	//fmt.Println("WriteAt start")
	//buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//for i,_ := range(buf1) {
	//	buf1[i] = 'A'
	//}
	//n2, err2 := dli.WriteAt(buf1, 0)
	//fmt.Printf("Write byte n = %d\n", n2)
	//fmt.Println(err2)
	//
	//fmt.Printf("ReadAt test\n")
	//buffer := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//n, err4 := dli.ReadAt(buffer, 0)
	//fmt.Printf("Read byte n = %d\n", n)
	//fmt.Println(buffer)
	//fmt.Println(err4)
	////
	//buffer2 := make([]byte, C.VIXDISKLIB_SECTOR_SIZE * 2)
	//n, err1 := dli.ReadAt(buffer2, 0)
	//fmt.Printf("Read byte n = %d\n", n)
	//fmt.Println(buffer2)
	//fmt.Println(err1)

	//// Test 2
	//buffer = make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//vixE := gDiskLib.ReadMetadata(dli,"adapterType", buffer, gDiskLib.VIXDISKLIB_SECTOR_SIZE, 128)
	//fmt.Println(buffer)
	//fmt.Println(vixE)
	//
	//buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//er := gDiskLib.GetMetadataKeys(dli, buffer2, gDiskLib.VIXDISKLIB_SECTOR_SIZE, 128)
	//fmt.Println(buffer2)
	//fmt.Println(er)
	//
	//mode := gDiskLib.GetTransportMode(dli);
	//fmt.Println(mode)
	//
	//modes := gDiskLib.ListTransportModes()
	//fmt.Println(modes)

	error1 := gDiskLib.Close(dli)
	fmt.Println(error1)
	error2 := gDiskLib.Disconnect(conn)
	fmt.Println(error2)
	error3 := gDiskLib.EndAccess(params)
	fmt.Println(error3)
	////gDiskLib.Exit()

}
