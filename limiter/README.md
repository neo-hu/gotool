# NewWriter
限制速率的write， bytesPerSec: 每秒最多写入bytesPerSec，如果bytesPerSec <=0 不限制速率
`
w := NewWriter(fd, 512, 10*1024)
// 每秒最大写入512个字节
`
