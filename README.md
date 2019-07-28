## SeqSrv
### Getting Started
1.Run server
```
./seq_srv -node-id=1
```
2.Request http://127.0.0.1:8181/seq/next

|     parameter   |   description    |
| --------------- | ---------:|
|-bind-addr|SeqSrv listen address(default:0.0.0:8181)|
|-node-id|This node id(default:1, note:this id must be globally unique)|
|-seq-file|Sequence number persistence file path(default:max_seq)|
|-init-cur-seq|Sequence number initial value(default:0)|
|-init-max-seq|Sequence number max value(default:1000)|
|-step|Max sequence number increment increasing step(default:1000)|