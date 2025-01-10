# Creating a Plugin in Korin

Korin allows you to extend its functionality by creating custom plugins.
Plugins can analyze, modify, or extend the code dynamically during preprocessing. 
This guide will walk you through the process of creating a plugin for Korin.

Before getting started, to understand how plugins work, we have to understand first the labeling mechanism and how 
Korin understands Golang code.

Korin uses a labeling system to identify and categorize different parts of the code during preprocessing. 
This system helps plugins to analyze and manipulate the code efficiently. Here is a list of all the labels that currently 
exist in Korin, please check [`klabels/declarations.go`](../pkg/klabels/declarations.go) for a more updated list.

### Comprehensive List of Labels

1. **VariableKind**:
    - **Purpose**: Identifies variable declarations.
    - **Usage**: Used to label lines where variables are declared, allowing plugins to process these variables.
    - **Note**: A `VariableKind` may point to a `const` or a `var`, and these can be identified by checking the labels to see if there is another 
   label of either `ConstKind` or `VarKind`. Additionally, you can figure out whether this is a multi-line const or var by using the `AnalysisHelper`'s 
   `CheckMultiLineConstOrVar` function as seen on [k:env](../pkg/kplugins/plugin_env.go).
    - **Data Type**: `[]klabels.VariableDeclaration`

2. **FunctionKind**:
    - **Purpose**: Identifies function declarations.
    - **Usage**: Used to label lines where functions are declared, enabling plugins to analyze function scopes and signatures.
    - **Note**: Anonymous functions are identified by having the `Name` as an empty string.
    - **Data Type**: `klabels.FunctionDeclaration`

3. **ConstScopeBeginKind**:
    - **Purpose**: Marks the beginning of a constant scope.
    - **Usage**: Helps plugins to identify and process constant declarations within a specific scope.
    - **Data Type**: `nil`

4. **ConstScopeEndKind**:
    - **Purpose**: Marks the end of a constant scope.
    - **Usage**: Helps plugins to identify the end of a constant declaration scope.
    - **Data Type**: `nil`

5. **PackageKind**:
    - **Purpose**: Identifies package declarations.
    - **Usage**: Used to label lines where package statements are declared, enabling plugins to process package information.
    - **Data Type**: `klabels.PackageDeclaration`

6. **TypeDeclarationKind**:
    - **Purpose**: Identifies type declarations.
    - **Usage**: Used to label lines where types are declared, allowing plugins to process these types.
    - **Note**: You can identify whether it's a struct, interface, or a type alias by checking the `Kind` field of the `TypeDeclaration`.
    - **Data Type**: `klabels.TypeDeclaration`
   
7. **FieldDeclarationKind**:
    - **Purpose**: Identifies field declarations, these fields may be from struct or interface.
    - **Usage**: Used to label lines where fields are declared, allowing plugins to process these fields.
    - **Data Type**: `[]klabels.FieldDeclaration`

8. **CommentKind**:
    - **Purpose**: Identifies comments within the code.
    - **Usage**: Used to label lines where comments are present, enabling plugins to process or ignore comments.
    - **Data Type**: `string`

9. **ReturnKind**:
    - **Purpose**: Identifies return statements.
    - **Usage**: Used to label lines where return statements are present, allowing plugins to process these return values.
    - **Data Type**: `klabels.ReturnStatement`

10. **ScopeBeginKind**:
    - **Purpose**: Marks the beginning of a scope.
    - **Usage**: Helps plugins to identify the start of a scope.
    - **Data Type**: `nil`

11. **ScopeEndKind**:
    - **Purpose**: Marks the end of a scope.
    - **Usage**: Helps plugins to identify the end of a scope.
    - **Data Type**: `nil`

12. **VarScopeBeginKind**:
    - **Purpose**: Marks the beginning of a variable scope.
    - **Usage**: Helps plugins to identify the start of a variable declaration scope.
    - **Data Type**: `nil`

13. **VarScopeEndKind**:
    - **Purpose**: Marks the end of a variable scope.
    - **Usage**: Helps plugins to identify the end of a variable declaration scope.
    - **Data Type**: `nil`

14. **ConstDeclarationKind**:
    - **Purpose**: Identifies constant declarations.
    - **Usage**: Used to label lines where constants are declared, allowing plugins to process these constants.
    - **Note**: This is always accompanied by a `VariableKind` which has all the information about the constant.
    - **Data Type**: `nil`

15. **VarDeclarationKind**:
    - **Purpose**: Identifies var declarations.
    - **Usage**: Used to label lines where constants are declared, allowing plugins to process these vars.
    - **Note**: This is always accompanied by a `VariableKind` which has all the information about the vars.
    - **Data Type**: `nil`


By understanding and utilizing these labels, you can create powerful plugins that can analyze and manipulate code effectively during preprocessing.

### Labeling System

