package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	VERSION = "Version 1.0.0"
)

type Flashcard struct {
	Q    string `json:"question"`
	A    string `json:"answer"`
	Help string `json:"help"`
	// QMedia
	// AMedia
	LearningStage int `json:"learningStage"`
}

type Set struct {
	// Encrypted bool
	Flashcards []Flashcard `json:"flashcards"`
}

func openFile(name string) (Set, error) {
	var set Set
	file, err := os.Open(name)
	if err != nil {
		return set, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&set)
	/*if set.Encrypted {
		// Decrypt
	}*/
	return set, err
}

func randList(length int) []int {
	list := make([]int, length)
	for i := 0; i < length; i++ {
		list[i] = i
	}
	for i := 0; i < length; i++ {
		j := rand.Intn(length)
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func edit(set Set) Set {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		for i := 0; i < len(set.Flashcards); i++ {
			fmt.Println(i + 1)
			fmt.Println("Question: ")
			fmt.Println(set.Flashcards[i].Q)
			fmt.Println("Answer: ")
			fmt.Println(set.Flashcards[i].A)
			fmt.Println("Help: ")
			fmt.Println(set.Flashcards[i].Help)
			fmt.Println()
		}
		fmt.Println("Select flashcard, 0 to add new, -1 to exit: ")
		scanner.Scan()
		selectionStr := scanner.Text()
		selection, err := strconv.Atoi(selectionStr)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid integer.")
			continue
		}
		if selection == -1 {
			break
		}
		if selection == 0 {
			var flashcard Flashcard
			fmt.Println("Question: ")
			scanner.Scan()
			flashcard.Q = scanner.Text()
			fmt.Println("Answer: ")
			scanner.Scan()
			flashcard.A = scanner.Text()
			fmt.Println("Help: ")
			scanner.Scan()
			flashcard.Help = scanner.Text()
			set.Flashcards = append(set.Flashcards, flashcard)
		} else {
			fmt.Println("Select action [E]dit/[d]elete: ")
			scanner.Scan()
			action := scanner.Text()
			if action == "d" || action == "D" {
				set.Flashcards = append(set.Flashcards[:selection-1], set.Flashcards[selection:]...)
				continue
			}
			fmt.Println("Select field [q]uestion/[a]nswer/[h]elp/[Enter]all: ")
			scanner.Scan()
			field := scanner.Text()
			if field == "q" || field == "Q" {
				fmt.Println("Question: ")
				scanner.Scan()
				set.Flashcards[selection-1].Q = scanner.Text()
			} else if field == "a" || field == "A" {
				fmt.Println("Answer: ")
				scanner.Scan()
				set.Flashcards[selection-1].A = scanner.Text()
			} else if field == "h" || field == "H" {
				fmt.Println("Help: ")
				scanner.Scan()
				set.Flashcards[selection-1].Help = scanner.Text()
			} else {
				fmt.Println("Question: ")
				scanner.Scan()
				set.Flashcards[selection-1].Q = scanner.Text()
				fmt.Println("Answer: ")
				scanner.Scan()
				set.Flashcards[selection-1].A = scanner.Text()
				fmt.Println("Help: ")
				scanner.Scan()
				set.Flashcards[selection-1].Help = scanner.Text()
			}
		}
	}
	return set
}

func checkAnswer(answer string, correct string) bool {
	answer = strings.ReplaceAll(strings.TrimSpace(answer), " ", "")
	correct = strings.ReplaceAll(strings.TrimSpace(correct), " ", "")
	if answer == correct {
		return true
	}
	if strings.Contains(correct, "(") && strings.Contains(correct, ")") {
		correctWithoutBrackets := strings.Split(correct, "(")[0] + strings.Split(correct, ")")[1]
		if answer == correctWithoutBrackets {
			return true
		}
	}
	if strings.Contains(correct, "/") {
		for _, correctAnswer := range strings.Split(correct, "/") {
			if answer == correctAnswer {
				return true
			}
		}
	}
	if strings.Contains(correct, "|") {
		for _, correctAnswer := range strings.Split(correct, "|") {
			if answer == correctAnswer {
				return true
			}
		}
	}
	return false
}

