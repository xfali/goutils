/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/28
 * @time 17:32
 * @version V1.0
 * Description: 
 */

package id

import "time"

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

type SnowFlake struct {
    /* 上次生产id时间戳 */
    lastTimestamp time.Duration
    // 0，并发控制
    sequence int64
    //工作节点id
    workerId int64
    // 数据中心ID
    datacenterId int64
}
//
//func NewSnowFlake() *SnowFlake {
//    return &SnowFlake {
//        lastTimestamp: -1,
//        sequence: 0,
//        workerId:
//    }
//}
//
//
//public IdUtil(){
//this.datacenterId = getDatacenterId(maxDatacenterId);
//this.workerId = getMaxWorkerId(datacenterId, maxWorkerId);
//}

/**
 * @param workerId
 *            工作机器ID
 * @param datacenterId
 *            序列号
 */
//public IdUtil(long workerId, long datacenterId) {
//if (workerId > maxWorkerId || workerId < 0) {
//throw new IllegalArgumentException(String.format("worker Id can't be greater than %d or less than 0", maxWorkerId));
//}
//if (datacenterId > maxDatacenterId || datacenterId < 0) {
//throw new IllegalArgumentException(String.format("datacenter Id can't be greater than %d or less than 0", maxDatacenterId));
//}
//this.workerId = workerId;
//this.datacenterId = datacenterId;
//}

/**
 * 获取下一个ID
 *
 * @return
 */
