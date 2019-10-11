package bolt_test
//bolt 数据库写的操作
import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"testing"
)

func TestBolt(t *testing.T)  {

	//1.open bolt database
	db , err := bolt.Open("test.db",0644,nil)
	defer db.Close()
	if err !=nil{
		log.Panic("open database failed")
	}
	//2.fing bucket if it does not exist then create a database
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			//bucket does not exist create
			bucket , err =tx.CreateBucket([]byte("b1"))  //变量b1 通常用一个文件管理

			if err != nil {
				log.Panic("创建bucket失败了")
			}
		}
		bucket.Put([]byte("11111"),[]byte("hello"))
		bucket.Put([]byte("22222"),[]byte("world"))


		return nil
	})

	//读数据
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			log.Panic("bucket b1 不应该为空，请检查")
		}
		v1 := bucket.Get([]byte("11111"))
		v2 := bucket.Get([]byte("22222"))
		fmt.Printf("v1 : %s \n", v1)
		fmt.Printf("v2 : %s \n", v2)


		return nil
		//数据库打开 用defer 关闭
	})







}