func Create(name string) {
	if name[len(name)-5:] != ".ofmp" {
		name += ".ofmp"
	}
	var set Set
	set.Flashcards = []Flashcard{}
	file, err := os.Create(name)
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.Encode(set)
		fmt.Println("Created " + name)
	} else {
		fmt.Println("Error creating " + name)
		fmt.Println(err)
	}
}
func Import() {
	fmt.Println("Import not yet implemented")
}
func Edit(name string) {
	file, err := os.Open(name)
	if err == nil {
		var set Set
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&set)
		if err == nil {
			set = edit(set)
			file, err = os.Create(name)
			if err == nil {
				defer file.Close()
				encoder := json.NewEncoder(file)
				encoder.Encode(set)
				fmt.Println("Edited " + name)
			} else {
				fmt.Println("Error creating " + name)
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error decoding " + name)
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error opening " + name)
		fmt.Println(err)
	}
}
func Read(name string) {
	set, err := openFile(name)
	if err == nil {
		for i := 0; i < len(set.Flashcards); i++ {
			fmt.Println(set.Flashcards[i].Q)
			fmt.Println()
			if set.Flashcards[i].Help != "" {
				fmt.Println(set.Flashcards[i].Help)
				fmt.Println()
			}
			fmt.Println(set.Flashcards[i].A)
			fmt.Println()
			fmt.Println()
		}
	} else {
		fmt.Println("Error opening " + name)
		fmt.Println(err)
	}
}
func Learn(name string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Select mode [W]rite/[v]iew: ")
	scanner.Scan()
	mode := scanner.Text()
	set, err := openFile(name)
	if err == nil {
		if mode == "v" || mode == "V" {
			nums := randList(len(set.Flashcards))
			for i := 0; i < len(set.Flashcards); i++ {
				fmt.Println(set.Flashcards[nums[i]].Q)
				if set.Flashcards[nums[i]].Help != "" {
					fmt.Println("\nShow help? [y]es/[N]o")
					scanner.Scan()
					help := scanner.Text()
					if help == "y" || help == "Y" {
						fmt.Println(set.Flashcards[nums[i]].Help)
					}
				}
				fmt.Println("\nPress enter to show answer")
				scanner.Scan()
				fmt.Println(set.Flashcards[nums[i]].A)
				fmt.Println("\nPress enter to continue")
				scanner.Scan()
			}
		} else {
			fmt.Println("Select learning stage [1]st/[2]nd/[3]rd/[4]th/[5]th: ")
			scanner.Scan()
			stageStr := scanner.Text()
			stage, convErr := strconv.Atoi(stageStr)
			if convErr == nil {
				if stage < 1 || stage > 5 {
					fmt.Println("Invalid stage")
					return
				}
				var fullBool bool
				fmt.Println("Full answers? [y]es/[N]o")
				scanner.Scan()
				full := scanner.Text()
				if full == "y" || full == "Y" {
					fullBool = true
				} else {
					fullBool = false
				}
				var cards int
				nums := randList(len(set.Flashcards))
				for i := 0; i < len(set.Flashcards); i++ {
					helpUsed := false
					if set.Flashcards[nums[i]].LearningStage == stage-1 {
						cards++
						fmt.Println("\n" + set.Flashcards[nums[i]].Q)
						if set.Flashcards[nums[i]].Help != "" {
							fmt.Println("\nShow help? [y]es/[N]o")
							scanner.Scan()
							help := scanner.Text()
							if help == "y" || help == "Y" {
								helpUsed = true
								fmt.Println(set.Flashcards[nums[i]].Help)
							}
						}
						fmt.Println("\nWrite answer and press enter")
						scanner.Scan()
						answer := scanner.Text()
						if answer == set.Flashcards[nums[i]].A {
							if !helpUsed {
								set.Flashcards[nums[i]].LearningStage++
								fmt.Println("Correct")
							} else {
								fmt.Println("Correct, but you used help")
							}
						} else if checkAnswer(answer, set.Flashcards[nums[i]].A) && !fullBool {
							if !helpUsed {
								set.Flashcards[nums[i]].LearningStage++
								fmt.Println("Correct")
								fmt.Println("Full answer: " + set.Flashcards[nums[i]].A)
							} else {
								fmt.Println("Correct, but you used help")
								fmt.Println("Full answer: " + set.Flashcards[nums[i]].A)
							}
						} else {
							fmt.Println("Wrong")
							fmt.Println("Correct answer: " + set.Flashcards[nums[i]].A)
						}
						fmt.Println("\nPress enter to continue")
						scanner.Scan()
					}
				}
				if cards == 0 {
					fmt.Println("There are no cards left in this learning stage.")
				} else {
					fmt.Println("Congratulations! You went through all cards in this learning stage.")
				}
				file, openErr := os.OpenFile(name, os.O_WRONLY, 0644)
				if openErr == nil {
					defer file.Close()
					encoder := json.NewEncoder(file)
					encoder.Encode(set)
					fmt.Println("Saved " + name)
				} else {
					fmt.Println("Error saving " + name)
					fmt.Println(err)
				}
			} else {
				fmt.Println("Invalid stage")
			}
		}
	} else {
		fmt.Println("Error opening " + name)
		fmt.Println(err)
	}
}
func Delete(name string) {
	err := os.Remove(name)
	if err == nil {
		fmt.Println("Deleted " + name)
	} else {
		fmt.Println("Error deleting " + name)
		fmt.Println(err)
	}
}
func Reset(name string) {
	set, err := openFile(name)
	if err == nil {
		for i := 0; i < len(set.Flashcards); i++ {
			set.Flashcards[i].LearningStage = 0
		}
		file, openErr := os.OpenFile(name, os.O_WRONLY, 0644)
		if openErr == nil {
			defer file.Close()
			encoder := json.NewEncoder(file)
			encoder.Encode(set)
			fmt.Println("Saved " + name)
		} else {
			fmt.Println("Error saving " + name)
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error opening " + name)
		fmt.Println(err)
	}
}
func Version() {
	fmt.Println(VERSION)
}
func Help() {
	fmt.Println("-c\tCreate a new flashcard set")
	fmt.Println("-i\tImport flashcards from a file")
	fmt.Println("-e\tEdit a flashcard set")
	fmt.Println("-r\tRead a flashcard set")
	fmt.Println("-l\tLearn a flashcard set")
	fmt.Println("-d\tDelete a flashcard set")
	fmt.Println("-v\tShow the version of OFMP")
	fmt.Println("-h\tShow help")
	fmt.Println("--reset\tReset Flashcard set learning stage")
	fmt.Println("--version\tShow the version of OFMP")
	fmt.Println("--help\tShow extended help")
	fmt.Println("--info\tShow info about OFMP")
}
func ExtendedHelp() {
	fmt.Println("OFMP is an open source flashcard management program")
	fmt.Println("Syntax: ofmp [command] [options]")
	fmt.Println("Commands:")
	fmt.Println("\t-c\tCreate a new flashcard set\n\t\tofmp -c [name]")
	fmt.Println("\t-i\tImport flashcards from a file\n\t\tofmp -i [file]")
	fmt.Println("\t-e\tEdit a flashcard set\n\t\tofmp -e [name]")
	fmt.Println("\t-r\tRead a flashcard set\n\t\tofmp -r [name]")
	fmt.Println("\t-l\tLearn a flashcard set\n\t\tofmp -l [name]")
	fmt.Println("\t-d\tDelete a flashcard set\n\t\tofmp -d [name]")
	fmt.Println("\t-v\tShow the version of OFMP\n\t\tofmp -v")
	fmt.Println("\t-h\tShow help\n\t\tofmp -h")
	fmt.Println("\t--reset\tReset Flashcard set learning stage\n\t\tofmp --reset [name]")
	fmt.Println("\t--version\tShow the version of OFMP\n\t\tofmp --version")
	fmt.Println("\t--help\tShow extended help\n\t\tofmp --help")
	fmt.Println("\t--info\tShow info about OFMP\n\t\tofmp --info")
}
func Info() {
	fmt.Println("Open Flashcard Management Program by Jeroen Leuenberger")
	fmt.Println(VERSION)
	fmt.Println("https://github.com/Jero075/OFMP")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command specified\nUse '-h' for help")
		return
	}
	switch os.Args[1] {
	case "-c":
		if len(os.Args) < 3 {
			fmt.Println("No file specified\nUse '-h' for help")
			return
		}
		Create(os.Args[2])
	case "-i":
		Import()
	case "-e":
		if len(os.Args) < 3 {
			fmt.Println("No file specified\nUse '-h' for help")
			return
		}
		Edit(os.Args[2])
	case "-r":
		if len(os.Args) < 3 {
			fmt.Println("No file specified\nUse '-h' for help")
			return
		}
		Read(os.Args[2])
	case "-l":
		if len(os.Args) < 3 {
			fmt.Println("No file specified\nUse '-h' for help")
			return
		}
		Learn(os.Args[2])
	case "-d":
		if len(os.Args) < 3 {
			fmt.Println("No file specified\nUse '-h' for help")
			return
		}
		Delete(os.Args[2])
	case "-v":
		Version()
	case "-h":
		Help()
	case "--reset":
		if len(os.Args) < 3 {
			fmt.Println("No file specified\nUse '-h' for help")
			return
		}
		Reset(os.Args[2])
	case "--version":
		Version()
	case "--help":
		ExtendedHelp()
	case "--info":
		Info()
	default:
		fmt.Println("Not a valid command\nUse '-h' for help")
	}
}
