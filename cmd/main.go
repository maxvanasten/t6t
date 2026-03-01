package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"io"

	p "github.com/maxvanasten/gscp/parser"
)

type Ast struct {
	Nodes []p.Node `json:"ast"`
}

type FunctionSignature struct {
	Name      string
	Arguments []string
}

type FunctionCall struct {
	Name      string
	Arguments []string
}

type Output struct {
	FunctionSignatures  []FunctionSignature
	FunctionCalls       []FunctionCall
	VariableAssignments []string
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
		fmt.Fprintln(os.Stderr, "Usage: t6t -f ast.json")
		os.Exit(1)
	}

	flag := args[0]

	output := Output{}

	switch flag {
	case "-f":
		for _, fd := range GetAll(ast.Nodes, "function_declaration") {
			functionSignature := FunctionSignature{fd.Data.FunctionName, ArgumentStrings(fd.Children[0].Children)}

			output.FunctionSignatures = append(output.FunctionSignatures, functionSignature)
		}
		for _, fc := range GetAll(ast.Nodes, "function_call") {
			functionCall := FunctionCall{fc.Data.FunctionName, ArgumentStrings(fc.Children)}

			output.FunctionCalls = append(output.FunctionCalls, functionCall)
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

func GetAll(nodes []p.Node, identifier string) []p.Node {
	result := []p.Node{}

	for _, node := range nodes {
		if node.Type == identifier {
			result = append(result, node)
		}
		if len(node.Children) > 0 {
			result = append(result, GetAll(node.Children, identifier)...)
		}
	}

	return result
}
