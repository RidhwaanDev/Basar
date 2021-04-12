package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Job struct {
	JobStatus int // 0 for waiting, 1 for running, 2 for done
	FileName  string
	FileData  []byte
}

type ClientUpdate struct {
	Status int
}

type Ticket struct {
	Id       string
	FileName string
}

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.FlushAll(ctx)

	// ping redis to make sure its up
}

// key string, value JSON of job struct
func SubmitJob(key string, job Job) error {
	value, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func MarkAsComplete(key string) {
	jobBytes, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		fmt.Printf("%s does not exist\n", key)
		return
	} else if err != nil {
		panic(err)
	} else {
		var someJob Job
		err := json.Unmarshal(jobBytes, &someJob)
		if err != nil {
			panic(err)
		}
		someJob.JobStatus = 2
		SubmitJob(key, someJob)
	}
}

func GetJob(key string) *Job {
	jobBytes, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		fmt.Printf("%s does not exist\n", key)
		return nil
	} else if err != nil {
		panic(err)
	} else {
		var finalJob Job
		err := json.Unmarshal(jobBytes, &finalJob)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%s : %v\n", key, finalJob)
		return &finalJob
	}
}
