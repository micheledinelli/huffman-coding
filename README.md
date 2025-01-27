# huffman-coding

<div align="center">

![go](https://img.shields.io/badge/Go-00ADD8.svg?style=plain&logo=Go&logoColor=white)
![go](https://img.shields.io/github/go-mod/go-version/micheledinelli/huffman-coding?style=flat)

</div>

Command line application to try Huffman code for lossless file compression.

> [!IMPORTANT]
> This is still a work in progress. The whole decompression has to be added and there is room for performance improvement.
> 
> The cli still supports `txt` files only.


Huffman code is a type of optimal prefix code.
The process of finding or using such a code is Huffman coding, an algorithm developed by [David A. Huffman](https://en.wikipedia.org/wiki/David_A._Huffman). 
Huffman algorithm is a greedy algorithm that uses a frequency-sorted binary tree and generates "prefix-free codes", in other words the bit string representing some particular symbol is never a prefix of the bit string representing any other symbol.

### How it works

Suppose that we have to compress a file that contains the following
sentence: "Nel mezzo del cammin", the characters are:

- z: 2 (11.76%)
- o: 1 (5.88%)
- l: 2 (11.76%)
- d: 1 (5.88%)
- e: 3 (17.65%)
- c: 1 (5.88%)
- a: 1 (5.88%)
- i: 1 (5.88%)
- n: 2 (11.76%)
- m: 3 (17.65%)

Characters are lower cased and the space character is not included, this is only for example purposes.

Note that

- ASCII encoding takes 8 bit for each char (8n bit)
- Alphabet encoding takes 4 bit for each char (4n bit)
  * 0000: z
  * 0001: o
  * and so on
- Huffman?

<div align="center">

![huffman](img/huffman-dark.webp)

</div>

The example leads to the following dictionary:
- i: 000
- l: 001
- z: 010
- n: 011
- d: 1000
- o: 1001
- c: 1010
- a: 1011
- e: 110
- m: 111

The number of bits required is

```math
\begin{aligned}
 ((11.76\% *  3) + (5.88\% * 4) + (11.76\% * 3) + (5.88\% * 4) &+ \\ (17.65\% * 3) + (5.88\% * 4) + (5.88\% * 4) + (5.88\% * 3) &+ \\
(11.76\% * 3) + (17.65\% * 3))*n &= \\
   &= 3.2346n \\
   3.2346n < 4n < 8n
\end{aligned}
```

### How to try

```console
git clone https://github.com/micheledinelli/huffman-coding.git
```

```console
go build -o huffman
```

Once built list all the available commands with
```console
./huffman --help
```

Along with the source code there is also a file named `divine-comedy-it.txt` which is the divine comedy by Dante Alighieri in Italian in `.txt` format. It can be used to test the Huffman code.

```console
./huffman c divine-comedy-it.txt
```

`divine-comedy-it.txt` is compressed in `divine-comedy-it.txt.bin` which is 280KB almost half of the size of the original one (543KB).
![results](img/results.webp)
The cli also outputs metadata (`divine-comedy-it.txt.metadata`), the unique encoding for a specific file in order to be able to decompress it later.
