package bittorrentfile

type BitTorrentInfoFile struct {
	Length int
	Path   []string
}

type BittorrentHash = []byte

type BitTorrentInfo struct {
	Name       string
	Files      []BitTorrentInfoFile
	PieceLengh int
	Pieces     []BittorrentHash
}

type BitTorrent struct {
	Announce     string
	AnnounceList []string
	Comment      string
	CreatedBy    string
	CreationDate int
	Encoding     string
	Publisher    string
	PublisherUrl string
	Info         BitTorrentInfo
}

func GetBitTorrent(dict map[string]interface{}) BitTorrent {
	var result BitTorrent
	announce, ok := dict["announce"]
	if ok {
		result.Announce = getStringFromInterface(announce)
	}

	announceListInterface, ok := dict["announce-list"]
	if ok {
		list := getListFromInterface(announceListInterface)
		announceList := make([]string, len(list))
		for index, announceItem := range list {
			innerList := getListFromInterface(announceItem)
			announceList[index] = getStringFromInterface(innerList[0])
		}

		result.AnnounceList = announceList
	}

	comment, ok := dict["comment"]
	if ok {
		result.Comment = getStringFromInterface(comment)
	}

	encoding, ok := dict["encoding"]
	if ok {
		result.Encoding = getStringFromInterface(encoding)
	}

	createdBy, ok := dict["created by"]
	if ok {
		result.CreatedBy = getStringFromInterface(createdBy)
	}

	createionDate, ok := dict["creation date"]
	if ok {
		result.CreationDate = getIntFromInterface(createionDate)
	}

	publisher, ok := dict["publisher"]
	if ok {
		result.Publisher = getStringFromInterface(publisher)
	}

	publisherUrl, ok := dict["publisher-url"]
	if ok {
		result.PublisherUrl = getStringFromInterface(publisherUrl)
	}

	info, ok := dict["info"]
	if ok {
		result.Info = getInfoFromDictionary(getStringDictionary(info))
	}

	return result
}

func getInfoFromDictionary(infoDictionary map[string]interface{}) BitTorrentInfo {
	var info BitTorrentInfo
	name, ok := infoDictionary["name"]
	if ok {
		info.Name = getStringFromInterface(name)
	}

	files, ok := infoDictionary["files"]
	if ok {
		fileList := getListFromInterface(files)
		info.Files = getFilesFromList(fileList)
	}

	pieceLength, ok := infoDictionary["piece length"]
	if ok {
		info.PieceLengh = getIntFromInterface(pieceLength)
	}

	pieces, ok := infoDictionary["pieces"]
	if ok {
		info.Pieces = getPeacesFromInterface(pieces)
	}

	return info
}

func getPeacesFromInterface(inter interface{}) []BittorrentHash {

	list, ok := inter.([]byte)
	if !ok {
		panic("Not a list of bytes")

	}

	length := len(list) / 20
	pieces := make([]BittorrentHash, length)
	for i := 0; i < length; i++ {
		pieces[i] = list[i*20 : (i+1)*20]
	}

	return pieces
}

func getStringDictionary(inter interface{}) map[string]interface{} {
	dict, ok := inter.(map[string]interface{})
	if ok {
		return dict
	}

	panic("Not a dictionary")
}

func getStringFromInterface(inter interface{}) string {
	bytes, ok := inter.([]byte)
	if ok {
		return string(bytes)
	}

	panic("Not a byte array")
}

func getIntFromInterface(inter interface{}) int {
	integer, ok := inter.(int)
	if ok {
		return integer
	}

	panic("Not an integer")
}

func getListFromInterface(inter interface{}) []interface{} {
	list, ok := inter.([]interface{})
	if ok {
		return list
	}

	panic("Not a list")
}

func getFilesFromList(list []interface{}) []BitTorrentInfoFile {
	resultList := make([]BitTorrentInfoFile, len(list))

	for index, element := range list {
		dict := getStringDictionary(element)
		var file BitTorrentInfoFile
		path, ok := dict["path"]
		if ok {
			paths := getListFromInterface(path)
			stringPaths := make([]string, len(paths))
			for index, path := range paths {
				stringPaths[index] = getStringFromInterface(path)
			}
			file.Path = stringPaths
		}

		lenght, ok := dict["length"]
		if ok {
			file.Length = getIntFromInterface(lenght)
		}

		resultList[index] = file
	}

	return resultList
}
