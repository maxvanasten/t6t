package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	p "github.com/maxvanasten/gscp/parser"
)

type Ast struct {
	Nodes []p.Node `json:"ast"`
}

type Function struct {
	Name      string
	Arguments []string
}

type Output struct {
	FunctionSignatures []Function
	FunctionCalls      []Function
}

func main() {
	args := os.Args[1:]
	var ast Ast
	if len(args) == 1 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		err = json.Unmarshal(data, &ast)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling json: %v\n", err)
			os.Exit(1)
		}
	} else if len(args) == 2 {
		filePath := args[1]
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&ast); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Usage: t6t -f ast.json OR gscp -p input.gsc | t6t -f")
		os.Exit(1)
	}

	flag := args[0]

	output := Output{}

	switch flag {
	case "-f":
		for _, f := range GetAllFunctions(ast.Nodes) {
			switch f.Type {
			case "function_call":
				output.FunctionCalls = append(output.FunctionCalls, Function{f.Data.FunctionName, ArgumentStrings(f.Children)})
			case "function_declaration":
				output.FunctionSignatures = append(output.FunctionSignatures, Function{f.Data.FunctionName, ArgumentStrings(f.Children[0].Children)})
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown flag: %v\nFlags:\n-f (Full analysis)\n", flag)
		os.Exit(1)
	}

	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Fprintf(os.Stdout, "Error encoding json: %v\n", err)
		os.Exit(1)
	}
}

func ArgumentStrings(arguments []p.Node) []string {
	output := []string{}
	if len(arguments) > 0 {
		for _, a := range arguments {
			str := strings.Builder{}
			switch a.Type {
			case "variable_reference":
				str.WriteString(a.Data.VarName)
			case "number":
				str.WriteString(a.Data.Content)
			case "string":
				str.WriteRune('"')
				str.WriteString(a.Data.Content)
				str.WriteRune('"')
			default:
				str.WriteString("<<<")
				str.WriteString(a.Type)
				str.WriteString(">>>")
			}
			output = append(output, str.String())
		}
	}
	return output
}

func GetAllFunctions(nodes []p.Node) []p.Node {
	functionNames := make(map[string]bool)
	result := []p.Node{}

	for _, node := range nodes {
		if node.Type == "function_call" || node.Type == "function_declaration" {
			if !functionNames[node.Data.FunctionName] {
				functionNames[node.Data.FunctionName] = true
				result = append(result, node)
			}
		}
		if len(node.Children) > 0 {
			result = append(result, GetAllFunctions(node.Children)...)
		}
	}

	return result
}
