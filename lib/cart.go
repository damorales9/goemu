package lib

import (
	"fmt"
	"os"
	"unsafe"
)

// Type aliases for unsigned integers
type u8 uint8
type u16 uint16
type u32 uint32
type u64 uint64

// Structs
type romHeader struct {
	entry          [4]u8
	logo           [0x30]u8
	title          [16]byte
	newLicCode     u16
	sgbFlag        u8
	typeCode       u8
	romSize        u8
	ramSize        u8
	destCode       u8
	licCode        u8
	version        u8
	checksum       u8
	globalChecksum u16
}

type cartContext struct {
	filename string
	romSize  u32
	romData  []u8
	header   *romHeader
}

// Global variable
var ctx cartContext

// ROM Types
var ROM_TYPES = []string{
	"ROM ONLY", "MBC1", "MBC1+RAM", "MBC1+RAM+BATTERY", "0x04 ???",
	"MBC2", "MBC2+BATTERY", "0x07 ???", "ROM+RAM 1", "ROM+RAM+BATTERY 1",
	"0x0A ???", "MMM01", "MMM01+RAM", "MMM01+RAM+BATTERY", "0x0E ???",
	"MBC3+TIMER+BATTERY", "MBC3+TIMER+RAM+BATTERY 2", "MBC3", "MBC3+RAM 2",
	"MBC3+RAM+BATTERY 2", "0x14 ???", "0x15 ???", "0x16 ???", "0x17 ???",
	"0x18 ???", "MBC5", "MBC5+RAM", "MBC5+RAM+BATTERY", "MBC5+RUMBLE",
	"MBC5+RUMBLE+RAM", "MBC5+RUMBLE+RAM+BATTERY", "0x1F ???", "MBC6",
	"0x21 ???", "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
}

// License Codes
var LIC_CODE = map[u8]string{
	0x00: "None", 0x01: "Nintendo R&D1", 0x08: "Capcom", 0x13: "Electronic Arts",
	0x18: "Hudson Soft", 0x19: "b-ai", 0x20: "kss", 0x22: "pow", 0x24: "PCM Complete",
	0x25: "san-x", 0x28: "Kemco Japan", 0x29: "seta", 0x30: "Viacom", 0x31: "Nintendo",
	0x32: "Bandai", 0x33: "Ocean/Acclaim", 0x34: "Konami", 0x35: "Hector", 0x37: "Taito",
	0x38: "Hudson", 0x39: "Banpresto", 0x41: "Ubi Soft", 0x42: "Atlus", 0x44: "Malibu",
	0x46: "angel", 0x47: "Bullet-Proof", 0x49: "irem", 0x50: "Absolute", 0x51: "Acclaim",
	0x52: "Activision", 0x53: "American sammy", 0x54: "Konami", 0x55: "Hi tech entertainment",
	0x56: "LJN", 0x57: "Matchbox", 0x58: "Mattel", 0x59: "Milton Bradley", 0x60: "Titus",
	0x61: "Virgin", 0x64: "LucasArts", 0x67: "Ocean", 0x69: "Electronic Arts", 0x70: "Infogrames",
	0x71: "Interplay", 0x72: "Broderbund", 0x73: "sculptured", 0x75: "sci", 0x78: "THQ",
	0x79: "Accolade", 0x80: "misawa", 0x83: "lozc", 0x86: "Tokuma Shoten Intermedia",
	0x87: "Tsukuda Original", 0x91: "Chunsoft", 0x92: "Video system", 0x93: "Ocean/Acclaim",
	0x95: "Varie", 0x96: "Yonezawa/s’pal", 0x97: "Kaneko", 0x99: "Pack in soft", 0xA4: "Konami (Yu-Gi-Oh!)",
}

func cartLicName() string {
	if ctx.header.newLicCode <= 0xA4 {
		return LIC_CODE[u8(ctx.header.newLicCode)]
	}
	return "UNKNOWN"
}

func cartTypeName() string {
	if ctx.header.typeCode <= 0x22 {
		return ROM_TYPES[ctx.header.typeCode]
	}
	return "UNKNOWN"
}

func cartLoad(cart string) bool {
	ctx.filename = cart

	fp, err := os.Open(cart)
	if err != nil {
		fmt.Printf("Failed to open: %s\n", cart)
		return false
	}
	defer fp.Close()

	fmt.Printf("Opened: %s\n", ctx.filename)

	fi, err := fp.Stat()
	if err != nil {
		fmt.Printf("Failed to get file info: %s\n", cart)
		return false
	}
	ctx.romSize = u32(fi.Size())

	ctx.romData = make([]u8, ctx.romSize)
	byteData := make([]byte, ctx.romSize) // Temporary byte slice for reading data

	_, err = fp.Read(byteData)
	if err != nil {
		fmt.Printf("Failed to read ROM data: %s\n", cart)
		return false
	}

	// Copy data from byteData to ctx.romData
	for i := range byteData {
		ctx.romData[i] = u8(byteData[i])
	}

	ctx.header = (*romHeader)(unsafe.Pointer(&ctx.romData[0x100]))
	ctx.header.title[15] = 0 // Null-terminate the title

	fmt.Printf("Cartridge Loaded:\n")
	fmt.Printf("\t Title    : %s\n", ctx.header.title)
	fmt.Printf("\t Type     : %2.2X (%s)\n", ctx.header.typeCode, cartTypeName())
	fmt.Printf("\t ROM Size : %d KB\n", 32<<ctx.header.romSize)
	fmt.Printf("\t RAM Size : %2.2X\n", ctx.header.ramSize)
	fmt.Printf("\t LIC Code : %2.2X (%s)\n", ctx.header.licCode, cartLicName())
	fmt.Printf("\t ROM Vers : %2.2X\n", ctx.header.version)

	var x u16
	for i := 0x0134; i <= 0x014C; i++ {
		x = x - u16(ctx.romData[i]) - 1
	}

	fmt.Printf("\t Checksum : %2.2X (%s)\n", ctx.header.checksum, func() string {
		if x&0xFF != 0 {
			return "FAILED"
		}
		return "PASSED"
	}())

	return true
}
