/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package idUtil

import (
    "errors"
    "fmt"
    "hash/crc32"
    "net"
    "os"
    "strconv"
    "sync"
    "time"
)

const (
    // 时间起始标记点，作为基准，一般取系统的最近时间（一旦确定不能变动）
    twepoch = 1512443299165
    // 机器标识位数
    workerIdBits = 5
    // 数据中心标识位数
    datacenterIdBits = 5
    // 机器ID最大值
    maxWorkerId = -1 ^ (-1 << workerIdBits)
    // 数据中心ID最大值
    maxDatacenterId = -1 ^ (-1 << datacenterIdBits)
    // 毫秒内自增位
    sequenceBits = 12
    // 机器ID偏左移12位
    workerIdShift = sequenceBits
    // 数据中心ID左移17位
    datacenterIdShift = sequenceBits + workerIdBits
    // 时间毫秒左移22位
    timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
    // 最大毫秒自增
    sequenceMask = -1 ^ (-1 << sequenceBits)
)


/**
 * @author Xiongfa Li
 * 唯一ID生成器，从2017年12月5日开始，能够使用68年左右，最大占19位字符
 */
type SnowFlake struct {
    /* 上次生产id时间戳 */
    lastTimestamp int64
    // 0，并发控制
    sequence int64
    //工作节点id
    workerId int64
    // 数据中心ID
    datacenterId int64

    lock sync.Mutex
}

type SFId int64
type SFStrId string

func NewSnowFlake() *SnowFlake {
    ret := &SnowFlake{
        lastTimestamp: -1,
        sequence:      0,
    }
    ret.datacenterId = getDatacenterId(maxDatacenterId)
    ret.workerId = getMaxWorkerId(ret.datacenterId, maxWorkerId)
    return ret
}

/**
 * @param workerId
 *            工作机器ID
 * @param datacenterId
 *            序列号
 */
func NewSnowFlakeWithId(workerId int64, dataCenterId int64) *SnowFlake {
    if workerId > maxWorkerId || workerId < 0 {
        panic(fmt.Sprintf("worker Id can't be greater than %d or less than 0", maxWorkerId))
    }
    if dataCenterId > maxDatacenterId || dataCenterId < 0 {
        panic(fmt.Sprintf("datacenter Id can't be greater than %d or less than 0", maxDatacenterId))
    }
    return &SnowFlake{
        lastTimestamp: -1,
        sequence:      0,
        workerId:      workerId,
        datacenterId:  dataCenterId,
    }
}

/**
 * 获取下一个ID
 *
 * @return
 */
func (sf *SnowFlake) NextId() (SFId, error) {
    sf.lock.Lock()
    defer sf.lock.Unlock()

    timestamp := sf.timeGen()
    if timestamp < sf.lastTimestamp {
        return -1, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate idUtil for %d milliseconds", sf.lastTimestamp-timestamp))
    }

    if sf.lastTimestamp == timestamp {
        // 当前毫秒内，则+1
        sf.sequence = (sf.sequence + 1) & sequenceMask
        if sf.sequence == 0 {
            // 当前毫秒内计数满了，则等待下一秒
            timestamp = sf.tilNextMillis(sf.lastTimestamp)
        }
    } else {
        sf.sequence = 0
    }
    sf.lastTimestamp = timestamp
    // ID偏移组合生成最终的ID，并返回ID
    nextId := ((timestamp - twepoch) << timestampLeftShift) | (sf.datacenterId << datacenterIdShift) | (sf.workerId << workerIdShift) | sf.sequence

    return SFId(nextId), nil
}

func (sf *SnowFlake) tilNextMillis(lastTimestamp int64) int64 {
    timestamp := sf.timeGen()
    for timestamp <= lastTimestamp {
        timestamp = sf.timeGen()
    }
    return timestamp
}

func (sf *SnowFlake) timeGen() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

/**
* <p>
* 获取 maxWorkerId
* </p>
*/
func getMaxWorkerId(dataCenterId int64, maxWorkerId int64) int64 {
    pid := os.Getpid()
    str := fmt.Sprintf("%d%d", dataCenterId, pid)
    return (int64(crc32.ChecksumIEEE([]byte(str))) & 0xffff) % (maxWorkerId + 1)
}

//
///**
// * <p>
// * 数据标识id部分
// * </p>
// */
func getDatacenterId(maxDatacenterId int) int64 {
    intes, err := net.Interfaces()
    if err == nil {
        for i := 0; i < len(intes); i++ {
            mac := intes[i].HardwareAddr
            macLen := len(mac)
            if macLen == 0 {
                continue
            }
            id := ((0x000000FF & int64(mac[macLen-1])) | (0x0000FF00 & (int64(mac[macLen-2]) << 8))) >> 6
            return id
        }
    }
    return 1
}

func (id SFId)Int64() int64 {
    return int64(id)
}

