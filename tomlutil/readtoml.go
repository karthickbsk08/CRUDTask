package tomlutil

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	gtoml "github.com/pelletier/go-toml"
	"golang.org/x/term"
)

func GetKeyVal(pData any, key string) string {
	return fmt.Sprintf("%v", pData.(map[string]any)[key])
}

func ReadTomlConfig(fileName string) (f any) {
	// Read the TOML file
	_, err := toml.DecodeFile(fileName, &f)
	if err != nil {
		log.Println(err)
		return f
	}
	return f
}

func ReadTomlMapinAnyType(fileName string, pType *any) {
	// Read the TOML file
	_, err := toml.DecodeFile(fileName, &pType)
	if err != nil {
		log.Println(err)
	}
}

// DecodeTOMLWithTypeCheck decodes TOML file into 'out' if it is a supported pointer type.
func DecodeTOMLWithTypeCheck(filepath string, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("out argument must be a non-nil pointer")
	}

	// Use a type switch on the concrete pointer type.
	switch out.(type) {
	case *struct{}, *map[string]any:
		// These two cover most generic uses (struct or map)
		// Decode directly:
		if _, err := toml.DecodeFile(filepath, out); err != nil {
			return err
		}
		return nil

	// Add cases for pointer to specific types if needed, e.g.,
	case *string, *int, *bool, *float64:
		if _, err := toml.DecodeFile(filepath, out); err != nil {
			return err
		}
		return nil

	default:
		// Unsupported type
		return fmt.Errorf("unsupported type: %T", out)
	}
}

func WriteTomlFile(u, p string, pFilePath string) error {

	if !promptPassword() {
		fmt.Println("Access denied ❌")
		return fmt.Errorf(" Access denied")
	}

	// Prompt for key and value
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter key to update: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	fmt.Print("Enter new value: ")
	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)

	// Load the TOML file
	tree, err := gtoml.LoadFile(pFilePath)
	if err != nil {
		fmt.Println("Error loading file:", err)
		return err
	}
	old := tree.Get(key)

	// Update the Port value
	tree.Set(key, value)

	// Save changes back to file
	err = os.WriteFile(pFilePath, []byte(tree.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Printf(` %s key was updated successfully! into old (%s) ==>  new (%s) `, key, old, value)
	return nil
}

const correctPassword = "123" // 👈 you can load from env or prompt

func promptPassword() bool {
	fmt.Print("Enter password to update config: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		return false
	}
	input := strings.TrimSpace(string(bytePassword))
	return input == correctPassword
}
