package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║                       CSV-DIFF                        ║")
	fmt.Println("╠═══════════════════════════════════════════════════════╣")
    	fmt.Println("║  Das Programm erwartet zwei csv Dateien als Eingabe.  ║")
	fmt.Println("║ Die Differenz wird im selben Verzeichnis gespeichert. ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	
	reader := bufio.NewReader(os.Stdin)

	file1, err := getFilePath(reader, "Erste csv-Datei: ")
	if err != nil {
		log.Fatal(err)
	}

	file2, err := getFilePath(reader, "Zweite csv-Datei: ")
	if err != nil {
		log.Fatal(err)
	}

	outputFile := "output_" + time.Now().Format("2006-01-02_15-04-05") + ".csv"

	records1, err := readCSV(file1)
	if err != nil {
		log.Fatal(err)
	}

	records2, err := readCSV(file2)
	if err != nil {
		log.Fatal(err)
	}

	err = writeCSV(outputFile, findUniqueRecords(records1, records2))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Vergleich abgeschlossen.\nErgebnisse in", outputFile, "gespeichert.\n")

	// Damit das Fenster beim Ausführen über Explorer offen bleibt.
	fmt.Println("\nNow please press the Any-Key and get me some coffee!")
	reader.ReadString('\n')
	os.Exit(0)
}

func getFilePath(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	path, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	path = strings.TrimSpace(path)

	if path == "" {
		return "", fmt.Errorf("Bitte geben Sie einen gültigen Dateipfad ein\n")
	}

	if !isValidPath(path) {
		return "", fmt.Errorf("Der angegebene Pfad ist ungültig\n")
	}

	return path, nil
}

func isValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Öffnen von %s: %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Lesen von %s: %v", filename, err)
	}

	return records, nil
}

func findUniqueRecords(records1, records2 [][]string) [][]string {
	var uniqueRecords [][]string

	for _, record := range records1 {
		if !containsRecord(records2, record) {
			uniqueRecords = append(uniqueRecords, record)
		}
	}

	for _, record := range records2 {
		if !containsRecord(records1, record) {
			uniqueRecords = append(uniqueRecords, record)
		}
	}

	return uniqueRecords
}

func containsRecord(records [][]string, record []string) bool {
	for _, r := range records {
		if equalRecords(r, record) {
			return true
		}
	}
	return false
}

func equalRecords(record1, record2 []string) bool {
	if len(record1) != len(record2) {
		return false
	}

	for i, value := range record1 {
		if value != record2[i] {
			return false
		}
	}

	return true
}

func writeCSV(filename string, records [][]string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("Fehler beim Erstellen von %s: %v", filename, err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    writer.UseCRLF = false
    err = writer.WriteAll(records)
    if err != nil {
        return fmt.Errorf("Fehler beim Schreiben in %s: %v", filename, err)
    }

    return nil
}
