package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"torrent_client/bittorrentfile"
)

type Element uint

const (
	Dictionary Element = iota
	List
	String
	Integer
	End
)

func getLiteralElement(reader *bufio.Reader) []byte {
	lenStr, err := reader.ReadString(':')
	if err != nil {
		panic(err)
	}

	len, err := strconv.Atoi(strings.TrimRight(lenStr, ":"))
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, len)
	n, err := io.ReadFull(reader, buffer)
	if err != nil {
		panic(err)
	}

	if n != len {
		panic(n)
	}
	return buffer
}

func getIntegerElement(reader *bufio.Reader) int {
	intStr, err := reader.ReadString('e')
	if err != nil {
		panic(err)
	}

	intValue, err := strconv.Atoi(strings.TrimRight(intStr, "e"))
	if err != nil {
		panic(err)
	}

	return intValue
}

func getDictionaryElement(reader *bufio.Reader) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for {
		elementType := getNextElementType(reader)
		if elementType == End {
			break
		}

		key := string(getLiteralElement(reader)[:])

		elementType = getNextElementType(reader)

		var value interface{}
		if elementType == String {
			value = getLiteralElement(reader)
		} else if elementType == Integer {
			value = getIntegerElement(reader)
		} else if elementType == List {
			value = getListElement(reader)
		} else if elementType == Dictionary {
			value = getDictionaryElement(reader)
		}

		resultMap[key] = value
	}

	return resultMap
}

func getListElement(reader *bufio.Reader) []interface{} {
	resultList := make([]interface{}, 0, 10)
	for {
		elementType := getNextElementType(reader)
		if elementType == End {
			break
		}
		if elementType == String {
			resultList = append(resultList, getLiteralElement(reader))
		} else if elementType == Integer {
			resultList = append(resultList, getIntegerElement(reader))
		} else if elementType == List {
			resultList = append(resultList, getListElement(reader))
		} else if elementType == Dictionary {
			resultList = append(resultList, getDictionaryElement(reader))
		}
	}

	return resultList
}

func getNextElementType(reader *bufio.Reader) Element {
	nextItem, err := reader.Peek(1)
	if err != nil {
		panic(err)
	}

	if nextItem[0] > 47 && nextItem[0] < 58 {
		return String
	}

	r, _, err := reader.ReadRune()
	if r == 'd' {
		return Dictionary
	}

	if r == 'l' {
		return List
	}

	if r == 'i' {
		return Integer
	}

	if r == 'e' {
		return End
	}

	panic("Unknown type of element")
}

func main() {
	f, err := os.Open("example.torrent")
	if err != nil {
		panic(err)
	}

	br := bufio.NewReader(f)

	elementType := getNextElementType(br)
	if elementType != Dictionary {
		panic("Root element MUST be a dictionary")
	}

	dict := getDictionaryElement(br)
	bitTorrent := bittorrentfile.GetBitTorrent(dict)
	fmt.Println(bitTorrent)
}
