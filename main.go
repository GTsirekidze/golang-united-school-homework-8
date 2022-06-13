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
	"strconv"
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
		err = ReadAndWriteFile(writer)
	}

	if args["operation"] == "add" {
		err = AddToFile(args, writer)
		if err != nil {
			return err
		}
		err = ReadAndWriteFile(writer)
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
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	json.Unmarshal(bytes, &users)
	ans := "[{"
	for _, v := range users {
		ans += "\"id\":" + v.Id + ",\"email\":\"" + v.Email + "\",\"age\":" + strconv.Itoa(v.Age) + "},"
		if v.Id == args["id"] {
			b, err := json.Marshal(v)
			if err != nil {
				fmt.Println("error:", err)
			}
			writer.Write(b)
		}
	}
	ans += "}]"
	fmt.Println(ans)
	if err = file.Truncate(0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
	file.Seek(0, 0)

	file.Write([]byte(ans))
	return nil
}

func findById(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	json.Unmarshal(bytes, &users)
	for _, v := range users {
		if v.Id == args["id"] {
			b, err := json.Marshal(v)
			if err != nil {
				fmt.Println("error:", err)
			}
			writer.Write(b)
			return nil
		}
	}
	return nil
}

func AddToFile(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0755)
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
	ans := "[{"
	isEmpty := len(users) == 0
	for _, v := range users {
		ans += "\"id\":" + v.Id + ",\"email\":\"" + v.Email + "\",\"age\":" + strconv.Itoa(v.Age) + "},"
		if v.Id == userToAdd.Id {
			return errors.New("Item with id " + v.Id + " already exists")
		}
	}
	if !isEmpty {
		ans += "{"
	}
	ans += "\"id\":" + userToAdd.Id + ",\"email\":\"" + userToAdd.Email + "\",\"age\":" + strconv.Itoa(userToAdd.Age) + "}]"
	fmt.Println(ans)
	if err = file.Truncate(0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
	file.Seek(0, 0)

	file.Write([]byte(ans))

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func ReadAndWriteFile(writer io.Writer) error {
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	writer.Write(bytes)

	if err = file.Close(); err != nil {
		return err
	}
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
