package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Arguments map[string]string
type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	err := ValidateInputs(args)

	if err != nil {
		return err
	}

	if args["operation"] == "list" {
		err = ReadAndWriteFile(args, writer)
	}

	if args["operation"] == "add" {
		err := AddToFile(args, writer)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if args["operation"] == "findById" {
		err = findById(args, writer)
		if err != nil {
			return err
		}
	}

	if args["operation"] == "remove" {
		err = removeById(args, writer)
		if err != nil {
			return err
		}
	}
	return nil
}
func removeById(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users, ansUsers []User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return err
	}
	fmt.Println(users)
	isFound := false
	for _, v := range users {

		if v.Id == args["id"] {
			b, err := json.Marshal(v)
			if err != nil {
				fmt.Println("error:", err)
			}
			writer.Write(b)
			isFound = true
			continue
		}
		ansUsers = append(ansUsers, v)
	}

	if !isFound {
		writer.Write([]byte("Item with id " + args["id"] + " not found"))
	}
	if err = file.Truncate(0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
	file.Seek(0, 0)

	b, err := json.Marshal(ansUsers)
	if err != nil {
		return err
	}
	file.Write(b)
	return nil
}

func findById(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	var users []User
	json.Unmarshal(bytes, &users)
	for _, v := range users {
		if v.Id == args["id"] {
			b, err := json.Marshal(v)
			if err != nil {
				return nil
			}
			writer.Write(b)
			return nil
		}
	}

	return nil
}

func AddToFile(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	var userToAdd User
	json.Unmarshal(bytes, &users)
	json.Unmarshal([]byte(args["item"]), &userToAdd)

	for _, v := range users {
		if v.Id == userToAdd.Id {
			writer.Write([]byte("Item with id " + v.Id + " already exists"))
			return nil
		}
	}

	users = append(users, userToAdd)

	b, err := json.Marshal(users)
	if err != nil {
		return nil
	}

	if err = file.Truncate(0); err != nil {
		return err
	}
	file.Seek(0, 0)

	file.Write(b)
	writer.Write(b)

	return nil
}

func ReadAndWriteFile(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	writer.Write(bytes)

	return nil
}

func ValidateInputs(args Arguments) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}
	if args["operation"] != "list" && args["operation"] != "add" && args["operation"] != "remove" && args["operation"] != "findById" {
		return errors.New("Operation " + args["operation"] + " not allowed!")
	}
	if args["operation"] == "add" && args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}
	if (args["operation"] == "findById" || args["operation"] == "remove") && args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}
	return nil
}

func main() {
	var operationFlag = flag.String("operation", "", "help message for flag operation")
	var itemFlag = flag.String("item", "", "help message for flag item")
	var fileNameFlag = flag.String("fileName ", "", "help message for flag file name")
	flag.Parse()
	args := parseArgs(*operationFlag, *itemFlag, *fileNameFlag)
	err := Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs(operation, item, fileName string) Arguments {
	ans := Arguments{"operation": operation, "item": item, "filename": fileName}

	return ans
}
