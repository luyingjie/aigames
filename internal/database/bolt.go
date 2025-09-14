package database

import (
	"aigames/pkg/constants"
	"aigames/pkg/logger"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// DB 数据库连接包装器
type DB struct {
	conn *bolt.DB
	path string
}

// NewDB 创建新的数据库连接
func NewDB(dbPath string) (*DB, error) {
	// 打开数据库文件
	conn, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	db := &DB{
		conn: conn,
		path: dbPath,
	}

	// 初始化数据库存储桶
	if err := db.initBuckets(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	logger.Info("数据库连接成功: %s", dbPath)
	return db, nil
}

// initBuckets 初始化数据库存储桶
func (db *DB) initBuckets() error {
	buckets := []string{
		constants.BucketUsers,
		constants.BucketGames,
		constants.BucketRooms,
		constants.BucketAIPlayers,
		constants.BucketChats,
		constants.BucketConfigs,
	}

	return db.conn.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("创建存储桶 %s 失败: %w", bucket, err)
			}
		}
		return nil
	})
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	if db.conn != nil {
		logger.Info("关闭数据库连接")
		return db.conn.Close()
	}
	return nil
}

// Get 获取数据
func (db *DB) Get(bucket, key string, result interface{}) error {
	return db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("键 %s 不存在", key)
		}

		return json.Unmarshal(data, result)
	})
}

// Put 存储数据
func (db *DB) Put(bucket, key string, data interface{}) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		encoded, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("序列化数据失败: %w", err)
		}

		return b.Put([]byte(key), encoded)
	})
}

// Delete 删除数据
func (db *DB) Delete(bucket, key string) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		return b.Delete([]byte(key))
	})
}

// Exists 检查键是否存在
func (db *DB) Exists(bucket, key string) bool {
	exists := false
	db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			data := b.Get([]byte(key))
			exists = (data != nil)
		}
		return nil
	})
	return exists
}

// List 列出存储桶中的所有键值对
func (db *DB) List(bucket string, result interface{}) error {
	return db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		items := make(map[string]interface{})
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item interface{}
			if err := json.Unmarshal(v, &item); err != nil {
				logger.Warn("反序列化数据失败: %v", err)
				continue
			}
			items[string(k)] = item
		}

		// 将map转换为目标类型
		encoded, err := json.Marshal(items)
		if err != nil {
			return fmt.Errorf("序列化结果失败: %w", err)
		}

		return json.Unmarshal(encoded, result)
	})
}

// ListKeys 列出存储桶中的所有键
func (db *DB) ListKeys(bucket string) ([]string, error) {
	var keys []string

	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})

	return keys, err
}

// Count 统计存储桶中的记录数
func (db *DB) Count(bucket string) (int, error) {
	count := 0

	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}
		return nil
	})

	return count, err
}

// BatchPut 批量存储数据
func (db *DB) BatchPut(bucket string, items map[string]interface{}) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		for key, data := range items {
			encoded, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("序列化数据失败 (key: %s): %w", key, err)
			}

			if err := b.Put([]byte(key), encoded); err != nil {
				return fmt.Errorf("存储数据失败 (key: %s): %w", key, err)
			}
		}

		return nil
	})
}

// BatchDelete 批量删除数据
func (db *DB) BatchDelete(bucket string, keys []string) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		for _, key := range keys {
			if err := b.Delete([]byte(key)); err != nil {
				return fmt.Errorf("删除数据失败 (key: %s): %w", key, err)
			}
		}

		return nil
	})
}

// Transaction 执行事务
func (db *DB) Transaction(fn func(*bolt.Tx) error) error {
	return db.conn.Update(fn)
}

// ViewTransaction 执行只读事务
func (db *DB) ViewTransaction(fn func(*bolt.Tx) error) error {
	return db.conn.View(fn)
}

// GetBucketStats 获取存储桶统计信息
func (db *DB) GetBucketStats(bucket string) (bolt.BucketStats, error) {
	var stats bolt.BucketStats

	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}
		stats = b.Stats()
		return nil
	})

	return stats, err
}

// GetStats 获取数据库统计信息
func (db *DB) GetStats() bolt.Stats {
	return db.conn.Stats()
}

// Backup 备份数据库
func (db *DB) Backup(dest string) error {
	return db.conn.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(dest, 0600)
	})
}

// CreateIndex 创建简单的索引(通过前缀扫描实现)
func (db *DB) CreateIndex(bucket, indexBucket, field string, getValue func(interface{}) string) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		// 创建索引存储桶
		indexB, err := tx.CreateBucketIfNotExists([]byte(indexBucket))
		if err != nil {
			return err
		}

		// 获取原数据存储桶
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		// 遍历所有数据，建立索引
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item interface{}
			if err := json.Unmarshal(v, &item); err != nil {
				continue
			}

			indexKey := getValue(item)
			if indexKey != "" {
				// 存储索引: indexKey -> 原始key
				if err := indexB.Put([]byte(indexKey), k); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetByIndex 通过索引查询数据
func (db *DB) GetByIndex(bucket, indexBucket, indexKey string, result interface{}) error {
	return db.conn.View(func(tx *bolt.Tx) error {
		// 获取索引存储桶
		indexB := tx.Bucket([]byte(indexBucket))
		if indexB == nil {
			return fmt.Errorf("索引存储桶 %s 不存在", indexBucket)
		}

		// 通过索引获取原始key
		originalKey := indexB.Get([]byte(indexKey))
		if originalKey == nil {
			return fmt.Errorf("索引 %s 不存在", indexKey)
		}

		// 获取原始数据
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("存储桶 %s 不存在", bucket)
		}

		data := b.Get(originalKey)
		if data == nil {
			return fmt.Errorf("数据不存在")
		}

		return json.Unmarshal(data, result)
	})
}

// GetBoltDB 获取底层的bolt.DB实例
func (db *DB) GetBoltDB() *bolt.DB {
	return db.conn
}
