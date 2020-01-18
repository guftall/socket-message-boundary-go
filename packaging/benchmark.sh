# go test -bench=FillCurrentPacketDataAllDataInFrame -benchmem
# go test -bench=FillCurrentPacketDataNotCompleteFrame -benchmem

go test -run=NONE -benchtime=15000000x -bench=FillCurrentPacketDataAllDataInFrame -memprofile=mem.log -cpuprofile=cpu.log
go tool pprof -text -nodecount=10 ./packaging.test mem.log
go tool pprof -text -nodecount=10 ./packaging.test cpu.log