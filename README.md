ethdiff
=======

Finds last common block between two Ethereum-compatible blockchains. This block number can be used to reset corrupted node strayed from consensus.

## Usage example

#### Find latest common block between Ethereum and Ethereum Classic.

```
[21:54:28] dt1:~> ~/go/bin/ethdiff  https://cloudflare-eth.com https://www.ethercluster.com/etc
2022/02/26 21:54:29.932755 diff.go:62: highestCommonBlock = 0xdf0d69 (14617961)
2022/02/26 21:54:29.932955 diff.go:64: highestCommonBlock (safe value) = 0xdf0ca1 (14617761)
2022/02/26 21:54:30.392783 diff.go:77: searchFunc(0x6f8651) = true
2022/02/26 21:54:30.655261 diff.go:77: searchFunc(0x37c328) = true
2022/02/26 21:54:30.894623 diff.go:77: searchFunc(0x1be194) = false
2022/02/26 21:54:31.215508 diff.go:77: searchFunc(0x29d25e) = true
2022/02/26 21:54:31.518864 diff.go:77: searchFunc(0x22d9f9) = true
2022/02/26 21:54:32.010745 diff.go:77: searchFunc(0x1f5dc7) = true
2022/02/26 21:54:32.551594 diff.go:77: searchFunc(0x1d9fae) = true
2022/02/26 21:54:32.900448 diff.go:77: searchFunc(0x1cc0a1) = false
2022/02/26 21:54:33.431563 diff.go:77: searchFunc(0x1d3028) = false
2022/02/26 21:54:33.736333 diff.go:77: searchFunc(0x1d67eb) = true
2022/02/26 21:54:34.263017 diff.go:77: searchFunc(0x1d4c0a) = true
2022/02/26 21:54:34.796005 diff.go:77: searchFunc(0x1d3e19) = false
2022/02/26 21:54:35.760275 diff.go:77: searchFunc(0x1d4512) = false
2022/02/26 21:54:36.073528 diff.go:77: searchFunc(0x1d488e) = false
2022/02/26 21:54:36.367861 diff.go:77: searchFunc(0x1d4a4c) = false
2022/02/26 21:54:36.712629 diff.go:77: searchFunc(0x1d4b2b) = false
2022/02/26 21:54:37.211520 diff.go:77: searchFunc(0x1d4b9b) = false
2022/02/26 21:54:37.732155 diff.go:77: searchFunc(0x1d4bd3) = false
2022/02/26 21:54:38.244955 diff.go:77: searchFunc(0x1d4bef) = false
2022/02/26 21:54:39.534426 diff.go:77: searchFunc(0x1d4bfd) = false
2022/02/26 21:54:40.017071 diff.go:77: searchFunc(0x1d4c04) = true
2022/02/26 21:54:40.279099 diff.go:77: searchFunc(0x1d4c01) = true
2022/02/26 21:54:40.447962 diff.go:77: searchFunc(0x1d4bff) = false
2022/02/26 21:54:40.626779 diff.go:77: searchFunc(0x1d4c00) = true
0x1d4bff
```

#### Check if Polygon node is in line with official RPC endpoint

```
[21:57:17] dt1:~> ~/go/bin/ethdiff https://polygon-rpc.com http://matic-full-node-3.mysterium.network:8545/
2022/02/26 21:57:42.740567 diff.go:62: highestCommonBlock = 0x183072f (25364271)
2022/02/26 21:57:42.740611 diff.go:64: highestCommonBlock (safe value) = 0x1830667 (25364071)
2022/02/26 21:57:42.837655 diff.go:77: searchFunc(0xc18334) = false
2022/02/26 21:57:43.006258 diff.go:77: searchFunc(0x12244ce) = false
2022/02/26 21:57:43.189835 diff.go:77: searchFunc(0x152a59b) = false
2022/02/26 21:57:43.379055 diff.go:77: searchFunc(0x16ad602) = false
2022/02/26 21:57:43.585609 diff.go:77: searchFunc(0x176ee35) = false
2022/02/26 21:57:43.880756 diff.go:77: searchFunc(0x17cfa4f) = false
2022/02/26 21:57:44.011970 diff.go:77: searchFunc(0x180005c) = false
2022/02/26 21:57:44.265298 diff.go:77: searchFunc(0x1818362) = false
2022/02/26 21:57:44.357502 diff.go:77: searchFunc(0x18244e5) = false
2022/02/26 21:57:44.558605 diff.go:77: searchFunc(0x182a5a7) = false
2022/02/26 21:57:44.627808 diff.go:77: searchFunc(0x182d608) = false
2022/02/26 21:57:45.050091 diff.go:77: searchFunc(0x182ee38) = false
2022/02/26 21:57:45.226193 diff.go:77: searchFunc(0x182fa50) = false
2022/02/26 21:57:45.557199 diff.go:77: searchFunc(0x183005c) = false
2022/02/26 21:57:45.973571 diff.go:77: searchFunc(0x1830362) = false
2022/02/26 21:57:46.077013 diff.go:77: searchFunc(0x18304e5) = false
2022/02/26 21:57:46.274309 diff.go:77: searchFunc(0x18305a7) = false
2022/02/26 21:57:46.469935 diff.go:77: searchFunc(0x1830608) = false
2022/02/26 21:57:46.580897 diff.go:77: searchFunc(0x1830638) = false
2022/02/26 21:57:47.035597 diff.go:77: searchFunc(0x1830650) = false
2022/02/26 21:57:47.180541 diff.go:77: searchFunc(0x183065c) = false
2022/02/26 21:57:47.308357 diff.go:77: searchFunc(0x1830662) = false
2022/02/26 21:57:47.453436 diff.go:77: searchFunc(0x1830665) = false
2022/02/26 21:57:48.474402 diff.go:77: searchFunc(0x1830667) = false
0x1830667
```

If last common block number is not equal to latest tested block (`highestCommonBlock (safe value)`) and node syncronization is stuck, you can try to fix it like this:

1. Login into corresponding node server shell.
2. Attach to bor like this: `bor attach /mnt/data/.bor/data/bor.ipc`
3. Rewind blockchain to latest common block: `debug.setHead("0x1830667")`

## Synopsis

```
ethdiff [options...] <left RPC address> <right RPC address>

Options:
  -offset uint
    	head backward offset for safe block retrieval (default 200)
  -retries uint
    	number of retries for RPC calls (default 3)
  -total-timeout duration
    	whole operation timeout (default 1m0s)
  -version
    	show program version and exit
```
