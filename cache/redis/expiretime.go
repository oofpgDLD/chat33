package redis

//key过期时间 所有key都设置了过期时间
const EXPIRETimeHour = 60 * 60 //一小时

const EXPIRETimeHalfDay = 12 * 60 * 60 //半天

const EXPIRETimeDay = 24 * 60 * 60 //一天

const EXPIRETimeWeek = 7 * 24 * 60 * 60 //一周

const EXPIRETimeMonth = 30 * 24 * 60 * 60 //一周
//锁超时时间
const LockTimeout = 3 //3秒
