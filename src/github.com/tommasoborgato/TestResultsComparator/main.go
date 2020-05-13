package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type failedTestsStruct struct {
	lines []string
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getGroupsFromLines(lines []string) map[string][]string {
	var groups = make(map[string][]string)
	var testNameArr []string

	reGroupName := regexp.MustCompile(`^jdk=.*`)
	reTestName := regexp.MustCompile(`org\.[\.\w]+`)
	groupName := ""
	testName := ""

	for _, line := range lines {
		// find the group name
		if reGroupName.FindString(line) != "" {
			groupName = reGroupName.FindString(line)
		}

		testNameArr = reTestName.FindStringSubmatch(line)
		if len(testNameArr) > 0 {
			testName = testNameArr[0]

			if groupName != "" {
				//fmt.Println(i, testNameArr)
				groups[groupName] = append(groups[groupName], testName)
			} else {
				log.Fatalf("Test with no group: %s", testName)
			}
		}
	}

	return groups
}

func main() {
	file1 := "/home/tborgato/projects/TestResultsComparator/workspace/FilesToCompare/XP.txt"
	file2 := "/home/tborgato/projects/TestResultsComparator/workspace/FilesToCompare/7_3_1.txt"
	fmt.Printf("%s VS %s\n", file1, file2)

	// read file 1
	lines1, err := readLines(file1)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// read file 2
	lines2, err := readLines(file2)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// group failed tests

	// XP failed tests
	group1 := getGroupsFromLines(lines1)
	// 7.3.1 failed tests
	group2 := getGroupsFromLines(lines2)

	fmt.Println(len(group1))
	fmt.Println(len(group2))

	// find XP tests that dind't fail in 7.3.1
	fmt.Println("=========================================================")
	for k1, v1 := range group1 {
		for _, tst1 := range v1 {
			tstFound := false
			for k2, v2 := range group2 {
				// same permutation
				if k1 == k2 {
					for _, tst2 := range v2 {
						if tst1 == tst2 {
							tstFound = true
							//break
						}
					}
					//break
				}
			}
			if !tstFound {
				fmt.Printf("XP %s %s\n", k1, tst1)
			}
		}
	}
	fmt.Println("=========================================================")
}