//public synchronized long nextId() {
//long timestamp = timeGen();
//if (timestamp < lastTimestamp) {
//throw new RuntimeException(String.format(
//"Clock moved backwards.  Refusing to generate id for %d milliseconds", lastTimestamp - timestamp));
//}
//
//if (lastTimestamp == timestamp) {
//// 当前毫秒内，则+1
//sequence = (sequence + 1) & sequenceMask;
//if (sequence == 0) {
//// 当前毫秒内计数满了，则等待下一秒
//timestamp = tilNextMillis(lastTimestamp);
//}
//} else {
//sequence = 0L;
//}
//lastTimestamp = timestamp;
//// ID偏移组合生成最终的ID，并返回ID
//long nextId = ((timestamp - twepoch) << timestampLeftShift) | (datacenterId << datacenterIdShift)
//| (workerId << workerIdShift) | sequence;
//
//return nextId;
//}
//
//private long tilNextMillis(final long lastTimestamp) {
//long timestamp = this.timeGen();
//while (timestamp <= lastTimestamp) {
//timestamp = this.timeGen();
//}
//return timestamp;
//}
//
//private long timeGen() {
////		Calendar cal = Calendar.getInstance();
////		long time = cal.getTime().getTime();
////		cal.add(Calendar.YEAR, 69);
////		time = cal.getTime().getTime();
////		return time;
//return System.currentTimeMillis();
//}
//
///**
// * <p>
// * 获取 maxWorkerId
// * </p>
// */
//protected static long getMaxWorkerId(long datacenterId, long maxWorkerId) {
//StringBuffer mpid = new StringBuffer();
//mpid.append(datacenterId);
//String name = ManagementFactory.getRuntimeMXBean().getName();
//if (!name.isEmpty()) {
///*
// * GET jvmPid
// */
//mpid.append(name.split("@")[0]);
//}
///*
// * MAC + PID 的 hashcode 获取16个低位
// */
//return (mpid.toString().hashCode() & 0xffff) % (maxWorkerId + 1);
//}
//
///**
// * <p>
// * 数据标识id部分
// * </p>
// */
//protected static long getDatacenterId(long maxDatacenterId) {
//long id = 0L;
//try {
//InetAddress ip = InetAddress.getLocalHost();
//NetworkInterface network = NetworkInterface.getByInetAddress(ip);
//if (network == null) {
//id = 1L;
//} else {
//byte[] mac = network.getHardwareAddress();
//id = ((0x000000FF & (long) mac[mac.length - 1])
//| (0x0000FF00 & (((long) mac[mac.length - 2]) << 8))) >> 6;
//id = id % (maxDatacenterId + 1);
//}
//} catch (Exception e) {
//System.out.println(" getDatacenterId: " + e.getMessage());
//}
//return id;
//}
//
//public static String id2String(long id, int bit) {
//String bitStr = "%0" + bit + "d";
//return String.format(bitStr, id);
//}
//
//public static long trans2Timestamp(long id) {
//return twepoch + (id >> timestampLeftShift);
//}
//
//private static final char compress_digit[] = {
//'0', '1', '2', '3', '4','5', '6', '7', '8', '9',
//'A', 'B', 'C', 'D', 'E','F', 'G', 'H', 'I', 'J',
//'K', 'L', 'M', 'N', 'O','P', 'Q', 'R', 'S', 'T',
//'U', 'V', 'W', 'X', 'Y','Z'
//};
//private static final int uncompress_digit[] = {
//0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
//0, 0, 0, 0, 0,0, 0, //占位
//10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
//20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
//30, 31, 32, 33, 34, 35
//};
//
//private static final char compress_lower_digit[] = {
//'0', '1', '2', '3', '4','5', '6', '7', '8', '9',
//'A', 'B', 'C', 'D', 'E','F', 'G', 'H', 'I', 'J',
//'K', 'L', 'M', 'N', 'O','P', 'Q', 'R', 'S', 'T',
//'U', 'V', 'W', 'X', 'Y','Z',
//'a', 'b', 'c', 'd','e', 'f', 'g', 'h', 'i','j',
//'k', 'l', 'm', 'n','o', 'p', 'q', 'r', 's','t',
//'u', 'v', 'w', 'x','y', 'z'
//};
//private static final int uncompress_lower_digit[] = {
//0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
//0, 0, 0, 0, 0,0, 0, //占位
//10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
//20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
//30, 31, 32, 33, 34, 35,
//0, 0,0, 0, 0, 0,
//36, 37, 38, 39, 40, 41, 42, 43, 44, 45,
//46, 47, 48, 49, 50, 51, 52, 53, 54, 55,
//56, 57, 58, 59, 60, 61
//};
//
///**
// * 将id压缩成包含数字和字母的id，最多13位（69年）
// * @param id
// * @return
// */
//public static String innerCompress2String(char[]compressDigit, long id) {
//return innerCompress2String(compressDigit, id, 0);
//}
//
///**
// * 将id压缩成包含数字和大写字母的id，自动补齐size位数（注意大于size时不会截断）
// * @param id
// * @param size
// * @return
// */
//public static String innerCompress2String(char[]compressDigit, long id, int size) {
//StringBuffer buf = new StringBuffer();
//long quotient, compl, old = id;
//do {
//quotient = old / compressDigit.length;
//compl = old % compressDigit.length;
//buf.insert(0, compressDigit[(int)compl]);
//old = quotient;
//} while (quotient > 0);
//int lenth = buf.length();
//while (size-- > lenth) {
//buf.insert(0, compressDigit[0]);
//}
//return buf.toString();
//}
//
///**
// * 将压缩的字符id转换为数字id
// * @param strId
// * @return
// */
//public static long innerUncompress2Long(char[]compressDigit, int[]uncompressDigit, String strId) {
//int len = strId.length(), cur = 1;
//long value = 0, mul = 1;
//do {
//value += uncompressDigit[(int)strId.charAt(len - cur) - compressDigit[0]] * mul;
//mul *= compressDigit.length;
//} while (len > cur++);
//return value;
//}
//
///**
// * 将id压缩成包含数字和大写字母的id，最多13位（69年）
// * @param id
// * @return
// */
//public static String compress2String(long id) {
//return innerCompress2String(compress_digit, id);
//}
//
///**
// * 将id压缩成包含数字和大写字母的id，自动补齐size位数（注意大于size时不会截断）
// * @param id
// * @param size
// * @return
// */
//public static String compress2String(long id, int size) {
//return innerCompress2String(compress_digit, id, size);
//}
//
///**
// * 将压缩的包含数字和大写字母的id转换为数字id
// * @param strId
// * @return
// */
//public static long uncompress2Long(String strId) {
//return innerUncompress2Long(compress_digit, uncompress_digit, strId);
//}
//
///**
// * 将id压缩成包含数字和大小写字母的id，最多11位（69年）
// * @param id
// * @return
// */
//public static String compress2StringUL(long id) {
//return innerCompress2String(compress_lower_digit, id);
//}
//
///**
// * 将id压缩成包含数字和大小写字母的id，自动补齐size位数（注意大于size时不会截断）
// * @param id
// * @param size
// * @return
// */
//public static String compress2StringUL(long id, int size) {
//return innerCompress2String(compress_lower_digit, id, size);
//}
//
///**
// * 将压缩成包含数字和大小写字母的id转换为数字id
// * @param strId
// * @return
// */
//public static long uncompress2LongUL(String strId) {
//return innerUncompress2Long(compress_lower_digit, uncompress_lower_digit, strId);
//}
