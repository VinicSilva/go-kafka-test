package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"github.com/VinicSilva/esquenta-imersao-go-course/infra/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	repository2 "github.com/VinicSilva/esquenta-imersao-go-course/infra/repository"
	usecase2 "github.com/VinicSilva/esquenta-imersao-go-course/usecase"
	_ "github.com/go-sql-driver/mysql"
)

func main()  {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/fullcycle")
	if err != nil {
		log.Fatalln(err)
	}
	repository := repository2.CourseMySQLRepository{Db: db}
	usecase := usecase2.CreateCourse{Repository: repository}

	var msgChan = make(chan *ckafka.Message)
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9094",
		"group.id": "appgo",
	}
	topics := []string{"courses"}
	consumer := kafka.NewConsumer(configMapConsumer, topics)
	go consumer.Consume(msgChan)

	for msg := range msgChan {
		var input usecase2.CreateCourseInputDto
		json.Unmarshal(msg.Value, &input)
		output, err := usecase.Execute(input)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println(output)
		}
	}
}

// {"name": "Curso Vnx", "description":"vnx 2.0","status": "Pending"}