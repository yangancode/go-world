# 分布式锁的实现

## 简介
mysql_lock：使用MySQL的乐观锁和悲观锁实现的分布式锁

redis_lock：使用Redis set方法实现的分布式锁

redlock：使用Redis redlock算法实现的分布式锁

ectd_lock：使用etcd实现的分布式锁