func (id SFId)LimitString(bit int) string {
    bitStr := "%0" + strconv.Itoa(bit) + "d"
    return fmt.Sprintf(bitStr, id)
}

func (id SFId)String() string {
    return strconv.FormatInt(int64(id), 10)
}

func (id SFId)Timestamp() time.Duration {
    return time.Duration(twepoch + (id >> timestampLeftShift)) * time.Millisecond
}

func (id SFId)Compress() SFStrId {
    return SFStrId(Compress2StringUL(int64(id)))
}

func (sid SFStrId)String() string {
    return string(sid)
}

func (sid SFStrId)UnCompress() SFId {
    return SFId(Uncompress2LongUL(string(sid)))
}

func (id SFId)Time() time.Time {
    t := int64(id.Timestamp())
    return time.Unix(0, t)
}

var compress_digit = []byte{
    '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
    'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
    'U', 'V', 'W', 'X', 'Y', 'Z',
}

var uncompress_digit = []byte{
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    0, 0, 0, 0, 0, 0, 0, //占位
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
    20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
    30, 31, 32, 33, 34, 35,
}

var compress_lower_digit = []byte{
    '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
    'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
    'U', 'V', 'W', 'X', 'Y', 'Z',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
    'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
    'u', 'v', 'w', 'x', 'y', 'z',
}

var uncompress_lower_digit = []byte{
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    0, 0, 0, 0, 0, 0, 0, //占位
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
    20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
    30, 31, 32, 33, 34, 35,
    0, 0, 0, 0, 0, 0,
    36, 37, 38, 39, 40, 41, 42, 43, 44, 45,
    46, 47, 48, 49, 50, 51, 52, 53, 54, 55,
    56, 57, 58, 59, 60, 61,
}

/**
* 将id压缩成包含数字和字母的id，最多13位（69年）
* @param idUtil
* @return
*/
func innerCompress2String(compressDigit []byte, id int64) string {
    return innerCompress2String2(compressDigit, id, 0)
}

/**
* 将id压缩成包含数字和大写字母的id，自动补齐size位数（注意大于size时不会截断）
* @param idUtil
* @param size
* @return
*/
func innerCompress2String2(compressDigit []byte, id int64, size int) string {
    //int64最大值为19位，如果compressDigit为2进制，则位数转换为 19*10/2 【19乘以（10进制/2进制）】=95，所以使用100足够满足需求
    buf := make([]byte, 100)
    var quotient, compl, old int64 = 0, 0, id
    clen := int64(len(compressDigit))
    index := len(buf) - 1
    for {
        quotient = old / clen
        compl = old % clen
        buf[index] = compressDigit[int(compl)]
        index--
        old = quotient
        if quotient == 0 {
            break
        }
    }
    lenth := len(buf) - index - 1
    for size > lenth {
        size--
        buf[index] = compressDigit[0]
        index--
    }
    buf = buf[index+1:]
    return string(buf)
}

/**
* 将压缩的字符id转换为数字id
* @param strId
* @return
*/
func innerUncompress2Long(compressDigit []byte, uncompressDigit []byte, strId string) int64 {
    strbyte := []byte(strId)
    byteLen := len(strbyte)
    cur := 1
    var value int64 = 0
    var mul int64 = 1
    for {
        v1 := int64(uncompressDigit[int(strbyte[byteLen-cur]-compressDigit[0])])
        value += v1 * mul
        mul *= int64(len(compressDigit))
        if byteLen <= cur {
            break
        }
        cur++
    }
    return value
}

/**
* 将id压缩成包含数字和大写字母的id，最多13位（69年）
* @param idUtil
* @return
*/
func Compress2String(id int64) string {
    return innerCompress2String(compress_digit, id)
}

/**
* 将id压缩成包含数字和大写字母的id，自动补齐size位数（注意大于size时不会截断）
* @param idUtil
* @param size
* @return
*/
func Compress2String2(id int64, size int) string {
    return innerCompress2String2(compress_digit, id, size)
}

/**
* 将压缩的包含数字和大写字母的id转换为数字id
* @param strId
* @return
*/
func Uncompress2Long(strId string) int64 {
    return innerUncompress2Long(compress_digit, uncompress_digit, strId)
}

/**
* 将id压缩成包含数字和大小写字母的id，最多11位（69年）
* @param idUtil
* @return
*/
func Compress2StringUL(id int64) string {
    return innerCompress2String(compress_lower_digit, id)
}

/**
* 将id压缩成包含数字和大小写字母的id，自动补齐size位数（注意大于size时不会截断）
* @param idUtil
* @param size
* @return
*/
func Compress2StringUL2(id int64, size int) string {
    return innerCompress2String2(compress_lower_digit, id, size)
}

/**
* 将压缩成包含数字和大小写字母的id转换为数字id
* @param strId
* @return
*/
func Uncompress2LongUL(strId string) int64 {
    return innerUncompress2Long(compress_lower_digit, uncompress_lower_digit, strId)
}
