package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const (
	VERSION = "Version 0.0.1"
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

func learningStageList(flashcards []Flashcard) []int {
	list := make([]int, 5)
	for i := 0; i < len(flashcards); i++ {
		if flashcards[i].LearningStage < 5 {
			list[flashcards[i].LearningStage]++
		}
	}
	return list
}

func Create() {
	fmt.Println("Create")
}
func Import() {
	fmt.Println("Import")
}
func Edit() {
	fmt.Println("Edit")
}
func Read(name string) {
	set, err := openFile(name)
	if err == nil {
		for i := 0; i < len(set.Flashcards); i++ {
			fmt.Println(set.Flashcards[i].Q)
			fmt.Println()
			fmt.Println(set.Flashcards[i].Help)
			fmt.Println()
			fmt.Println(set.Flashcards[i].A)
		}
	} else {
		fmt.Println("Error opening " + name)
		fmt.Println(err)
	}
}
func Learn(name string) {
	fmt.Println("Select mode [W]rite/[v]iew: ")
	var mode string
	fmt.Scanln(&mode)
	set, err := openFile(name)
	if err == nil {
		if mode == "v" || mode == "V" {
			nums := randList(len(set.Flashcards))
			for i := 0; i < len(set.Flashcards); i++ {
				fmt.Println(set.Flashcards[nums[i]].Q)
				fmt.Println("\nShow help? [y]/[N]")
				var help string
				fmt.Scanln(&help)
				if help == "y" || help == "Y" {
					fmt.Println(set.Flashcards[nums[i]].Help)
				}
				fmt.Println("\nPress enter to show answer")
				fmt.Scanln()
				fmt.Println(set.Flashcards[nums[i]].A)
				fmt.Println("\nPress enter to continue")
				fmt.Scanln()
			}
		} else {
			for i := 0; i < 5; i++ {
				nums := randList(len(set.Flashcards))
				for j := 0; j < len(set.Flashcards); j++ {
					if set.Flashcards[nums[i]].LearningStage == i {
						fmt.Println(set.Flashcards[nums[j]].Q)
						fmt.Println("\nShow help? [y]/[N]")
						var help string
						fmt.Scanln(&help)
						if help == "y" || help == "Y" {
							fmt.Println(set.Flashcards[nums[j]].Help)
						}
						fmt.Println("\nWrite answer and press enter")
						var answer string
						fmt.Scanln(&answer)
						if answer == set.Flashcards[nums[j]].A {
							set.Flashcards[nums[j]].LearningStage++
							fmt.Println("Correct")
						} else {
							fmt.Println("Wrong")
						}
						fmt.Println("\nPress enter to continue")
						fmt.Scanln()
					}
				}
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
		fmt.Println("Reset " + name)
	} else {
		fmt.Println("Error resetting " + name)
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
		Create()
	case "-i":
		Import()
	case "-e":
		Edit()
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