Korin's labeling system works by analyzing each line of code and assigning appropriate labels based on the content. 
These labels are then used by plugins to perform specific actions. 
Here is a brief overview of how the labeling system works:

1. **Analysis**:
    - Each line of code is analyzed to determine its type (e.g., variable declaration, function declaration).
    - Labels are assigned based on the analysis, although multiple labels can be assigned to a single line, such as 
   if there is a comment, or if the line is a `const` or a `var` declaration.

2. **Label Data**:
    - Labels contain data relevant to the type of code they represent. 
    - For example, a `VariableKind` label will be of type `[]VariableDeclarations` which contains information 
   about the variable's name, type, and value. A more comprehensive list of the `label.Kind` and it's associated `label.Data`'s 
   data type can be found under [Comprehensive List of Labels](#comprehensive-list-of-labels).

By understanding and utilizing these labels, you can create powerful plugins that can analyze and manipulate code effectively during preprocessing.

## Plugin Structure

A Korin plugin is a Go package that implements the `Plugin` interface. 
The interface requires the `Process` method, which is called for each line of code during preprocessing.

### Example Plugin

Here is a basic example of a Korin plugin:

```go
package myplugin

import (
    "github.com/ShindouMihou/korin/pkg/klabels"
    "github.com/ShindouMihou/korin/pkg/kplugins"
)

type MyPlugin struct {
	Plugin
}

func (p MyPlugin) Name() string {
	return "PluginName"
}

func (p MyPlugin) Group() string {
	return "group.plugin"
}

func (p MyPlugin) Version() string {
	return "1.0.0"
}

// Context is used to keep some information within the same file, this is important for plugins that requires some sort of 
// multi-line knowledge, such as, plugins that annotates an entire type struct. 
//
// A best example of this is found in the `k:named` plugin where it annotates the entire struct with the specified tags.
func (p MyPlugin) Context(file string) *any {
	return nil
}

func (p *MyPlugin) FreeContext(file string) {
	// Free the context here which is done when the file is finished.
}

func (p MyPlugin) Process(line string, index int, headers *kplugins.Headers, stack []klabels.Analysis) (string, error) {
    // Your plugin logic here
    return "", nil // Returning "" will keep the line unchanged, recommended to use over returning line itself.
}
```

## Using Helpers

Korin provides several helper types to assist with plugin development:

### `SyntaxHelper`

`SyntaxHelper` provides utility functions for working with Go syntax. It helps in identifying and manipulating different parts of the code.

### `AnalysisHelper`

`AnalysisHelper` assists in analyzing the code. It provides functions to extract and interpret various code elements, making it easier to understand the structure and content of the code being processed.

### `ReadHelper`

`ReadHelper` offers functions to read and interpret labels. It helps in identifying specific labels within the code, which can be used to trigger specific plugin actions.

## Example: Using Helpers in a Plugin

Here is an example of a plugin that uses `ReadHelper` to identify and process specific labels:

```go
package myplugin

import (
    "github.com/ShindouMihou/korin/pkg/klabels"
    "github.com/ShindouMihou/korin/pkg/kplugins"
)

type MyPlugin struct{}

func (p MyPlugin) Process(line string, index int, headers *kplugins.Headers, stack []klabels.Analysis) (string, error) {
	analysis := stack[index]
	return kplugins.ReadHelper.Require(klabels.VariableKind, analysis.Labels, func(label klabels.Label) (string, error) {
		variables := (label).Data.([]klabels.VariableDeclaration)
		// Process the variables that were labeled.
        return line, nil // Only return `line` when there is changes to prevent performance issues.
    })
}
```

To learn more about the plugins, simply check your autocomplete for the `kplugins` package:
- `kplugins.ReadHelper`
- `kplugins.SyntaxHelper`
- `kplugins.AnalysisHelper`

You can also read upon the source code of the native plugins itself:
- [k:env](../pkg/kplugins/plugin_env.go): Automatically sets the value of a variable to its corresponding environment value.
- [k:float](../pkg/kplugins/plugin_error_propogation.go): Automatically floats the error up the stack.
- [k:println](../pkg/kplugins/plugin_println.go): Automatically prints the value of a variable to the console.
- [k:named](../pkg/kplugins/plugin_annotate_serializers.go): Automatically annotates the struct fields with the specified tags for naming, used for JSON, YAML, etc.


## Registering Your Plugin

To use your plugin, you need to register it with Korin. This is done by adding your plugin to the list of plugins in the `main` function:

```go
package main

import (
    "github.com/ShindouMihou/korin/pkg/korin"
    "path/to/your/plugin"
)

func main() {
    korin := korin.New()
    korin.Plugins = append(korin.Plugins, &plugin.MyPlugin{})
    korin.Run("cmd/app.go")  // Replace with your actual entry point
}
```

## Conclusion

Creating a plugin in Korin is straightforward and allows you to extend the functionality of the preprocessor. 
By leveraging the provided helpers, you can efficiently analyze and manipulate the code during preprocessing. Happy coding!