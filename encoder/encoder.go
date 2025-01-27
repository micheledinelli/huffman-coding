package encoder

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	"huffman/queue"
	"huffman/workers"
	"log"
	"os"
	"strings"
)

func Encode(filename string, outputFilename *string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s", filename)
	}
	defer file.Close()

	dictionaries := workers.ProcessFile(file)
	dict := make(map[string]int)
	for d := range dictionaries {
		for k, v := range d {
			dict[k] += v
		}
	}

	codes := Huffman(dict)
	Dump(codes, file, outputFilename)
}

func Decode(filename, metadata string, outputFilename *string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s", filename)
	}
	defer file.Close()

	paddingLength := make([]byte, 1)
	_, err = file.Read(paddingLength)
	if err != nil {
		log.Fatalf("Error reading padding length: %v", err)
	}

	encodedData, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading encoded data: %v", err)
	}

	metadataBytes, err := os.ReadFile(metadata)
	if err != nil {
		log.Fatalf("Error reading metadata: %v", err)
	}

	var codes map[string]string
	err = json.Unmarshal(metadataBytes, &codes)
	if err != nil {
		log.Fatalf("Error unmarshalling metadata: %v", err)
	}

	reverseCodes := make(map[string]string)
	for char, code := range codes {
		reverseCodes[code] = char
	}

	var decodedData strings.Builder
	var bitstreamBuffer strings.Builder
	for _, b := range encodedData[1:] {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			bitstreamBuffer.WriteString(fmt.Sprintf("%d", bit))
		}
	}

	bitstream := bitstreamBuffer.String()[:len(bitstreamBuffer.String())-int(paddingLength[0])]

	var codeBuffer strings.Builder
	for _, bit := range bitstream {
		codeBuffer.WriteString(string(bit))
		if char, found := reverseCodes[codeBuffer.String()]; found {
			decodedData.WriteString(char)
			codeBuffer.Reset()
		}
	}

	if outputFilename == nil || *outputFilename == "" {
		defaultOutput := "output.txt"
		outputFilename = &defaultOutput
	}

	outputFile, err := os.Create(*outputFilename)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(decodedData.String())
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
}

func Dump(codes map[string]string, 	file *os.File, outputFilename *string) error {
	if outputFilename == nil || *outputFilename == "" {
		outfile := file.Name() + ".bin"
		outputFilename = &outfile
	}

	f, err := os.OpenFile(*outputFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(file)
	var encodedBits strings.Builder
	file.Seek(0, 0)
	for scanner.Scan() {
		line := scanner.Text()
		line += "\n"
		for _, c := range line {
			encodedBits.WriteString(codes[string(c)])
		}
	}

	bitstream := encodedBits.String()
	paddingLength := (8 - len(bitstream)%8) % 8
	for i := 0; i < paddingLength; i++ {
		bitstream += "0"
	}

	var byteArray bytes.Buffer
	for i := 0; i < len(bitstream); i += 8 {
		byteChunk := bitstream[i : i+8]
		var byteValue byte
		for j := 0; j < 8; j++ {
			byteValue = byteValue << 1
			if byteChunk[j] == '1' {
				byteValue |= 1
			}
		}
		byteArray.WriteByte(byteValue)
	}

	_, err = f.Write([]byte{byte(paddingLength)})
	if err != nil {
		return fmt.Errorf("error writing padding length: %v", err)
	}

	_, err = f.Write(byteArray.Bytes())
	if err != nil {
		return fmt.Errorf("error writing encoded bits: %v", err)
	}

	metaOutfile := file.Name() + ".metadata"
	metaFile, err := os.Create(metaOutfile)
	if err != nil {
		return fmt.Errorf("error creating metadata file: %v", err)
	}
	defer metaFile.Close()

	metadataBytes, err := json.Marshal(codes)
	if err != nil {
		return fmt.Errorf("error marshalling metadata: %v", err)
	}

	_, err = metaFile.Write(metadataBytes)
	if err != nil {
		return fmt.Errorf("error writing metadata: %v", err)
	}

	return nil
}

func Huffman(dic map[string]int) map[string]string {
	pq := make(queue.PriorityQueue, len(dic))

	i := 0
	for key, value := range dic {
		pq[i] = &queue.Node{
			Value:    key,
			Priority: value,
			Index:    i,
			Left:     nil,
			Right:    nil,
		}
		i++
	}
	heap.Init(&pq)

	for i := 0; i < len(dic)-1; i++ {
		node1 := heap.Pop(&pq).(*queue.Node)
		node2 := heap.Pop(&pq).(*queue.Node)

		node := &queue.Node{
			Value:    "@@@", // internal node with arbitrary value
			Priority: node1.Priority + node2.Priority,
			Index:    pq.Len(),
			Left:     node1,
			Right:    node2,
		}

		heap.Push(&pq, node)
	}

	codes := encodeTree(pq.Pop().(*queue.Node))
	return codes
}

func encodeTree(root *queue.Node) map[string]string {
	codes := make(map[string]string)
	var visit func(root *queue.Node, code string)
	visit = func(root *queue.Node, code string) {
		if root == nil {
			return
		}

		// We reached a leaf
		if root.Left == nil && root.Right == nil {
			codes[root.Value] = code
		}

		visit(root.Left, code+"0")
		visit(root.Right, code+"1")
	}

	visit(root, "")
	return codes
}
